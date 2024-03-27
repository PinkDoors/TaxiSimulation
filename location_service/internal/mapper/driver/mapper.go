package driver

import (
	openapi "location_service/internal/httpadapter/generate"
	"location_service/internal/model"
)

func ToResponseDriver(driver model.Driver) (responseDriver openapi.Driver) {
	return openapi.Driver{
		Id:  &driver.DriverId,
		Lat: driver.Lat,
		Lng: driver.Lng,
	}
}
