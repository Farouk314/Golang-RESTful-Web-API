package main

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"mux"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/vladimiroff/jwt-go"

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

// In memory data
var certificates []Certificate              //Certificates
var users []User                            //Users
var mySigningKey = []byte("mockSigningKey") //mock signing key

func lookUpUserIDByName(userName string) (string, error) {
	for _, item := range users {
		if item.Name == userName {
			return item.ID, nil
		}
	}
	return "", errors.New("Could not get user by name" + userName)
}

//Homehandler for basic testing
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("HOMEPAGE"); err != nil {
		panic(err)
	}
}

// Initiate some in memory data
func initInMemoryData() {
	userA, userB := 0, 1
	users = append(users,
		User{
			ID:    "A",
			Email: "A@mail.com",
			Name:  "userA",
		},
		User{
			ID:    "B",
			Email: "B@mail.com",
			Name:  "userB",
		},
	)
	//userA's certificates:
	certificates = append(certificates,
		Certificate{
			ID:        "1",
			Title:     "A certificate title",
			CreatedAt: "do something for dates",
			OwnerID:   users[userA].ID,
			Year:      2018,
			Note:      "",
			Transfer: &Transfer{
				To:     "0002@mail.com",
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
	//userB's certificates
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

// Route handler GET certificates method
func getCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(certificates); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// Route handler GET certificates/{id} method
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

//Route handler get users certificates method
func getUsersCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usersCertificates []Certificate
	params := mux.Vars(r)
	for _, item := range certificates {
		if item.OwnerID == params["userId"] {
			usersCertificates = append(usersCertificates, item)
		}
	}
	if err := json.NewEncoder(w).Encode(usersCertificates); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// Route handler POST certificate method
func createCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var certificate Certificate

	// TODO
	user, _, _ := r.BasicAuth()
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&certificate); err != nil {
		fmt.Fprintf(w, err.Error())
	}
	if certificate.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"title": "Title must be populated"})
		return
	}

	certificate.OwnerID, _ = lookUpUserIDByName(user)
	certificate.OwnerID = "2"
	certificate.ID = strconv.Itoa(rand.Intn(1000000)) //Mock ID - not safe
	certificate.CreatedAt = "TODO Date"
	certificate.Year = time.Now().Year()
	certificates = append(certificates, certificate)
	if err := json.NewEncoder(w).Encode(certificate); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// Route handler PUT certificate/{id} method
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

// Route Handler DELETE certificate/{id} method
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
	user.ID = strconv.Itoa(rand.Intn(1000000)) //Mock ID - not safe
	users = append(users, user)
}

func createTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range certificates {
		if item.ID == params["id"] && item.OwnerID == params["userId"] {
			var certificate Certificate
			//certificates = append(certificates[:index], certificates[index+1:]...)
			if err := json.NewDecoder(r.Body).Decode(&certificate); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			item.Transfer.To = certificate.Transfer.To
			item.Transfer.Status = "Pending"
		}
	}
}

func acceptTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] {
			//need to check if userId's email corresponds to transfer.To email....
			certificate := item
			certificate.OwnerID = params["userId"]
			// certificate.Transfer.To = ""
			// certificate.Transfer.Status = ""
			certificates = append(certificates[index:], certificates[index+1:]...)
			certificates = append(certificates, certificate)
			if err := json.NewEncoder(w).Encode(certificate); err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			break
		}
	}

}

//JWT generator
func generateJWT() (string, error) {
	fmt.Println("generateJWT...")
	//Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	//Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	claims["admin"] = true
	claims["name"] = "Rob Wind"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	//Sign the token
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

//Get Token Handler
var getTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	validToken, err := generateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// fmt.Fprintf(w, validToken)
	if err = json.NewEncoder(w).Encode(validToken); err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//Write token to responsewriter
	//json.NewEncoder(w).Encode(tokenString)
	//w.Write([]byte(tokenString))
})

//Basic Authentication
func basicAuth(handler http.HandlerFunc, userName string, password string, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, pw, ok := r.BasicAuth()
		fmt.Println("u:" + u + "pw:" + pw)
		if !ok || subtle.ConstantTimeCompare([]byte(u), []byte(userName)) != 1 || subtle.ConstantTimeCompare([]byte(pw), []byte(password)) != 1 {
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
	//initInMemoryData()
	data := "userEmail:pw" //must be format user:password to pass BasicAuth Basic
	dataEncoded := base64.StdEncoding.EncodeToString([]byte(data))
	// Will print base64 encoded of data. Header request format: Basic <b64EncodedString>
	fmt.Printf("Data encoded: %v", dataEncoded)
	fmt.Println("")

	s := &http.Server{
		Addr:    ":8000",
		Handler: initRoutes(r),
	}

	log.Fatal(s.ListenAndServe())
}

// Router
func initRoutes() *Router {
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/certificates", getCertificates).Methods("GET")
	r.HandleFunc("/users/{id}/certificate", createCertificate).Methods("POST")
	r.HandleFunc("/certificate/{id}", getCertificate).Methods("GET")
	r.HandleFunc("/users/{userId}/certificate/{id}", patchCertificate).Methods("PATCH")
	r.HandleFunc("/certificate/{id}", deleteCertificate).Methods("DELETE")
	r.HandleFunc("/users/{userId}/certificates", getUsersCertificates).Methods("GET")
	r.HandleFunc("/users/{userId}/certificates/{id}/transfers", createTransfer).Methods("POST")                                          // PUT OR PATCH?
	r.HandleFunc("/users/{userId}/certificates/{id}/transfers", basicAuth(acceptTransfer, "userEmail", "pw", "my-realm")).Methods("PUT") //PUT or PATCH?
	r.HandleFunc("/api/signup", createUser).Methods("POST")
	r.HandleFunc("/get-token", getTokenHandler).Methods("GET")
	return r
}
