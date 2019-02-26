package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// InitHandlers : Initiate route handlers
func (a *App) InitHandlers() http.Handler {
	r := mux.NewRouter()
	//Endpoints
	r.HandleFunc("/certificates/{id}", a.CreateCertificate).Methods("POST")
	r.HandleFunc("/certificates/{id}", a.UpdateCertificate).Methods("PATCH")
	r.HandleFunc("/certificates/{id}", a.DeleteCertificate).Methods("DELETE")
	r.HandleFunc("/users/{userId}/certificates", a.GetUsersCertificates).Methods("GET")
	r.HandleFunc("/certificates/{id}/transfers", a.CreateTransfer).Methods("PATCH")
	r.HandleFunc("/certificates/{id}/transfers", a.AcceptTransfer).Methods("PUT")
	return r
}
