package rest

import (
	"errors"
	"log/slog"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

type RestServer struct {
	logger    *slog.Logger
	dbPool    *pgxpool.Pool
	dbQueries *rbsdb.Queries
	api       huma.API
}

func RegisterAllRoutes(logger *slog.Logger, dbPool *pgxpool.Pool, dbQueries *rbsdb.Queries, api huma.API) {

	restServer := &RestServer{

		logger:    logger,
		dbPool:    dbPool,
		dbQueries: dbQueries,
		api:       api,
	}

	restServer.registerCountryRoutes()
	restServer.registerCustodianRoutes()
	restServer.registerExchangeRoutes()
	restServer.registerCustodian2ExchangeRoutes()
}

//-----------------------------------------------------------------------------

func describeEndpoint(operationID string, summaryAndDescription string) func(op *huma.Operation) {

	return func(op *huma.Operation) {

		op.OperationID = operationID
		op.Summary = summaryAndDescription
		op.Description = summaryAndDescription
	}
}

//-----------------------------------------------------------------------------

func mapDBError(err error) error {

	if err == nil {

		return nil
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		case pgerrcode.UniqueViolation:

			return huma.Error409Conflict("Conflict in data processing")

		case pgerrcode.ForeignKeyViolation:

			return huma.Error422UnprocessableEntity("Data processing failed due to dependencies")

		case pgerrcode.NotNullViolation:

			return huma.Error400BadRequest("Data processing not possible due to incorrect input data")
		}
	}

	return huma.Error500InternalServerError("Internal server error")
}

//-----------------------------------------------------------------------------
// mapping helper from db to api
//-----------------------------------------------------------------------------

func mapFromNullInt2(data pgtype.Int2) *int16 {

	if !data.Valid {

		return nil
	}

	return &data.Int16
}

func mapToNullInt2(data *int16) pgtype.Int2 {

	if data != nil {

		return pgtype.Int2{

			Int16: *data,
			Valid: true,
		}
	}

	return pgtype.Int2{Valid: false}
}

//-----------------------------------------------------------------------------

func mapFromNullString(data pgtype.Text) *string {

	if !data.Valid {

		return nil
	}

	return &data.String
}

func mapToNullString(data *string) pgtype.Text {

	if data != nil {

		return pgtype.Text{

			String: *data,
			Valid:  true,
		}
	}

	return pgtype.Text{Valid: false}
}

//-----------------------------------------------------------------------------
