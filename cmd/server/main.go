package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ForTheTrashBin/RbsClone/internal/osspecific"
	"github.com/ForTheTrashBin/RbsClone/internal/rbsdb"
	"github.com/ForTheTrashBin/RbsClone/internal/serverconfig"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/sync/errgroup"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {

	//-------------------------------------------------------------------------
	// Striktes HTTP -> HTTPS Redirect
	//-------------------------------------------------------------------------

	host := r.Host

	if h, p, err := net.SplitHostPort(r.Host); err == nil {

		if p == serverconfig.GetHttpPort() {

			host = net.JoinHostPort(h, serverconfig.GetHttpsPort())
		}
	} else {

		host = net.JoinHostPort(r.Host, serverconfig.GetHttpsPort())
	}

	target := "https://" + host + r.URL.RequestURI()

	http.Redirect(w, r, target, http.StatusTemporaryRedirect) // could be http.StatusMovedPermanently too
}

func httpsHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Secure work completed!"))
}

func helper(log *slog.Logger) {
	/*
		databaseHost := "skylax-dkt-01-docker"
		databasePort := "5432"

		databaseUsername := "admin%40example.com"
		databasePassword := "eiterbatzen123"

		databaseScheme := "postgres"
		databaseName := "my_database"

		dsn := url.URL{

			Scheme: databaseScheme,
			User:   url.UserPassword(databaseUsername, databasePassword),
			Host:   fmt.Sprintf("%s:%s", databaseHost, databasePort),
			Path:   databaseName,
		}

		q := dsn.Query()
		q.Add("sslmode", "disable")

		// dsn.RawQuery = q.Encode()
	*/
	dsn := "postgres://postgres:eiterbatzen123@skylax-dkt-01-docker:5432/my_database?sslmode=disable"

	fmt.Println("The connectionstring:", dsn)

	//-------------------------------------------------------------------------

	dbPool, err := sql.Open("pgx", dsn)

	if err != nil {

		log.Error("Can't open database", "error", err)

		return
	}

	defer dbPool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	defer cancel()

	if err := dbPool.PingContext(ctx); err != nil {

		log.Error("Can't ping database", "error", err)

		return
	}

	//-------------------------------------------------------------------------
	// pool tuning

	dbPool.SetMaxOpenConns(25)
	dbPool.SetMaxIdleConns(5)
	dbPool.SetConnMaxLifetime(5 * time.Minute)

	//-------------------------------------------------------------------------

	baseQueries := rbsdb.New(dbPool)
	/*
		log.Info("Insert START")

		for idx := 0; idx < 10000; idx++ {

			shortcode := fmt.Sprintf("XXX%d", idx+1)

			_, err = baseQueries.CreateExchange(ctx,
				rbsdb.CreateExchangeParams{
					Shortcode:   shortcode,
					Lastname:    "Last",
					Firstname:   sql.NullString{String: "First", Valid: true},
					Statuscode:  88,
					Scorepoints: 99,
				})

			if err != nil {

				log.Error("CreateExchange", "error", err)

				return
			}
		}

		log.Info("Insert END")
	*/

	newExchange, err := baseQueries.GetExchangeByShortcode(ctx, "XXX23")

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {

			log.Info("GetExchange: No records")
		} else {

			log.Error("GetExchange", "error", err)

			return
		}
	} else {

		log.Info("DB_Record", "newExchange", newExchange)
	}
	/*
	   exchanges, err := baseQueries.ListExchanges(ctx)

	   if err != nil {

	   		log.Error("ListExchanges", "error", err)

	   		return
	   	}

	   for _, ex := range exchanges {

	   		fmt.Println("Ex: ", ex)
	   	}
	*/
}

func main() {

	//-------------------------------------------------------------------------
	// Create a new structured logger that writes JSON to stdout.
	//-------------------------------------------------------------------------

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	helper(log)

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

	serverHTTP := &http.Server{

		Handler:           http.HandlerFunc(httpHandler),
		Addr:              ":" + serverconfig.GetHttpPort(),
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       3 * time.Second,
	}

	httpListener, err := net.Listen("tcp", serverHTTP.Addr)

	if err != nil {

		log.Error("failed to bind http", "error", err)

		return
	}

	//-------------------------------------------------------------------------
	// Create a linstener to listen on the https port
	//-------------------------------------------------------------------------

	serverHTTPS := &http.Server{

		Addr:              ":" + serverconfig.GetHttpsPort(),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           http.HandlerFunc(httpsHandler),
	}

	httpsListener, err := net.Listen("tcp", serverHTTPS.Addr)

	if err != nil {

		httpListener.Close()

		log.Error("failed to bind https", "error", err)

		return
	}

	//-------------------------------------------------------------------------

	cert, err := tls.LoadX509KeyPair("certs/cert.pem", "certs/cert-key.pem")

	if err != nil {

		httpListener.Close()
		httpsListener.Close()

		log.Error("failed to load TLS certificate", "error", err)

		return
	}

	httpsListener = tls.NewListener(httpsListener, &tls.Config{Certificates: []tls.Certificate{cert}})

	//-------------------------------------------------------------------------
	// Log successful bindings
	//-------------------------------------------------------------------------

	log.Info("http server listening", "addr", httpListener.Addr().String())
	log.Info("https server listening", "addr", httpsListener.Addr().String())

	//-------------------------------------------------------------------------
	// http-server for automatic redirect from http to https
	//-------------------------------------------------------------------------

	errGrp.Go(func() error {

		if err := serverHTTP.Serve(httpListener); err != nil {

			if errors.Is(err, http.ErrServerClosed) {

				log.Info("http-server closed (gracefully)")

				return nil
			} else {

				log.Error("http server returned error", "error", err)

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

				log.Info("https-server closed (gracefully)")

				return nil
			} else {

				log.Error("https server returned error", "error", err)

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

		log.Info("Shutdown signal received or error detected. Servers are shutting down...")

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

		log.Error("Server terminated with an error", "error", err)
	} else {

		log.Info("Server shut down normally")
	}
}
