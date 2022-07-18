package main

import (
	"USATUKirill96/gridgo/locations/internal"
	"USATUKirill96/gridgo/locations/pkg/location"
	pb "USATUKirill96/gridgo/protobuf"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
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

	cfg := elasticsearch.Config{Addresses: []string{os.Getenv("ELASTICSEARCH_URL")}}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	locations := location.NewRepository(es, ElasticsearchIndex)
	ls := location.Service{Locations: locations}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_GRPC")))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterLocationsServer(s, &internal.Server{LocationService: ls})
	log.Printf("Server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	app := internal.Application{
		LocationService: location.Service{Locations: locations},
	}

	r := mux.NewRouter()
	r.HandleFunc("/", app.GetDistance)
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_PORT")),
	}

	fmt.Printf("Server started and running at %v", srv.Addr)
	err = srv.ListenAndServe()
	fmt.Println(err)

}
