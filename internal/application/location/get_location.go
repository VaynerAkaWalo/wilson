package usecase_location

import (
	"context"
	"golang-template/internal/domain/location"
)

type (
	Repository interface {
		Get(context.Context, location.Id) (location.Location, error)
		GetAll(ctx context.Context) ([]location.Location, error)
	}

	GetLocationService struct {
		repository Repository
	}
)

func NewGetLocationService(repository Repository) *GetLocationService {
	return &GetLocationService{
		repository: repository,
	}
}

func (service *GetLocationService) GetLocation(ctx context.Context, id location.Id) (location.Location, error) {
	return service.repository.Get(ctx, id)
}

func (service *GetLocationService) GetAllLocation(ctx context.Context) ([]location.Location, error) {
	return service.repository.GetAll(ctx)
}
