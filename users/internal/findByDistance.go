package internal

import (
	"USATUKirill96/gridgo/users/pkg/pagination"
	"USATUKirill96/gridgo/users/pkg/user"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type FindByDistanceInput struct {
	Username   string                `json:"username"`
	Distance   int                   `json:"distance"`
	Pagination pagination.Pagination `json:"-"`
}

func (input FindByDistanceInput) validate() map[string]interface{} {

	errorsMap := make(map[string]interface{})

	if len([]rune(input.Username)) < 4 {
		errorsMap["username"] = "Username is too short. Make it 4 symbols or longer"
	}
	if len([]rune(input.Username)) > 16 {
		errorsMap["username"] = "Username is too long. Make it 16 symbols or shorter"
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(input.Username) {
		errorsMap["username"] = "Username contains forbidden characters. Only letters and numbers allowed"
	}
	if input.Distance < 0 {
		errorsMap["distance"] = "Distance must be an integer larger than zero"
	}
	err := input.Pagination.Validate()
	if err != nil {
		errorsMap["pagination"] = err.(pagination.Error).Map
	}

	return errorsMap
}

func ParseFindByDistanceInput(r *http.Request) (*FindByDistanceInput, map[string]interface{}) {
	errorsMap := make(map[string]interface{})
	q := r.URL.Query()
	dst, err := strconv.Atoi(q.Get("distance"))
	if err != nil {
		errorsMap["distance"] = "incorrect format. Integer values allowed"
	}

	pg, err := pagination.FromQuery(q)
	if err != nil {
		errorsMap["pagination"] = err.(pagination.Error).Map
	}

	i := &FindByDistanceInput{
		Username:   q.Get("username"),
		Distance:   dst,
		Pagination: pg,
	}
	return i, errorsMap
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
	users, err := app.UserService.FindByDistance(i.Username, i.Distance, i.Pagination)
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
	// if no users found, there will be an empty array in result.
	// Do not change to "var" declaration, in this case a null will be passed to result instead of empty array
	res := []string{}
	for _, u := range users {
		res = append(res, u.Username)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"result": res})
}
