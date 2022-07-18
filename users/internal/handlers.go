package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
)

type UpdateLocationInput struct {
	Username  string  `json:"username"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (input UpdateLocationInput) Validate() map[string]string {
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

	errs := input.Validate()
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
