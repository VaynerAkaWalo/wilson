package main

import (
	"context"
	"github.com/VaynerAkaWalo/go-toolkit/xevent"
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xlog"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/caarlos0/env/v11"
	"golang-template/internal/adapters"
	"golang-template/internal/adapters/action"
	"golang-template/internal/adapters/location"
	"golang-template/internal/adapters/profile"
	"golang-template/internal/application/action"
	"golang-template/internal/application/location"
	"golang-template/internal/domain/action"
	"golang-template/internal/domain/profile"
	"golang-template/internal/domain/transaction"
	"log"
	"log/slog"
)

type appConfig struct {
	AwsAccessKey string `env:"DDB_ACCESS_KEY"`
	AwsSecretKey string `env:"DDB_ACCESS_SECRET_KEY"`
}

func main() {
	slog.SetDefault(slog.New(xlog.NewPreConfiguredHandler(transaction.ContextKey, profile.ContextKey)))

	cfg, err := env.ParseAs[appConfig]()
	if err != nil {
		log.Fatal("unable to load env config")
	}

	cp := credentials.NewStaticCredentialsProvider(cfg.AwsAccessKey, cfg.AwsSecretKey, "")

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(cp), config.WithRegion("eu-north-1"))
	if err != nil {
		log.Fatal(err)
	}

	broker := xevent.NewBroker(action.Event{}, transaction.GoldChangeEvent{})

	profileStore := adapter_profile.NewDDBProfileStore(awsCfg)
	locationStore := adapters.NewLocationStore()

	profileHandler := adapter_profile.HttpHandler{
		Service: profile.Service{
			ProfileRepository:  profileStore,
			LocationRepository: locationStore,
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

	authN := xhttp.NewAuthenticator(authProvider)

	getLocation := usecase_location.NewGetLocationService(locationStore)

	locationHandler := adapter_location.HttpHandler{
		Service: *getLocation,
	}

	httpServer := xhttp.Server{
		Addr:     ":8787",
		Handlers: []xhttp.RouteHandler{profileHandler, locationHandler},
		AuthN:    authN,
	}

	log.Fatal(httpServer.ListenAndServe())
}
