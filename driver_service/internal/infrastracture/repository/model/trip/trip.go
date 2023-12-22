package trip

import (
	"driver_service/internal/domain/models"
	"driver_service/internal/domain/models/trip"
)

type Trip struct {
	Id         string               `bson:"id"`
	DriverId   string               `bson:"driverId"`
	From       models.LatLngLiteral `bson:"from"`
	To         models.LatLngLiteral `bson:"to"`
	Price      models.Money         `bson:"price"`
	TripStatus trip.TripStatus      `bson:"tripStatus"`
}
