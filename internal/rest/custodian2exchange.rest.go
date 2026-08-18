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

type Custodian2ExchangeNoPK struct {
	Flags   int16  `json:"flags" format:"int16" minimum:"0" doc:"Some binary encoded flags for this data (see external documentation)"`
	Value01 string `json:"value01" minLength:"1" maxLength:"80" doc:"A value01 for this data"`
	Value02 int16  `json:"value02" format:"int16" minimum:"0" doc:"A value01 for this data"`
}

type Custodian2Exchange struct {
	Idexchange  uuid.UUID `json:"idexchange" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
	Idcustodian uuid.UUID `json:"idcustodian" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
	Custodian2ExchangeNoPK
}

//-----------------------------------------------------------------------------

type Custodian2ExchangeRequestIdExchange struct {
	Idexchange uuid.UUID `path:"idexchange" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
}

type Custodian2ExchangeRequestIdCustodian struct {
	Idcustodian uuid.UUID `path:"idcustodian" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
}

//-----------------------------------------------------------------------------

type MapCustodian2Exchange struct {
	Idcustodian uuid.UUID `json:"idcustodian" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
	Custodian2ExchangeNoPK
}

type MapCustodian2ExchangeRequestIdExchange struct {
	Idexchange uuid.UUID `path:"idexchange" doc:"This is one of the two parts of the unique identifier of this data"`

	Body []MapCustodian2Exchange
}

//-----------------------------------------------------------------------------

type MapExchange2Custodian struct {
	Idexchange uuid.UUID `json:"idexchange" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`
	Custodian2ExchangeNoPK
}

type MapCustodian2ExchangeRequestIdCustodian struct {
	Idcustodian uuid.UUID `json:"idcustodian" format:"uuid" doc:"This is one of the two parts of the unique identifier of this data"`

	Body []MapExchange2Custodian
}

//-----------------------------------------------------------------------------

type Custodian2ExchangeResponse struct {
	Body []Custodian2Exchange
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerCustodian2ExchangeRoutes() {

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Custodian2Exchange"},
		OperationID: "getCustodian2ExchangeByIdCustodian",
		Summary:     "Get a list of all mappings by idcustodian supplied",
		Description: "Get a list of all mappings by idcustodian supplied",
		Method:      http.MethodGet,
		Path:        "/custodian2exchange/idcustodian/{idcustodian}",
	}, func(ctx context.Context, request *Custodian2ExchangeRequestIdCustodian) (*Custodian2ExchangeResponse, error) {

		rs.logger.Info("GetCustodian2ExchangeByIdCustodian", "IdCustodian", request.Idcustodian)

		queries := rbsdb.New(rs.db)

		dbSlice, err := queries.GetCustodian2ExchangeByIdCustodian(ctx, request.Idcustodian)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.Error("ReadCustodian2ExchangeById", "error", err)

				return nil, mapDBError(err)
			}
		}

		result := make([]Custodian2Exchange, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APICustodian2Exchange(dbData)
		}

		return &Custodian2ExchangeResponse{Body: result}, nil
	})

	huma.Register(rs.api, huma.Operation{
		Tags:        []string{"Custodian2Exchange"},
		OperationID: "getCustodian2ExchangeByIdExchange",
		Summary:     "Get a list of mappings by idexchange supplied",
		Description: "Get a list of mappings by idexchange supplied",
		Method:      http.MethodGet,
		Path:        "/custodian2exchange/idexchange/{idexchange}",
	}, func(ctx context.Context, request *Custodian2ExchangeRequestIdExchange) (*Custodian2ExchangeResponse, error) {

		rs.logger.Info("GetCustodian2ExchangeByIdExchange", "IdExchange", request.Idexchange)

		queries := rbsdb.New(rs.db)

		dbSlice, err := queries.GetCustodian2ExchangeByIdExchange(ctx, request.Idexchange)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {

				return nil, huma.Error404NotFound("No data found")

			} else {

				rs.logger.Error("ReadCustodian2ExchangeById", "error", err)

				return nil, mapDBError(err)
			}
		}

		result := make([]Custodian2Exchange, len(dbSlice))

		for idx, dbData := range dbSlice {

			result[idx] = mapDB2APICustodian2Exchange(dbData)
		}

		return &Custodian2ExchangeResponse{Body: result}, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Custodian2Exchange"},
		OperationID:   "mapCustodians2Exchange",
		Summary:       "Modify the mapping of multiple custodians to a single exchange",
		Description:   "Modify the mapping of multiple custodians to a single exchange",
		Method:        http.MethodPost,
		Path:          "/custodian2exchange/{idexchange}",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, request *MapCustodian2ExchangeRequestIdExchange) (*struct{}, error) {

		rs.logger.Info("MapCustodian2ExchangeRequestIdExchange")

		tx, err := rs.db.Begin(ctx)

		if err != nil {

			rs.logger.Error("MapCustodian2ExchangeRequestIdExchange", "error", err)

			return nil, mapDBError(err)
		}

		defer tx.Rollback(ctx)

		queries := rbsdb.New(tx)

		//---------------------------------------------------------------------
		// Step A: Load the current status from the database
		//---------------------------------------------------------------------

		dbSlice, err := queries.GetCustodian2ExchangeByIdExchange(ctx, request.Idexchange)

		if err != nil {

			rs.logger.Error("MapCustodian2ExchangeRequestIdExchange", "error", err)

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		//---------------------------------------------------------------------

		dbMap := make(map[uuid.UUID]rbsdb.Custodian2exchange)

		for _, dbRecord := range dbSlice {

			dbMap[dbRecord.Idcustodian] = dbRecord
		}

		//---------------------------------------------------------------------

		targetMap := make(map[uuid.UUID]MapCustodian2Exchange)

		for _, target := range request.Body {

			targetMap[target.Idcustodian] = target
		}

		//---------------------------------------------------------------------
		// Step B:
		//---------------------------------------------------------------------

		for _, target := range request.Body {

			existing, exists := dbMap[target.Idcustodian]

			if !exists {

				err := queries.InsertCustodian2Exchange(ctx, rbsdb.InsertCustodian2ExchangeParams{

					Idexchange:  request.Idexchange,
					Idcustodian: target.Idcustodian,
					Flags:       target.Flags,
					Value01:     target.Value01,
					Value02:     target.Value02,
				})

				if err != nil {

					rs.logger.Error("MapCustodian2ExchangeRequestIdExchange", "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			} else {

				if existing.Flags != target.Flags || existing.Value01 != target.Value01 || existing.Value02 != target.Value02 {

					_, err := queries.UpdateCustodian2Exchange(ctx, rbsdb.UpdateCustodian2ExchangeParams{

						Idexchange:  request.Idexchange,
						Idcustodian: target.Idcustodian,
						Flags:       target.Flags,
						Value01:     target.Value01,
						Value02:     target.Value02,
					})

					if err != nil {

						rs.logger.Error("MapCustodian2ExchangeRequestIdExchange", "error", err)

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

				_, err := queries.DeleteCustodian2Exchange(ctx, rbsdb.DeleteCustodian2ExchangeParams{

					Idexchange:  request.Idexchange,
					Idcustodian: existing.Idcustodian,
				})

				if err != nil {

					rs.logger.Error("MapCustodian2ExchangeRequestIdExchange", "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			}
		}

		return nil, nil
	})

	//-------------------------------------------------------------------------

	huma.Register(rs.api, huma.Operation{
		Tags:          []string{"Custodian2Exchange"},
		OperationID:   "mapExchanges2Custodian",
		Summary:       "Modify the mapping of multiple exchanges to a single custodian",
		Description:   "Modify the mapping of multiple exchanges to a single custodian",
		Method:        http.MethodPost,
		Path:          "/exchange2custodian/{idcustodian}",
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, request *MapCustodian2ExchangeRequestIdCustodian) (*struct{}, error) {

		rs.logger.Info("MapCustodian2ExchangeRequestIdCustodian")

		tx, err := rs.db.Begin(ctx)

		if err != nil {

			rs.logger.Error("MapCustodian2ExchangeRequestIdCustodian", "error", err)

			return nil, mapDBError(err)
		}

		defer tx.Rollback(ctx)

		queries := rbsdb.New(tx)

		//---------------------------------------------------------------------
		// Step A: Load the current status from the database
		//---------------------------------------------------------------------

		dbSlice, err := queries.GetCustodian2ExchangeByIdCustodian(ctx, request.Idcustodian)

		if err != nil {

			rs.logger.Error("MapCustodian2ExchangeRequestIdCustodian", "error", err)

			return nil, huma.Error500InternalServerError("Internal server error")
		}

		//---------------------------------------------------------------------

		dbMap := make(map[uuid.UUID]rbsdb.Custodian2exchange)

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

				err := queries.InsertCustodian2Exchange(ctx, rbsdb.InsertCustodian2ExchangeParams{

					Idcustodian: request.Idcustodian,
					Idexchange:  target.Idexchange,
					Flags:       target.Flags,
					Value01:     target.Value01,
					Value02:     target.Value02,
				})

				if err != nil {

					rs.logger.Error("MapCustodian2ExchangeRequestIdCustodian", "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			} else {

				if existing.Flags != target.Flags || existing.Value01 != target.Value01 || existing.Value02 != target.Value02 {

					_, err := queries.UpdateCustodian2Exchange(ctx, rbsdb.UpdateCustodian2ExchangeParams{

						Idcustodian: request.Idcustodian,
						Idexchange:  target.Idexchange,
						Flags:       target.Flags,
						Value01:     target.Value01,
						Value02:     target.Value02,
					})

					if err != nil {

						rs.logger.Error("MapCustodian2ExchangeRequestIdCustodian", "error", err)

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

				_, err := queries.DeleteCustodian2Exchange(ctx, rbsdb.DeleteCustodian2ExchangeParams{

					Idcustodian: request.Idcustodian,
					Idexchange:  existing.Idexchange,
				})

				if err != nil {

					rs.logger.Error("MapCustodian2ExchangeRequestIdCustodian", "error", err)

					return nil, huma.Error500InternalServerError("Internal server error")
				}
			}
		}

		return nil, nil
	})
}

//-----------------------------------------------------------------------------

func mapDB2APICustodian2Exchange(record rbsdb.Custodian2exchange) Custodian2Exchange {

	return Custodian2Exchange{

		Idexchange:  record.Idexchange,
		Idcustodian: record.Idcustodian,

		Custodian2ExchangeNoPK: Custodian2ExchangeNoPK{
			Flags:   record.Flags,
			Value01: record.Value01,
			Value02: record.Value02,
		},
	}
}
