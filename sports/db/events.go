package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/LightningDev/entain-master/sports/proto/sports"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
)

// EventsRepo provides repository access to sports events.
type EventsRepo interface {
	// Init will initialise our events repository.
	Init() error

	// List will return a list of sports events.
	List(filter *sports.ListEventsRequestFilter, orderby string) ([]*sports.Event, error)

	// GetByID will return a single sports event.
	GetByID(id int64) (*sports.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the events repository dummy data.
func (r *eventsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

// List all events
func (r *eventsRepo) List(filter *sports.ListEventsRequestFilter, orderBy string) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[eventsList]

	query, args = r.applyFilter(query, filter)

	// Apply ORDER BY as a last step in SQL query.
	query = r.applySort(query, orderBy)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

// Implement GetByID.
func (r *eventsRepo) GetByID(id int64) (*sports.Event, error) {
	var (
		err   error
		query string
	)

	query = getEventQueries()[eventByID]

	row := r.db.QueryRow(query, id)
	if err != nil {
		return nil, err
	}

	return r.scanEvent(row)
}

func (r *eventsRepo) applySort(query string, orderby string) string {
	if orderby == "" {
		return query
	}

	return query + " ORDER BY " + orderby
}

func (r *eventsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if filter.Name != "" {
		clauses = append(clauses, "name = ?")
		args = append(args, filter.Name)
	}

	if filter.Location != "" {
		clauses = append(clauses, "location = ?")
		args = append(args, filter.Location)
	}

	if filter.Sport != "" {
		clauses = append(clauses, "sport = ?")
		args = append(args, filter.Sport)
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

// Scan list of sports events.
func (m *eventsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.Name, &event.Location, &event.Sport, &advertisedStart, &event.Status); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts

		events = append(events, &event)
	}

	return events, nil
}

// Scan a single sports event row.
func (m *eventsRepo) scanEvent(
	row *sql.Row,
) (*sports.Event, error) {
	var event sports.Event
	var advertisedStart time.Time

	if err := row.Scan(&event.Id, &event.Name, &event.Location, &event.Sport, &advertisedStart, &event.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	ts, err := ptypes.TimestampProto(advertisedStart)
	if err != nil {
		return nil, err
	}

	event.AdvertisedStartTime = ts

	return &event, nil
}
