package location

import "testing"

func getFakeLocationService() Service {
	repository := NewFakeRepository()
	service := Service{repository}
	return service
}

func TestService_InsertLocation(t *testing.T) {
	service := getFakeLocationService()
	var (
		username          = "ServiceTest"
		longitude float32 = 15.3451
		latitude  float32 = 19.5312
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
