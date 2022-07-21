package internal

import (
	"USATUKirill96/gridgo/users/pkg/location"
	"USATUKirill96/gridgo/users/pkg/user"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplication_FindByDistance(t *testing.T) {
	ur := user.NewFakeRepository()
	lr := location.NewFakeLocationClient()
	s := user.Service{Users: &ur, Locations: lr}
	app := Application{UserService: s}

	// load fixtures
	us := []*user.User{
		{0, "User1", 5.000, 3.000},
		{0, "User2", 5.000, 4.000},
	}
	for _, u := range us {
		_, err := ur.Insert(*u)
		if err != nil {
			t.Error(err)
		}
	}

	cases := []struct {
		Username       string
		Distance       interface{}
		ExpectedStatus int
	}{
		{"User1", 100, http.StatusOK},
		{"User12345678901234567", 100, http.StatusBadRequest},
		{"Us", 100, http.StatusBadRequest},
		{"User@", 100, http.StatusBadRequest},
		{"User1", -250, http.StatusBadRequest},
		{"NotExisting", 100, http.StatusNotFound},
		{"User1", "NotANumber", http.StatusBadRequest},
	}

	for _, tc := range cases {

		w := httptest.NewRecorder()
		params := fmt.Sprintf("username=%v&distance=%v", tc.Username, tc.Distance)
		r, _ := http.NewRequest(
			"GET",
			fmt.Sprintf("http://localhost:8080/users/?%v", params),
			nil,
		)
		app.FindByDistance(w, r)

		if w.Code != tc.ExpectedStatus {
			t.Errorf(
				"Unexpected status code from FindByDistance. Expected: %v, got: %v", tc.ExpectedStatus, w.Code,
			)
		}
	}
}
