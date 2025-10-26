package transaction

import (
	"context"
	"fmt"
	"golang-template/internal/domain/action"
	"golang-template/pkg/ievent"
	"log/slog"
)

type (
	ActionHandler struct {
		service           *Service
		eventOrchestrator *ievent.Orchestrator[action.Event]
	}
)

func NewActionHandler(service *Service, eventOrchestrator *ievent.Orchestrator[action.Event]) *ActionHandler {
	return &ActionHandler{
		service:           service,
		eventOrchestrator: eventOrchestrator,
	}
}

func (handler *ActionHandler) StartEventConsumption(ctx context.Context) error {
	eventChannel := handler.eventOrchestrator.RegisterListener(ctx)
	slog.InfoContext(ctx, "started consumption for action events in transaction domain")

	for {
		select {
		case event := <-eventChannel:

			trans, err := New(string(event.ProfileId), event.GoldReward)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Failed to create transaction for action %s", event.Id))
				break
			}

			err = handler.service.Perform(ctx, trans)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Failed to process transaction %s for event %s", trans.Id, event.Id))
				break
			}
		}
	}
}
