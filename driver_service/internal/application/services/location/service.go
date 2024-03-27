package location

import (
	"context"
	"driver_service/internal/domain/models"
)

type Service interface {
	GetDriversWithLocations(ctx context.Context, latLngLiteral models.LatLngLiteral) ([]models.Driver, error)
}
