package main

import (
	"github.com/VaynerAkaWalo/go-toolkit/xhttp"
	"github.com/VaynerAkaWalo/go-toolkit/xlog"
	"log"
	"log/slog"
)

func main() {
	slog.SetDefault(slog.New(xlog.NewPreConfiguredHandler()))

	authProvider, err := xhttp.NewAuthenticationProvider()
	if err != nil {
		log.Fatal(err)
	}

	authN := xhttp.NewAuthenticator(authProvider)

	httpServer := xhttp.Server{
		Addr:     ":8000",
		Handlers: []xhttp.RouteHandler{},
		AuthN:    authN,
	}

	log.Fatal(httpServer.ListenAndServe())
}
