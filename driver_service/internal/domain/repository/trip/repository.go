package trip

import (
	"context"
	"driver_service/internal/domain/models/trip"
	"github.com/google/uuid"
)

type Repository interface {
	GetTrips(ctx context.Context) ([]trip.Trip, error)
	GetTrip(ctx context.Context, tripId uuid.UUID) (*trip.Trip, error)

	CreateTrip(ctx context.Context, trip trip.Trip) error
	AcceptTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error)
	StartTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error)
	EndTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error)
	CancelTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error)
}
