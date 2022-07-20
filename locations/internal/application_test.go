package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var run = flag.Bool("integration", false, "If integral testing required")

// TestApplication_GetDistance is an E2E test, so it'll be only run if specific flag provided
func TestApplication_GetDistance(t *testing.T) {

	if !*run {
		fmt.Println("Skipping TestApplication_GetDistance")
		return
	}

	ti := location.NewTestIndex()
	r := ti.NewRepository()
	defer ti.TearDown()

	ls := []location.Location{
		{"0", "User0", 10.54321, -10.54321, day1},
		{"1", "User0", 12.54321, -12.54321, day1},
		{"2", "User1", 12.54321, -12.54321, day1},
		{"3", "User1", 15.54321, -15.54321, day2},
		{"4", "User1", 17.54321, -17.54321, day3},
	}

	for _, l := range ls {
		r.Insert(l)
	}

	params := fmt.Sprintf("username=%v&after=%v", "User1", day1.Format(time.RFC3339))
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost/?%v", params),
		nil,
	)
	w := httptest.NewRecorder()

	s := location.Service{Locations: r}
	app := Application{LocationService: s}
	app.GetDistance(w, req)

	res := w.Result()
	defer res.Body.Close()

	if w.Code != http.StatusOK {
		t.Errorf("Unexpected status. Expected: %v, got: %v", http.StatusOK, w.Code)
	}
	response := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&response)
	if response["distance"] != float64(773) {
		t.Errorf("Value doesn't match. Expected 773, got %v", response["distance"])
	}
}
