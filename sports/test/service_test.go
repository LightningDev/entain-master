package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/LightningDev/entain-master/sports/proto/sports"
	"github.com/LightningDev/entain-master/sports/service"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

// MockEventsRepo is a mock implementation of the db.EventsRepo interface.
type MockEventsRepo struct{}

func (m *MockEventsRepo) List(filter *sports.ListEventsRequestFilter, orderBy string) ([]*sports.Event, error) {
	// Return hardcoded data instead of querying a database.
	ts, _ := ptypes.TimestampProto(time.Now())
	return []*sports.Event{
		{
			Id:                  1,
			Name:                "Mock Event 1",
			Location:            "Brisbane",
			Sport:               "Bear Fighting",
			AdvertisedStartTime: ts,
			Status:              "Open",
		},
	}, nil
}

// Don't need Init in our mock implementation.
func (m *MockEventsRepo) Init() error {
	return nil
}

func (m *MockEventsRepo) GetByID(id int64) (*sports.Event, error) {
	// Return hardcoded data instead of querying a database.
	ts, _ := ptypes.TimestampProto(time.Now())
	return &sports.Event{
		Id:                  1,
		Name:                "Mock Event 1",
		Location:            "Brisbane",
		Sport:               "Bear Fighting",
		AdvertisedStartTime: ts,
		Status:              "Open",
	}, nil
}

func TestListEvents(t *testing.T) {
	repo := &MockEventsRepo{}
	s := service.NewSportsService(repo)

	req := &sports.ListEventsRequest{
		Filter: &sports.ListEventsRequestFilter{
			Name:     "Mock Event 1",
			Location: "Brisbane",
			Sport:    "Bear Fighting",
		},
		OrderBy: "advertised_start_time",
	}

	events, err := s.ListEvents(context.Background(), req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.Equal(t, 1, len(events.Events), "Expected only one event")
	assert.Equal(t, int64(1), events.Events[0].Id)
	assert.Equal(t, "Mock Event 1", events.Events[0].Name)
	assert.Equal(t, "Brisbane", events.Events[0].Location)
	assert.Equal(t, "Bear Fighting", events.Events[0].Sport)
	assert.Equal(t, "Open", events.Events[0].Status)
}

func TestGetEvent(t *testing.T) {
	repo := &MockEventsRepo{}
	s := service.NewSportsService(repo)

	req := &sports.GetEventRequest{Id: 1}

	event, err := s.GetEvent(context.Background(), req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.NotNil(t, event.Event, "Expected event not to be nil")
	assert.Equal(t, int64(1), event.Event.Id)
	assert.Equal(t, "Mock Event 1", event.Event.Name)
	assert.Equal(t, "Brisbane", event.Event.Location)
	assert.Equal(t, "Bear Fighting", event.Event.Sport)
	assert.Equal(t, "Open", event.Event.Status)
}
