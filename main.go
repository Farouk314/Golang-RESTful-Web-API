package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Certificate Struct
type Certificate struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt string    `json:"createdAt"` //should be type date?
	OwnerID   string    `json:"ownerId"`
	Year      int       `json:"year"`
	Note      string    `json:"note"`
	Transfer  *Transfer `json:"transfer"`
}

// User Struct
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Transfer Struct
type Transfer struct {
	To     string `json:"to"` //should be type email?
	Status string `json:"status"`
}

// In memory data
var certificates []Certificate //Certificates
var users []User               //Users

// Initiate some certificates
func initCertificates() {
	certificates = append(certificates, Certificate{
		ID:        "1",
		Title:     "A certificate title",
		CreatedAt: "do something for dates",
		OwnerID:   "guid1",
		Year:      2019,
		Note:      "First Certificate",
		Transfer: &Transfer{
			To:     "",
			Status: "",
		},
	})
	certificates = append(certificates, Certificate{
		ID:        "2",
		Title:     "A certificate title",
		CreatedAt: "do something for dates",
		OwnerID:   "guid2",
		Year:      2015,
		Note:      "Second Certificate",
		Transfer: &Transfer{
			To:     "",
			Status: "",
		},
	})
}

// Route handler GET certificates method
func getCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Encoding to the certificates ResponseWriter...")
	json.NewEncoder(w).Encode(certificates)
}

// Route handler GET certificates/{id} method
func getCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println("Looping through all certificates...")
	for _, item := range certificates {
		if item.ID == params["id"] {
			fmt.Println("[getCertificate] Writing to response writer certificate{id}=" + item.ID)
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// Route handler POST certificate method
func createCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var certificate Certificate
	_ = json.NewDecoder(r.Body).Decode(&certificate)  //
	certificate.ID = strconv.Itoa(rand.Intn(1000000)) //Mock ID - not safe
	certificates = append(certificates, certificate)
	json.NewEncoder(w).Encode(certificate)
}

// Route handler PATCH certificate/{id} method
func patchCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//To do: PATCH method can update individual fields in certificate
}

// Route handler PUT certificate/{id} method
func putCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			var certificate Certificate
			_ = json.NewDecoder(r.Body).Decode(&certificate) //PUT method: implement error if missing fields
			certificate.ID = params["id"]
			certificates = append(certificates, certificate)
			json.NewEncoder(w).Encode(certificate)
		}
	}
}

// Route Handler DELETE certificate/{id} method
func deleteCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println("Looping through all certificates...")
	for index, item := range certificates {
		if item.ID == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			fmt.Println("Deleted certificate with id=" + item.ID)
			break
		}
	}
	json.NewEncoder(w).Encode(certificates)
}

// Main
func main() {
	fmt.Println("Running...")

	initCertificates()

	r := mux.NewRouter()
	r.HandleFunc("/api/certificates", getCertificates).Methods("GET")
	r.HandleFunc("/api/certificates/{id}", getCertificate).Methods("GET")
	r.HandleFunc("/api/certificates", createCertificate).Methods("POST")
	r.HandleFunc("/api/certificates/{id}", putCertificate).Methods("PUT")
	r.HandleFunc("/api/certificates/{id}", patchCertificate).Methods("PATCH")
	r.HandleFunc("/api/certificates/{id}", deleteCertificate).Methods("DELETE")

	s := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	log.Fatal(s.ListenAndServe())
}
