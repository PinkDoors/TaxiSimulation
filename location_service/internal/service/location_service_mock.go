package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"location_service/internal/model"
)

type MockLocationService struct {
	mock.Mock
}

func (m *MockLocationService) GetDrivers(ctx context.Context, radius, lat, lng float32) ([]*model.Driver, error) {
	args := m.Called(ctx, radius, lat, lng)
	return args.Get(0).([]*model.Driver), args.Error(1)
}

func (m *MockLocationService) UpdateDriver(ctx context.Context, driverId string, lat, lng float32) error {
	args := m.Called(ctx, driverId, lat, lng)
	return args.Error(0)
}
