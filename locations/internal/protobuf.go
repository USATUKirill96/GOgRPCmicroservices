package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	pb "USATUKirill96/gridgo/protobuf"
	"context"
)

type Server struct {
	pb.UnimplementedLocationsServer
	LocationService location.Service
}

func (s *Server) Insert(_ context.Context, l *pb.NewLocation) (*pb.Empty, error) {
	err := s.LocationService.InsertLocation(l.Username, l.Longitude, l.Latitude)
	return &pb.Empty{}, err
}
