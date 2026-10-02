package rest

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

//-----------------------------------------------------------------------------
//
//	GET		Get by Idcustodian	"/exchange2custodian/idcustodian/{idcustodian}"
//	GET		Get by Idexchange	"/exchange2custodian/idexchange/{idexchange}"
//	PUT		Update				"/exchange2custodian/idcustodian/{idcustodian}"
//	PUT		Update				"/exchange2custodian/idexchange/{idexchange}"
//
//-----------------------------------------------------------------------------

type Exchange2CustodianListItem struct {
	Idexchange uuid.UUID `json:"idexchange" format:"uuid" doc:"This is the unique identifier of this data"`
	Value1     int16     `json:"value1" format:"int16" minimum:"0" maximum:"255" doc:"This is the first special payload for testing"`
	Value2     int16     `json:"value2" format:"int16" minimum:"0" maximum:"255" doc:"This is the second special payload for testing"`
}

type Exchange2CustodianDefaultList struct {
	Idexchangedef uuid.UUID `json:"idexchangedef" format:"uuid" doc:"The id of the default-exchange"`

	Exchanges []Exchange2CustodianListItem `json:"exchanges"`
}

//-----------------------------------------------------------------------------

type Exchange2CustodianGetRequest struct {
	Idcustodian uuid.UUID `path:"idcustodian" format:"uuid" doc:"This is the unique identifier a this data"`
}

type Exchange2CustodianGetResponse struct {
	Body Exchange2CustodianDefaultList
}

//-----------------------------------------------------------------------------

type Exchange2CustodianPutRequest struct {
	Idcustodian uuid.UUID `path:"idcustodian" format:"uuid" doc:"This is the unique identifier a this data"`
	Body        Exchange2CustodianDefaultList
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerExchange2CustodianRoutes() {

	group := huma.NewGroup(rs.api, "/exchange2custodian")

	group.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {

		op.Tags = append(op.Tags, "Exchange2Custodian")

		next(op)
	})

	//-------------------------------------------------------------------------

	huma.Get(group, "/{idcustodian}", func(ctx context.Context, request *Exchange2CustodianGetRequest) (*Exchange2CustodianGetResponse, error) {

		dbSlice, err := rs.dbQueries.GetExchange2CustodianByIdcustodian(ctx, request.Idcustodian)

		if err != nil {

			rs.logger.ErrorContext(ctx, "GetExchange2CustodianByIdcustodian", "error", err)

			return nil, mapDBError(err)
		}

		var Idexchangedef uuid.UUID = uuid.New()

		if len(dbSlice) > 0 {
			Idexchangedef = dbSlice[0].Idexchange
		}

		result := Exchange2CustodianDefaultList{
			Idexchangedef: Idexchangedef,
			Exchanges:     make([]Exchange2CustodianListItem, 0),
		}

		for idx, dbData := range dbSlice {

			result.Exchanges[idx] = Exchange2CustodianListItem{
				Idexchange: dbData.Idexchange,
				Value1:     dbData.Value1,
				Value2:     dbData.Value2,
			}
		}

		return &Exchange2CustodianGetResponse{Body: result}, nil

	}, describeEndpoint("getExchange2CustodianByIdcustodian", "Get a list of all mappings by idcustodian supplied"))

	//-------------------------------------------------------------------------

	huma.Put(group, "/{idcustodian}", func(ctx context.Context, request *Exchange2CustodianPutRequest) (*struct{}, error) {

		tx, err := rs.dbPool.Begin(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, "MapExchange2CustodianRequestIdcustodian", "error", err)

			return nil, mapDBError(err)
		}

		defer tx.Rollback(ctx)

		queries := rs.dbQueries.WithTx(tx)

		//---------------------------------------------------------------------
		// Step A: Delete all records for given 'Idcustodian'
		//---------------------------------------------------------------------

		_, err = queries.DeleteExchange2CustodianByIdcustodian(ctx, request.Idcustodian)

		if err != nil {

			rs.logger.ErrorContext(ctx, "PutExchange2CustodianByIdcustodian", "error", err)

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		if len(request.Body.Exchanges) > 0 {
			//-----------------------------------------------------------------
			// Step B: Find default-exchange in collection
			//-----------------------------------------------------------------

			index := slices.IndexFunc(request.Body.Exchanges,
				func(entry Exchange2CustodianListItem) bool {
					return entry.Idexchange == request.Idcustodian
				})

			if index < 0 {
				rs.logger.ErrorContext(ctx, "PutExchange2CustodianByIdcustodian",
					"error", fmt.Errorf("no default-exchange found in array: %s", request.Body.Idexchangedef))

				return nil, huma.Error500InternalServerError("Internal server error")
			}

			//-----------------------------------------------------------------
			// Step C: Insert default-exchange at position 0 into db
			//-----------------------------------------------------------------

			err = queries.InsertExchange2Custodian(ctx, rbsdb.InsertExchange2CustodianParams{
				Idcustodian: request.Idcustodian,
				Idexchange:  request.Body.Exchanges[index].Idexchange,
				Sequenceno:  0,
				Value1:      request.Body.Exchanges[index].Value1,
				Value2:      request.Body.Exchanges[index].Value2})

			if err != nil {

				rs.logger.ErrorContext(ctx, "PutExchange2CustodianByIdcustodian", "error", err)

				return nil, huma.Error500InternalServerError("Internal server error")
			}

			//-----------------------------------------------------------------
			// Step D: Remove default-exchange from collection
			//-----------------------------------------------------------------

			request.Body.Exchanges = slices.Delete(request.Body.Exchanges, index, index+1)

			//-----------------------------------------------------------------
			// Step E: Insert all remaining echanges to db
			//-----------------------------------------------------------------

			for sequenceno, exchange := range request.Body.Exchanges {
				err = queries.InsertExchange2Custodian(ctx, rbsdb.InsertExchange2CustodianParams{
					Idcustodian: request.Idcustodian,
					Idexchange:  exchange.Idexchange,
					Sequenceno:  int32(sequenceno + 1),
					Value1:      exchange.Value1,
					Value2:      exchange.Value2})

				if err != nil {

					rs.logger.ErrorContext(ctx, "PutExchange2CustodianByIdcustodian", "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			}
		}

		tx.Commit(ctx)

		return nil, nil

	}, describeEndpoint("putExchange2CustodianByIdcustodian", "Update the mapping of multiple exchanges to a single custodian"), defaultStatus(http.StatusAccepted))
}
