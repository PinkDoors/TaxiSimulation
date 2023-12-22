package trip

import (
	"context"
	config "driver_service/configs"
	"driver_service/internal/domain/models"
	"driver_service/internal/domain/models/trip"
	domainRepo "driver_service/internal/domain/repository/trip"
	trip2 "driver_service/internal/infrastracture/repository/model/trip"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var _ domainRepo.Repository = &Repository{}

type Repository struct {
	config *config.Config
	logger *zap.Logger
}

func NewRepository(
	config *config.Config,
	logger *zap.Logger,
) *Repository {
	return &Repository{
		config: config,
		logger: logger,
	}
}

func (r *Repository) GetTrips(ctx context.Context) ([]trip.Trip, error) {
	client, err := r.getMongoClient(ctx)
	if err != nil {
		r.logger.Error("Failed to create Mongo client.", zap.Error(err))
		return nil, err
	}

	defer func() {
		if mongoDisconnectErr := client.Disconnect(ctx); mongoDisconnectErr != nil {
			r.logger.Error("Failed to disconnect from Mongo client.", zap.Error(err))
		}
	}()

	col := client.Database("driver-service").Collection("trips")

	strUUID := "550e8400-e29b-41d4-a716-446612340000"
	parsedUUID, err := uuid.Parse(strUUID)
	fmt.Println("Inserting 1 documents...")
	_, err2 := col.InsertOne(ctx, trip.Trip{
		Id:         parsedUUID,
		DriverId:   "test",
		From:       models.LatLngLiteral{},
		To:         models.LatLngLiteral{},
		Price:      models.Money{},
		TripStatus: "DRIVER_SEARCH",
	})
	if err2 != nil {
		// log.Fatal(err)
	}

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) == false {
			r.logger.Error("Failed to read trips.", zap.Error(err))
		}
	}

	var foundTrips []trip.Trip
	err = cursor.All(ctx, &foundTrips)
	if err != nil {
		r.logger.Error("Failed to decode trips.", zap.Error(err))
		return nil, err
	}

	return foundTrips, nil
}

func (r *Repository) GetTrip(ctx context.Context, tripId uuid.UUID) (*trip.Trip, error) {
	client, err := r.getMongoClient(ctx)
	if err != nil {
		r.logger.Error("Failed to create Mongo client.", zap.Error(err))
		return nil, err
	}

	defer func() {
		if mongoDisconnectErr := client.Disconnect(ctx); mongoDisconnectErr != nil {
			r.logger.Error("Failed to disconnect from Mongo client.", zap.Error(err))
		}
	}()

	col := client.Database("driver-service").Collection("trips")
	result := col.FindOne(ctx, bson.M{"id": tripId.String()})

	findErr := result.Err()
	if findErr != nil {
		if errors.Is(findErr, mongo.ErrNoDocuments) {
			return nil, nil
		}
		r.logger.Error("Failed to read trip.", zap.Error(findErr))
		return nil, findErr
	}

	var foundTrip trip.Trip
	if err = result.Decode(&foundTrip); err != nil {
		r.logger.Error("Failed to decode trip", zap.Error(err))
		return nil, err
	}

	return &foundTrip, nil
}

func (r *Repository) CreateTrip(ctx context.Context, trip trip.Trip) error {
	client, err := r.getMongoClient(ctx)
	if err != nil {
		r.logger.Error("Failed to create Mongo client.", zap.Error(err))
		return err
	}

	defer func() {
		if mongoDisconnectErr := client.Disconnect(ctx); mongoDisconnectErr != nil {
			r.logger.Error("Failed to disconnect from Mongo client.", zap.Error(err))
		}
	}()

	col := client.Database("driver-service").Collection("trips")

	mongoEntity := trip2.Trip{
		Id:       trip.Id.String(),
		DriverId: "",
		From: models.LatLngLiteral{
			Lat: trip.From.Lat,
			Lng: trip.From.Lng,
		},
		To: models.LatLngLiteral{
			Lat: trip.To.Lat,
			Lng: trip.To.Lng,
		},
		Price: models.Money{
			Amount:   trip.Price.Amount,
			Currency: trip.Price.Currency,
		},
		TripStatus: trip.TripStatus,
	}

	_, err = col.InsertOne(ctx, mongoEntity)
	if err != nil {
		r.logger.Error("Failed to create new trip", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) AcceptTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error) {
	return r.updateTripStatus(ctx, tripId, trip.DriverFound)
}

func (r *Repository) StartTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error) {
	return r.updateTripStatus(ctx, tripId, trip.Started)
}

func (r *Repository) EndTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error) {
	return r.updateTripStatus(ctx, tripId, trip.Ended)
}

func (r *Repository) CancelTrip(ctx context.Context, tripId uuid.UUID) (tripFound bool, err error) {
	return r.updateTripStatus(ctx, tripId, trip.Canceled)
}

func (r *Repository) updateTripStatus(ctx context.Context, tripId uuid.UUID, tripStatus trip.TripStatus) (tripFound bool, err error) {
	client, err := r.getMongoClient(ctx)
	if err != nil {
		r.logger.Error("Failed to create Mongo client.", zap.Error(err))
		return false, err
	}

	defer func() {
		if mongoDisconnectErr := client.Disconnect(ctx); mongoDisconnectErr != nil {
			r.logger.Error("Failed to disconnect from Mongo client.", zap.Error(err))
		}
	}()

	col := client.Database("driver-service").Collection("trips")

	filter := bson.D{{"id", tripId.String()}}
	update := bson.D{{"$set", bson.D{{"tripStatus", tripStatus}}}}

	result, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Error("Failed to update trip.", zap.Error(err))
		return false, err
	}

	return result.MatchedCount > 0, nil
}

func (r *Repository) getMongoClient(ctx context.Context) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(r.config.Db.Uri)
	opts.SetTimeout(time.Duration(r.config.Db.Timeout) * time.Second)
	optsAuth := options.Credential{
		Username:   r.config.Db.Username,
		Password:   r.config.Db.Password,
		AuthSource: r.config.Db.AuthSource,
	}

	opts.SetAuth(optsAuth)

	opts.Monitor = otelmongo.NewMonitor()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func checkFindErr(err error) {
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return
		}
		log.Fatal(err)
	}
}
