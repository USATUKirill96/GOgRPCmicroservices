package location

type ServiceRepository interface {
	Insert(location Location) (Location, error)
}

type Service struct {
	Locations ServiceRepository
}

func (s Service) InsertLocation(username string, longitude, latitude float32) error {
	l := Location{
		Username:  username,
		Longitude: longitude,
		Latitude:  latitude,
	}
	_, err := s.Locations.Insert(l)
	return err
}
