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
//	GET		Get List 			"/custodian"
//	GET		Get by Id			"/custodian/id/{id}"
//	GET		Get by Shortcode	"/custodian/shortcode/{shortcode}"
//	POST	Create				"/custodian"
//	DELETE	Delete				"/custodian/{id}"
//	PUT		Update				"/custodian/{id}"
//
//-----------------------------------------------------------------------------

type CustodianListItem struct {
	Id        uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"5" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"30" doc:"A longer more descriptive description of this data"`
}

type CustodianNoPK struct {
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"5" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"30" doc:"A longer more descriptive description of this data"`
	Flags     int16     `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
	Idcountry uuid.UUID `json:"idcountry" format:"uuid" doc:"A reference to a country, where the custodion is in"`
	Depotno   *string   `json:"depotno,omitempty" maxLength:"10" doc:"This dopot numer assocciated with this custodian"`
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

type CustodianRequestCreate struct {
	Body CustodianNoPK
}

type CustodianRequestUpdate struct {
	Id   uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Body CustodianNoPK
}

//-----------------------------------------------------------------------------

type CustodianResponseCreate struct {
	Id uuid.UUID `header:"id" format:"uuid" doc:"Generated id for newly created data"`
}
type CustodianResponse struct {
	Body Custodian
}

type CustodianListResponse struct {
	Body []CustodianListItem
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerCustodianRoutes() {

	group := huma.NewGroup(rs.api, "/custodian")

	group.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {

		op.Tags = append(op.Tags, "Custodian")

		next(op)
	})

	//-------------------------------------------------------------------------

	huma.Get(group, "", func(ctx context.Context, request *struct{}) (*CustodianListResponse, error) {

		dbSlice, err := rs.dbQueries.GetCustodians(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, "ListCustodians", "error", err)

			return nil, mapDBError(err)
		}

		result := make([]CustodianListItem, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APICustodianListItem(dbData)
		}

		return &CustodianListResponse{Body: result}, nil

	}, describeEndpoint("getCustodians", "Get a list of all custodians"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/id/{id}", func(ctx context.Context, request *CustodianRequestId) (*CustodianResponse, error) {

		dbresult, err := rs.dbQueries.GetCustodianByID(ctx, request.Id)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, "ReadCustodianById", "error", err)

				return nil, mapDBError(err)
			}
		}

		custodian := mapDB2APICustodian(dbresult)

		return &CustodianResponse{Body: custodian}, nil

	}, describeEndpoint("getCustodianById", "Get a single custodian based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/shortcode/{shortcode}", func(ctx context.Context, request *CustodianRequestShortcode) (*CustodianResponse, error) {

		dbresult, err := rs.dbQueries.GetCustodianByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, "ReadCustodianById", "error", err)

				return nil, mapDBError(err)
			}
		}

		custodian := mapDB2APICustodian(dbresult)

		return &CustodianResponse{Body: custodian}, nil

	}, describeEndpoint("getCustodianByShortcode", "Get a single custodian based on the shortcode supplied"))

	//-------------------------------------------------------------------------

	huma.Post(group, "", func(ctx context.Context, request *CustodianRequestCreate) (*CustodianResponseCreate, error) {

		insertParams := rbsdb.InsertCustodianParams{

			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
			Flags:     request.Body.Flags,
			Idcountry: request.Body.Idcountry,
			Depotno:   mapToNullString(request.Body.Depotno),
		}

		result, err := rs.dbQueries.InsertCustodian(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "CreateCustodian", "error", err)

			return nil, mapDBError(err)
		}

		return &CustodianResponseCreate{Id: result}, nil

	}, describeEndpoint("createCustodian", "Create a new custodian"))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *CustodianRequestId) (*struct{}, error) {

		result, err := rs.dbQueries.DeleteCustodian(ctx, request.Id)

		if err != nil {

			rs.logger.ErrorContext(ctx, "DeleteCustodian", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("deleteCustodian", "Delete a single custodian based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Put(group, "/{id}", func(ctx context.Context, request *CustodianRequestUpdate) (*struct{}, error) {

		updateParams := rbsdb.UpdateCustodianParams{

			Idcustodian: request.Id,
			Shortcode:   request.Body.Shortcode,
			Name:        request.Body.Name,
			Flags:       request.Body.Flags,
			Idcountry:   request.Body.Idcountry,
			Depotno:     mapToNullString(request.Body.Depotno),
		}

		result, err := rs.dbQueries.UpdateCustodian(ctx, updateParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "UpdateCustodian", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("updateCustodian", "Update an existing custodian based on the id supplied"))
}

//-----------------------------------------------------------------------------

func mapDB2APICustodianListItem(record rbsdb.Custodian) CustodianListItem {

	return CustodianListItem{

		Id:        record.Idcustodian,
		Shortcode: record.Shortcode,
		Name:      record.Name,
	}
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
