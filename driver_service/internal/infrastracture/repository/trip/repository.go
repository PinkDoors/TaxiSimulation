package trip

import (
	"context"
	config "driver_service/configs"
	"driver_service/internal/domain/models"
	"driver_service/internal/domain/models/trip"
	domainRepo "driver_service/internal/domain/repository/trip"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
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

func (r *Repository) getMongoClient(ctx context.Context) (*mongo.Client, error) {
	opts := options.Client()
	opts.ApplyURI(r.config.DB.URI)
	opts.SetTimeout(time.Duration(r.config.DB.TIMEOUT) * time.Second)
	optsAuth := options.Credential{
		Username:   r.config.DB.USERNAME,
		Password:   r.config.DB.PASSWORD,
		AuthSource: r.config.DB.AUTHSOURCE,
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
	checkFindErr(err)

	var foundTrips []trip.Trip
	err = cursor.All(ctx, &foundTrips)
	if err != nil {
		r.logger.Error("Failed to read trips from Mongo.", zap.Error(err))
		return nil, err
	}

	return foundTrips, nil
}

func (r *Repository) GetTrip(ctx context.Context, tripId uuid.UUID) (trip.Trip, error) {
	client, err := r.getMongoClient(ctx)
	if err != nil {
		r.logger.Error("Failed to create Mongo client.", zap.Error(err))
		return trip.Trip{}, err
	}

	defer func() {
		if mongoDisconnectErr := client.Disconnect(ctx); mongoDisconnectErr != nil {
			r.logger.Error("Failed to disconnect from Mongo client.", zap.Error(err))
		}
	}()

	col := client.Database("driver-service").Collection("trips")
	result := col.FindOne(ctx, bson.M{"id": tripId})
	checkFindErr(result.Err())

	var foundTrip trip.Trip
	if err = result.Decode(&foundTrip); err != nil {
		r.logger.Error("Failed to read trips from Mongo.", zap.Error(err))
		return trip.Trip{}, err
	}

	return foundTrip, nil
}
