package main

import "time"

// Certificate Struct
type Certificate struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"` // Will be a type date
	OwnerID   string    `json:"ownerId"`
	Year      int       `json:"year"`
	Note      string    `json:"note"`
	Transfer  *Transfer `json:"transfer"`
}
