package main

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xlog"
	"golang-template/internal/adapters"
	"golang-template/internal/adapters/action"
	"golang-template/internal/adapters/profile"
	"golang-template/internal/application/action"
	"golang-template/internal/domain/action"
	"golang-template/internal/domain/profile"
	"golang-template/internal/domain/transaction"
	"golang-template/pkg/ievent"
	"log"
	"log/slog"
)

func main() {
	slog.SetDefault(slog.New(xlog.NewPreConfiguredHandler(transaction.ContextKey)))

	actionEventOrchestrator := ievent.NewOrchestrator[action.Event]()
	transactionEventOrchestrator := ievent.NewOrchestrator[transaction.GoldChangeEvent]()

	profileStore := adapters.NewRepository()

	profileHandler := adapter_profile.HttpHandler{
		Service: profile.Service{
			ProfileRepository: profileStore,
		},
		ActionEventOrchestrator:     actionEventOrchestrator,
		GoldChangeEventOrchestrator: transactionEventOrchestrator,
	}

	authProvider, err := xhttp.NewAuthenticationProvider()
	if err != nil {
		log.Fatal(err)
	}

	transactionService := &transaction.Service{
		BalanceStore:      profileStore,
		EventOrchestrator: transactionEventOrchestrator,
	}

	transactionActionHandler := transaction.NewActionHandler(transactionService, actionEventOrchestrator)
	go transactionActionHandler.StartEventConsumption(context.TODO())

	actionHandler := adapter_action.ActionHandler{
		Service: usecase_action.PerformActionService{
			ProfileRepository:  profileStore,
			LocationRepository: adapter_action.LocationStore{},
			EventOrchestrator:  actionEventOrchestrator,
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
