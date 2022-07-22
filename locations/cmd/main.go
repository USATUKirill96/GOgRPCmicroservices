package main

import (
	"USATUKirill96/gridgo/locations/internal"
	"USATUKirill96/gridgo/locations/pkg/location"
	"USATUKirill96/gridgo/tools/logging"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/joho/godotenv"
	"log"
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

	gRPCSrv := internal.Server{
		Logger:          logger,
		LocationService: ls,
	}
	go gRPCSrv.Serve()

	app := internal.Application{
		LocationService: location.Service{Locations: locations},
		Logger:          logger,
	}
	err = app.Serve()
	if err != nil {
		logger.ERROR(err)
	}
}
