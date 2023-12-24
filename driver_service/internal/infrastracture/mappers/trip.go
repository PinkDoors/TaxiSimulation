package mappers

import (
	"driver_service/internal/domain/models"
	domainModels "driver_service/internal/domain/models/trip"
	"driver_service/internal/infrastracture/repository/dto/trip"
	"github.com/google/uuid"
)

func TripDtosToModels(dtos []trip.Trip) ([]domainModels.Trip, error) {
	var tripModels []domainModels.Trip

	for _, tripDTO := range dtos {
		tripModel, err := TripDtoToModel(tripDTO)
		if err != nil {
			return nil, err
		}
		tripModels = append(tripModels, tripModel)
	}

	return tripModels, nil
}

func TripDtoToModel(dto trip.Trip) (domainModels.Trip, error) {
	tripId, err := uuid.Parse(dto.Id)
	if err != nil {
		return domainModels.Trip{}, err
	}

	return domainModels.Trip{
		Id:       tripId,
		DriverId: dto.DriverId,
		From: models.LatLngLiteral{
			Lat: float32(dto.From.Lat),
			Lng: float32(dto.From.Lng),
		},
		To: models.LatLngLiteral{
			Lat: float32(dto.To.Lat),
			Lng: float32(dto.To.Lng),
		},
		Price: models.Money{
			Amount:   dto.Price.Amount,
			Currency: dto.Price.Currency,
		},
		TripStatus: dto.TripStatus,
	}, nil
}
