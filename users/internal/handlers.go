package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
)

type UpdateLocationInput struct {
	Username  string  `json:"username"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (input UpdateLocationInput) validate() map[string]string {
	errs := make(map[string]string)

	if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(input.Username) {
		errs["username"] = "Username contains forbidden characters. Only letters and numbers allowed"
	}
	if math.Abs(input.Latitude) > 90 {
		errs["latitude"] = "Latitude is incorrect. Values within -90 and 90 allowed"
	}
	if math.Abs(input.Longitude) > 180 {
		errs["longitude"] = "Longitude is incorrect. Values within -180 and 180 allowed"
	}

	return errs
}

func (app Application) UpdateLocation(w http.ResponseWriter, r *http.Request) {

	var input UpdateLocationInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Incorrect input", err)))
		return
	}

	errs := input.validate()
	if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, _ := json.Marshal(errs)
		w.Write(jsonResp)
		return
	}
	err = app.UserService.UpdateLocation(input.Username, input.Longitude, input.Latitude)

	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte("Location updated"))
}

type FindByDistanceInput struct {
	Username string `json:"username"`
	Distance int    `json:"distance"`
}

func (input FindByDistanceInput) validate() map[string]string {

	errs := make(map[string]string)

	if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(input.Username) {
		errs["username"] = "Username contains forbidden characters. Only letters and numbers allowed"
	}
	if input.Distance < 0 {
		errs["distance"] = "Distance must be an integer larger than zero"
	}
	return errs
}

func ParseFindByDistanceInput(r *http.Request) (*FindByDistanceInput, error) {
	q := r.URL.Query()
	dst, err := strconv.Atoi(q.Get("distance"))
	if err != nil {
		return nil, err
	}
	i := &FindByDistanceInput{
		Username: q.Get("username"),
		Distance: dst,
	}
	return i, nil
}

func (app Application) FindByDistance(w http.ResponseWriter, r *http.Request) {
	i, err := ParseFindByDistanceInput(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	errs := i.validate()
	if len(errs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(errs))) // -------------------------------------- FIXME ------------------
		return
	}
	users, err := app.UserService.FindByDistance(i.Username, i.Distance)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var res []string
	for _, u := range users {
		res = append(res, u.Username)
	}
	w.Write([]byte(fmt.Sprint(res)))
}
