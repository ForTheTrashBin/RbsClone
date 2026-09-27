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
//	GET		Get List 			"/custodian"
//	GET		Get by ID			"/custodian/id/{id}"
//	GET		Get by Shortcode	"/custodian/shortcode/{shortcode}"
//	POST	Create				"/custodian"
//	DELETE	Delete				"/custodian/{id}"
//	PUT		Update				"/custodian/{id}"
//
//-----------------------------------------------------------------------------

type CustodianListItem struct {
	ID        uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"5" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
}

type CustodianNoPK struct {
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"5" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
	Flags     int16     `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
	Idcountry uuid.UUID `json:"idcountry" format:"uuid" doc:"A reference to a country, where the custodion is in"`
	Depotno   *string   `json:"depotno,omitempty" maxLength:"10" doc:"This dopot numer assocciated with this custodian"`
}

type Custodian struct {
	ID        uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"5" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
	Flags     int16     `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
	Idcountry uuid.UUID `json:"idcountry" format:"uuid" doc:"A reference to a country, where the custodion is in"`
	Depotno   *string   `json:"depotno,omitempty" maxLength:"10" doc:"This dopot numer assocciated with this custodian"`
}

//-----------------------------------------------------------------------------

type CustodianRequestId struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type CustodianRequestShortcode struct {
	Shortcode string `path:"shortcode" minLength:"1" maxLength:"5" doc:"This is the unique identifier a this data"`
}

type CustodianRequestCreate struct {
	Body CustodianNoPK
}

type CustodianRequestUpdate struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Body CustodianNoPK
}

//-----------------------------------------------------------------------------

type CustodianResponseCreate struct {
	Header struct {
		Location string `header:"Location" doc:"URL of the newly created entity"`
	}

	Body Custodian
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

		dbresult, err := rs.dbQueries.GetCustodianById(ctx, request.ID)

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

		dbResult, err := rs.dbQueries.InsertCustodian(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "CreateCustodian", "error", err)

			return nil, mapDBError(err)
		}

		response := CustodianResponseCreate{}

		response.Header.Location = fmt.Sprintf("/custodian/id/%s", dbResult.ID.String())

		response.Body = mapDB2APICustodian(dbResult)

		return &response, nil

	}, describeEndpoint("createCustodian", "Create a new custodian"), defaultStatus(http.StatusCreated))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *CustodianRequestId) (*struct{}, error) {

		dbResult, err := rs.dbQueries.DeleteCustodian(ctx, request.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, "DeleteCustodian", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := dbResult.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("deleteCustodian", "Delete a single custodian based on the id supplied")) // huma-Defaultstatus = http.StatusNoContent

	//-------------------------------------------------------------------------

	huma.Put(group, "/{id}", func(ctx context.Context, request *CustodianRequestUpdate) (*CustodianResponse, error) {

		updateParams := rbsdb.UpdateCustodianParams{

			ID:        request.ID,
			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
			Flags:     request.Body.Flags,
			Idcountry: request.Body.Idcountry,
			Depotno:   mapToNullString(request.Body.Depotno),
		}

		dbResult, err := rs.dbQueries.UpdateCustodian(ctx, updateParams)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("")
			}

			rs.logger.ErrorContext(ctx, "UpdateCountry", "error", err)

			return nil, mapDBError(err)
		}

		response := CustodianResponse{}

		response.Body = mapDB2APICustodian(dbResult)

		return &response, nil

	}, describeEndpoint("updateCustodian", "Update an existing custodian based on the id supplied")) // huma-Defaultstatus = http.StatusOk
}

//-----------------------------------------------------------------------------

func mapDB2APICustodianListItem(record rbsdb.Custodian) CustodianListItem {

	return CustodianListItem{

		ID:        record.ID,
		Shortcode: record.Shortcode,
		Name:      record.Name,
	}
}

func mapDB2APICustodian(record rbsdb.Custodian) Custodian {

	return Custodian{

		ID:        record.ID,
		Shortcode: record.Shortcode,
		Name:      record.Name,
		Flags:     record.Flags,
		Idcountry: record.Idcountry,
		Depotno:   mapFromNullString(record.Depotno),
	}
}
