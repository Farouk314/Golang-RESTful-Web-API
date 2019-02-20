package main

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

//Homehandler for basic testing
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("HOMEPAGE")
}

// Initiate some in memory data
func initInMemoryData() {
	userA := User{
		ID:    "0001",
		Email: "0001@mail.com",
		Name:  "Joseph",
	}

	userB := User{
		ID:    "0002",
		Email: "0002@mail.com",
		Name:  "Roberto",
	}
	//userA's certificates:
	certificates = append(certificates, Certificate{
		ID:        "1",
		Title:     "A certificate title",
		CreatedAt: "do something for dates",
		OwnerID:   userA.ID,
		Year:      2018,
		Note:      "",
		Transfer: &Transfer{
			To:     "0002@mail.com",
			Status: "",
		},
	})
	certificates = append(certificates, Certificate{
		ID:        "3",
		Title:     "A certificate title",
		CreatedAt: "do something for dates",
		OwnerID:   userA.ID,
		Year:      2019,
		Note:      "",
		Transfer: &Transfer{
			To:     "",
			Status: "",
		},
	})
	//userB's certificates
	certificates = append(certificates, Certificate{
		ID:        "2",
		Title:     "A certificate title",
		CreatedAt: "do something for dates",
		OwnerID:   userB.ID,
		Year:      2015,
		Note:      "",
		Transfer: &Transfer{
			To:     "",
			Status: "",
		},
	})
	certificates = append(certificates, Certificate{
		ID:        "7",
		Title:     "A certificate title",
		CreatedAt: "do something for dates",
		OwnerID:   userB.ID,
		Year:      2000,
		Note:      "",
		Transfer: &Transfer{
			To:     "",
			Status: "",
		},
	})
}

// Route handler GET certificates method
func getCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(certificates)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// Route handler GET certificates/{id} method
func getCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range certificates {
		if item.ID == params["id"] {
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			return
		}
	}
}

func getUsersCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usersCertificates []Certificate
	params := mux.Vars(r)
	for _, item := range certificates {
		if item.OwnerID == params["userId"] {
			usersCertificates = append(usersCertificates, item)
		}
	}
	err := json.NewEncoder(w).Encode(usersCertificates)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

// Route handler POST certificate method
func createCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var certificate Certificate
	if err := json.NewDecoder(r.Body).Decode(&certificate); err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//authHeader := r.Header().Set("")
	//fmt.Printf("Auth: %v", authHeader)
	certificate.ID = strconv.Itoa(rand.Intn(1000000)) //Mock ID - not safe
	certificates = append(certificates, certificate)
	if err := json.NewEncoder(w).Encode(certificate); err != nil {
		panic(err)
	}
}

// Route handler PATCH certificate/{id} method
func patchCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Todo (Farouk): PATCH method can update individual fields in certificate
}

// Route handler PUT certificate/{id} method
func putCertificate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range certificates {
		if item.ID == params["id"] {
			certificates = append(certificates[:index], certificates[index+1:]...)
			var certificate Certificate
			//TODO (Farouk): PUT method: implement error if missing fields
			err := json.NewDecoder(r.Body).Decode(&certificate)
			if err != nil {
				fmt.Fprintf(w, err.Error())
			}
			certificate.ID = params["id"]
			certificates = append(certificates, certificate)
			err = json.NewEncoder(w).Encode(certificate)
			if err != nil {
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
	err := json.NewEncoder(w).Encode(certificates)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ID = strconv.Itoa(rand.Intn(1000000)) //Mock ID - not safe
	users = append(users, user)
	//jwt?
	//json.NewEncoder(w).Encode("")
}

func createTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range certificates {
		if item.ID == params["id"] && item.OwnerID == params["userId"] {
			var certificate Certificate
			//certificates = append(certificates[:index], certificates[index+1:]...)
			err := json.NewDecoder(r.Body).Decode(&certificate)
			if err != nil {
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
			certificate.Transfer.To = ""
			certificate.Transfer.Status = ""
			certificates = append(certificates[index:], certificates[index+1:]...)
			certificates = append(certificates, certificate)
			err := json.NewEncoder(w).Encode(certificate)
			if err != nil {
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
	err = json.NewEncoder(w).Encode(validToken)
	if err != nil {
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

	initInMemoryData()
	data := "user:pw" //must be format user:pw to pass BasicAuth
	dataEncoded := base64.StdEncoding.EncodeToString([]byte(data))
	fmt.Println(dataEncoded) //will print base64 encoded of data. Header request format: Basic <b64EncodedString>

	//Router
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/certificate", createCertificate).Methods("POST")
	r.HandleFunc("/certificates", getCertificates).Methods("GET")
	r.HandleFunc("/certificate/{id}", getCertificate).Methods("GET")
	r.HandleFunc("/certificate/{id}", putCertificate).Methods("PUT")
	r.HandleFunc("/certificate/{id}", patchCertificate).Methods("PATCH")
	r.HandleFunc("/certificate/{id}", deleteCertificate).Methods("DELETE")
	r.HandleFunc("/users/{userId}/certificates", getUsersCertificates).Methods("GET")
	r.HandleFunc("/users/{userId}/certificates/{id}/transfers", createTransfer).Methods("POST")                                     // PUT OR PATCH?
	r.HandleFunc("/users/{userId}/certificates/{id}/transfers", basicAuth(acceptTransfer, "user", "pw", "my-realm")).Methods("PUT") //PUT or PATCH?
	r.HandleFunc("/api/signup", createUser).Methods("POST")
	r.HandleFunc("/get-token", getTokenHandler).Methods("GET")

	s := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	log.Fatal(s.ListenAndServe())
}
