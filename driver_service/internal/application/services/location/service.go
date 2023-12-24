package location

import (
	"context"
	"driver_service/internal/domain/models"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
}

type LocationService interface {
	GetDriversWithLocations(ctx context.Context, latLngLiteral models.LatLngLiteral) ([]models.Driver, error)
}
