package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Insecure work completed!"))
	/*
	   // Striktes HTTP -> HTTPS Redirect

	   target := "https://" + r.Host + r.URL.Path

	   if len(r.URL.RawQuery) > 0 {

	   		target += "?" + r.URL.RawQuery
	   	}

	   http.Redirect(w, r, target, http.StatusMovedPermanently) // could be http.StatusTemporaryRedirect too
	*/
}

func httpsHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Secure work completed!"))
}

func main() {

	//-------------------------------------------------------------------------
	// Create a new structured logger that writes JSON to stdout.
	//-------------------------------------------------------------------------

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

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

	errGrp, ctx := errgroup.WithContext(ctx)

	//-------------------------------------------------------------------------
	// http-server for automatic redirect from http to https
	//-------------------------------------------------------------------------

	serverHTTP := &http.Server{

		Addr:              ":8080",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           http.HandlerFunc(httpHandler),
	}

	errGrp.Go(func() error {

		if err := serverHTTP.ListenAndServe(); err != nil {

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

	serverHTTPS := &http.Server{

		Addr:              ":8443",
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           http.HandlerFunc(httpsHandler),
	}

	errGrp.Go(func() error {

		// TODO: Zertifikate bei "letsencrypt"

		// if err := serverHTTPS.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		if err := serverHTTPS.ListenAndServe(); err != nil {

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
		// cause this function to return.
		//-------------------------------------------------------------------------

		<-ctx.Done()

		log.Info("Shutdown signal received. Servers are shutting down...")

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
