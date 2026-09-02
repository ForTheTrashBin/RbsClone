package rest

import (
	"context"
	"database/sql"
	"errors"

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
	Id        uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode string    `json:"shortcode" minLength:"2" maxLength:"2" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"30" doc:"A longer more descriptive description of this data"`
}

type CountryNoPK struct {
	Shortcode  string `json:"shortcode" minLength:"2" maxLength:"2" doc:"A unique short name for this data"`
	Name       string `json:"name" minLength:"1" maxLength:"30" doc:"A longer more descriptive description of this data"`
	Flags      int16  `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
	Ibanlength *int16 `json:"ibanlenth,omitempty" format:"int16" minimum:"0" doc:"The exact length of the IBAN required in that country"`
	Risktype   int16  `json:"risktype" format:"int16" minimum:"0" doc:"This risk profile of this country"`
}

type Country struct {
	Id uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	CountryNoPK
}

//-----------------------------------------------------------------------------

type CountryRequestId struct {
	Id uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type CountryRequestShortcode struct {
	Shortcode string `path:"shortcode" minLength:"2" maxLength:"2" doc:"This is the unique identifier a this data"`
}

type CountryRequestCreate struct {
	Body CountryNoPK
}

type CountryRequestUpdate struct {
	Id   uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Body CountryNoPK
}

//-----------------------------------------------------------------------------

type CountryResponseCreate struct {
	Id uuid.UUID `header:"id" format:"uuid" doc:"Generated id for newly created data"`
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

		dbresult, err := rs.dbQueries.GetCountryByID(ctx, request.Id)

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

		insertParams := rbsdb.InsertCountryParams{

			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
			Flags:      request.Body.Flags,
			Ibanlength: mapToNullInt2(request.Body.Ibanlength),
			Risktype:   request.Body.Risktype,
		}

		result, err := rs.dbQueries.InsertCountry(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "CreateCountry", "error", err)

			return nil, mapDBError(err)
		}

		return &CountryResponseCreate{Id: result}, nil

	}, describeEndpoint("createCountry", "Create a new country"))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *CountryRequestId) (*struct{}, error) {

		result, err := rs.dbQueries.DeleteCountry(ctx, request.Id)

		if err != nil {

			rs.logger.ErrorContext(ctx, "DeleteCountry", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("deleteCountry", "Delete a single country based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Put(group, "/{id}", func(ctx context.Context, request *CountryRequestUpdate) (*struct{}, error) {

		updateParams := rbsdb.UpdateCountryParams{

			Idcountry:  request.Id,
			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
			Flags:      request.Body.Flags,
			Ibanlength: mapToNullInt2(request.Body.Ibanlength),
			Risktype:   request.Body.Risktype,
		}

		result, err := rs.dbQueries.UpdateCountry(ctx, updateParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "UpdateCountry", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("updateCountry", "Update an existing country based on the id supplied"))
}

//-----------------------------------------------------------------------------

func mapDB2APICountryListItem(record rbsdb.Country) CountryListItem {

	return CountryListItem{

		Id:        record.Idcountry,
		Shortcode: record.Shortcode,
		Name:      record.Name,
	}
}

func mapDB2APICountry(record rbsdb.Country) Country {

	return Country{

		Id: record.Idcountry,
		CountryNoPK: CountryNoPK{
			Shortcode:  record.Shortcode,
			Name:       record.Name,
			Flags:      record.Flags,
			Ibanlength: mapFromNullInt2(record.Ibanlength),
			Risktype:   record.Risktype,
		},
	}
}
