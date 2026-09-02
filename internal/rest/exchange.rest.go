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
//	GET		Get List 			"/exchange"
//	GET		Get by Id			"/exchange/id/{id}"
//	GET		Get by Shortcode	"/exchange/shortcode/{shortcode}"
//	POST	Create				"/exchange"
//	DELETE	Delete				"/exchange/{id}"
//	PUT		Update				"/exchange/{id}"
//
//-----------------------------------------------------------------------------

type ExchangeListItem struct {
	Id        uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"8" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
}

type ExchangeNoPK struct {
	Shortcode string `json:"shortcode" minLength:"1" maxLength:"8" doc:"A unique short name for this data"`
	Name      string `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
	Flags     int16  `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
}

type Exchange struct {
	Id uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	ExchangeNoPK
}

//-----------------------------------------------------------------------------

type ExchangeRequestId struct {
	Id uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type ExchangeRequestShortcode struct {
	Shortcode string `path:"shortcode" minLength:"1" maxLength:"8" doc:"This is the unique identifier a this data"`
}

type ExchangeRequestCreate struct {
	Body ExchangeNoPK
}

type ExchangeRequestUpdate struct {
	Id   uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Body ExchangeNoPK
}

//-----------------------------------------------------------------------------

type ExchangeResponseCreate struct {
	Id uuid.UUID `header:"id" format:"uuid" doc:"Generated id for newly created data"`
}
type ExchangeResponse struct {
	Body Exchange
}

type ExchangeListResponse struct {
	Body []ExchangeListItem
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerExchangeRoutes() {

	group := huma.NewGroup(rs.api, "/exchange")

	group.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {

		op.Tags = append(op.Tags, "Exchange")

		next(op)
	})

	//-------------------------------------------------------------------------

	huma.Get(group, "", func(ctx context.Context, request *struct{}) (*ExchangeListResponse, error) {

		dbSlice, err := rs.dbQueries.GetExchanges(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, "ListExchanges", "error", err)

			return nil, mapDBError(err)
		}

		result := make([]ExchangeListItem, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APIExchangeListItem(dbData)
		}

		return &ExchangeListResponse{Body: result}, nil

	}, describeEndpoint("getExchanges", "Get a list of all exchanges"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/id/{id}", func(ctx context.Context, request *ExchangeRequestId) (*ExchangeResponse, error) {

		dbresult, err := rs.dbQueries.GetExchangeByID(ctx, request.Id)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, "ReadExchangeById", "error", err)

				return nil, mapDBError(err)
			}
		}

		exchange := mapDB2APIExchange(dbresult)

		return &ExchangeResponse{Body: exchange}, nil

	}, describeEndpoint("getExchangeById", "Get a single exchange based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Get(group, "/shortcode/{shortcode}", func(ctx context.Context, request *ExchangeRequestShortcode) (*ExchangeResponse, error) {

		dbresult, err := rs.dbQueries.GetExchangeByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.ErrorContext(ctx, "ReadExchangeById", "error", err)

				return nil, mapDBError(err)
			}
		}

		exchange := mapDB2APIExchange(dbresult)

		return &ExchangeResponse{Body: exchange}, nil

	}, describeEndpoint("getExchangeByShortcode", "Get a single exchange based on the shortcode supplied"))

	//-------------------------------------------------------------------------

	huma.Post(group, "", func(ctx context.Context, request *ExchangeRequestCreate) (*ExchangeResponseCreate, error) {

		insertParams := rbsdb.InsertExchangeParams{

			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
			Flags:     request.Body.Flags,
		}

		result, err := rs.dbQueries.InsertExchange(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "CreateExchange", "error", err)

			return nil, mapDBError(err)
		}

		return &ExchangeResponseCreate{Id: result}, nil

	}, describeEndpoint("createExchange", "Create a new exchange"))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *ExchangeRequestId) (*struct{}, error) {

		result, err := rs.dbQueries.DeleteExchange(ctx, request.Id)

		if err != nil {

			rs.logger.ErrorContext(ctx, "DeleteExchange", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("deleteExchange", "Delete a single exchange based on the id supplied"))

	//-------------------------------------------------------------------------

	huma.Put(group, "/{id}", func(ctx context.Context, request *ExchangeRequestUpdate) (*struct{}, error) {

		updateParams := rbsdb.UpdateExchangeParams{

			Idexchange: request.Id,
			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
			Flags:      request.Body.Flags,
		}

		result, err := rs.dbQueries.UpdateExchange(ctx, updateParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "UpdateExchange", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("updateExchange", "Update an existing exchange based on the id supplied"))
}

//-----------------------------------------------------------------------------

func mapDB2APIExchangeListItem(record rbsdb.Exchange) ExchangeListItem {

	return ExchangeListItem{

		Id:        record.Idexchange,
		Shortcode: record.Shortcode,
		Name:      record.Name,
	}
}

func mapDB2APIExchange(record rbsdb.Exchange) Exchange {

	return Exchange{

		Id: record.Idexchange,
		ExchangeNoPK: ExchangeNoPK{
			Shortcode: record.Shortcode,
			Name:      record.Name,
			Flags:     record.Flags,
		},
	}
}
