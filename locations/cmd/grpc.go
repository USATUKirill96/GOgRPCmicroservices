package main

import (
	"USATUKirill96/gridgo/locations/internal"
	"USATUKirill96/gridgo/locations/pkg/location"
	pb "USATUKirill96/gridgo/protobuf"
	"USATUKirill96/gridgo/tools/logging"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
)

// rungRPS parses .env variables and starts gRPC server
// Warning: the function is blocking. Only run in a separated goroutine
func rungRPS(ls location.Service, log logging.Logger) {

	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_GRPC")),
	)

	if err != nil {
		log.ERROR(err)
	}

	s := grpc.NewServer()
	pb.RegisterLocationsServer(s, &internal.Server{LocationService: ls, Logger: log})
	log.INFO(fmt.Sprintf("gRPC server listening at %v \n", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		log.ERROR(errors.New(fmt.Sprintf("failed to serve: %v", err)))
	}
}
