package service_test

import (
	"testing"

	"database/sql"
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
)

// Test list races with visiblity filter
func TestListRaces_VisibilityFilter(t *testing.T) {
	database, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer database.Close()

	repo := db.NewRacesRepo(database)

	// Setup memory data
	err = repo.Init()
	assert.NoError(t, err)

	reqFilter := &racing.ListRacesRequestFilter{
		Visibility: &racing.ListRacesRequestFilter_Visible{
			Visible: true,
		},
	}

	// Get races
	races, err := repo.List(reqFilter, "")

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result slice is not empty
	assert.NotEmpty(t, races)

	// Assert that all returned races are visible
	for _, race := range races {
		assert.True(t, race.Visible)
	}
}

// Test list races with no visiblity filter
func TestListRaces_No_VisibilityFilter(t *testing.T) {
	database, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer database.Close()

	repo := db.NewRacesRepo(database)

	// Setup memory data
	err = repo.Init()
	assert.NoError(t, err)

	reqFilter := &racing.ListRacesRequestFilter{
		Visibility: &racing.ListRacesRequestFilter_Visible{},
	}

	// Get races
	races, err := repo.List(reqFilter, "")

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result slice is not empty
	assert.NotEmpty(t, races)

	// Assert that it returns all races regardless of the visiblity
	for _, race := range races {
		assert.NotNil(t, race)
	}
}

// Test list races with ORDER BY advertised_start_time
func TestListRaces_OrderByAdvertisedStartTime(t *testing.T) {
	database, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer database.Close()

	repo := db.NewRacesRepo(database)

	// Setup memory data
	err = repo.Init()
	assert.NoError(t, err)

	reqFilter := &racing.ListRacesRequestFilter{
		Visibility: &racing.ListRacesRequestFilter_Visible{},
	}

	// Get races
	races, err := repo.List(reqFilter, "advertised_start_time")

	// Assert that there was no error
	assert.NoError(t, err)

	// Assert that the result slice is not empty
	assert.NotEmpty(t, races)

	// Database ORDER BY ASC by default
	// So assert that if races are ordered by advertised_start_time correctly
	// race[i] time should >= race[i-1] time
	for i := 1; i < len(races); i++ {
		assert.True(t, races[i].AdvertisedStartTime.Seconds >= races[i-1].AdvertisedStartTime.Seconds)
	}
}
