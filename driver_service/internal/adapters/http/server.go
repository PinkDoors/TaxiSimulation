package http

import (
	"context"
	generated "driver_service/internal/adapters/http/generate"
	"driver_service/internal/application/services/driver"
	"driver_service/internal/application/services/trip"
	trip2 "driver_service/internal/mappers/trip"
	"fmt"
	"time"
)

type DriverServer struct {
	driverService *driver.Service
	tripService   *trip.Service
}

func NewDriverServer(
	driverService *driver.Service,
	tripService *trip.Service,
) *DriverServer {
	return &DriverServer{
		driverService: driverService,
		tripService:   tripService,
	}
}

func (s *DriverServer) GetTrips(ctx context.Context, request generated.GetTripsRequestObject) (generated.GetTripsResponseObject, error) {
	const maxAttempts = 10
	attempts := 0

	for attempts <= maxAttempts {
		foundTrips, err := s.driverService.GetTripsForDriver(ctx, request.Params.UserId)
		if err != nil {
			return generated.GetTrips500Response{}, err
		}

		if len(foundTrips) > 0 {
			var targetTrips []generated.Trip
			for _, sourceTrip := range foundTrips {
				targetTrip := trip2.ToResponseTrip(sourceTrip)
				targetTrips = append(targetTrips, targetTrip)
			}
			return generated.GetTrips200JSONResponse(targetTrips), nil
		}

		attempts++
		time.Sleep(time.Second)
	}

	return generated.GetTrips504Response{}, nil
}

func (s *DriverServer) GetTripByID(ctx context.Context, request generated.GetTripByIDRequestObject) (generated.GetTripByIDResponseObject, error) {
	fmt.Println("Handling GetTripByID request")

	foundTrip, err := s.tripService.GetTrip(ctx, request.TripId)
	if err != nil {
		return generated.GetTripByID500Response{}, err
	}

	if foundTrip == nil {
		return generated.GetTripByID404Response{}, err
	}

	targetTrip := trip2.ToResponseTrip(*foundTrip)

	return generated.GetTripByID200JSONResponse(targetTrip), nil
}

func (s *DriverServer) AcceptTrip(ctx context.Context, request generated.AcceptTripRequestObject) (generated.AcceptTripResponseObject, error) {
	userId := request.Params.UserId

	tripFound, err := s.driverService.AcceptTrip(ctx, request.TripId, userId)
	if err != nil {
		return generated.AcceptTrip500Response{}, err
	}

	if tripFound == false {
		return generated.AcceptTrip404Response{}, nil
	}

	return generated.AcceptTrip200Response{}, nil
}

func (s *DriverServer) CancelTrip(ctx context.Context, request generated.CancelTripRequestObject) (generated.CancelTripResponseObject, error) {
	fmt.Println("Handling CancelTrip request")

	userId := request.Params.UserId
	reason := request.Params.Reason

	tripFound, err := s.driverService.CancelTrip(ctx, request.TripId, userId, *reason)
	if err != nil {
		return generated.CancelTrip500Response{}, err
	}

	if tripFound == false {
		return generated.CancelTrip404Response{}, nil
	}

	return generated.CancelTrip200Response{}, nil
}

func (s *DriverServer) StartTrip(ctx context.Context, request generated.StartTripRequestObject) (generated.StartTripResponseObject, error) {
	userId := request.Params.UserId

	tripFound, err := s.driverService.StartTrip(ctx, request.TripId, userId)
	if err != nil {
		return generated.StartTrip500Response{}, err
	}

	if tripFound == false {
		return generated.StartTrip404Response{}, nil
	}

	return generated.StartTrip200Response{}, nil
}

func (s *DriverServer) EndTrip(ctx context.Context, request generated.EndTripRequestObject) (generated.EndTripResponseObject, error) {
	userId := request.Params.UserId

	tripFound, err := s.driverService.EndTrip(ctx, request.TripId, userId)
	if err != nil {
		return generated.EndTrip500Response{}, err
	}

	if tripFound == false {
		return generated.EndTrip404Response{}, nil
	}

	return generated.EndTrip200Response{}, nil
}
