package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type GetDistanceInput struct {
	Username string    `json:"username"`
	After    time.Time `json:"after"`
	Before   time.Time `json:"before"`
}

func ParseGetDistanceInput(r *http.Request) (GetDistanceInput, error) {
	q := r.URL.Query()

	var (
		after  time.Time
		before time.Time
	)
	var err error

	user := q.Get("username")
	rawAfter := q.Get("after")
	if rawAfter != "" {
		after, err = time.Parse("2021-09-02T11:26:18+00:00", rawAfter)
		if err != nil {
			return GetDistanceInput{}, err
		}
	}

	rawBefore := q.Get("before")
	if rawBefore != "" {
		before, err = time.Parse("2021-09-02T11:26:18+00:00", rawAfter)
		if err != nil {
			return GetDistanceInput{}, err
		}
	}
	return GetDistanceInput{user, after, before}, nil
}

func (input GetDistanceInput) Validate() map[string]string {
	errs := make(map[string]string)
	if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(input.Username) {
		errs["username"] = "Username contains forbidden characters. Only letters and numbers allowed"
	}
	if input.Before.After(input.After) {
		errs["range"] = "Time label 'Before' follows 'After'"
	}
	return errs
}

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

	input, err := ParseGetDistanceInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, _ := json.Marshal(err)
		w.Write(jsonResp)
	}
	errs := input.Validate()
	if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, _ := json.Marshal(errs)
		w.Write(jsonResp)
		return
	}
	input = input.SetDefaultRange()

	d, err := app.LocationService.GetDistance(input.Username, input.After, input.Before)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	distance := strconv.Itoa(d)
	w.Write([]byte(distance))

}
