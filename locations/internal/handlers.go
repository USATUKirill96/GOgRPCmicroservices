package internal

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type GetDistanceInput struct {
	Username string    `json:"username"`
	After    time.Time `json:"after"`
	Before   time.Time `json:"before"`
}

type ParseError struct {
	err map[string]interface{}
}

func (e ParseError) Error() string {
	return fmt.Sprint(e.err)
}

// ParseGetDistanceInput transforms raw *http.Request into GetDistanceInput value
func ParseGetDistanceInput(r *http.Request) (GetDistanceInput, *ParseError) {
	q := r.URL.Query()

	var (
		after  time.Time
		before time.Time
	)
	var err error

	user := q.Get("username")
	rawAfter := q.Get("after")
	if rawAfter != "" {
		// HACK: + symbol stands for space in ASCII, so we need to manually get it back
		// TODO: Check if gorilla mux has a solution for it
		rawAfter = strings.Replace(rawAfter, " ", "+", 1)
		after, err = time.Parse(time.RFC3339, rawAfter)
		if err != nil {
			return GetDistanceInput{},
				&ParseError{map[string]interface{}{"after": "incorrect format. ISO 8601 allowed"}}
		}
	}

	rawBefore := q.Get("before")
	if rawBefore != "" {
		rawBefore = strings.Replace(rawBefore, " ", "+", 1)
		before, err = time.Parse(time.RFC3339, rawBefore)
		if err != nil {
			return GetDistanceInput{},
				&ParseError{map[string]interface{}{"before": "incorrect format. ISO 8601 allowed"}}

		}
	}
	return GetDistanceInput{user, after, before}, nil
}

// Validate ensures that GetDistanceInput contains only valid values and returns a map of errors if any
func (input GetDistanceInput) Validate() map[string]string {
	errs := make(map[string]string)
	if input.Username == "" {
		errs["username"] = "Username must be provided"
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(input.Username) {
		errs["username"] = "Username contains forbidden characters. Only letters and numbers allowed"
	}
	if input.Before.After(input.After) {
		errs["range"] = "Time label 'Before' follows 'After'"
	}
	return errs
}

// SetDefaultRange sets default values such as GetDistanceInput.After and GetDistanceInput.Before
func (input GetDistanceInput) SetDefaultRange() GetDistanceInput {
	if input.After.IsZero() {
		input.After = time.Now().Add(-24 * time.Hour)
	}
	if input.Before.IsZero() {
		input.Before = time.Now()
	}
	return input
}

func (app Application) GetDistance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	input, perr := ParseGetDistanceInput(r)
	if perr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": perr.err})
		return
	}
	errs := input.Validate()
	if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}
	input = input.SetDefaultRange()

	d, err := app.LocationService.GetDistance(input.Username, input.After, input.Before)
	if err != nil {
		if errors.Is(err, location.NotEnoughLocations) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"errors": err})
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"distance": d})

}
