package main

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xevent"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xlog"
	"golang-template/internal/adapters"
	"golang-template/internal/adapters/action"
	"golang-template/internal/adapters/profile"
	"golang-template/internal/application/action"
	"golang-template/internal/domain/action"
	"golang-template/internal/domain/profile"
	"golang-template/internal/domain/transaction"
	"log"
	"log/slog"
)

func main() {
	slog.SetDefault(slog.New(xlog.NewPreConfiguredHandler(transaction.ContextKey, profile.ContextKey)))

	broker := xevent.NewBroker(action.Event{}, transaction.GoldChangeEvent{})

	profileStore := adapters.NewRepository()

	profileHandler := adapter_profile.HttpHandler{
		Service: profile.Service{
			ProfileRepository: profileStore,
		},
		Broker: broker,
	}

	authProvider, err := xhttp.NewAuthenticationProvider()
	if err != nil {
		log.Fatal(err)
	}

	transactionService := &transaction.Service{
		BalanceStore: profileStore,
		Broker:       broker,
	}

	transactionActionHandler := transaction.NewActionHandler(transactionService, broker)
	go transactionActionHandler.StartEventConsumption(context.TODO())

	actionHandler := adapter_action.ActionHandler{
		Service: usecase_action.PerformActionService{
			ProfileRepository:  profileStore,
			LocationRepository: adapter_action.LocationStore{},
			Broker:             broker,
		},
	}

	actionHandler.StartActionLoop()

	authN := xhttp.NewAuthenticator(authProvider, "GET /event")

	httpServer := xhttp.Server{
		Addr:     ":8787",
		Handlers: []xhttp.RouteHandler{profileHandler},
		AuthN:    authN,
	}

	log.Fatal(httpServer.ListenAndServe())
}
