package location

import (
	"sync"
	"testing"
	"time"
)

// Use https://www.movable-type.co.uk/scripts/latlong.html to check calculated distances

func getFakeLocationService(locations ...Location) Service {
	repository := NewFakeRepository()

	if len(locations) > 0 {
		repository.locations = locations
	}

	service := Service{repository}
	return service
}

func TestService_InsertLocation(t *testing.T) {
	service := getFakeLocationService()
	var (
		username          = "ServiceTest"
		longitude float64 = 15.3451
		latitude  float64 = 19.5312
	)

	err := service.InsertLocation(username, longitude, latitude)
	if err != nil {
		t.Errorf("Insert location: %v", err)
	}

	// Assert location
	var found bool
	for _, l := range service.Locations.(*FakeRepository).locations {
		if l.Username == username && l.Longitude == longitude && l.Latitude == latitude {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Insert location: location wasn't saved")
	}

}

func TestService_GetDistance(t *testing.T) {

	var day1 time.Time = time.Date(2021, 8, 10, 0, 0, 0, 100, time.Local)
	var day2 time.Time = time.Date(2021, 8, 20, 0, 0, 0, 100, time.Local)
	var day3 time.Time = time.Date(2021, 8, 30, 0, 0, 0, 100, time.Local)

	ls := []Location{
		{"0", "User0", 10.54321, -10.54321, day1},
		{"1", "User0", 12.54321, -12.54321, day1},
		{"2", "User1", 12.54321, -12.54321, day1},
		{"3", "User1", 15.54321, -15.54321, day2},
		{"4", "User1", 17.54321, -17.54321, day3},
	}
	service := getFakeLocationService(ls...)

	cases := []struct {
		username string
		after    time.Time
		before   time.Time
		expected int
	}{
		{"User0", day1, day2, 311}, // exclude User1
		{"User1", day1, day3, 773}, // include not-equal time
		{"User1", day2, day3, 308}, // exclude day 1
		{"User1", day3, day3, 0},   // only one location
	}

	for _, tc := range cases {
		distance, err := service.GetDistance(tc.username, tc.after, tc.before)
		if err != nil {
			t.Error(err)
		}
		if distance != tc.expected {
			t.Errorf("GetDistance: expected %v, got %v", tc.expected, distance)
		}
	}
}

func TestService_calculateDistance(t *testing.T) {

	cases := []struct {
		from     Location
		to       Location
		expected int
	}{
		{
			Location{"2", "User1", 12.54321, -12.54321, time.Now()},
			Location{"3", "User1", 15.54321, -15.54321, time.Now()},
			465249,
		},
		{
			Location{"3", "User1", 15.54321, -15.54321, time.Now()},
			Location{"4", "User1", 17.54321, -17.54321, time.Now()},
			308399,
		},
	}

	for _, tc := range cases {
		ch := make(chan float64, 1)
		wg := sync.WaitGroup{}
		wg.Add(1)
		go calculateDistance(tc.from, tc.to, ch, &wg)
		wg.Wait()
		res := <-ch
		close(ch)
		if int(res) != tc.expected {
			t.Errorf("calculateDistance: expected %v, got %v", tc.expected, int(res))
		}
	}
}

func TestService_sumDistance(t *testing.T) {
	inputs := []float64{149.5, 149.5, 200, 300, 600, 1}
	expected := 1400

	ich := make(chan float64)
	och := make(chan int)

	go sumDistances(ich, och)

	for _, v := range inputs {
		ich <- v
	}
	close(ich)
	res := <-och
	if res != expected {
		t.Errorf("sumDistance: expected %v, got %v", expected, res)
	}
}
