package location

import (
	"fmt"
	"math/rand"
	"time"
)

// NewFakeRepository creates a mocked repository which stores data in memory.
// Optional parameter err might be provided. It contains an error, which must be returned from the repository
// (for error-handling cases)
func NewFakeRepository(err ...error) *FakeRepository {
	var e error
	if len(err) > 0 {
		e = err[0]
	} else {
		e = nil
	}
	return &FakeRepository{err: e}
}

type FakeRepository struct {
	Locations []Location
	err       error
}

func (r *FakeRepository) Insert(l Location) (Location, error) {
	if r.err != nil {
		return l, r.err
	}
	l.ID = fmt.Sprint(rand.Uint64())
	r.Locations = append(r.Locations, l)
	return l, nil
}

func (r *FakeRepository) Find(u string, after, before time.Time) ([]Location, error) {
	if r.err != nil {
		return []Location{}, r.err
	}
	var locations []Location
	for _, l := range r.Locations {
		if l.Username == u &&
			(l.Updated.After(after) || l.Updated.Equal(after)) &&
			(l.Updated.Before(before) || l.Updated.Equal(before)) {
			locations = append(locations, l)
		}
	}
	return locations, nil
}
