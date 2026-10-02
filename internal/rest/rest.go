package rest

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//-----------------------------------------------------------------------------

type LocationHeader struct {
	Location string `header:"Location" doc:"URL of the newly created entity"`
}

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

	restServer.registerUtilities()

	restServer.registerCountryRoutes()
	restServer.registerCustodianRoutes()
	restServer.registerExchangeRoutes()
}

//-----------------------------------------------------------------------------

func (rs *RestServer) registerUtilities() {

	huma.Register(rs.api, huma.Operation{
		OperationID:   "get-ping",
		Method:        http.MethodGet,
		Path:          "/ping",
		Summary:       "Connection-Test",
		Description:   "Test of the connection to this server",
		Tags:          []string{"Utilities"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, request *struct{}) (*struct{}, error) {

		return nil, nil
	})

	huma.Register(rs.api, huma.Operation{
		OperationID:   "get-health",
		Method:        http.MethodGet,
		Path:          "/health",
		Summary:       "State of services and components (Health)",
		Description:   "Check the state of services and components",
		Tags:          []string{"Utilities"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, request *struct{}) (*struct{}, error) {

		if err := rs.dbPool.Ping(ctx); err != nil {

			rs.logger.ErrorContext(ctx, "DB-Pingtest", "error", err)

			return nil, huma.Error503ServiceUnavailable("")
		} else {

			return nil, nil
		}
	})
}

//-----------------------------------------------------------------------------

type operationIdType string

const operationIdKey operationIdType = "MyOperationIdKey"

func getOperationIdFromContext(ctx context.Context) string {

	if operationIdFromContext, ok := ctx.Value(operationIdKey).(string); ok {

		return cases.Title(language.Und, cases.NoLower).String(operationIdFromContext)
	} else {

		return "No operationId set as value of context!"
	}
}

func describeEndpoint(operationID string, summaryAndDescription string) func(op *huma.Operation) {

	return func(op *huma.Operation) {

		op.OperationID = operationID

		op.Summary = summaryAndDescription
		op.Description = summaryAndDescription

		//---------------------------------------------------------------------

		op.Middlewares = append(op.Middlewares, func(humaCtx huma.Context, next func(huma.Context)) {

			newCtx := context.WithValue(humaCtx.Context(), operationIdKey, op.OperationID)

			humaCtx = huma.WithContext(humaCtx, newCtx)

			next(humaCtx)
		})
	}
}

func defaultStatus(defaultStatus int) func(op *huma.Operation) {
	return func(op *huma.Operation) {

		op.DefaultStatus = defaultStatus
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
