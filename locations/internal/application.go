package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	"USATUKirill96/gridgo/tools/logging"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	LocationService location.Service
	Logger          logging.Logger
}

func (app Application) Routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/distance", app.GetDistance).Methods("GET")
	r.Use(app.LogRequests)
	r.Use(app.RecoverPanic)
	return r
}

func (app Application) Serve() error {
	srv := &http.Server{
		Handler: app.Routes(),
		Addr:    fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_PORT")),

		IdleTimeout:  time.Minute,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 5,
	}

	shutdownError := make(chan error)

	go func() {

		// Catch an exit signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		// Log the information
		app.Logger.INFO(fmt.Sprintf("Shutting down the server. Signal: %v", s.String()))

		// 5 seconds timeout to exit
		//Why 5 seconds: https://github.com/golang/go/issues/33191
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// if were any errors, send to the channel
		shutdownError <- srv.Shutdown(ctx)

		os.Exit(0)
	}()

	app.Logger.INFO(fmt.Sprintf("HTTP Server started and running at %v \n", srv.Addr))
	err := srv.ListenAndServe()
	// Calling `Shutdown()` will immediately return http.ErrServerClosed. There is no reason to log it as an error
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	return nil
}

func (app Application) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger.INFO("", r)
		next.ServeHTTP(w, r)
	})
}

func (app Application) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]interface{}{"errors": "Internal server error"})
				app.Logger.ERROR(err.(error))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
