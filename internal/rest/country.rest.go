package rest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

//-----------------------------------------------------------------------------
//
//	GET		Get List 			"/country"
//	GET		Get by Id			"/country/id/{id}"
//	GET		Get by Shortcode	"/country/shortcode/{shortcode}"
//	POST	Create				"/country"
//	DELETE	Delete				"/country/{id}"
//	PUT		Update				"/country/{id}"
//
//-----------------------------------------------------------------------------

type CountryListItem struct {
	ID        uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode string    `json:"shortcode" minLength:"2" maxLength:"2" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
}

type CountryNoPK struct {
	Shortcode  string `json:"shortcode" minLength:"2" maxLength:"2" doc:"A unique short name for this data"`
	Name       string `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
	Flags      int16  `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
	Ibanlength *int16 `json:"ibanlenth,omitempty" format:"int16" minimum:"0" doc:"The exact length of the IBAN required in that country"`
	Risktype   int16  `json:"risktype" format:"int16" minimum:"0" doc:"This risk profile of this country"`
}

type Country struct {
	ID         uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode  string    `json:"shortcode" minLength:"2" maxLength:"2" doc:"A unique short name for this data"`
	Name       string    `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
	Flags      int16     `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
	Ibanlength *int16    `json:"ibanlenth,omitempty" format:"int16" minimum:"0" doc:"The exact length of the IBAN required in that country"`
	Risktype   int16     `json:"risktype" format:"int16" minimum:"0" doc:"This risk profile of this country"`
}

//-----------------------------------------------------------------------------

type CountryRequestId struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type CountryRequestShortcode struct {
	Shortcode string `path:"shortcode" minLength:"2" maxLength:"2" doc:"This is the unique identifier a this data"`
}

type CountryRequestCreate struct {
	Body CountryNoPK
}

type CountryRequestUpdate struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Body CountryNoPK
}

//-----------------------------------------------------------------------------

type CountryResponseCreate struct {
	Header struct {
		Location string `header:"Location" doc:"URL of the newly created entity"`
	}

	Body Country
}

type CountryResponse struct {
	Body Country
}

type CountryListResponse struct {
	Body []CountryListItem
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerCountryRoutes() {

	group := huma.NewGroup(rs.api, "/country")

	group.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {

		op.Tags = append(op.Tags, "Country")

		next(op)
	})

	//-------------------------------------------------------------------------

	huma.Get(group, "", func(ctx context.Context, request *struct{}) (*CountryListResponse, error) {

		dbSlice, err := rs.dbQueries.GetCountries(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, "ListCountries", "error", err)

			return nil, mapDBError(err)
		}

		result := make([]CountryListItem, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APICountryListItem(dbData)
		}

		return &CountryListResponse{Body: result}, nil

	}, describeEndpoint("getCountries", "Get a list of all countries"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/id/{id}", func(ctx context.Context, request *CountryRequestId) (*CountryResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		dbresult, err := rs.dbQueries.GetCountryById(ctx, request.ID)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, "ReadCountryById", "error", err)

				return nil, mapDBError(err)
			}
		}

		country := mapDB2APICountry(dbresult)

		return &CountryResponse{Body: country}, nil

	}, describeEndpoint("getCountryById", "Get a single country based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/shortcode/{shortcode}", func(ctx context.Context, request *CountryRequestShortcode) (*CountryResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		dbresult, err := rs.dbQueries.GetCountryByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, "ReadCountryById", "error", err)

				return nil, mapDBError(err)
			}
		}

		country := mapDB2APICountry(dbresult)

		return &CountryResponse{Body: country}, nil

	}, describeEndpoint("getCountryByShortcode", "Get a single country based on the shortcode supplied"))

	//-------------------------------------------------------------------------

	huma.Post(group, "", func(ctx context.Context, request *CountryRequestCreate) (*CountryResponseCreate, error) {

		// time.Sleep(2000 * time.Millisecond)

		insertParams := rbsdb.InsertCountryParams{

			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
			Flags:      request.Body.Flags,
			Ibanlength: mapToNullInt2(request.Body.Ibanlength),
			Risktype:   request.Body.Risktype,
		}

		dbResult, err := rs.dbQueries.InsertCountry(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "CreateCountry", "error", err)

			return nil, mapDBError(err)
		}

		response := CountryResponseCreate{}

		response.Header.Location = fmt.Sprintf("/country/id/%s", dbResult.ID.String())

		response.Body = mapDB2APICountry(dbResult)

		return &response, nil

	}, describeEndpoint("createCountry", "Create a new country"), defaultStatus(http.StatusCreated))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *CountryRequestId) (*struct{}, error) {

		// time.Sleep(2000 * time.Millisecond)

		dbResult, err := rs.dbQueries.DeleteCountry(ctx, request.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, "DeleteCountry", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := dbResult.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("deleteCountry", "Delete a single country based on the id supplied")) // huma-Defaultstatus = http.StatusNoContent

	//-------------------------------------------------------------------------

	huma.Put(group, "/{id}", func(ctx context.Context, request *CountryRequestUpdate) (*CountryResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		updateParams := rbsdb.UpdateCountryParams{

			ID:         request.ID,
			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
			Flags:      request.Body.Flags,
			Ibanlength: mapToNullInt2(request.Body.Ibanlength),
			Risktype:   request.Body.Risktype,
		}

		dbResult, err := rs.dbQueries.UpdateCountry(ctx, updateParams)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("")
			}

			rs.logger.ErrorContext(ctx, "UpdateCountry", "error", err)

			return nil, mapDBError(err)
		}

		response := CountryResponse{}

		response.Body = mapDB2APICountry(dbResult)

		return &response, nil

	}, describeEndpoint("updateCountry", "Update an existing country based on the id supplied")) // huma-Defaultstatus = http.StatusOk
}

//-----------------------------------------------------------------------------

func mapDB2APICountryListItem(record rbsdb.Country) CountryListItem {

	return CountryListItem{

		ID:        record.ID,
		Shortcode: record.Shortcode,
		Name:      record.Name,
	}
}

func mapDB2APICountry(record rbsdb.Country) Country {

	return Country{

		ID:         record.ID,
		Shortcode:  record.Shortcode,
		Name:       record.Name,
		Flags:      record.Flags,
		Ibanlength: mapFromNullInt2(record.Ibanlength),
		Risktype:   record.Risktype,
	}
}
