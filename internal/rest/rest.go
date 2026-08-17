package rest

import (
	"errors"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

type RestServer struct {
	logger *slog.Logger
	db     *pgxpool.Pool
	api    huma.API
}

func RegisterAllRoutes(logger *slog.Logger, db *pgxpool.Pool, api huma.API) {

	restServer := &RestServer{

		logger: logger,
		db:     db,
		api:    api,
	}

	restServer.registerCountryRoutes()
	restServer.registerCustodianRoutes()
	restServer.registerExchangeRoutes()
	restServer.registerCustodian2ExchangeRoutes()
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

func mapDBError(err error) error {

	if err == nil {

		return nil
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {

		switch pgErr.Code {

		case "23505":

			return huma.Error409Conflict("Conflict in data processing")

		case "23503":

			return huma.Error422UnprocessableEntity("Data processing failed due to dependencies")

		case "23502":

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
