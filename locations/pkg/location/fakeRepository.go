package location

func NewFakeRepository() *FakeRepository { return &FakeRepository{} }

type FakeRepository struct {
	locations []Location
	lastID    int
}

func (r *FakeRepository) Insert(l Location) (Location, error) {
	r.lastID += 1
	l.ID = r.lastID
	r.locations = append(r.locations, l)
	return l, nil
}
