package locationrepo

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"location_service/internal/model"
	"location_service/internal/repo"
)

type locationRepo struct {
	pgxPool *pgxpool.Pool
	logger  *zap.Logger
}

func (r *locationRepo) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
		return tx
	}

	return r.pgxPool
}

func (r *locationRepo) UpdateDriverLocation(ctx context.Context, driverId string, lat, lng float32) error {
	conn := r.conn(ctx)

	r.logger.Info(
		"Try to update driver location",
		zap.String("DriverId", driverId),
		zap.Float32("Latitude", lat),
		zap.Float32("Longitude", lng),
	)
	query := `UPDATE drivers SET lat=$1, lng=$2 WHERE id=$3`
	_, err := conn.Exec(ctx, query, lat, lng, driverId)
	if err != nil {
		return err
	}
	r.logger.Info(
		"Successful update of driver",
		zap.String("DriverId", driverId),
	)

	return nil
}

func (r *locationRepo) ListDriver(ctx context.Context, radius, lat, lng float32) ([]*model.Driver, error) {
	conn := r.conn(ctx)

	r.logger.Info(
		"Try to found drivers inside the circle",
		zap.Float32("Latitude", lat),
		zap.Float32("Longitude", lng),
		zap.Float32("Radius", radius),
	)
	query := `SELECT id, lat, lng FROM drivers WHERE earth_box(ll_to_earth($1, $2), $3) @> ll_to_earth(lat, lng)`
	rows, err := conn.Query(ctx, query, lat, lng, radius)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	drivers := make([]*model.Driver, 0)
	for rows.Next() {
		driver := new(model.Driver)
		err := rows.Scan(&driver.DriverId, &driver.Lat, &driver.Lng)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, driver)
	}
	r.logger.Info(
		"Drivers found",
		zap.Int("count", len(drivers)),
	)

	return drivers, nil
}

func New(pgxPool *pgxpool.Pool, logger *zap.Logger) repo.Driver {
	return &locationRepo{
		pgxPool: pgxPool,
		logger:  logger,
	}
}
