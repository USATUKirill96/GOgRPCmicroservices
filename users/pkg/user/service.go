package user

import (
	pb "USATUKirill96/gridgo/protobuf"
	"context"
	"errors"
	"time"
)

type ServiceRepository interface {
	ByUsername(string) (*User, error)
	Insert(User) (*User, error)
	Update(User) (*User, error)
}

type Service struct {
	Users     ServiceRepository
	Locations pb.LocationsClient
}

// UpdateLocation updated location of a user by its username
// SIDE EFFECT: creates a new user if username doesn't match any existing
// TODO: discuss risks of the side effect, possibly use a separated endpoint/function for it
func (s Service) UpdateLocation(username string, longitude, latitude float64) error {
	u, err := s.Users.ByUsername(username)
	if err != nil && !errors.Is(NotFound, err) {
		return err
	}

	if u != nil { // User already exists
		u.Longitude = longitude
		u.Latitude = latitude
		u, err = s.Users.Update(*u)
		if err != nil {
			return err
		}
	} else { // Need to create a new user
		u, err = s.Users.Insert(User{Username: username, Longitude: longitude, Latitude: latitude})
		if err != nil {
			return err
		}
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
