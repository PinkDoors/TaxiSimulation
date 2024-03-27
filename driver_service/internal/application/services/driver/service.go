package driver

import (
	"context"
	"driver_service/internal/application/eventbus/producer/trip_outbound"
	"driver_service/internal/application/services/location"
	"driver_service/internal/application/services/trip"
	domainModels "driver_service/internal/domain/models/trip"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	tripService     *trip.Service
	locationService location.Service
	producer        trip_outbound.Producer
	logger          *zap.Logger
}

func NewService(
	tripService *trip.Service,
	locationService location.Service,
	producer trip_outbound.Producer,
	logger *zap.Logger,
) *Service {
	return &Service{
		tripService:     tripService,
		locationService: locationService,
		producer:        producer,
		logger:          logger,
	}
}

func (s *Service) GetTripsForDriver(ctx context.Context, driverId uuid.UUID) ([]domainModels.Trip, error) {
	var tripsForDriver []domainModels.Trip

	createdTrips, err := s.tripService.GetCreatedTrips(ctx)
	if err != nil {
		return nil, err
	}

	for _, createdTrip := range createdTrips {
		drivers, getDriversErr := s.locationService.GetDriversWithLocations(ctx, createdTrip.From)
		if getDriversErr != nil {
			continue
		}

		for _, driver := range drivers {
			println(driver.Id.String())
			if driver.Id == driverId {
				tripsForDriver = append(tripsForDriver, createdTrip)
				break
			}
		}
	}

	return tripsForDriver, nil
}

func (s *Service) AcceptTrip(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID) (tripFound bool, err error) {
	return s.produceChangeTripStatusCommand(ctx, tripId, driverId, domainModels.DriverFound)
}

func (s *Service) CancelTrip(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID, reason string) (tripFound bool, err error) {
	return s.produceChangeTripStatusCommand(ctx, tripId, driverId, domainModels.Canceled)
}

func (s *Service) StartTrip(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID) (tripFound bool, err error) {
	return s.produceChangeTripStatusCommand(ctx, tripId, driverId, domainModels.Started)
}

func (s *Service) EndTrip(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID) (tripFound bool, err error) {
	return s.produceChangeTripStatusCommand(ctx, tripId, driverId, domainModels.Ended)
}

func (s *Service) produceChangeTripStatusCommand(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID, newStatus domainModels.Status) (tripFound bool, err error) {
	foundTrip, err := s.tripService.GetTrip(ctx, tripId)
	if err != nil {
		return false, err
	}

	if foundTrip == nil {
		return false, nil
	}

	err = s.producer.Produce(ctx, tripId, driverId, time.Now(), newStatus)
	if err != nil {
		return true, err
	}

	return true, nil
}
