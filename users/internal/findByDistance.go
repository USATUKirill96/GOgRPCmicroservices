package internal

import (
	"USATUKirill96/gridgo/users/pkg/user"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type FindByDistanceInput struct {
	Username string `json:"username"`
	Distance int    `json:"distance"`
}

func (input FindByDistanceInput) validate() map[string]string {

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
	if input.Distance < 0 {
		errs["distance"] = "Distance must be an integer larger than zero"
	}
	return errs
}

func ParseFindByDistanceInput(r *http.Request) (*FindByDistanceInput, map[string]string) {
	errors := make(map[string]string)
	q := r.URL.Query()
	dst, err := strconv.Atoi(q.Get("distance"))
	if err != nil {
		errors["distance"] = "incorrect format. Integer values allowed"
	}
	i := &FindByDistanceInput{
		Username: q.Get("username"),
		Distance: dst,
	}
	return i, errors
}

func (app Application) FindByDistance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	i, errrs := ParseFindByDistanceInput(r)
	if len(errrs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errrs})
		return
	}

	errrs = i.validate()
	if len(errrs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": errrs})
		return
	}
	users, err := app.UserService.FindByDistance(i.Username, i.Distance)
	if err != nil {
		if errors.Is(err, user.NotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{"errors": "User doesn't exist"})
			return
		}
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": "Internal error. Please try again later"})
		return
	}
	var res []string
	for _, u := range users {
		res = append(res, u.Username)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"result": res})
}
