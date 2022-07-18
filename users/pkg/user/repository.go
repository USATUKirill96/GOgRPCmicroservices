package user

import (
	"database/sql"
	"errors"
)

func NewRepository(db *sql.DB) Repository { return Repository{db} }

type Repository struct {
	db *sql.DB
}

func (r Repository) Insert(u User) (*User, error) {
	stmt := `
	   INSERT INTO app_user (username, longitude, latitude) 
	   VALUES ($1, $2, $3) 
	RETURNING id
	`
	var id int
	err := r.db.QueryRow(stmt, u.Username, u.Longitude, u.Latitude).Scan(&id)
	if err != nil {
		return nil, err
	}
	u.ID = id

	return &u, nil
}

func (r Repository) ByUsername(username string) (*User, error) {
	stmt := `
	SELECT id, username, longitude, latitude
	  FROM app_user
	 WHERE username = $1
	`

	u := &User{}
	err := r.db.QueryRow(stmt, username).Scan(&u.ID, &u.Username, &u.Longitude, &u.Latitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (r Repository) Update(u User) (*User, error) {
	if u.ID == 0 {
		return nil, CannotUpdate
	}

	stmt := `
	UPDATE app_user 
	   SET longitude = $2, latitude = $3
	 WHERE id = $1
	`
	err := r.db.QueryRow(stmt, u.ID, u.Longitude, u.Latitude).Err()
	if err != nil {
		return nil, err
	}
	return &u, nil
}
