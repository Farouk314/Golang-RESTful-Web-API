package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// InitHandlers : Initiate route handlers
func (a *App) InitHandlers() http.Handler {
	r := mux.NewRouter()
	userA, userB := 0, 1
	r.HandleFunc("/Home", a.HomeHandler).Methods("GET")
	r.HandleFunc("/certificates", a.GetCertificates).Methods("GET")
	r.HandleFunc("/certificates/{id}", a.CreateCertificate).Methods("POST")
	r.HandleFunc("/certificates/{id}", a.GetCertificate).Methods("GET")
	r.HandleFunc("/certificates/{id}", a.UpdateCertificate).Methods("PATCH")
	r.HandleFunc("/certificates/{id}", a.DeleteCertificate).Methods("DELETE")
	r.HandleFunc("/users/{userId}/certificates", a.GetUsersCertificates).Methods("GET")

	// TODO (Farouk): PUT OR PATCH?
	r.HandleFunc("/certificates/{id}/transfers",
		basicAuth(a.CreateTransfer, users[userA].Email, "userApw", "my-realm")).Methods("PATCH")
	// TODO (Farouk): PUT or PATCH?
	r.HandleFunc("/certificates/{id}/transfers",
		basicAuth(a.AcceptTransfer, users[userB].Email, "userBpw", "my-realm")).Methods("PUT")
	return r
}
