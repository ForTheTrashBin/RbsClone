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

type CountryId struct {
	ID uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type CountryIdPath struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type CountryShortcode struct {
	Shortcode string `json:"shortcode" minLength:"2" maxLength:"2" doc:"A unique short name for this data"`
}

type CountryShortcodePath struct {
	Shortcode string `path:"shortcode" minLength:"2" maxLength:"2" doc:"A unique short name for this data"`
}

type CountryName struct {
	Name string `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
}

type CountryFlags struct {
	Flags int16 `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
}

type CountryIbanlength struct {
	Ibanlength *int16 `json:"ibanlenth,omitempty" format:"int16" minimum:"8" maximum:"34" doc:"The exact length of the IBAN required in that country"`
}

type CountryRisktype struct {
	Risktype int16 `json:"risktype" format:"int16" minimum:"0" doc:"This risk profile of this country"`
}

//-----------------------------------------------------------------------------

type CountryListItem struct {
	CountryId
	CountryShortcode
	CountryName
}

type CountryNoPK struct {
	CountryShortcode
	CountryName
	CountryFlags
	CountryIbanlength
	CountryRisktype
}

type Country struct {
	CountryId
	CountryShortcode
	CountryName
	CountryFlags
	CountryIbanlength
	CountryRisktype
}

//-----------------------------------------------------------------------------

type CountryRequestId struct {
	CountryIdPath
}

type CountryRequestShortcode struct {
	CountryShortcodePath
}

type CountryRequestCreate struct {
	Body CountryNoPK
}

type CountryRequestUpdate struct {
	CountryIdPath
	Body CountryNoPK
}

//-----------------------------------------------------------------------------

type CountryResponseList struct {
	Body []CountryListItem
}

type CountryResponseCreate struct {
	Header LocationHeader
	Body   Country
}

type CountryResponse struct {
	Body Country
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerCountryRoutes() {

	group := huma.NewGroup(rs.api, "/country")

	group.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {

		op.Tags = append(op.Tags, "Country")

		next(op)
	})

	//-------------------------------------------------------------------------

	huma.Get(group, "", func(ctx context.Context, request *struct{}) (*CountryResponseList, error) {

		OperationId := getOperationIdFromContext(ctx)

		dbSlice, err := rs.dbQueries.GetCountries(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		result := make([]CountryListItem, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APICountryListItem(dbData)
		}

		return &CountryResponseList{Body: result}, nil

	}, describeEndpoint("getCountries", "Get a list of all countries"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/id/{id}", func(ctx context.Context, request *CountryRequestId) (*CountryResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		dbresult, err := rs.dbQueries.GetCountryById(ctx, request.ID)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, OperationId, "error", err)

				return nil, mapDBError(err)
			}
		}

		country := mapDB2APICountry(dbresult)

		return &CountryResponse{Body: country}, nil

	}, describeEndpoint("getCountryById", "Get a single country based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/shortcode/{shortcode}", func(ctx context.Context, request *CountryRequestShortcode) (*CountryResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		dbresult, err := rs.dbQueries.GetCountryByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, OperationId, "error", err)

				return nil, mapDBError(err)
			}
		}

		country := mapDB2APICountry(dbresult)

		return &CountryResponse{Body: country}, nil

	}, describeEndpoint("getCountryByShortcode", "Get a single country based on the shortcode supplied"))

	//-------------------------------------------------------------------------

	huma.Post(group, "", func(ctx context.Context, request *CountryRequestCreate) (*CountryResponseCreate, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		insertParams := rbsdb.InsertCountryParams{

			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
			Flags:      request.Body.Flags,
			Ibanlength: mapToNullInt2(request.Body.Ibanlength),
			Risktype:   request.Body.Risktype,
		}

		dbResult, err := rs.dbQueries.InsertCountry(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		return &CountryResponseCreate{
			Header: LocationHeader{Location: fmt.Sprintf("/country/id/%s", dbResult.ID.String())},
			Body:   mapDB2APICountry(dbResult),
		}, nil

	}, describeEndpoint("createCountry", "Create a new country"), defaultStatus(http.StatusCreated))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *CountryRequestId) (*struct{}, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		dbResult, err := rs.dbQueries.DeleteCountry(ctx, request.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

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

		OperationId := getOperationIdFromContext(ctx)

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

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		return &CountryResponse{Body: mapDB2APICountry(dbResult)}, nil

	}, describeEndpoint("updateCountry", "Update an existing country based on the id supplied")) // huma-Defaultstatus = http.StatusOk
}

//-----------------------------------------------------------------------------

func mapDB2APICountryListItem(record rbsdb.Country) CountryListItem {

	return CountryListItem{

		CountryId:        CountryId{ID: record.ID},
		CountryShortcode: CountryShortcode{Shortcode: record.Shortcode},
		CountryName:      CountryName{Name: record.Name},
	}
}

func mapDB2APICountry(record rbsdb.Country) Country {

	return Country{

		CountryId:         CountryId{ID: record.ID},
		CountryShortcode:  CountryShortcode{Shortcode: record.Shortcode},
		CountryName:       CountryName{Name: record.Name},
		CountryFlags:      CountryFlags{Flags: record.Flags},
		CountryIbanlength: CountryIbanlength{Ibanlength: mapFromNullInt2(record.Ibanlength)},
		CountryRisktype:   CountryRisktype{Risktype: record.Risktype},
	}
}
