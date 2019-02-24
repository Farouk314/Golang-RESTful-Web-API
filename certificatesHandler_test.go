package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCertificates(t *testing.T) {
	req, err := http.NewRequest("GET", "localhost:8000/api/certificates", nil)
	if err != nil {
		t.Fatalf("Could not create get request: %v", err)
	}
	rec := httptest.NewRecorder()

	GetCertificates(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK;  got %v", res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
}
