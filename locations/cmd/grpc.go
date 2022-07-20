package main

import (
	"USATUKirill96/gridgo/locations/internal"
	"USATUKirill96/gridgo/locations/pkg/location"
	pb "USATUKirill96/gridgo/protobuf"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// rungRPS parses .env variables and starts gRPC server
// Warning: the function is blocking. Only run in a separated goroutine
func rungRPS(ls location.Service) {

	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_GRPC")),
	)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterLocationsServer(s, &internal.Server{LocationService: ls})
	log.Printf("gRPC server listening at %v \n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
