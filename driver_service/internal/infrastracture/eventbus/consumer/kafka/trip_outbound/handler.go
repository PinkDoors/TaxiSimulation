package trip_outbound

import (
	"context"
	"driver_service/internal/application/services/trip"
	"driver_service/internal/domain/models"
	domainModels "driver_service/internal/domain/models/trip"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TripInboundMessageHandler struct {
	tripService *trip.Service
	logger      *zap.Logger
}

func NewTripInboundMessageHandler(
	tripService *trip.Service,
	logger *zap.Logger,
) *TripInboundMessageHandler {
	return &TripInboundMessageHandler{
		tripService: tripService,
		logger:      logger,
	}
}

func (timh *TripInboundMessageHandler) Handle(ctx context.Context, message kafka.Message) error {
	var event Event
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		timh.logger.Warn("Error decoding message", zap.Error(err))
		// skipping this message
		return nil
	}

	parsedUUID, err := uuid.Parse(event.Data.TripID)
	if err != nil {
		timh.logger.Error("Error decoding uuid.", zap.Error(err))
		return nil
	}

	switch event.Type {
	case "trip.event.accepted":
		_, err := timh.tripService.AcceptTrip(ctx, parsedUUID)
		if err != nil {
			return err
		}
	case "trip.event.canceled":
		_, err := timh.tripService.CancelTrip(ctx, parsedUUID)
		if err != nil {
			return err
		}
	case "trip.event.created":
		tripModel := domainModels.Trip{
			Id:       parsedUUID,
			DriverId: "",
			From: models.LatLngLiteral{
				Lat: event.Data.From.Lat,
				Lng: event.Data.From.Lng,
			},
			To: models.LatLngLiteral{
				Lat: event.Data.To.Lat,
				Lng: event.Data.To.Lng,
			},
			Price: models.Money{
				Amount:   event.Data.Price.Amount,
				Currency: event.Data.Price.Currency,
			},
			TripStatus: domainModels.DriverSearch,
		}
		err := timh.tripService.CreateTrip(ctx, tripModel)
		if err != nil {
			return err
		}
	case "trip.event.ended":
		_, err := timh.tripService.EndTrip(ctx, parsedUUID)
		if err != nil {
			return err
		}
	case "trip.event.started":
		_, err := timh.tripService.StartTrip(ctx, parsedUUID)
		if err != nil {
			return err
		}
	default:
		timh.logger.Error("Unknown message type - "+event.Type, zap.Error(err))
		// skipping this message
		return nil
	}

	return nil
}
