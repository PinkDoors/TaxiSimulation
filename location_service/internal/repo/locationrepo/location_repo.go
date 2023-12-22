package locationrepo

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"location_service/internal/model"
	"location_service/internal/repo"
)

type locationRepo struct {
	pgxPool *pgxpool.Pool
}

func (r *locationRepo) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
		return tx
	}

	return r.pgxPool
}

func (r *locationRepo) UpdateDriverLocation(ctx context.Context, driverId string, lat, lng float32) error {
	conn := r.conn(ctx)

	query := `UPDATE drivers SET lat=$1, lng=$2 WHERE id=$3`
	_, err := conn.Exec(ctx, query, lat, lng, driverId)
	if err != nil {
		return err
	}

	return nil
}

func (r *locationRepo) ListDriver(ctx context.Context, radius, lat, lng float32) ([]*model.Driver, error) {
	conn := r.conn(ctx)

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

	return drivers, nil
}

func New(pgxPool *pgxpool.Pool) repo.Driver {
	return &locationRepo{
		pgxPool: pgxPool,
	}
}
