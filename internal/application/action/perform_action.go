package usecase_action

import (
	"context"
	"golang-template/internal/domain/action"
)

type (
	ProfileRepository interface {
		Get(context.Context, action.ProfileId) (action.Profile, error)
		GetAll(context.Context) ([]action.Profile, error)
		UpdateGold(context.Context, action.Profile) error
	}

	LocationRepository interface {
		Get(context.Context, action.LocationId) (action.Location, error)
	}

	PerformActionService struct {
		ProfileRepository  ProfileRepository
		LocationRepository LocationRepository
	}
)

func (service PerformActionService) Execute(ctx context.Context, profileId action.ProfileId) error {
	profile, err := service.ProfileRepository.Get(ctx, profileId)
	if err != nil {
		return err
	}

	location, err := service.LocationRepository.Get(ctx, profile.Location)
	if err != nil {
		return err
	}

	act := action.New(profile, location)
	err = profile.ConsumeAction(*act)
	if err != nil {
		return err
	}

	return service.ProfileRepository.UpdateGold(ctx, profile)
}

func (service PerformActionService) GetEligibleProfiles(ctx context.Context) ([]action.ProfileId, error) {
	profiles, err := service.ProfileRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]action.ProfileId, 0)

	for _, prof := range profiles {
		ids = append(ids, prof.Id)
	}

	return ids, nil
}
