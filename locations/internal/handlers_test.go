package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var day1 time.Time = time.Date(2021, 8, 10, 0, 0, 0, 100, time.Local)
var day2 time.Time = time.Date(2021, 8, 20, 0, 0, 0, 100, time.Local)
var day3 time.Time = time.Date(2021, 8, 30, 0, 0, 0, 100, time.Local)

func NewFakeApplication() Application {
	ls := []location.Location{
		{"0", "User0", 10.54321, -10.54321, day1},
		{"1", "User0", 12.54321, -12.54321, day1},
		{"2", "User1", 12.54321, -12.54321, day1},
		{"3", "User1", 15.54321, -15.54321, day2},
		{"4", "User1", 17.54321, -17.54321, day3},
	}
	r := location.NewFakeRepository()
	s := location.Service{Locations: r}
	for _, l := range ls {
		r.Insert(l)
	}
	return Application{LocationService: s}
}

func TestApplication_GetDistance_Correct(t *testing.T) {
	app := NewFakeApplication()
	params := fmt.Sprintf("username=%v&after=%v", "User1", day1.Format(time.RFC3339))
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost/?%v", params),
		nil,
	)
	w := httptest.NewRecorder()
	app.GetDistance(w, req)
	res := w.Result()
	defer res.Body.Close()

	if w.Code != http.StatusOK {
		t.Errorf("Unexpected status. Expected: %v, got: %v", http.StatusOK, w.Code)
	}
	response := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&response)
	if response["distance"] != 773 {
		t.Errorf("Value doesn't match. Expected 773, got %v", response["distance"])
	}
}

func TestApplication_GetDistance_IncorrectUsername(t *testing.T) {
	app := NewFakeApplication()
	params := fmt.Sprintf("username=%v&after=%v", "Unexisted", day1.Format(time.RFC3339))
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost/?%v", params),
		nil,
	)
	w := httptest.NewRecorder()
	app.GetDistance(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Unexpected status. Expected: %v, got: %v", http.StatusNotFound, w.Code)
	}
}

func TestApplication_GetDistance_IncorrectTimeFormat(t *testing.T) {
	app := NewFakeApplication()
	params := fmt.Sprintf("username=%v&after=%v", "User1", "2021-09 -02T11: 26: 18+00:00")
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost/?%v", params),
		nil,
	)
	w := httptest.NewRecorder()
	app.GetDistance(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Unexpected status. Expected: %v, got: %v", http.StatusNotFound, w.Code)
	}
}
