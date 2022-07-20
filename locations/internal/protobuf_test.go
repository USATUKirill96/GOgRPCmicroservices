package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	pb "USATUKirill96/gridgo/protobuf"
	"context"
	"errors"
	"testing"
)

func TestServer_Insert_Correct(t *testing.T) {
	r := location.NewFakeRepository()
	service := location.Service{Locations: r}
	s := Server{LocationService: service}
	l := &pb.NewLocation{Username: "TestServer", Longitude: 15.391, Latitude: 48.358}
	_, err := s.Insert(context.Background(), l)
	if err != nil {
		t.Error(err)
	}
	if len(r.Locations) != 1 {
		t.Errorf("Unexpected number of locations saved. Expected: 1, got: %v", len(r.Locations))
	}
	nl := r.Locations[0]

	if nl.Username != l.Username {
		t.Errorf("Username doesn't match. Expected %v, got %v", l.Username, nl.Username)
	}
	if nl.Longitude != l.Longitude {
		t.Errorf("Longitude doesn't match. Expected %v, got %v", l.Longitude, nl.Longitude)
	}
	if nl.Latitude != l.Latitude {
		t.Errorf("Latitude doesn't match. Expected %v, got %v", l.Latitude, nl.Latitude)
	}
}

func TestServer_Insert_Incorrect(t *testing.T) {
	r := location.NewFakeRepository(location.CouldNotInsertLocation)
	service := location.Service{Locations: r}
	s := Server{LocationService: service}
	l := &pb.NewLocation{Username: "TestServer", Longitude: 15.391, Latitude: 48.358}
	_, err := s.Insert(context.Background(), l)
	if !errors.Is(err, location.CouldNotInsertLocation) {
		t.Errorf("Unexpected return. Expected: %v, got: %v", location.CouldNotInsertLocation, err)
	}
}
