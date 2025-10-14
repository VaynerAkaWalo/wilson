package adapter_action

import (
	"context"
	"golang-template/internal/domain/action"
)

type (
	LocationStore struct{}
)

func (l LocationStore) Get(ctx context.Context, id action.LocationId) (action.Location, error) {
	return action.Location{
		Id:         id,
		Multiplier: 0.35,
	}, nil
}
