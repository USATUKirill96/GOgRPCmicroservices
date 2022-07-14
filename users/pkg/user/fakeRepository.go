package user

import (
	"sync"
)

// NewFakeRepository creates a Fake repository for testing
func NewFakeRepository() FakeRepository {

	userFixtures := map[string]User{
		"User0": User{0, "User0", 12.4321, -28.1635},
		"User1": User{1, "User1", 25.1628, 17.4351},
		"User2": User{2, "User3", -45.0256, -8.4321},
	}
	return FakeRepository{
		users:      userFixtures,
		lastUserID: 2,
	}
}

type FakeRepository struct {
	users      map[string]User
	lastUserID int
	mu         sync.Mutex
}

// Insert adds new services to the fake database by his username
func (r *FakeRepository) Insert(username string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Check if services already exists
	_, exists := r.users[username]
	if exists {
		return nil, AlreadyExists
	}
	// Increment ID and add to storage
	r.lastUserID += 1
	newUser := User{Username: username, ID: r.lastUserID}
	r.users[username] = newUser

	return &newUser, nil
}

// Get returns a User value if exists in database
func (r *FakeRepository) Get(username string) (*User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, NotFound
	}
	return &user, nil
}

func (r *FakeRepository) UpdateLocation(u *User, long float32, lat float32) (*User, error) {
	_, exists := r.users[u.Username]
	if !exists {
		return nil, NotFound
	}
	u.Latitude = lat
	u.Longitude = long
	r.users[u.Username] = *u
	return u, nil
}
