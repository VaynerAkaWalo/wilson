package usecase_action

import (
	"context"
	"golang-template/internal/domain/action"
	"golang-template/pkg/ievent"
)

type (
	ProfileRepository interface {
		Get(context.Context, action.ProfileId) (action.Profile, error)
		GetAll(context.Context) ([]action.Profile, error)
	}

	LocationRepository interface {
		Get(context.Context, action.LocationId) (action.Location, error)
	}

	PerformActionService struct {
		ProfileRepository  ProfileRepository
		LocationRepository LocationRepository
		EventOrchestrator  *ievent.Orchestrator[action.Event]
	}
)

func (service PerformActionService) Execute(ctx context.Context, id action.ProfileId) error {
	profile, err := service.ProfileRepository.Get(ctx, id)
	if err != nil {
		return err
	}

	location, err := service.LocationRepository.Get(ctx, profile.Location)
	if err != nil {
		return err
	}

	act := action.New(profile, location)
	return service.EventOrchestrator.PublishEvent(ctx, act.CreateEvent())
}

func (service PerformActionService) GetEligibleProfiles(ctx context.Context) ([]action.Profile, error) {
	profiles, err := service.ProfileRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}
