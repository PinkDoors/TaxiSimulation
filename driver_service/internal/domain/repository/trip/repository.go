package trip

import (
	"context"
	"driver_service/internal/domain/models/trip"
	"github.com/google/uuid"
)

type Repository interface {
	GetTrips(ctx context.Context) ([]trip.Trip, error)
	GetTrip(ctx context.Context, tripId uuid.UUID) (trip.Trip, error)
}
