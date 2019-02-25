package main

import (
	"errors"
	"time"
)

// Certificates
var certificates []Certificate

// Users
var users []User

// LookUpUserIDByEmail returns a the UsserID for a specified email address
func LookUpUserIDByEmail(userEmail string) (string, error) {
	for _, item := range users {
		if item.Email == userEmail {
			return item.ID, nil
		}
	}
	return "", errors.New("Could not get user by email: " + userEmail)
}

// InitInMemoryData : initiates some in memory data for testing
func (a *App) InitInMemoryData() {
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
	//userA's Certificates
	certificates = append(certificates,
		Certificate{
			ID:        "1",
			Title:     "A certificate title",
			CreatedAt: time.Now(),
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
			CreatedAt: time.Now(),
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
			CreatedAt: time.Now(),
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
			CreatedAt: time.Now(),
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
