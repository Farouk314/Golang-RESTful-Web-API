package main

import (
	"log"
	"net/http"
)

// App represents the application and it's state
type App struct {
	Handler http.Handler
}

// Initialise the RESTful API
func (a *App) Initialise() {
	a.InitInMemoryData()
	a.Handler = a.InitCors(a.InitHandlers)
}

// Run initiates the application
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Handler))
}
