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

// HomeHandler default route for testing
func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode("HOMEPAGE"); err != nil {
	// 	fmt.Fprintf(w, err.Error())
	// }
	text := r.FormValue("v")
	if text == "" {
		http.Error(w, "Missing value", http.StatusBadRequest)
		return
	}

	v, err := strconv.Atoi(text)
	if err != nil {
		http.Error(w, "not a number: "+text, http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, v*2)
}

// GetCertificates route handler GET certificates method
func (a *App) GetCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(certificates); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// GetCertificate route handler GET certificates/{id} method
func (a *App) GetCertificate(w http.ResponseWriter, r *http.Request) {
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

// GetUsersCertificates route handler get user's certificates method
func (a *App) GetUsersCertificates(w http.ResponseWriter, r *http.Request) {
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

// CreateCertificate route handler POST certificate method
func (a *App) CreateCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var certificate Certificate
	userEmail, _, _ := r.BasicAuth()
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&certificate); err != nil {
		fmt.Fprintf(w, err.Error())
		return
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
	//Go lang time formatting must reference: 2006-01-02T15:04:05Z07:00 (RF3339 Layout)
	// cca, err := time.Parse(time.ANSIC, certificate.CreatedAt)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(map[string]interface{}{
	// 		"createdAt": "Date must be in format ANSIC (Mon Jan _2 15:04:05 2006)",
	// 	})
	// }
	// certificate.CreatedAt = time.Now()
	certificate.Year = time.Now().Year()
	certificates = append(certificates, certificate)
	if err := json.NewEncoder(w).Encode(certificate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"title": "Title must be populated",
		})
	}
}

// UpdateCertificate route handler PATCH certificate/{id} method
func (a *App) UpdateCertificate(w http.ResponseWriter, r *http.Request) {
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

// DeleteCertificate route Handler DELETE certificate/{id} method
func (a *App) DeleteCertificate(w http.ResponseWriter, r *http.Request) {
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

// CreateTransfer route handler certificates/{id}/transfers PATCH method
func (a *App) CreateTransfer(w http.ResponseWriter, r *http.Request) {
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

// AcceptTransfer route handler certificates/{id}/transfers transfer PUT method
func (a *App) AcceptTransfer(w http.ResponseWriter, r *http.Request) {
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
