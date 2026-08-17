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

type CustodianNoPK struct {
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"5" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"30" doc:"A longer more descriptive description of this data"`
	Flags     int16     `json:"flags" format:"int16" minimum:"0" doc:"Some binary encodes flags for this data"`
	Idcountry uuid.UUID `json:"idcountry" format:"uuid" doc:"A reference to a country, where the custodion is in"`
	Depotno   *string   `json:"depotna,omitempty" maxLength:"10" doc:"This dopot numer assocciated with this custodian"`
}

type Custodian struct {
	Id uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	CustodianNoPK
}

//-----------------------------------------------------------------------------

type CustodianRequestId struct {
	Id uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type CustodianRequestShortcode struct {
	Shortcode string `path:"shortcode" minLength:"1" maxLength:"5" doc:"This is the unique identifier a this data"`
}

type CustodianRequestNoPK struct {
	Body CustodianNoPK
}

type CustodianRequestUpdate struct {
	Id   uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Body CustodianNoPK
}

//-----------------------------------------------------------------------------

type CustodianResponseId struct {
	Id uuid.UUID `header:"id" format:"uuid" doc:"Generated id for newly created data"`
}
type CustodianResponse struct {
	Body Custodian
}

type CustodiansResponse struct {
	Body []Custodian
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerCustodianRoutes() {

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Custodian"},
		OperationID: "getCustodians",
		Summary:     "Get a list of all custodians",
		Description: "Get a list of all custodians",
		Method:      http.MethodGet,
		Path:        "/custodians",
	}, func(ctx context.Context, input *struct{}) (*CustodiansResponse, error) {

		rs.logger.Info("GetCustodians")

		queries := rbsdb.New(rs.db)

		dbSlice, err := queries.GetCustodians(ctx)

		if err != nil {

			rs.logger.Error("ListCustodians", "Error", err)

			return nil, mapDBError(err)
		}

		result := make([]Custodian, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APICustodian(dbData)
		}

		return &CustodiansResponse{Body: result}, nil
	})

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Custodian"},
		OperationID: "getCustodianById",
		Summary:     "Get a single custodian based on the id supplied",
		Description: "Get a single custodian based on the id supplied",
		Method:      http.MethodGet,
		Path:        "/custodian/id/{id}",
	}, func(ctx context.Context, request *CustodianRequestId) (*CustodianResponse, error) {

		rs.logger.Info("GetCustodianById", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		dbresult, err := queries.GetCustodianByID(ctx, request.Id)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.Error("ReadCustodianById", "Error", err)

				return nil, mapDBError(err)
			}
		}

		custodian := mapDB2APICustodian(dbresult)

		return &CustodianResponse{Body: custodian}, nil
	})

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Custodian"},
		OperationID: "getCustodianByShortcode",
		Summary:     "Get a single custodian based on the shortcode supplied",
		Description: "Get a single custodian based on the shortcode supplied",
		Method:      http.MethodGet,
		Path:        "/custodian/shortcode/{shortcode}",
	}, func(ctx context.Context, request *CustodianRequestShortcode) (*CustodianResponse, error) {

		rs.logger.Info("GetCustodianByShortcode", "Shortcode", request.Shortcode)

		queries := rbsdb.New(rs.db)

		dbresult, err := queries.GetCustodianByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.Error("ReadCustodianById", "Error", err)

				return nil, mapDBError(err)
			}
		}

		custodian := mapDB2APICustodian(dbresult)

		return &CustodianResponse{Body: custodian}, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Custodian"},
		OperationID:   "createCustodian",
		Summary:       "Create a new custodian",
		Description:   "Create a new custodian",
		Method:        http.MethodPost,
		Path:          "/custodian",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, request *CustodianRequestNoPK) (*CustodianResponseId, error) {

		rs.logger.Info("CreateCustodian")

		queries := rbsdb.New(rs.db)

		insertParams := rbsdb.InsertCustodianParams{

			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
			Flags:     request.Body.Flags,
			Idcountry: request.Body.Idcountry,
			Depotno:   mapToNullString(request.Body.Depotno),
		}

		result, err := queries.InsertCustodian(ctx, insertParams)

		if err != nil {

			rs.logger.Error("CreateCustodian", "Error", err)

			return nil, mapDBError(err)
		}

		return &CustodianResponseId{Id: result}, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Custodian"},
		OperationID:   "deleteCustodian",
		Summary:       "Delete a single custodian based on the id supplied",
		Description:   "Delete a single custodian based on the id supplied",
		Method:        http.MethodDelete,
		Path:          "/custodian/{id}",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, request *CustodianRequestId) (*struct{}, error) {

		rs.logger.Info("DeleteCustodian", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		result, err := queries.DeleteCustodian(ctx, request.Id)

		if err != nil {

			rs.logger.Error("DeleteCustodian", "Error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Custodian"},
		OperationID:   "updateCustodian",
		Summary:       "Update an existing custodian based on the id supplied",
		Description:   "Update an existing custodian based on the id supplied",
		Method:        http.MethodPut,
		Path:          "/custodian/{id}",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, request *CustodianRequestUpdate) (*struct{}, error) {

		rs.logger.Info("UpdateCustodian", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		updateParams := rbsdb.UpdateCustodianParams{

			Idcustodian: request.Id,
			Shortcode:   request.Body.Shortcode,
			Name:        request.Body.Name,
			Flags:       request.Body.Flags,
			Idcountry:   request.Body.Idcountry,
			Depotno:     mapToNullString(request.Body.Depotno),
		}

		result, err := queries.UpdateCustodian(ctx, updateParams)

		if err != nil {

			rs.logger.Error("UpdateCustodian", "Error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil
	})
}

func mapDB2APICustodian(record rbsdb.Custodian) Custodian {

	return Custodian{

		Id: record.Idcustodian,
		CustodianNoPK: CustodianNoPK{
			Shortcode: record.Shortcode,
			Name:      record.Name,
			Flags:     record.Flags,
			Idcountry: record.Idcountry,
			Depotno:   mapFromNullString(record.Depotno),
		},
	}
}
