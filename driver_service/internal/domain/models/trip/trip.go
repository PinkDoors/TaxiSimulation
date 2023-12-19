package trip

import (
	"driver_service/internal/domain/models"
	"github.com/google/uuid"
)

type Trip struct {
	Id         uuid.UUID            `bson:"id"`
	DriverId   string               `bson:"driverId"`
	From       models.LatLngLiteral `bson:"from"`
	To         models.LatLngLiteral `bson:"to"`
	Price      models.Money         `bson:"price"`
	TripStatus TripStatus           `bson:"tripStatus"`
}
