package service

import (
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)

	// GetRace will return a single race.
	GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.GetRaceResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter, in.OrderBy)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}

// Implement GetRace to satisfy the Racing interface.
func (s *racingService) GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.GetRaceResponse, error) {
	race, err := s.racesRepo.GetByID(in.Id)
	if err != nil {
		return nil, err
	}

	// Returns 404 if race not found.
	if race == nil {
		return nil, status.Error(codes.NotFound, "Race not found")
	}

	return &racing.GetRaceResponse{Race: race}, nil
}
