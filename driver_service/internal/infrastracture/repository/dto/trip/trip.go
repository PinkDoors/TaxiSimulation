package trip

import (
	"driver_service/internal/domain/models/trip"
	"driver_service/internal/infrastracture/repository/dto"
)

type Trip struct {
	Id         string            `bson:"id"`
	DriverId   string            `bson:"driverId"`
	From       dto.LatLngLiteral `bson:"from"`
	To         dto.LatLngLiteral `bson:"to"`
	Price      dto.Money         `bson:"price"`
	TripStatus trip.Status       `bson:"tripStatus"`
}
