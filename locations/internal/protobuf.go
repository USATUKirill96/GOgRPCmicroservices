package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	pb "USATUKirill96/gridgo/protobuf"
	"USATUKirill96/gridgo/tools/logging"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
)

type Server struct {
	pb.UnimplementedLocationsServer
	LocationService location.Service
	Logger          logging.Logger
}

func (s *Server) Serve() {

	lis, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%v", os.Getenv("LOCATION_SERVICE_GRPC")),
	)

	if err != nil {
		s.Logger.ERROR(err)
	}

	srv := grpc.NewServer()
	pb.RegisterLocationsServer(srv, s)
	s.Logger.INFO(fmt.Sprintf("gRPC server listening at %v \n", lis.Addr()))
	if err := srv.Serve(lis); err != nil {
		s.Logger.ERROR(errors.New(fmt.Sprintf("failed to serve: %v", err)))
	}
}

func (s *Server) Insert(_ context.Context, l *pb.NewLocation) (*pb.Empty, error) {
	err := s.LocationService.InsertLocation(l.Username, l.Longitude, l.Latitude)
	if err != nil {
		s.Logger.ERROR(errors.New(fmt.Sprintf("Inserting location: %v \n, %v", l, err)))
	}
	return &pb.Empty{}, err
}
