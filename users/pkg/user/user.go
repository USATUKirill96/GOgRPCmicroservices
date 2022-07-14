package user

import "errors"

var AlreadyExists error = errors.New("database: User already exists")
var NotFound error = errors.New("database: User not found")

type User struct {
	ID        int
	Username  string
	Longitude float32
	Latitude  float32
}
