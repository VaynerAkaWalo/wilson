package transaction

import (
	"context"
	"errors"
	"fmt"
	"github.com/VaynerAkaWalo/go-toolkit/xevent"
	"github.com/VaynerAkaWalo/go-toolkit/xuuid"
	"golang-template/internal/domain/profile"
	"log/slog"
)

type (
	Service struct {
		BalanceStore BalanceRepository
		Broker       *xevent.Broker
	}

	BalanceRepository interface {
		GetBalance(ctx context.Context, profile string) (Balance, error)
		UpdateBalance(ctx context.Context, balance Balance) (Balance, error)
	}
)

func (service *Service) Perform(ctx context.Context, transaction Transaction) error {
	ctx = context.WithValue(ctx, ContextKey, string(transaction.Id))
	ctx = context.WithValue(ctx, profile.ContextKey, transaction.Profile)
	slog.InfoContext(ctx, "Attempting to perform transaction")

	for attempt := range 3 {
		balance, err := service.BalanceStore.GetBalance(ctx, transaction.Profile)
		if err != nil {
			return err
		}

		balance.Gold += transaction.BalanceChange

		balance, err = service.BalanceStore.UpdateBalance(ctx, balance)
		if err == nil {
			slog.InfoContext(ctx, "Transaction completed successfully")

			changeEvent := GoldChangeEvent{
				Id:          xuuid.UUID(),
				Profile:     balance.Profile,
				GoldBalance: balance.Gold,
			}

			xevent.PublishEvent[GoldChangeEvent](service.Broker, ctx, changeEvent)
			return err
		}

		if !errors.Is(err, VersionMismatchError{}) {
			slog.ErrorContext(ctx, fmt.Sprintf("Transaction failed because of unknown error %s", err.Error()))
			return err
		}

		slog.WarnContext(ctx, fmt.Sprintf("Transaction attempt %d failed", attempt))
	}

	return fmt.Errorf("all transaction %s attempts failed", transaction.Id)
}
