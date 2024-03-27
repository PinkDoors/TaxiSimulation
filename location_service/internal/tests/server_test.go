package tests

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	server "location_service/internal/httpadapter"
	generated "location_service/internal/httpadapter/generate"
	"location_service/internal/model"
	"location_service/internal/service"
	"testing"
)

func TestGetDrivers(t *testing.T) {
	// Создаем mock объекты
	mockLocationService := new(service.MockLocationService)
	getDriversCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "get_drivers_request", Help: "Increment for /drivers endpoint",
		},
	)
	updateDriverCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "update_driver_request", Help: "Increment for /drivers/{driver_id}/location endpoint",
		},
	)
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any

	var expectedDrivers []*model.Driver
	mockLocationService.On("GetDrivers", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedDrivers, nil)

	locationServer := server.New(mockLocationService, logger, getDriversCounter, updateDriverCounter)

	req := generated.GetDriversRequestObject{
		Params: generated.GetDriversParams{
			Lat:    50.0,
			Lng:    50.0,
			Radius: 10,
		},
	}

	drivers, err := locationServer.GetDrivers(context.Background(), req)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if drivers == nil {
		t.Error("drivers is nil or empty")
	}

	mockLocationService.AssertExpectations(t)
}
