package service

import (
	"context"
	"go.uber.org/zap"
	"location_service/internal/model"
	"location_service/internal/repo"
)

type LocationService struct {
	repo   repo.Driver
	logger *zap.Logger
}

func New(
	repo repo.Driver,
	logger *zap.Logger,
) *LocationService {
	return &LocationService{
		repo:   repo,
		logger: logger,
	}
}

func (l *LocationService) GetDrivers(ctx context.Context, radius, lat, lng float32) ([]*model.Driver, error) {
	return l.repo.ListDriver(ctx, radius, lat, lng)
}

func (l *LocationService) UpdateDriver(ctx context.Context, driverId string, lat, lng float32) error {
	return l.repo.UpdateDriverLocation(ctx, driverId, lat, lng)
}
