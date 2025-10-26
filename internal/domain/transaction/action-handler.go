package transaction

import (
	"context"
	"fmt"
	"github.com/VaynerAkaWalo/go-toolkit/xevent"
	"golang-template/internal/domain/action"
	"log/slog"
)

type (
	ActionHandler struct {
		service *Service
		broker  *xevent.Broker
	}
)

func NewActionHandler(service *Service, broker *xevent.Broker) *ActionHandler {
	return &ActionHandler{
		service: service,
		broker:  broker,
	}
}

func (handler *ActionHandler) StartEventConsumption(ctx context.Context) error {
	eventChannel := xevent.RegisterListener[action.Event](handler.broker, ctx)
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
