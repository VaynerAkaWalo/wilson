package adapter_action

import (
	"context"
	"golang-template/internal/application/action"
	"log/slog"
	"time"
)

type (
	ActionHandler struct {
		Service usecase_action.PerformActionService
	}
)

func (handler ActionHandler) StartActionLoop() {
	ticker := time.NewTicker(6 * time.Second)

	go func() {
		for {
			handler.performActions()
			<-ticker.C
		}
	}()
}

func (handler ActionHandler) performActions() {
	ctx := context.TODO()
	profiles, err := handler.Service.GetEligibleProfiles(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	for _, prof := range profiles {
		slog.InfoContext(ctx, string("performing action for profile "+prof))
		handler.Service.Execute(ctx, prof)
	}
}
