package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// homehandler default route for testing
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("HOMEPAGE"); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// getCertificates route handler GET certificates method
func getCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(certificates); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// getCertificate route handler GET certificates/{id} method
func getCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range certificates {
		if item.ID == params["id"] {
			if err := json.NewEncoder(w).Encode(item); err != nil {
				fmt.Fprintf(w, err.Error())
			}
			return
		}
	}
}

// getUsersCertificates route handler get user's certificates method
func getUsersCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userEmail, _, _ := r.BasicAuth()
	var usersCertificates []Certificate
	for _, item := range certificates {
		userID, err := LookUpUserIDByEmail(userEmail)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		if item.OwnerID == userID && params["userId"] == userID {
			usersCertificates = append(usersCertificates, item)
		}
		if userID != params["userId"] {
			fmt.Fprintf(w, "You cannot access user "+params["userId"]+"'s certificates")
			return
		}
	}
	if err := json.NewEncoder(w).Encode(usersCertificates); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// createCertificate route handler POST certificate method
func createCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var certificate Certificate
	userEmail, _, _ := r.BasicAuth()
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&certificate); err != nil {
		fmt.Fprintf(w, err.Error())
	}
	if certificate.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"title": "Title must be populated",
		})
		return
	}

	certificate.OwnerID, _ = LookUpUserIDByEmail(userEmail)
	// Todo (Farouk): Mock ID - not safe
	certificate.ID = strconv.Itoa(rand.Intn(1000000))
	// Todo (Farouk): Parse as date
	certificate.CreatedAt = "DO SOMETHING"
	certificate.Year = time.Now().Year()
	certificates = append(certificates, certificate)
	if err := json.NewEncoder(w).Encode(certificate); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// updateCertificate route handler PATCH certificate/{id} method
func updateCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range certificates {
		if item.ID == params["id"] {
			certificate := &item
			var c Certificate
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				fmt.Fprintf(w, err.Error())
			}
			(*certificate).Title = c.Title
			(*certificate).Note = c.Note
			if err := json.NewEncoder(w).Encode((*certificate)); err != nil {
				fmt.Fprintf(w, err.Error())
			}
			return
		}
	}
}

// deleteCertificate route Handler DELETE certificate/{id} method
func deleteCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			break
		}
	}
	if err := json.NewEncoder(w).Encode(certificates); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// createTransfer route handler certificates/{id}/transfers PATCH method
func createTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range certificates {
		if item.ID == params["id"] {
			certificate := &item
			var c Certificate
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			(*certificate).Transfer.To = c.Transfer.To
			(*certificate).Transfer.Status = "Pending transfer"
			json.NewEncoder(w).Encode(*certificate)
		}
		return
	}
}

// acceptTransfer route handler certificates/{id}/transfers transfer PUT method
func acceptTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userEmail, _, _ := r.BasicAuth()
	for _, item := range certificates {
		fmt.Println("Checking item " + item.ID + "against params ID: " + params["id"])
		if item.ID == params["id"] && item.Transfer.To == userEmail {
			certificate := &item
			var err error
			coid, err := LookUpUserIDByEmail(userEmail)
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			(*certificate).OwnerID = coid
			(*certificate).Transfer.To = ""
			(*certificate).Transfer.Status = ""
			if err := json.NewEncoder(w).Encode(*certificate); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			break
		}
	}
}
