package trip

import (
	"context"
	trip2 "driver_service/internal/domain/models/trip"
	"driver_service/internal/domain/repository/trip"
	"github.com/google/uuid"
)

type Service struct {
	tripRepository trip.Repository
}

func NewService(
	tripRepo trip.Repository,
) *Service {
	return &Service{
		tripRepository: tripRepo,
		// userMetrics:    metrics.NewUserMetrics(),
	}
}

func (s *Service) GetTrips(ctx context.Context) ([]trip2.Trip, error) {
	return s.tripRepository.GetTrips(ctx)
}

func (s *Service) GetTrip(ctx context.Context, tripId uuid.UUID) (trip2.Trip, error) {
	return s.tripRepository.GetTrip(ctx, tripId)
}
