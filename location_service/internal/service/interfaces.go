package service

import (
	"context"
	"location_service/internal/model"
)

type Service interface {
	GetDrivers(ctx context.Context, radius, lat, lng float32) ([]*model.Driver, error)
	UpdateDriver(ctx context.Context, driverId string, lat, lng float32) error
}
