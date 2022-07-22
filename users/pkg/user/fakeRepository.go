package user

import (
	"USATUKirill96/gridgo/locations/pkg/location"
	"USATUKirill96/gridgo/users/pkg/pagination"
	"sync"
)

// NewFakeRepository creates a Fake repository for testing
func NewFakeRepository() FakeRepository {

	return FakeRepository{
		users:      make(map[string]User),
		lastUserID: 0,
	}
}

type FakeRepository struct {
	users      map[string]User
	lastUserID int
	mu         sync.Mutex
}

// Insert adds new services to the fake database by his username
func (r *FakeRepository) Insert(u User) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Check if services already exists
	_, exists := r.users[u.Username]
	if exists {
		return nil, AlreadyExists
	}
	// Increment ID and add to storage
	r.lastUserID += 1
	u.ID = r.lastUserID
	r.users[u.Username] = u

	return &u, nil
}

func (r *FakeRepository) Update(u User) (*User, error) {
	_, exists := r.users[u.Username]
	if !exists {
		return nil, NotFound
	}
	r.users[u.Username] = u
	return &u, nil
}

// ByUsername returns a User value if exists in database
func (r *FakeRepository) ByUsername(username string) (*User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, NotFound
	}
	return &user, nil
}

func (r *FakeRepository) ByDistance(tgu User, dst int, pg pagination.Pagination) ([]*User, error) {

	var matches []*User
	for _, u := range r.users {
		d := location.Distance{
			FromLon: tgu.Longitude,
			ToLon:   u.Longitude,
			FromLat: tgu.Latitude,
			ToLAt:   u.Latitude,
		}
		if d.Meters() <= float64(dst*1000) && u.Username != tgu.Username {
			matches = append(matches, &u)
		}
	}
	if pg.Offset != 0 {
		if pg.Offset > len(matches) {
			matches = []*User{}
		}
		matches = matches[pg.Offset-1:]
	}
	if pg.Limit != 0 {
		if pg.Limit < len(matches) {
			matches = matches[:pg.Limit]
		}
	}
	return matches, nil
}
