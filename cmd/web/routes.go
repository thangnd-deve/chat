package main

import (
	"chat/internal/handlers"
	"github.com/bmizerany/pat"
	"net/http"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Hom))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	return mux
}
