package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
)

// Create user route Handler
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
