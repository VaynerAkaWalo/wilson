package adapter_action

import (
	"context"
	"fmt"
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
	parentCtx := context.Background()
	profiles, err := handler.Service.GetEligibleProfiles(parentCtx)
	if err != nil {
		slog.ErrorContext(parentCtx, err.Error())
	}

	for _, prof := range profiles {
		err = handler.Service.Execute(parentCtx, prof.Id)
		if err != nil {
			slog.ErrorContext(parentCtx, fmt.Sprintf("action failed %v", err.Error()))
		}
	}
}
