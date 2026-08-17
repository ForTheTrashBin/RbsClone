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

type ExchangeNoPK struct {
	Shortcode string `json:"shortcode" minLength:"1" maxLength:"8" doc:"A unique short name for this data"`
	Name      string `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
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

type ExchangeRequestNoPK struct {
	Body ExchangeNoPK
}

type ExchangeRequestUpdate struct {
	Id   uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Body ExchangeNoPK
}

//-----------------------------------------------------------------------------

type ExchangeResponseId struct {
	Id uuid.UUID `header:"id" format:"uuid" doc:"Generated id for newly created data"`
}
type ExchangeResponse struct {
	Body Exchange
}

type ExchangesResponse struct {
	Body []Exchange
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerExchangeRoutes() {

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Exchange"},
		OperationID: "getExchanges",
		Summary:     "Get a list of all exchanges",
		Description: "Get a list of all exchanges",
		Method:      http.MethodGet,
		Path:        "/exchanges",
	}, func(ctx context.Context, input *struct{}) (*ExchangesResponse, error) {

		rs.logger.Info("GetExchanges")

		queries := rbsdb.New(rs.db)

		dbSlice, err := queries.GetExchanges(ctx)

		if err != nil {

			rs.logger.Error("ListExchanges", "Error", err)

			return nil, mapDBError(err)
		}

		result := make([]Exchange, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APIExchange(dbData)
		}

		return &ExchangesResponse{Body: result}, nil
	})

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Exchange"},
		OperationID: "getExchangeById",
		Summary:     "Get a single exchange based on the id supplied",
		Description: "Get a single exchange based on the id supplied",
		Method:      http.MethodGet,
		Path:        "/exchange/id/{id}",
	}, func(ctx context.Context, request *ExchangeRequestId) (*ExchangeResponse, error) {

		rs.logger.Info("GetExchangeById", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		dbresult, err := queries.GetExchangeByID(ctx, request.Id)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.Error("ReadExchangeById", "Error", err)

				return nil, mapDBError(err)
			}
		}

		exchange := mapDB2APIExchange(dbresult)

		return &ExchangeResponse{Body: exchange}, nil
	})

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Exchange"},
		OperationID: "getExchangeByShortcode",
		Summary:     "Get a single exchange based on the shortcode supplied",
		Description: "Get a single exchange based on the shortcode supplied",
		Method:      http.MethodGet,
		Path:        "/exchange/shortcode/{shortcode}",
	}, func(ctx context.Context, request *ExchangeRequestShortcode) (*ExchangeResponse, error) {

		rs.logger.Info("GetExchangeByShortcode", "Shortcode", request.Shortcode)

		queries := rbsdb.New(rs.db)

		dbresult, err := queries.GetExchangeByShortcode(ctx, request.Shortcode)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.Error("ReadExchangeById", "Error", err)

				return nil, mapDBError(err)
			}
		}

		exchange := mapDB2APIExchange(dbresult)

		return &ExchangeResponse{Body: exchange}, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Exchange"},
		OperationID:   "createExchange",
		Summary:       "Create a new exchange",
		Description:   "Create a new exchange",
		Method:        http.MethodPost,
		Path:          "/exchange",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, request *ExchangeRequestNoPK) (*ExchangeResponseId, error) {

		rs.logger.Info("CreateExchange")

		queries := rbsdb.New(rs.db)

		insertParams := rbsdb.InsertExchangeParams{

			Shortcode: request.Body.Shortcode,
			Name:      request.Body.Name,
		}

		result, err := queries.InsertExchange(ctx, insertParams)

		if err != nil {

			rs.logger.Error("CreateExchange", "Error", err)

			return nil, mapDBError(err)
		}

		return &ExchangeResponseId{Id: result}, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Exchange"},
		OperationID:   "deleteExchange",
		Summary:       "Delete a single exchange based on the id supplied",
		Description:   "Delete a single exchange based on the id supplied",
		Method:        http.MethodDelete,
		Path:          "/exchange/{id}",
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, request *ExchangeRequestId) (*struct{}, error) {

		rs.logger.Info("DeleteExchange", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		result, err := queries.DeleteExchange(ctx, request.Id)

		if err != nil {

			rs.logger.Error("DeleteExchange", "Error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Exchange"},
		OperationID:   "updateExchange",
		Summary:       "Update an existing exchange based on the id supplied",
		Description:   "Update an existing exchange based on the id supplied",
		Method:        http.MethodPut,
		Path:          "/exchange/{id}",
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, request *ExchangeRequestUpdate) (*struct{}, error) {

		rs.logger.Info("UpdateExchange", "Id", request.Id)

		queries := rbsdb.New(rs.db)

		updateParams := rbsdb.UpdateExchangeParams{

			Idexchange: request.Id,
			Shortcode:  request.Body.Shortcode,
			Name:       request.Body.Name,
		}

		result, err := queries.UpdateExchange(ctx, updateParams)

		if err != nil {

			rs.logger.Error("UpdateExchange", "Error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := result.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil
	})
}

func mapDB2APIExchange(record rbsdb.Exchange) Exchange {

	return Exchange{

		Id: record.Idexchange,
		ExchangeNoPK: ExchangeNoPK{
			Shortcode: record.Shortcode,
			Name:      record.Name,
		},
	}
}
