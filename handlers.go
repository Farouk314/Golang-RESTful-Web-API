package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// InitHandlers : Initiate route handlers
func InitHandlers() http.Handler {
	r := mux.NewRouter()
	userA, userB := 0, 1
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/certificates", getCertificates).Methods("GET")
	r.HandleFunc("/certificate/{id}", createCertificate).Methods("POST")
	r.HandleFunc("/certificate/{id}", getCertificate).Methods("GET")
	r.HandleFunc("/certificate/{id}", updateCertificate).Methods("PATCH")
	r.HandleFunc("/certificate/{id}", deleteCertificate).Methods("DELETE")
	r.HandleFunc("/users/{userId}/certificates", getUsersCertificates).Methods("GET")

	// TODO (Farouk): PUT OR PATCH?
	r.HandleFunc("/certificates/{id}/transfers",
		basicAuth(createTransfer, users[userA].Email, "userApw", "my-realm")).Methods("PATCH")
	// TODO (Farouk): PUT or PATCH?
	r.HandleFunc("/certificates/{id}/transfers",
		basicAuth(acceptTransfer, users[userB].Email, "userBpw", "my-realm")).Methods("PUT")
	return r
}
