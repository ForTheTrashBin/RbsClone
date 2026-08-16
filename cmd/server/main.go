package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ForTheTrashBin/RbsClone/internal/osspecific"
	"github.com/ForTheTrashBin/RbsClone/internal/rest"
	"github.com/ForTheTrashBin/RbsClone/internal/serverconfig"
	"golang.org/x/sync/errgroup"
)

//-----------------------------------------------------------------------------
// Handler and wrapper for http messages
//-----------------------------------------------------------------------------

type HttpHandlerWrapper struct {
	config serverconfig.Config
}

func NewHttpHandlerWrapper(config serverconfig.Config) HttpHandlerWrapper {

	return HttpHandlerWrapper{

		config: config,
	}
}

//-----------------------------------------------------------------------------

func (hhw HttpHandlerWrapper) httpHandler(w http.ResponseWriter, r *http.Request) {

	newHost := r.Host

	httpPort := strconv.Itoa(int(hhw.config.GetHTTPPort()))
	httpsPort := strconv.Itoa(int(hhw.config.GetHTTPSPort()))

	if host, port, err := net.SplitHostPort(r.Host); err == nil {

		if port == httpPort {

			newHost = net.JoinHostPort(host, httpsPort)
		}
	} else {

		newHost = net.JoinHostPort(r.Host, httpsPort)
	}

	target := "https://" + newHost + r.URL.RequestURI()

	http.Redirect(w, r, target, http.StatusMovedPermanently) // could be http.StatusTemporaryRedirect too
}

//-----------------------------------------------------------------------------
// The main entry-point of this app
//-----------------------------------------------------------------------------

func main() {

	var config serverconfig.Config

	//-------------------------------------------------------------------------
	// Read content of .env-file into the 'local' environment of this proccess
	//-------------------------------------------------------------------------

	if err := config.InitGoDotEnv(); err != nil {

		panic(fmt.Errorf("Error from 'godotenv': %w", err))
	}

	//-------------------------------------------------------------------------
	// Parse the content of the 'local' environmet into the Config-struct
	//-------------------------------------------------------------------------

	if err := config.InitCaarlos0(); err != nil {

		panic(fmt.Errorf("Error from 'env': %w", err))
	}

	//-------------------------------------------------------------------------
	// Create a new structured logger that writes JSON to stdout.
	//-------------------------------------------------------------------------

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.LogLevel})) // Level: parseLogLevel(config.LogLevel)}))

	logger.Info("Application config", "DB_USER", config.DB.User)
	logger.Info("Application config", "DB_PASSWORD", config.DB.Password)
	logger.Info("Application config", "DB_HOST", config.DB.Host)
	logger.Info("Application config", "DB_PORT", config.DB.Port)
	logger.Info("Application config", "DB_DATABASE", config.DB.Database)

	logger.Info("Application config", "RT_HTTPPORT", config.GetHTTPPort())
	logger.Info("Application config", "RT_HTTPSPORT", config.GetHTTPSPort())

	//-------------------------------------------------------------------------
	// Connect to database
	//-------------------------------------------------------------------------

	// Create a pgxpool.Config for manual configuration

	pgxconfig, err := pgxpool.ParseConfig("")

	if err != nil {

		logger.Error("Can't parse postgres-config", "error", err)

		return
	}

	// Create a pgxpool.Config for manual configuration

	pgxconfig.ConnConfig.User = config.DB.User
	pgxconfig.ConnConfig.Password = config.DB.Password
	pgxconfig.ConnConfig.Host = config.DB.Host
	pgxconfig.ConnConfig.Port = config.DB.Port
	pgxconfig.ConnConfig.Database = config.DB.Database

	pgxconfig.ConnConfig.ConnectTimeout = 5 * time.Second

	// Connection pool tuning

	pgxconfig.MinConns = 5
	pgxconfig.MaxConns = 25
	pgxconfig.MaxConnLifetime = 1 * time.Hour
	pgxconfig.MaxConnIdleTime = 15 * time.Minute
	pgxconfig.HealthCheckPeriod = 1 * time.Minute

	//-------------------------------------------------------------------------
	// Connect to postgres database
	//-------------------------------------------------------------------------

	ctxDBConnect, cancelDBConnectTimeout := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancelDBConnectTimeout()

	//-------------------------------------------------------------------------

	dbPool, err := pgxpool.NewWithConfig(ctxDBConnect, pgxconfig)

	if err != nil {

		logger.Error("Can't open database", "error", err)

		return
	}

	defer dbPool.Close()

	//-------------------------------------------------------------------------

	if err := dbPool.Ping(ctxDBConnect); err != nil {

		logger.Error("Can't ping database", "error", err)

		return
	}

	//-------------------------------------------------------------------------

	cancelDBConnectTimeout()

	//-------------------------------------------------------------------------

	// tools.InitializeCountry(logger, dbPool)

	//-------------------------------------------------------------------------
	// Create a context that is canceled (ctx.Done()) when an interrupt signal is received
	//-------------------------------------------------------------------------

	ctx, stopSignaling := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	//-------------------------------------------------------------------------
	// NotifyContext has "redirected" various signals (e.g., SIGTERM) to its
	// "own" handler. Calling 'stopSignaling' reverses this, and to make sure
	// this isn't forgotten --> 'defer stopSignaling'
	//
	// Remember:
	//
	// NotifyContext changes the given context with 'WithCancelCause' and
	// returns a new context that is canceled when one of the signals is
	// received.
	//-------------------------------------------------------------------------

	defer stopSignaling() // Don't forget to stop the signal notification when done

	//-------------------------------------------------------------------------
	// Disable echo for the duration of the processing
	//-------------------------------------------------------------------------

	if err := osspecific.SetEcho(false); err == nil {

		// Remember: defered functions are executed in LIFO order, so FlushStdin will be called before SetEcho(true)

		defer osspecific.SetEcho(true) // Restore echo on exit
		defer osspecific.FlushStdin()  // Flush stdin on exit
	}

	//-------------------------------------------------------------------------

	errGrp, ctx := errgroup.WithContext(ctx)

	//-------------------------------------------------------------------------
	// Create a listener to listen on the http port
	//-------------------------------------------------------------------------

	httpHandlerWrapper := NewHttpHandlerWrapper(config)

	serverHTTP := &http.Server{

		Handler:           http.HandlerFunc(httpHandlerWrapper.httpHandler),
		Addr:              ":" + strconv.Itoa(int(config.GetHTTPPort())),
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       3 * time.Second,
	}

	httpListener, err := net.Listen("tcp", serverHTTP.Addr)

	if err != nil {

		logger.Error("failed to bind http", "error", err)

		return
	}

	//-------------------------------------------------------------------------
	// Create a listener to listen on the https port
	//-------------------------------------------------------------------------

	humaConfig := huma.DefaultConfig("API for Rbs", "0.1.0")

	humaConfig.CreateHooks = nil

	//-------------------------------------------------------------------------

	router := http.NewServeMux()

	api := humago.New(router, humaConfig)

	rest.RegisterAllRoutes(logger, dbPool, api)

	serverHTTPS := &http.Server{

		Handler:           router, // myhandler, // http.HandlerFunc(myhandler), // http.HandlerFunc(httpsHandler),
		Addr:              ":" + strconv.Itoa(int(config.GetHTTPSPort())),
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       3 * time.Second,
	}

	httpsListener, err := net.Listen("tcp", serverHTTPS.Addr)

	if err != nil {

		httpListener.Close()

		logger.Error("failed to bind https", "error", err)

		return
	}

	//-------------------------------------------------------------------------

	cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/cert-key.pem")

	if err != nil {

		httpListener.Close()
		httpsListener.Close()

		logger.Error("failed to load TLS certificate", "error", err)

		return
	}

	httpsListener = tls.NewListener(httpsListener, &tls.Config{Certificates: []tls.Certificate{cert}})

	//-------------------------------------------------------------------------
	// Log successful bindings
	//-------------------------------------------------------------------------

	logger.Info("http server listening", "addr", httpListener.Addr().String())
	logger.Info("https server listening", "addr", httpsListener.Addr().String())

	//-------------------------------------------------------------------------
	// http-server for automatic redirect from http to https
	//-------------------------------------------------------------------------

	errGrp.Go(func() error {

		if err := serverHTTP.Serve(httpListener); err != nil {

			if errors.Is(err, http.ErrServerClosed) {

				logger.Info("http-server closed (gracefully)")

				return nil
			} else {

				logger.Error("http server returned error", "error", err)

				return err
			}
		}

		return nil
	})

	//-------------------------------------------------------------------------
	// https-server to serve https requests
	//-------------------------------------------------------------------------

	errGrp.Go(func() error {

		if err := serverHTTPS.Serve(httpsListener); err != nil {

			if errors.Is(err, http.ErrServerClosed) {

				logger.Info("https-server closed (gracefully)")

				return nil
			} else {

				logger.Error("https server returned error", "error", err)

				return err
			}
		}

		return nil
	})

	errGrp.Go(func() error {

		//-------------------------------------------------------------------------
		// Wait until EITHER an OS signal is received OR a server throws an
		// error (ctx is terminated)
		//
		// 1. If an OS signal is received, the servers will be shut down gracefully.
		// 2. If a server throws an error, the other server will be shut down gracefully.
		//
		// A server-specific error is detected whenever one of the 'errGrp.Go'
		// functions returns an error. This will cancel the context too, which will
		// cause this function to run.
		//-------------------------------------------------------------------------

		<-ctx.Done()

		logger.Info("Shutdown signal received or error detected. Servers are shutting down...")

		//-------------------------------------------------------------------------
		// Set a separate timeout for shutdown (e.g., 10 seconds)
		//-------------------------------------------------------------------------

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		//-------------------------------------------------------------------------
		// Shut down both servers simultaneously (uses a sub-errgroup for the shutdown)
		//-------------------------------------------------------------------------

		shutdownGroup, _ := errgroup.WithContext(shutdownCtx)

		shutdownGroup.Go(func() error { return serverHTTP.Shutdown(shutdownCtx) })
		shutdownGroup.Go(func() error { return serverHTTPS.Shutdown(shutdownCtx) })

		return shutdownGroup.Wait()
	})

	if err := errGrp.Wait(); err != nil {

		logger.Error("Server terminated with an error", "error", err)
	} else {

		logger.Info("Server shut down normally")
	}
}
