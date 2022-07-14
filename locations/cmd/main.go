package main

import (
	"USATUKirill96/gridgo/locations/internal"
	"USATUKirill96/gridgo/locations/pkg/location"
	pb "USATUKirill96/gridgo/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	locations := location.NewFakeRepository()
	ls := location.Service{Locations: locations}

	lis, _ := net.Listen("tcp", ":8002")
	s := grpc.NewServer()
	pb.RegisterLocationsServer(s, &internal.Server{LocationService: ls})
	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
