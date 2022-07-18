package user

import (
	"USATUKirill96/gridgo/users/pkg/location"
	"log"
	"testing"
)

const (
	longitude = iota + 1
	latitude
)

var userServiceTestUser = User{
	ID:        4,
	Username:  "UserServiceTest",
	Longitude: 4.3191,
	Latitude:  3.4815,
}

func getFakeUserService() Service {
	users := NewFakeRepository()
	_, err := users.Insert(userServiceTestUser)
	if err != nil {
		log.Fatal(err)
	}

	locations := location.NewFakeLocationClient()
	service := Service{Users: &users, Locations: locations}
	return service
}

func TestUserService_UpdateLocation(t *testing.T) {
	service := getFakeUserService()
	coords := []float64{longitude: 31.42561, latitude: -18.24312}

	err := service.UpdateLocation(userServiceTestUser.Username, coords[longitude], coords[latitude])
	if err != nil {
		t.Errorf("Update location: %v", err)
	}

	u, err := service.Users.ByUsername(userServiceTestUser.Username)
	if err != nil {
		t.Errorf("Get services: %v", err)
	}

	assertUserLocationUpdated(u, coords, t)
	assertLocationInserted(service.Locations.(*location.FakeLocationClient), u, t)
}

func assertUserLocationUpdated(u *User, coords []float64, t *testing.T) {
	if u.Longitude != coords[longitude] {
		t.Errorf("Update longitude. %v, got %v", coords[longitude], u.Longitude)
	}
	if u.Latitude != coords[latitude] {
		t.Errorf("Update latitude. Excepted %v, got %v", coords[latitude], u.Latitude)
	}
}

func assertLocationInserted(locations *location.FakeLocationClient, u *User, t *testing.T) {
	for _, l := range locations.GetLocations() {
		if l.Username == u.Username && l.Latitude == u.Latitude && l.Longitude == u.Longitude {
			return
		}
	}
	t.Errorf("Insert location: New location is not found")

}
