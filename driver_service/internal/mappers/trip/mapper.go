package trip

import (
	"driver_service/internal/adapters/http/generate"
	"driver_service/internal/domain/models/trip"
	trip2 "driver_service/internal/domain/models/trip"
)

func ToResponseTrip(model trip2.Trip) (responseTrip openapi.Trip) {
	var tripStatus openapi.TripStatus = ToResponseStatus(model.TripStatus)

	return openapi.Trip{
		Id:       &model.Id,
		DriverId: &model.DriverId,
		From: &openapi.LatLngLiteral{
			Lat: model.From.Lat,
			Lng: model.From.Lng,
		},
		To: &openapi.LatLngLiteral{
			Lat: model.To.Lat,
			Lng: model.To.Lng,
		},
		Price: &openapi.Money{
			Amount:   model.Price.Amount,
			Currency: model.Price.Currency,
		},
		Status: &tripStatus,
	}
}

func ToResponseStatus(modelStatus trip2.TripStatus) (responseStatus openapi.TripStatus) {
	switch modelStatus {
	case trip.DriverSearch:
		return openapi.DRIVERSEARCH
	case trip.DriverFound:
		return openapi.DRIVERFOUND
	case trip.OnPosition:
		return openapi.ONPOSITION
	case trip.Started:
		return openapi.STARTED
	case trip.Ended:
		return openapi.ENDED
	case trip.Canceled:
		return openapi.CANCELED
	default:
		return ""
	}
}
