//go:build !integration

package user

import (
	"USATUKirill96/gridgo/users/pkg/location"
	"USATUKirill96/gridgo/users/pkg/pagination"
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

func TestService_FindByDistance(t *testing.T) {

	const (
		close1 = iota
		close2
		far1
		far2
	)

	neighbors := []*User{
		&User{0, "Close1", 5.000, 3.000}, // 96.63 km
		&User{0, "Close2", 5.000, 4.000}, // 95.05 km
		&User{0, "Far1", 6.000, 3.000},   // 194.4 km
		&User{0, "Far2", 6.000, 6.000},   // 335.7 km
	}

	service := getFakeUserService()

	for _, n := range neighbors {
		_, err := service.Users.Insert(*n)
		if err != nil {
			log.Fatal(err)
		}
	}

	cases := []struct {
		distance int
		expected []*User
	}{
		{50, []*User{}},
		{100, []*User{neighbors[close1], neighbors[close2]}},
		{200, []*User{neighbors[close1], neighbors[close2], neighbors[far1]}},
		{400, []*User{neighbors[close1], neighbors[close2], neighbors[far1], neighbors[far2]}},
	}

	for _, tc := range cases {
		result, err := service.FindByDistance(userServiceTestUser.Username, tc.distance, pagination.Pagination{})
		if err != nil {
			t.Error(err)
		}
		// TODO: make special checks if exact requested users return from service
		if len(result) != len(tc.expected) {
			t.Errorf("Unexpected returns. Expected: %v, got: %v", len(tc.expected), len(result))
		}
	}
}
