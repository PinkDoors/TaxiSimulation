package httpadapter

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	generated "location_service/internal/httpadapter/generate"
	mapper "location_service/internal/mapper/driver"
	"location_service/internal/service"
)

type LocationServer struct {
	locationService     service.Service
	logger              *zap.Logger
	getDriversCounter   prometheus.Counter
	updateDriverCounter prometheus.Counter
}

func New(
	locationService service.Service,
	logger *zap.Logger,
	getDriversCounter prometheus.Counter,
	updateDriverCounter prometheus.Counter,
) *LocationServer {
	return &LocationServer{
		locationService:     locationService,
		logger:              logger,
		getDriversCounter:   getDriversCounter,
		updateDriverCounter: updateDriverCounter,
	}
}

func (l *LocationServer) GetDrivers(
	ctx context.Context,
	request generated.GetDriversRequestObject,
) (generated.GetDriversResponseObject, error) {
	l.getDriversCounter.Inc()
	l.logger.Info("URI /drivers was called")
	driverList, err := l.locationService.GetDrivers(
		ctx,
		request.Params.Radius,
		request.Params.Lat,
		request.Params.Lng,
	)
	if err != nil {
		l.logger.Error(
			"Internal Server Error",
			zap.String("cause", err.Error()),
		)
		return generated.GetDrivers500Response{}, err
	}
	if len(driverList) == 0 {
		l.logger.Error("Not found Error")
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
	l.updateDriverCounter.Inc()
	l.logger.Info(
		"URI /drivers/{driver_id}/location was called",
		zap.String("DriverId", request.DriverId.String()),
	)
	err := l.locationService.UpdateDriver(
		ctx,
		fmt.Sprintf("%s", request.DriverId),
		request.Body.Lat,
		request.Body.Lng,
	)
	if err != nil {
		l.logger.Error(
			"Internal Server Error",
			zap.String("cause", err.Error()),
		)
		return generated.UpdateDriverLocation500Response{}, err
	}

	return generated.UpdateDriverLocation200Response{}, nil
}
