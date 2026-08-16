package rest

import (
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

type RESTServer struct {
	logger *slog.Logger
	db     *pgxpool.Pool
	api    huma.API
}

func RegisterAllRoutes(logger *slog.Logger, db *pgxpool.Pool, api huma.API) {

	srv := &RESTServer{

		logger: logger,
		db:     db,
		api:    api,
	}

	srv.registerCountryRoutes()
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
/*
func mapFromNullString(data sql.NullString) *string {

	if !data.Valid {

		return nil
	}

	return &data.String
}

func mapToNullString(data *string) sql.NullString {

	if data != nil {

		return sql.NullString{

			String: *data,
			Valid:  true,
		}
	}

	return sql.NullString{Valid: false}
}
*/
//-----------------------------------------------------------------------------
