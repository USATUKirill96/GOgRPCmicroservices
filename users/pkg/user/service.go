package user

import (
	pb "USATUKirill96/gridgo/protobuf"
	"context"
	"errors"
	"time"
)

type ServiceRepository interface {
	Get(string) (*User, error)
	Insert(string) (*User, error)
	UpdateLocation(*User, float32, float32) (*User, error)
}

type Service struct {
	Users     ServiceRepository
	Locations pb.LocationsClient
}

// UpdateLocation updated location of a user by its username
// SIDE EFFECT: creates a new user if username doesn't match any existing
// TODO: discuss risks of the side effect, possibly use a separated endpoint/service for it
func (s Service) UpdateLocation(username string, longitude float32, latitude float32) error {
	u, err := s.Users.Get(username)
	if err != nil {
		if !errors.Is(NotFound, err) {
			return err
		}

		u, err = s.Users.Insert(username)
		if err != nil {
			return err
		}
	}
	_, err = s.Users.UpdateLocation(u, longitude, latitude)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	data := &pb.NewLocation{
		Username:  u.Username,
		Latitude:  latitude,
		Longitude: longitude,
	}
	_, err = s.Locations.Insert(ctx, data)
	if err != nil {
		return err
	}
	return nil
}
