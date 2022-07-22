package user

import (
	pb "USATUKirill96/gridgo/protobuf"
	"USATUKirill96/gridgo/users/pkg/pagination"
	"context"
	"errors"
	"time"
)

const UpdateLocationRetries int = 5

type ServiceRepository interface {
	ByUsername(string) (*User, error)
	ByDistance(User, int, pagination.Pagination) ([]*User, error)
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
	go func() {
		// Network might be unstable.
		// TODO: research and implement a rollback mechanism
		for i := 0; i < UpdateLocationRetries; i++ {
			_, err = s.Locations.Insert(ctx, data)
			if err == nil {
				return
			}
			time.Sleep(2 * time.Second)
		}
	}()
	return nil
}

func (s Service) FindByDistance(username string, distance int, pg pagination.Pagination) ([]*User, error) {
	u, err := s.Users.ByUsername(username)
	if err != nil {
		return nil, err
	}
	neighbors, err := s.Users.ByDistance(*u, distance, pg)
	return neighbors, err
}
