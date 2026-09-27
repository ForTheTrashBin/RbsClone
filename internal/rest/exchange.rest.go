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

type ExchangeListItem struct {
	ID        uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"8" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
}

type ExchangeNoPK struct {
	Shortcode string `json:"shortcode" minLength:"1" maxLength:"8" doc:"A unique short name for this data"`
	Name      string `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
	Flags     int16  `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
}

type Exchange struct {
	ID        uuid.UUID `json:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Shortcode string    `json:"shortcode" minLength:"1" maxLength:"8" doc:"A unique short name for this data"`
	Name      string    `json:"name" minLength:"1" maxLength:"80" doc:"A longer more descriptive description of this data"`
	Flags     int16     `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
}

//-----------------------------------------------------------------------------

type ExchangeRequestId struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
}

type ExchangeRequestShortcode struct {
	Shortcode string `path:"shortcode" minLength:"1" maxLength:"8" doc:"This is the unique identifier a this data"`
}

type ExchangeRequestCreate struct {
	Body ExchangeNoPK
}

type ExchangeRequestUpdate struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"This is the unique identifier a this data"`
	Body ExchangeNoPK
}

//-----------------------------------------------------------------------------

type ExchangeResponseCreate struct {
	Header struct {
		Location string `header:"Location" doc:"URL of the newly created entity"`
	}

	Body Exchange
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

		dbresult, err := rs.dbQueries.GetExchangeById(ctx, request.ID)

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

		dbResult, err := rs.dbQueries.InsertExchange(ctx, insertParams)

		if err != nil {

			rs.logger.ErrorContext(ctx, "CreateExchange", "error", err)

			return nil, mapDBError(err)
		}

		response := ExchangeResponseCreate{}

		response.Header.Location = fmt.Sprintf("/exchange/id/%s", dbResult.ID.String())

		response.Body = mapDB2APIExchange(dbResult)

		return &response, nil

	}, describeEndpoint("createExchange", "Create a new exchange"), defaultStatus(http.StatusCreated))

	//-------------------------------------------------------------------------

	huma.Delete(group, "/{id}", func(ctx context.Context, request *ExchangeRequestId) (*struct{}, error) {

		dbResult, err := rs.dbQueries.DeleteExchange(ctx, request.ID)

		if err != nil {

			rs.logger.ErrorContext(ctx, "DeleteExchange", "error", err)

			return nil, mapDBError(err)
		}

		if rowsAffected := dbResult.RowsAffected(); rowsAffected == 0 {

			return nil, huma.Error404NotFound("")
		}

		return nil, nil

	}, describeEndpoint("deleteExchange", "Delete a single exchange based on the id supplied")) // huma-Defaultstatus = http.StatusNoContent

	//-------------------------------------------------------------------------

	huma.Put(group, "/{id}", func(ctx context.Context, request *ExchangeRequestUpdate) (*ExchangeResponse, error) {

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

			rs.logger.ErrorContext(ctx, "UpdateCountry", "error", err)

			return nil, mapDBError(err)
		}

		response := ExchangeResponse{}

		response.Body = mapDB2APIExchange(dbResult)

		return &response, nil

	}, describeEndpoint("updateExchange", "Update an existing exchange based on the id supplied")) // huma-Defaultstatus = http.StatusOk
}

//-----------------------------------------------------------------------------

func mapDB2APIExchangeListItem(record rbsdb.Exchange) ExchangeListItem {

	return ExchangeListItem{

		ID:        record.ID,
		Shortcode: record.Shortcode,
		Name:      record.Name,
	}
}

func mapDB2APIExchange(record rbsdb.Exchange) Exchange {

	return Exchange{

		ID:        record.ID,
		Shortcode: record.Shortcode,
		Name:      record.Name,
		Flags:     record.Flags,
	}
}
