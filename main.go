package main

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Certificate Struct
type Certificate struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt string    `json:"createdAt"` // Will be a type date
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
	To     string `json:"to"`
	Status string `json:"status"`
}

// Certificates
var certificates []Certificate

// Users
var users []User

// lookUpUserIDByEmail returns a the UsserID for a specified email address
func lookUpUserIDByEmail(userEmail string) (string, error) {
	for _, item := range users {
		if item.Email == userEmail {
			return item.ID, nil
		}
	}
	return "", errors.New("Could not get user by email: " + userEmail)
}

// homehandler default route for testing
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("HOMEPAGE"); err != nil {
		panic(err)
	}
}

// initInMemoryData initiates some in memory data for testing
func initInMemoryData() {
	userA, userB := 0, 1
	users = append(users,
		User{
			ID:    "A",
			Email: "userAEmail",
			Name:  "userA",
		},
		User{
			ID:    "B",
			Email: "userBEmail",
			Name:  "userB",
		},
	)
	certificates = append(certificates,
		Certificate{
			ID:        "1",
			Title:     "A certificate title",
			CreatedAt: "do something for dates",
			OwnerID:   users[userA].ID,
			Year:      2018,
			Note:      "",
			Transfer: &Transfer{
				To:     "",
				Status: "",
			}},
		Certificate{
			ID:        "3",
			Title:     "A certificate title",
			CreatedAt: "do something for dates",
			OwnerID:   users[userA].ID,
			Year:      2019,
			Note:      "",
			Transfer: &Transfer{
				To:     "",
				Status: "",
			},
		},
	)
	certificates = append(certificates,
		Certificate{
			ID:        "2",
			Title:     "A certificate title",
			CreatedAt: "do something for dates",
			OwnerID:   users[userB].ID,
			Year:      2015,
			Note:      "",
			Transfer: &Transfer{
				To:     "",
				Status: "",
			}},
		Certificate{
			ID:        "7",
			Title:     "A certificate title",
			CreatedAt: "do something for dates",
			OwnerID:   users[userB].ID,
			Year:      2000,
			Note:      "",
			Transfer: &Transfer{
				To:     "",
				Status: "",
			},
		},
	)
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
		userID, err := lookUpUserIDByEmail(userEmail)
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

	certificate.OwnerID, _ = lookUpUserIDByEmail(userEmail)
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

// patchCertificate route handler PATCH certificate/{id} method
func patchCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] {
			certificate := certificates[index]
			certificates = append(certificates[:index], certificates[index+1:]...)

			if err := json.NewDecoder(r.Body).Decode(&certificate); err != nil {
				fmt.Fprintf(w, err.Error())
			}
			certificate.ID = params["id"]

			certificates = append(certificates, certificate)
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
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

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		panic(err)
	}
	//Todo (Farouk): Mock ID - not safe
	user.ID = strconv.Itoa(rand.Intn(1000000))
	users = append(users, user)
}

// createTransfer route handler create transfer PATCH method
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
			(*certificate).Transfer.Status = c.Transfer.Status
			json.NewEncoder(w).Encode(*certificate)
		}
		return
	}
}

// acceptTransfer route handler accept transfer PATCH method
func acceptTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userEmail, _, _ := r.BasicAuth()
	for index, item := range certificates {
		fmt.Println("Checking item " + item.ID + "againgst params ID: " + params["id"])
		if item.ID == params["id"] && item.Transfer.To == userEmail {
			fmt.Println("Certificate ID and userEmail match.")
			certificate := &item
			fmt.Println("Certificate.ID: " + (*certificate).ID)
			var err error
			coid, err := lookUpUserIDByEmail(userEmail)
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			(*certificate).OwnerID = coid
			(*certificate).Transfer.To = ""
			(*certificate).Transfer.Status = "accepted transfer"
			certificates = append(certificates[:index], certificates[index+1:]...)
			certificates = append(certificates, (*certificate))
			if err := json.NewEncoder(w).Encode(*certificate); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			break
		}
	}

}

// basicAuth basic http auth against email and password
func basicAuth(handler http.HandlerFunc, userEmail string, password string, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, pw, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(u), []byte(userEmail)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pw), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

// Main
func main() {
	fmt.Println("Running...")

	initInMemoryData()

	r := mux.NewRouter()
	userA, userB := 0, 1
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/certificates", getCertificates).Methods("GET")
	r.HandleFunc("/users/{id}/certificate", createCertificate).Methods("POST")
	r.HandleFunc("/certificate/{id}", getCertificate).Methods("GET")
	r.HandleFunc("/users/{userId}/certificate/{id}", patchCertificate).Methods("PATCH")
	r.HandleFunc("/certificate/{id}", deleteCertificate).Methods("DELETE")
	r.HandleFunc("/users/{userId}/certificates", getUsersCertificates).Methods("GET")
	r.HandleFunc("/certificates/{id}/transfers",
		basicAuth(createTransfer, users[userA].Email, "userApw", "my-realm")).Methods("PATCH") // PUT OR PATCH?
	r.HandleFunc("/certificates/{id}/transfers",
		basicAuth(acceptTransfer, users[userB].Email, "userBpw", "my-realm")).Methods("PUT") //PUT or PATCH?
	r.HandleFunc("/api/signup", createUser).Methods("POST")

	s := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	log.Fatal(s.ListenAndServe())
}
