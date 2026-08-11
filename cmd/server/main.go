package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ForTheTrashBin/RbsClone/internal/osspecific"
	"github.com/ForTheTrashBin/RbsClone/internal/serverconfig"
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

		Addr:              ":" + serverconfig.GetHttpPort(),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           http.HandlerFunc(httpHandler),
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

		// TODO: Zertificates from "letsencrypt"

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
