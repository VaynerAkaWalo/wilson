package adapters

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/ojrac/opensimplex-go"
	"golang-template/internal/domain/location"
	"golang-template/internal/domain/profile"
	"maps"
	"math/rand/v2"
	"net/http"
	"slices"
)

const (
	scale float64 = 6.5
)

type (
	LocationStore struct {
		locations map[string]location.Location
	}
)

func NewLocationStore() *LocationStore {
	locationMap := make(map[string]location.Location)
	noise := opensimplex.NewNormalized(6767)

	for x := range 32 {
		for y := range 32 {
			noiseLevel := noise.Eval2(float64(x)/scale, float64(y)/scale)

			loc := location.New(x, y, 1+rand.Float64(), noiseLevel)
			locationMap[string(loc.Id)] = *loc
		}
	}

	return &LocationStore{
		locations: locationMap,
	}
}

func (l LocationStore) Get(ctx context.Context, id location.Id) (location.Location, error) {
	loc, found := l.locations[string(id)]
	if !found {
		return location.Location{}, xhttp.NewError("location not found", http.StatusNotFound)
	}

	return loc, nil
}

func (l LocationStore) GetAll(ctx context.Context) ([]location.Location, error) {
	return slices.Collect(maps.Values(l.locations)), nil
}

func (l LocationStore) GetStartLocation(ctx context.Context) (profile.LocationId, error) {
	locations := slices.Collect(maps.Values(l.locations))

	for {
		randomLocation := locations[rand.IntN(len(locations))]

		if randomLocation.Type == location.BEACH {
			return profile.LocationId(randomLocation.Id), nil
		}
	}
}
