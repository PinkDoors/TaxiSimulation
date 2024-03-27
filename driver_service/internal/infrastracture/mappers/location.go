package mappers

import (
	"driver_service/internal/domain/models"
	"driver_service/internal/infrastracture/services/location/dto"
	"github.com/google/uuid"
)

func DriverWithLocationDtosToModels(dtos []dto.Driver) ([]models.Driver, error) {
	var drivers []models.Driver

	for _, driverDTO := range dtos {
		driver, err := driverWithLocationDtoToModel(driverDTO)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

func driverWithLocationDtoToModel(dto dto.Driver) (models.Driver, error) {
	id, err := uuid.Parse(dto.ID)
	if err != nil {
		return models.Driver{}, err
	}

	return models.Driver{
		Id: id,
		Location: models.LatLngLiteral{
			Lat: dto.Lat,
			Lng: dto.Lng,
		},
	}, nil
}
