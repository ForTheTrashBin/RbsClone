package rest

import (
	"context"

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

type Exchange2Custodian struct {
	IIdcustodian uuid.UUID `json:"idcustodian" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
	IIdexchange  uuid.UUID `json:"idexchange" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
	Sequenceno   int32     `json:"sequenceno" format:"int32" doc:"Determines the order of the stock exchamges"`
}

//-----------------------------------------------------------------------------

type Exchange2CustodianRequestIdcustodian struct {
	Idcustodian uuid.UUID `path:"idcustodian" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
}

//-----------------------------------------------------------------------------

type MapExchange2Custodian struct {
	Idexchange uuid.UUID `json:"idexchange" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
	Sequenceno int32     `json:"sequenceno" format:"int32" doc:"Determines the order of the stock exchamges"`
}

type MapExchange2CustodianRequestIdcustodian struct {
	Idcustodian uuid.UUID `path:"idcustodian" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`

	Body []MapExchange2Custodian
}

//-----------------------------------------------------------------------------

type Exchange2CustodianResponse struct {
	Body []Exchange2Custodian
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerExchange2CustodianRoutes() {

	group := huma.NewGroup(rs.api, "/exchange2custodian")

	group.UseModifier(func(op *huma.Operation, next func(*huma.Operation)) {

		op.Tags = append(op.Tags, "Exchange2Custodian")

		next(op)
	})

	//-------------------------------------------------------------------------

	huma.Get(group, "/idcustodian/{idcustodian}", func(ctx context.Context, request *Exchange2CustodianRequestIdcustodian) (*Exchange2CustodianResponse, error) {

		dbSlice, err := rs.dbQueries.GetExchange2CustodianByIdcustodian(ctx, request.Idcustodian)

		if err != nil {

			rs.logger.ErrorContext(ctx, "ReadExchange2CustodianByIdcustodian", "error", err)

			return nil, mapDBError(err)
		}

		result := make([]Exchange2Custodian, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APIExchange2Custodian(dbData)
		}

		return &Exchange2CustodianResponse{Body: result}, nil

	}, describeEndpoint("getExchange2CustodianByIdcustodian", "Get a list of all mappings by idcustodian supplied"))

	//-------------------------------------------------------------------------

	/*	huma.Put(group, "/idcustodian/{idcustodian}", func(ctx context.Context, request *MapExchange2CustodianRequestIdcustodian) (*struct{}, error) {

		tx, err := rs.dbPool.Begin(ctx)

		if err != nil {

			rs.logger.ErrorContext(ctx, "MapExchange2CustodianRequestIdcustodian", "error", err)

			return nil, mapDBError(err)
		}

		defer tx.Rollback(ctx)

		queries := rs.dbQueries.WithTx(tx)

		//---------------------------------------------------------------------
		// Step A: Load the current status from the database
		//---------------------------------------------------------------------

		dbSlice, err := queries.GetExchange2CustodianByIdcustodian(ctx, request.Idcustodian)

		if err != nil {

			rs.logger.ErrorContext(ctx, "MapExchange2CustodianRequestIdcustodian", "error", err)

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		//---------------------------------------------------------------------

		dbMap := make(map[uuid.UUID]rbsdb.Exchange2custodian)

		for _, dbRecord := range dbSlice {

			dbMap[dbRecord.Idcustodian] = dbRecord
		}

		//---------------------------------------------------------------------

		targetMap := make(map[uuid.UUID]MapExchange2Custodian)

		for _, target := range request.Body {

			targetMap[target.Idexchange] = target
		}

		//---------------------------------------------------------------------
		// Step B:
		//---------------------------------------------------------------------

		for _, target := range request.Body {

			existing, exists := dbMap[target.Idexchange]

			if !exists {

				err := queries.InsertExchange2Custodian(ctx, rbsdb.InsertExchange2CustodianParams{

					Idcustodian: request.Idcustodian,
					Idexchange:  target.Idexchange,
					Flags:       target.Flags,
					Value01:     target.Value01,
					Value02:     target.Value02,
				})

				if err != nil {

					rs.logger.ErrorContext(ctx, "MapExchange2CustodianRequestIdcustodian", "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			} else {

				if existing.Flags != target.Flags || existing.Value01 != target.Value01 || existing.Value02 != target.Value02 {

					_, err := queries.UpdateExchange2Custodian(ctx, rbsdb.UpdateExchange2CustodianParams{

						Idcustodian: request.Idcustodian,
						Idexchange:  target.Idexchange,
						Flags:       target.Flags,
						Value01:     target.Value01,
						Value02:     target.Value02,
					})

					if err != nil {

						rs.logger.ErrorContext(ctx, "MapExchange2CustodianRequestIdcustodian", "error", err)

						return nil, huma.Error500InternalServerError("Internal server error")
					}
				}
			}
		}

		//---------------------------------------------------------------------
		// Step C:
		//---------------------------------------------------------------------

		for _, existing := range dbMap {

			if _, exists := targetMap[existing.Idcustodian]; !exists {

				_, err := queries.DeleteExchange2Custodian(ctx, rbsdb.DeleteExchange2CustodianParams{

					Idcustodian: request.Idcustodian,
					Idexchange:  existing.Idexchange,
				})

				if err != nil {

					rs.logger.ErrorContext(ctx, "MapExchange2CustodianRequestIdcustodian", "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			}
		}

		tx.Commit(ctx)

		return nil, nil

	}, describeEndpoint("mapExchanges2Custodian", "Update the mapping of multiple exchanges to a single custodian")) */
}

//-----------------------------------------------------------------------------

func mapDB2APIExchange2Custodian(record rbsdb.Exchange2custodian) Exchange2Custodian {

	return Exchange2Custodian{

		IIdcustodian: record.Idcustodian,
		IIdexchange:  record.Idexchange,
		Sequenceno:   record.Sequenceno,
	}
}
