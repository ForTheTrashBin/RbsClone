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
//	GET		Get List 			"/exchange"
//	GET		Get by Id			"/exchange/id/{id}"
//	GET		Get by Shortcode	"/exchange/shortcode/{shortcode}"
//	POST	Create				"/exchange"
//	DELETE	Delete				"/exchange/{id}"
//	PUT		Update				"/exchange/{id}"
//
//-----------------------------------------------------------------------------

type ExchangeId struct {
	ID uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type ExchangeIdPath struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type ExchangeShortcode struct {
	Shortcode string `json:"shortcode" minLength:"1" maxLength:"8" doc:"A unique short name for this data"`
}

type ExchangeShortcodePath struct {
	Shortcode string `path:"shortcode" minLength:"1" maxLength:"8" doc:"A unique short name for this data"`
}

type ExchangeName struct {
	Name string `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
}

type ExchangeFlags struct {
	Flags int16 `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
}

//-----------------------------------------------------------------------------

type ExchangeListItem struct {
	ExchangeId
	ExchangeShortcode
	ExchangeName
}

type ExchangeNoPK struct {
	ExchangeShortcode
	ExchangeName
	ExchangeFlags
}

type Exchange struct {
	ExchangeId
	ExchangeShortcode
	ExchangeName
	ExchangeFlags
}

//-----------------------------------------------------------------------------

type ExchangeRequestId struct {
	ExchangeIdPath
}

type ExchangeRequestShortcode struct {
	ExchangeShortcodePath
}

type ExchangeRequestCreate struct {
	Body ExchangeNoPK
}

type ExchangeRequestUpdate struct {
	ExchangeIdPath
	Body ExchangeNoPK
}

//-----------------------------------------------------------------------------

type ExchangeResponseList struct {
	Body []ExchangeListItem
}

type ExchangeResponseCreate struct {
	Header LocationHeader
	Body   Exchange
}

type ExchangeResponse struct {
	Body Exchange
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerExchangeRoutes() {

	group := huma.NewGroup(rs.api, "/exchange")

	group.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {

		op.Tags = append(op.Tags, "Exchange")

		next(op)
	})

	//-------------------------------------------------------------------------

	huma.Get(group, "", func(ctx context.Context, request *struct{}) (*ExchangeResponseList, error) {

		OperationId := getOperationIdFromContext(ctx)

		dbSlice, err := rs.dbQueries.GetExchanges(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		result := make([]ExchangeListItem, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APIExchangeListItem(dbData)
		}

		return &ExchangeResponseList{Body: result}, nil

	}, describeEndpoint("getExchanges", "Get a list of all exchanges"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/id/{id}", func(ctx context.Context, request *ExchangeRequestId) (*ExchangeResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		dbresult, err := rs.dbQueries.GetExchangeById(ctx, request.ID)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, OperationId, "error", err)

				return nil, mapDBError(err)
			}
		}

		exchange := mapDB2APIExchange(dbresult)

		return &ExchangeResponse{Body: exchange}, nil

	}, describeEndpoint("getExchangeById", "Get a single exchange based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/shortcode/{shortcode}", func(ctx context.Context, request *ExchangeRequestShortcode) (*ExchangeResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		dbresult, err := rs.dbQueries.GetExchangeByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, OperationId, "error", err)

				return nil, mapDBError(err)
			}
		}

		exchange := mapDB2APIExchange(dbresult)

		return &ExchangeResponse{Body: exchange}, nil

	}, describeEndpoint("getExchangeByShortcode", "Get a single exchange based on the shortcode supplied"))

	//-------------------------------------------------------------------------

	huma.Post(group, "", func(ctx context.Context, request *ExchangeRequestCreate) (*ExchangeResponseCreate, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		insertParams := rbsdb.InsertExchangeParams{

			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
			Flags:     request.Body.Flags,
		}

		dbResult, err := rs.dbQueries.InsertExchange(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		return &ExchangeResponseCreate{
			Header: LocationHeader{Location: fmt.Sprintf("/exchange/id/%s", dbResult.ID.String())},
			Body:   mapDB2APIExchange(dbResult),
		}, nil

	}, describeEndpoint("createExchange", "Create a new exchange"), defaultStatus(http.StatusCreated))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *ExchangeRequestId) (*struct{}, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		dbResult, err := rs.dbQueries.DeleteExchange(ctx, request.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := dbResult.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("deleteExchange", "Delete a single exchange based on the id supplied")) // huma-Defaultstatus = http.StatusNoContent

	//-------------------------------------------------------------------------

	huma.Put(group, "/{id}", func(ctx context.Context, request *ExchangeRequestUpdate) (*ExchangeResponse, error) {

		// time.Sleep(2000 * time.Millisecond)

		OperationId := getOperationIdFromContext(ctx)

		updateParams := rbsdb.UpdateExchangeParams{

			ID:        request.ID,
			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
			Flags:     request.Body.Flags,
		}

		dbResult, err := rs.dbQueries.UpdateExchange(ctx, updateParams)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("")
			}

			rs.logger.ErrorContext(ctx, OperationId, "error", err)

			return nil, mapDBError(err)
		}

		return &ExchangeResponse{Body: mapDB2APIExchange(dbResult)}, nil

	}, describeEndpoint("updateExchange", "Update an existing exchange based on the id supplied")) // huma-Defaultstatus = http.StatusOk
}

//-----------------------------------------------------------------------------

func mapDB2APIExchangeListItem(record rbsdb.Exchange) ExchangeListItem {

	return ExchangeListItem{
		ExchangeId:        ExchangeId{ID: record.ID},
		ExchangeShortcode: ExchangeShortcode{Shortcode: record.Shortcode},
		ExchangeName:      ExchangeName{Name: record.Name},
	}
}

func mapDB2APIExchange(record rbsdb.Exchange) Exchange {

	return Exchange{

		ExchangeId:        ExchangeId{ID: record.ID},
		ExchangeShortcode: ExchangeShortcode{Shortcode: record.Shortcode},
		ExchangeName:      ExchangeName{Name: record.Name},
		ExchangeFlags:     ExchangeFlags{Flags: record.Flags},
	}
}
