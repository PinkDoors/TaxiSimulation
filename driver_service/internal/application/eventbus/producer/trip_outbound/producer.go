package trip_outbound

import (
	"context"
	"driver_service/internal/domain/models/trip"
	"github.com/google/uuid"
	"time"
)

type Producer interface {
	Produce(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID, time time.Time, newTripStatus trip.Status) error
}
