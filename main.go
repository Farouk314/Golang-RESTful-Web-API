package main

import (
	"fmt"
	"log"
	"net/http"
)

// Main
func main() {
	fmt.Println("Running...")

	InitInMemoryData()

	s := &http.Server{
		Addr:    ":8000",
		Handler: InitCors(InitHandlers),
	}

	log.Fatal(s.ListenAndServe())
}
