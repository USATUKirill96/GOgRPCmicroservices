package location

import (
	"errors"
	"sync"
	"time"
)

var CouldNotInsertLocation = errors.New("insert: could not insert location")

type IServiceRepository interface {
	Insert(location Location) (Location, error)
	Find(string, time.Time, time.Time) ([]Location, error)
}

type Service struct {
	Locations IServiceRepository
}

func (s Service) InsertLocation(username string, longitude, latitude float64) error {
	l := Location{
		Username:  username,
		Longitude: longitude,
		Latitude:  latitude,
		Updated:   time.Now(),
	}
	_, err := s.Locations.Insert(l)
	if err != nil {
		return CouldNotInsertLocation
	}
	return nil
}

var NotEnoughLocations error = errors.New("distance: not enough Locations to calculate distance")

func (s Service) GetDistance(username string, after, before time.Time) (int, error) {
	ls, err := s.Locations.Find(username, after, before)
	if err != nil {
		return 0, err
	}
	// Not enough Locations to calculate distance
	if len(ls) < 2 {
		return 0, NotEnoughLocations
	}

	wg := sync.WaitGroup{}
	dch := make(chan float64) //distances in meters
	sch := make(chan int)     // total calculated distance

	for i := 0; i < len(ls)-1; i++ {
		wg.Add(1)
		go calculateDistance(ls[i], ls[i+1], dch, &wg)
	}
	go sumDistances(dch, sch)

	wg.Wait()
	close(dch)
	totalMeters := <-sch
	totalKilometers := totalMeters / 1000
	return totalKilometers, nil
}

func calculateDistance(lfrom, lto Location, dch chan<- float64, wg *sync.WaitGroup) {
	defer wg.Done()

	d := Distance{
		FromLon: lfrom.Longitude,
		FromLat: lfrom.Latitude,
		ToLon:   lto.Longitude,
		ToLAt:   lto.Latitude,
	}
	dch <- d.Meters()
	return
}

func sumDistances(dch <-chan float64, sch chan<- int) {
	var total float64
	for {
		select {
		case d, ok := <-dch:
			if !ok {
				sch <- int(total)
				close(sch)
				return
			}
			total += d
		}
	}
}
