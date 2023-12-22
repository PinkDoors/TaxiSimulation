package httpadapter

import (
	"context"
	"fmt"
	generated "location_service/internal/httpadapter/generate"
	mapper "location_service/internal/mapper/driver"
	"location_service/internal/service"
)

type LocationServer struct {
	locationService service.LocationService
}

func New(
	locationService service.LocationService,
) *LocationServer {
	return &LocationServer{
		locationService: locationService,
	}
}

func (l *LocationServer) GetDrivers(
	ctx context.Context,
	request generated.GetDriversRequestObject,
) (generated.GetDriversResponseObject, error) {
	driverList, err := l.locationService.GetDrivers(
		ctx,
		request.Params.Radius,
		request.Params.Lat,
		request.Params.Lng,
	)
	if err != nil {
		return generated.GetDrivers500Response{}, err
	}
	if len(driverList) == 0 {
		return generated.GetDrivers404Response{}, nil
	}

	drivers := make([]generated.Driver, 0)
	for _, driver := range driverList {
		driverResult := mapper.ToResponseDriver(*driver)
		drivers = append(drivers, driverResult)
	}

	return generated.GetDrivers200JSONResponse(drivers), nil
}

func (l *LocationServer) UpdateDriverLocation(
	ctx context.Context,
	request generated.UpdateDriverLocationRequestObject,
) (generated.UpdateDriverLocationResponseObject, error) {
	err := l.locationService.UpdateDriver(
		ctx,
		fmt.Sprintf("%s", request.DriverId),
		request.Body.Lat,
		request.Body.Lng,
	)
	if err != nil {
		return generated.UpdateDriverLocation500Response{}, err
	}

	return generated.UpdateDriverLocation200Response{}, nil
}