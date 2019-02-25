package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// InitHandlers : Initiate route handlers
func (a *App) InitHandlers() http.Handler {
	r := mux.NewRouter()
	userA, userB := 0, 1
	//Endpoints
	r.HandleFunc("/certificates/{id}", a.CreateCertificate).Methods("POST")
	r.HandleFunc("/certificates/{id}", a.UpdateCertificate).Methods("PATCH")
	r.HandleFunc("/certificates/{id}", a.DeleteCertificate).Methods("DELETE")
	r.HandleFunc("/users/{userId}/certificates",
		basicAuth(a.GetUsersCertificates, users[userA].Email, "userApw", "my-realm")).Methods("GET")
	r.HandleFunc("/certificates/{id}/transfers",
		basicAuth(a.CreateTransfer, users[userA].Email, "userApw", "my-realm")).Methods("PATCH")
	r.HandleFunc("/certificates/{id}/transfers",
		basicAuth(a.AcceptTransfer, users[userB].Email, "userBpw", "my-realm")).Methods("PUT")
	return r
}
