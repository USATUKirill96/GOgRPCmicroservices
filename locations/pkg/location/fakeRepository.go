package location

import (
	"fmt"
	"math/rand"
	"time"
)

func NewFakeRepository() *FakeRepository { return &FakeRepository{} }

type FakeRepository struct {
	locations []Location
}

func (r *FakeRepository) Insert(l Location) (Location, error) {
	l.ID = fmt.Sprint(rand.Uint64())
	r.locations = append(r.locations, l)
	return l, nil
}

func (r *FakeRepository) Find(u string, after, before time.Time) ([]Location, error) {
	var locations []Location
	for _, l := range r.locations {
		if l.Username == u &&
			(l.Updated.After(after) || l.Updated.Equal(after)) &&
			(l.Updated.Before(before) || l.Updated.Equal(before)) {
			locations = append(locations, l)
		}
	}
	return locations, nil
}
