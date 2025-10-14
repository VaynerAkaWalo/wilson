package main

import (
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xlog"
	"golang-template/internal/adapters"
	"golang-template/internal/adapters/action"
	"golang-template/internal/adapters/profile"
	"golang-template/internal/application/action"
	"golang-template/internal/domain/profile"
	"log"
	"log/slog"
)

func main() {
	slog.SetDefault(slog.New(xlog.NewPreConfiguredHandler()))

	profileStore := adapters.NewRepository()

	profileHandler := adapter_profile.HttpHandler{
		Service: profile.Service{
			ProfileRepository: profileStore,
		},
	}

	authProvider, err := xhttp.NewAuthenticationProvider()
	if err != nil {
		log.Fatal(err)
	}

	actionHandler := adapter_action.ActionHandler{
		Service: usecase_action.PerformActionService{
			ProfileRepository:  profileStore,
			LocationRepository: adapter_action.LocationStore{},
		},
	}

	actionHandler.StartActionLoop()

	authN := xhttp.NewAuthenticator(authProvider)

	httpServer := xhttp.Server{
		Addr:     ":8787",
		Handlers: []xhttp.RouteHandler{profileHandler},
		AuthN:    authN,
	}

	log.Fatal(httpServer.ListenAndServe())
}
