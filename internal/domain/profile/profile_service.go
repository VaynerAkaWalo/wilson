package profile

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"log/slog"
	"net/http"
)

type (
	Repository interface {
		GetProfilesByOwner(context.Context, OwnerId) ([]Profile, error)
		Save(context.Context, *Profile) error
	}
	LocationRepository interface {
		GetStartLocation(context.Context) (LocationId, error)
	}

	Service struct {
		ProfileRepository  Repository
		LocationRepository LocationRepository
	}
)

func (service Service) GetProfilesByOwner(ctx context.Context, id OwnerId) ([]Profile, error) {
	return service.ProfileRepository.GetProfilesByOwner(ctx, id)
}

func (service Service) CreateProfile(ctx context.Context, name string) (*Profile, error) {
	ownerId, ok := ctx.Value(xhttp.UserId).(string)
	if !ok {
		return nil, xhttp.NewError("failed to get owner for profile", http.StatusInternalServerError)
	}

	startLocation, err := service.LocationRepository.GetStartLocation(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "unable to get start location for new profile")
		return nil, xhttp.NewError("unexpected error while creating profile", http.StatusInternalServerError)
	}

	newProfile, err := New(name, ownerId, startLocation)
	if err != nil {
		return nil, err
	}

	err = service.ProfileRepository.Save(ctx, newProfile)
	return newProfile, err
}
