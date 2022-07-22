package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	"USATUKirill96/gridgo/tools/logging"
	"encoding/json"
	"net/http"
)

type Application struct {
	LocationService location.Service
	Logger          logging.Logger
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
