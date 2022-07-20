package main

import (
	"USATUKirill96/gridgo/locations/internal"
	"USATUKirill96/gridgo/locations/pkg/location"
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
	locations := location.NewRepository(es, ElasticsearchIndex)
	ls := location.Service{Locations: locations}

	go rungRPS(ls) // Running server for internal communication

	app := internal.Application{
		LocationService: location.Service{Locations: locations},
	}

	r := mux.NewRouter()
	r.HandleFunc("/distance", app.GetDistance).Methods("GET")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_PORT")),
	}

	fmt.Printf("HTTP Server started and running at %v \n", srv.Addr)
	err = srv.ListenAndServe()
	fmt.Println(err)

}
