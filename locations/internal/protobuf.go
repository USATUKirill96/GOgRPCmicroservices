package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	pb "USATUKirill96/gridgo/protobuf"
	"USATUKirill96/gridgo/tools/logging"
	"context"
	"errors"
	"fmt"
)

type Server struct {
	pb.UnimplementedLocationsServer
	LocationService location.Service
	Logger          logging.Logger
}

func (s *Server) Insert(_ context.Context, l *pb.NewLocation) (*pb.Empty, error) {
	err := s.LocationService.InsertLocation(l.Username, l.Longitude, l.Latitude)
	if err != nil {
		s.Logger.ERROR(errors.New(fmt.Sprintf("Inserting location: %v \n, %v", l, err)))
	}
	return &pb.Empty{}, err
}
