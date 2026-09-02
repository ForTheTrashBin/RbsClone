package logging

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

type ctxKeyLogger struct{} // Private! It's nobody's business!

//-----------------------------------------------------------------------------

func SlogMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			ctx := context.WithValue(r.Context(), ctxKeyLogger{}, logger)

			defer func() {

				logMessage := fmt.Sprintf("%s %s %d", r.Method, r.URL.Path, ww.Status())

				logger.InfoContext(ctx, logMessage,
					slog.Duration("duration", time.Since(start)),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("remote", r.RemoteAddr),
				)
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}

func GetLogger(ctx context.Context) *slog.Logger {

	if logger, ok := ctx.Value(ctxKeyLogger{}).(*slog.Logger); ok {

		return logger
	}

	return slog.Default()
}
