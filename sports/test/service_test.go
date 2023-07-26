package service_test

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/LightningDev/entain-master/sports/db"
	"github.com/LightningDev/entain-master/sports/proto/sports"
	"github.com/stretchr/testify/assert"
)

// Test list all events should returns all available events
func TestListSports_ListEvents(t *testing.T) {
	database, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer database.Close()

	repo := db.NewEventsRepo(database)

	// Setup memory data
	err = repo.Init()
	assert.NoError(t, err)

	reqFilter := &sports.ListEventsRequestFilter{}

	// Get events
	events, err := repo.List(reqFilter, "")

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result slice is not empty
	assert.NotEmpty(t, events)

	// Assert that all returned sports are visible
	for _, event := range events {
		assert.NotNil(t, event)
	}
}

// Test list events with ORDER BY advertised_start_time
func TestListSports_OrderByAdvertisedStartTime(t *testing.T) {
	database, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer database.Close()

	repo := db.NewEventsRepo(database)

	// Setup memory data
	err = repo.Init()
	assert.NoError(t, err)

	reqFilter := &sports.ListEventsRequestFilter{}

	// Get events
	events, err := repo.List(reqFilter, "advertised_start_time")

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result slice is not empty
	assert.NotEmpty(t, events)

	// Database ORDER BY ASC by default
	// So assert that if events are ordered by advertised_start_time correctly
	// events[i] time should >= events[i-1] time
	for i := 1; i < len(events); i++ {
		assert.True(t, events[i].AdvertisedStartTime.Seconds >= events[i-1].AdvertisedStartTime.Seconds)
	}
}

// Test GetSport to return a single event
func TestGetSport_GetByID(t *testing.T) {
	database, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer database.Close()

	repo := db.NewEventsRepo(database)

	// Setup memory data
	err = repo.Init()
	assert.NoError(t, err)

	// Random id from 1-100 in seed data
	rand.Seed(time.Now().UnixNano())
	randomID := int64(rand.Intn(100) + 1)

	// Get sport
	sport, err := repo.GetByID(randomID)

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result slice is not empty
	assert.NotNil(t, sport)

	// Assert that it returns a sport with id equal to randomID
	assert.Equal(t, randomID, sport.Id)
}
