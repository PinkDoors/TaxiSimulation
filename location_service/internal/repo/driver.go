package repo

import (
	"context"
	"location_service/internal/model"
)

type Driver interface {
	UpdateDriverLocation(ctx context.Context, driverId string, lat, lng float32) error
	ListDriver(ctx context.Context, radius, lat, lng float32) ([]*model.Driver, error)
}
