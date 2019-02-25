package main

import (
	"net/http"

	"github.com/gorilla/handlers"
)

// InitCors initialises cors for a given handler func
func (a *App) InitCors(handlerFn func() http.Handler) http.Handler {
	headers := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Content-Type",
		"Authorization",
	})

	methods := handlers.AllowedMethods(
		[]string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"HEAD",
			"OPTIONS",
		})

	origin := handlers.AllowedOrigins([]string{"*"})

	return handlers.CORS(headers, methods, origin)(handlerFn())
}
