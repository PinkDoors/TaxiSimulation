package trip_outbound

import (
	"context"
	"driver_service/configs/kafka/producer"
	"driver_service/internal/domain/models/trip"
	eventbus "driver_service/internal/infrastracture/eventbus/producer"
	"driver_service/internal/infrastracture/eventbus/producer/kafka"
	"driver_service/pkg/kafka/trip_inbound"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type Producer struct {
	kafkaProducer eventbus.Producer
	config        producer.Config
	logger        *zap.Logger
}

func NewProducer(
	async *bool,
	config producer.Config,
	logger *zap.Logger,
) *Producer {
	kafkaProducer := kafka.NewProducer(async, config, logger)

	return &Producer{
		kafkaProducer: kafkaProducer,
		config:        config,
		logger:        logger,
	}
}

func (p *Producer) Produce(ctx context.Context, tripId uuid.UUID, driverId uuid.UUID, time time.Time, newTripStatus trip.Status) error {
	commandType, err := getCommandType(newTripStatus)
	if err != nil {
		return err
	}

	command := trip_inbound.Command{
		ID:              uuid.New(),
		Source:          "/driver",
		Type:            commandType,
		DataContentType: "application/json",
		Time:            time,
		Data: trip_inbound.Data{
			TripID:   tripId,
			DriverID: driverId,
		},
	}

	jsonData, err := json.Marshal(command)
	if err != nil {
		p.logger.Error("Error when decoding command", zap.Error(err))
	}

	err = p.kafkaProducer.Produce(ctx, jsonData)
	if err != nil {
		return err
	}

	return nil
}

func getCommandType(tripStatus trip.Status) (string, error) {
	switch tripStatus {
	case trip.DriverSearch:
		return "trip.command.create", nil
	case trip.DriverFound:
		return "trip.command.accept", nil
	case trip.Started:
		return "trip.command.started", nil
	case trip.Ended:
		return "trip.command.ended", nil
	case trip.Canceled:
		return "trip.command.canceled", nil
	default:
		logString := fmt.Sprintf("unsupported trip status - %s", tripStatus)
		return "nil", errors.New(logString)
	}
}
