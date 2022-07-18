package user

import "errors"

var AlreadyExists error = errors.New("database: User already exists")
var CannotUpdate error = errors.New("database: trying to update an unexciting user")
var NotFound error = errors.New("database: User not found")

type User struct {
	ID        int
	Username  string
	Longitude float64
	Latitude  float64
}
