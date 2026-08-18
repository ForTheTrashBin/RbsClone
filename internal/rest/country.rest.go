package rest

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

//-----------------------------------------------------------------------------

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

type CountriesResponse struct {
	Body []Country
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerCountryRoutes() {

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Country"},
		OperationID: "getCountries",
		Summary:     "Get a list of all countries",
		Description: "Get a list of all countries",
		Method:      http.MethodGet,
		Path:        "/countries",
	}, func(ctx context.Context, input *struct{}) (*CountriesResponse, error) {

		rs.logger.Info("GetCountries")

		queries := rbsdb.New(rs.db)

		dbSlice, err := queries.GetCountries(ctx)

		if err != nil {

			rs.logger.Error("ListCountries", "error", err)

			return nil, mapDBError(err)
		}

		result := make([]Country, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APICountry(dbData)
		}

		return &CountriesResponse{Body: result}, nil
	})

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Country"},
		OperationID: "getCountryById",
		Summary:     "Get a single country based on the id supplied",
		Description: "Get a single country based on the id supplied",
		Method:      http.MethodGet,
		Path:        "/country/id/{id}",
	}, func(ctx context.Context, request *CountryRequestId) (*CountryResponse, error) {

		rs.logger.Info("GetCountryById", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		dbresult, err := queries.GetCountryByID(ctx, request.Id)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.Error("ReadCountryById", "error", err)

				return nil, mapDBError(err)
			}
		}

		country := mapDB2APICountry(dbresult)

		return &CountryResponse{Body: country}, nil
	})

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Country"},
		OperationID: "getCountryByShortcode",
		Summary:     "Get a single country based on the shortcode supplied",
		Description: "Get a single country based on the shortcode supplied",
		Method:      http.MethodGet,
		Path:        "/country/shortcode/{shortcode}",
	}, func(ctx context.Context, request *CountryRequestShortcode) (*CountryResponse, error) {

		rs.logger.Info("GetCountryByShortcode", "Shortcode", request.Shortcode)

		queries := rbsdb.New(rs.db)

		dbresult, err := queries.GetCountryByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.Error("ReadCountryById", "error", err)

				return nil, mapDBError(err)
			}
		}

		country := mapDB2APICountry(dbresult)

		return &CountryResponse{Body: country}, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Country"},
		OperationID:   "createCountry",
		Summary:       "Create a new country",
		Description:   "Create a new country",
		Method:        http.MethodPost,
		Path:          "/country",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, request *CountryRequestCreate) (*CountryResponseCreate, error) {

		rs.logger.Info("CreateCountry")

		queries := rbsdb.New(rs.db)

		insertParams := rbsdb.InsertCountryParams{

			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
			Flags:      request.Body.Flags,
			Ibanlength: mapToNullInt2(request.Body.Ibanlength),
			Risktype:   request.Body.Risktype,
		}

		result, err := queries.InsertCountry(ctx, insertParams)

		if err != nil {

			rs.logger.Error("CreateCountry", "error", err)

			return nil, mapDBError(err)
		}

		return &CountryResponseCreate{Id: result}, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Country"},
		OperationID:   "deleteCountry",
		Summary:       "Delete a single country based on the id supplied",
		Description:   "Delete a single country based on the id supplied",
		Method:        http.MethodDelete,
		Path:          "/country/{id}",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, request *CountryRequestId) (*struct{}, error) {

		rs.logger.Info("DeleteCountry", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		result, err := queries.DeleteCountry(ctx, request.Id)

		if err != nil {

			rs.logger.Error("DeleteCountry", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Country"},
		OperationID:   "updateCountry",
		Summary:       "Update an existing country based on the id supplied",
		Description:   "Update an existing country based on the id supplied",
		Method:        http.MethodPut,
		Path:          "/country/{id}",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, request *CountryRequestUpdate) (*struct{}, error) {

		rs.logger.Info("UpdateCountry", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		updateParams := rbsdb.UpdateCountryParams{

			Idcountry:  request.Id,
			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
			Flags:      request.Body.Flags,
			Ibanlength: mapToNullInt2(request.Body.Ibanlength),
			Risktype:   request.Body.Risktype,
		}

		result, err := queries.UpdateCountry(ctx, updateParams)

		if err != nil {

			rs.logger.Error("UpdateCountry", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil
	})
}

//-----------------------------------------------------------------------------

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
