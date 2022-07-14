package location

import (
	pb "USATUKirill96/gridgo/protobuf"
	"context"
	"google.golang.org/grpc"
)

func NewFakeLocationClient() *FakeLocationClient { return &FakeLocationClient{} }

type FakeLocationClient struct {
	locations []Location
}

func (lc *FakeLocationClient) GetLocations() []Location { return lc.locations }

func (lc *FakeLocationClient) Insert(
	_ context.Context,
	input *pb.NewLocation,
	_ ...grpc.CallOption,
) (*pb.Empty, error) {
	lc.locations = append(lc.locations, Location{input.Username, input.Longitude, input.Latitude})
	return nil, nil
}
