package main

import (
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xlog"
	"golang-template/internal/adapters/aprofile"
	"golang-template/internal/domain/profile"
	"log"
	"log/slog"
)

func main() {
	slog.SetDefault(slog.New(xlog.NewPreConfiguredHandler()))

	profileStore := aprofile.NewRepository()

	profileHandler := aprofile.HttpHandler{
		Service: profile.Service{
			ProfileRepository: profileStore,
		},
	}

	authProvider, err := xhttp.NewAuthenticationProvider()
	if err != nil {
		log.Fatal(err)
	}

	authN := xhttp.NewAuthenticator(authProvider)

	httpServer := xhttp.Server{
		Addr:     ":8787",
		Handlers: []xhttp.RouteHandler{profileHandler},
		AuthN:    authN,
	}

	log.Fatal(httpServer.ListenAndServe())
}
