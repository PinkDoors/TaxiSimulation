package driver

import (
	"context"
	"driver_service/internal/application/models/trip_inbound"
	"driver_service/internal/application/trip"
	"driver_service/internal/infrastracture/eventbus/producer"
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	service  *trip.Service
	producer producer.Producer
	logger   *zap.Logger
}

func NewService(
	service *trip.Service,
	producer producer.Producer,
	logger *zap.Logger,
) *Service {
	return &Service{
		service:  service,
		producer: producer,
		logger:   logger,
	}
}

func (s *Service) AcceptTrip(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID) (tripFound bool, err error) {
	getTrip, err := s.service.GetTrip(ctx, tripId)
	if err != nil {
		return false, err
	}

	if getTrip == nil {
		return false, nil
	}

	command := trip_inbound.Command{
		ID:              uuid.New(),
		Source:          "/driver",
		Type:            "trip.command.accept",
		DataContentType: "application/json",
		Time:            time.Now(),
		Data: trip_inbound.Data{
			TripID:   tripId,
			DriverID: driverId,
		},
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		s.logger.Error("Error when decoding command", zap.Error(err))
	}

	err = s.producer.Produce(ctx, jsonData)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (s *Service) CancelTrip(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID) (tripFound bool, err error) {
	getTrip, err := s.service.GetTrip(ctx, tripId)
	if err != nil {
		return false, err
	}

	if getTrip == nil {
		return false, nil
	}

	command := trip_inbound.Command{
		ID:              uuid.New(),
		Source:          "/driver",
		Type:            "trip.command.cancel",
		DataContentType: "application/json",
		Time:            time.Now(),
		Data: trip_inbound.Data{
			TripID:   tripId,
			DriverID: driverId,
		},
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		s.logger.Error("Error when decoding command", zap.Error(err))
	}

	err = s.producer.Produce(ctx, jsonData)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (s *Service) StartTrip(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID) (tripFound bool, err error) {
	getTrip, err := s.service.GetTrip(ctx, tripId)
	if err != nil {
		return false, err
	}

	if getTrip == nil {
		return false, nil
	}

	command := trip_inbound.Command{
		ID:              uuid.New(),
		Source:          "/driver",
		Type:            "trip.command.start",
		DataContentType: "application/json",
		Time:            time.Now(),
		Data: trip_inbound.Data{
			TripID:   tripId,
			DriverID: driverId,
		},
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		s.logger.Error("Error when decoding command", zap.Error(err))
	}

	err = s.producer.Produce(ctx, jsonData)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (s *Service) EndTrip(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID) (tripFound bool, err error) {
	getTrip, err := s.service.GetTrip(ctx, tripId)
	if err != nil {
		return false, err
	}

	if getTrip == nil {
		return false, nil
	}

	command := trip_inbound.Command{
		ID:              uuid.New(),
		Source:          "/driver",
		Type:            "trip.command.end",
		DataContentType: "application/json",
		Time:            time.Now(),
		Data: trip_inbound.Data{
			TripID:   tripId,
			DriverID: driverId,
		},
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		s.logger.Error("Error when decoding command", zap.Error(err))
	}

	err = s.producer.Produce(ctx, jsonData)
	if err != nil {
		return true, err
	}

	return true, nil
}
