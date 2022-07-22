package internal

import (
	"USATUKirill96/gridgo/users/pkg/location"
	"USATUKirill96/gridgo/users/pkg/user"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplication_UpdateLocation(t *testing.T) {
	ur := user.NewFakeRepository()
	lr := location.NewFakeLocationClient()
	s := user.Service{Users: &ur, Locations: lr}
	app := Application{UserService: s}

	cases := []struct {
		Username       string  `json:"username"`
		Longitude      float64 `json:"longitude"`
		Latitude       float64 `json:"latitude"`
		ExpectedStatus int     `json:"-"`
	}{
		{"User1", 15.391, 18.248, 200},               // correct
		{"User123456789012345", 15.391, 18.248, 400}, // username is too long
		{"Us", 15.391, 18.248, 400},                  // username is too short
		{"User@", 15.391, 18.248, 400},               // forbidden characters
		{"User1", -191.391, 18.248, 400},             // incorrect longitude
		{"User1", 15.391, 180.248, 400},              // incorrect latitude
	}

	for _, tc := range cases {
		w := httptest.NewRecorder()
		data, _ := json.Marshal(tc)
		r, _ := http.NewRequest(
			"POST",
			"http://localhost:8080/location",
			bytes.NewBuffer(data),
		)
		app.UpdateLocation(w, r)

		if w.Code != tc.ExpectedStatus {
			t.Errorf(
				"Unexpected status code from UpdateLocation. Expected: %v, got: %v", tc.ExpectedStatus, w.Code,
			)
		}
	}
}
