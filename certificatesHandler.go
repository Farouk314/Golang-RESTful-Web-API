package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// CreateCertificate route handler POST certificate method
func (a *App) CreateCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userName, _, _ := r.BasicAuth()
	userID, err := LookUpUserIDByName(userName)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	var certificate Certificate
	var rc Certificate

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&rc); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	if rc.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Title must be populated")
		return
	}
	if rc.CreatedAt.Format(time.RFC3339) == "0001-01-01T00:00:00Z" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "CreatedAt must be populated")
		return
	}
	params := mux.Vars(r)
	//Set from req Body
	certificate.Title = rc.Title
	certificate.Note = rc.Note
	certificate.CreatedAt = rc.CreatedAt
	//Set by server
	certificate.OwnerID = userID
	certificate.ID = params["id"]
	certificate.Year = certificate.CreatedAt.Year()
	certificate.Transfer = &Transfer{To: "", Status: ""}
	//Add newly created certificate to in memory data
	certificates = append(certificates, certificate)
	if err := json.NewEncoder(w).Encode(certificate); err != nil {
		fmt.Fprintf(w, "Could not encode to JSON: %v", err)
	}
}

// GetUsersCertificates route handler get user's certificates method
func (a *App) GetUsersCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userName, _, _ := r.BasicAuth()
	userID, err := LookUpUserIDByName(userName)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	if userID != params["userId"] {
		fmt.Fprintf(w, "You cannot access user "+params["userId"]+"'s certificates")
		return
	}
	var usersCertificates []Certificate
	for _, item := range certificates {
		if item.OwnerID == userID && params["userId"] == userID {
			usersCertificates = append(usersCertificates, item)
		}
	}
	if err := json.NewEncoder(w).Encode(usersCertificates); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// UpdateCertificate route handler PATCH certificate/{id} method
func (a *App) UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userName, _, _ := r.BasicAuth()
	userID, err := LookUpUserIDByName(userName)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] && item.OwnerID == userID {
			var c Certificate

			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			if c.Title == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Title cannot be empty")
				return
			}
			certificates[index].Title = c.Title
			certificates[index].Note = c.Note
			if err := json.NewEncoder(w).Encode(certificates[index]); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
		}
	}
}

// DeleteCertificate route Handler DELETE certificate/{id} method
func (a *App) DeleteCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userName, _, _ := r.BasicAuth()
	userID, err := LookUpUserIDByName(userName)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] && item.OwnerID == userID {
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
	userName, _, _ := r.BasicAuth()
	userID, err := LookUpUserIDByName(userName)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] && item.OwnerID == userID {
			var c Certificate

			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			if c.Transfer.To == "" {
				fmt.Fprintf(w, "Transfer to must be populated")
				return
			}
			certificates[index].Transfer.To = c.Transfer.To
			certificates[index].Transfer.Status = "Pending transfer"
			if err := json.NewEncoder(w).Encode(certificates[index]); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			return
		}
	}
}

// AcceptTransfer route handler certificates/{id}/transfers transfer PUT method
func (a *App) AcceptTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userName, _, _ := r.BasicAuth()
	userID, err := LookUpUserIDByName(userName)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] && item.Transfer.To == userName {
			certificates[index].OwnerID = userID
			certificates[index].Transfer.To = ""
			certificates[index].Transfer.Status = ""
			if err := json.NewEncoder(w).Encode(certificates[index]); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			return
		}
	}
}
