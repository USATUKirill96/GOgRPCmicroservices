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

func (input UpdateLocationInput) validate() map[string]string {
	errs := make(map[string]string)

	if len([]rune(input.Username)) < 4 {
		errs["username"] = "Username is too short. Make it 4 symbols or longer"
	}
	if len([]rune(input.Username)) > 16 {
		errs["username"] = "Username is too long. Make it 16 symbols or shorter"
	}
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
	w.Header().Set("Content-Type", "application/json")

	var input UpdateLocationInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		message := map[string]string{err.(*json.UnmarshalTypeError).Field: "Incorrect value type"}
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": message})
		return
	}

	errs := input.validate()
	if len(errs) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errs})
		return
	}
	err = app.UserService.UpdateLocation(input.Username, input.Longitude, input.Latitude)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"errors": "Internal server problem. Please try again later"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"success": "Your location has been updated"})

}
