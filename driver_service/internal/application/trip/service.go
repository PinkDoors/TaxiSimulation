package trip

import (
	"context"
	trip2 "driver_service/internal/domain/models/trip"
	"driver_service/internal/domain/repository/trip"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	tripRepository trip.Repository
	logger         *zap.Logger
}

func NewService(
	tripRepo trip.Repository,
	logger *zap.Logger,
) *Service {
	return &Service{
		tripRepository: tripRepo,
		logger:         logger,
		// userMetrics:    metrics.NewUserMetrics(),
	}
}

func (s *Service) GetTrips(ctx context.Context) ([]trip2.Trip, error) {
	return s.tripRepository.GetTrips(ctx)
}

func (s *Service) GetTrip(ctx context.Context, tripId uuid.UUID) (*trip2.Trip, error) {
	return s.tripRepository.GetTrip(ctx, tripId)
}

func (s *Service) CreateTrip(ctx context.Context, tripId uuid.UUID) error {
	return s.tripRepository.CreateTrip(ctx, tripId)
}

func (s *Service) AcceptTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error) {
	return s.tripRepository.AcceptTrip(ctx, tripId)
}

func (s *Service) CancelTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error) {
	return s.tripRepository.CancelTrip(ctx, tripId)
}

func (s *Service) StartTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error) {
	return s.tripRepository.StartTrip(ctx, tripId)
}
