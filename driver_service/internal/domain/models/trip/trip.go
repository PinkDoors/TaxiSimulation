package trip

import (
	"driver_service/internal/domain/models"
	"github.com/google/uuid"
)

type Trip struct {
	Id         uuid.UUID
	DriverId   string
	From       models.LatLngLiteral
	To         models.LatLngLiteral
	Price      models.Money
	TripStatus TripStatus
}
