package adapter_location

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang-template/internal/domain/location"
	"golang-template/internal/domain/profile"
	"math/rand/v2"
)

type (
	LocationStore struct {
		conn *pgxpool.Pool
	}
)

func NewLocationStore(pool *pgxpool.Pool) *LocationStore {
	return &LocationStore{
		conn: pool,
	}
}

func (l LocationStore) Get(ctx context.Context, id location.Id) (location.Location, error) {
	rows, err := l.conn.Query(ctx, "select id, name, latitude, longitude, rewardMultiplier, type from locations where id=$1", string(id))
	if err != nil {
		return location.Location{}, err
	}
	return pgx.CollectOneRow[location.Location](rows, pgx.RowToStructByName[location.Location])
}

func (l LocationStore) GetAll(ctx context.Context) ([]location.Location, error) {
	rows, err := l.conn.Query(ctx, "select id, name, latitude, longitude, rewardMultiplier, type from locations")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows[location.Location](rows, pgx.RowToStructByName[location.Location])
}

func (l LocationStore) GetStartLocation(ctx context.Context) (profile.LocationId, error) {
	locations, err := l.GetAll(ctx)
	if err != nil {
		return "", err
	}

	for {
		randomLocation := locations[rand.IntN(len(locations))]

		if randomLocation.Type == location.BEACH {
			return profile.LocationId(randomLocation.Id), nil
		}
	}
}
