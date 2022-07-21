package main

import (
	"USATUKirill96/gridgo/locations/internal"
	"USATUKirill96/gridgo/locations/pkg/location"
	"USATUKirill96/gridgo/tools/logging"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var ElasticsearchIndex = "locations_1.0"

func main() {
	// Load .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Setup elasticsearch client and repository
	cfg := elasticsearch.Config{Addresses: []string{os.Getenv("ELASTICSEARCH_URL")}}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	logger := logging.NewLogger(es)
	locations := location.NewRepository(es, ElasticsearchIndex)
	ls := location.Service{Locations: locations}

	go rungRPS(ls, logger) // Running server for internal communication

	app := internal.Application{
		LocationService: location.Service{Locations: locations},
		Logger:          logger,
	}

	r := mux.NewRouter()
	r.HandleFunc("/distance", app.GetDistance).Methods("GET")
	r.Use(app.LogRequests)
	r.Use(app.RecoverPanic)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_PORT")),
	}

	logger.INFO(fmt.Sprintf("HTTP Server started and running at %v \n", srv.Addr))
	err = srv.ListenAndServe()

}
