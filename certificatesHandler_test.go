package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var a App

func newRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("Could not create %s request to url: %v, error: %v", method, url, err)
	}
	return req
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder() //Like a response writer plate
	a.Initialise()
	a.Handler.ServeHTTP(rec, req)

	return rec
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

func TestHomeHandler(t *testing.T) {
	req := newRequest(t, "GET", "http://localhost:8000/Home?v=2", nil)
	rec := executeRequest(req)

	a.HomeHandler(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)
	res := rec.Result()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Could not read response: %v", err)
	}
	i, err := strconv.Atoi(string(b[0]))
	if err != nil {
		t.Fatalf("expected an integer, got %v", string(b))
	}
	if i != 4 {
		t.Fatalf("Expected 4, got %v", i)
	}

	//Clean up in memory data
	certificates = certificates[:0]
	users = users[:0]
}

func TestGetCertificate(t *testing.T) {
	req := newRequest(t, "GET", "http://localhost:8000/certificates/1", nil)
	rec := executeRequest(req)
	a.Handler.ServeHTTP(rec, req)
	a.GetCertificate(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)

	res := rec.Result()
	defer res.Body.Close()

	// Certificate taken from resp body
	var certificate Certificate
	if err := json.NewDecoder(res.Body).Decode(&certificate); err != nil {
		t.Fatalf("Could not decode resp body: %v", err)
	}

	//Certificate in memory id=1
	cim := certificates[0]
	cimb, err := json.Marshal(cim)
	if err != nil {
		t.Fatalf("Could not marshal into JSON: %v", err)
	}
	cb, err := json.Marshal(certificate)
	if err != nil {
		t.Fatalf("Could not marshal into JSON: %v", err)
	}

	if string(cimb) != string(cb) {
		t.Fatalf("Expected certificate %+v, got %+v", string(cimb), string(cb))
	}

	//Clean up in memory data
	certificates = certificates[:0]
	users = users[:0]
}

func TestGetUsersCertificates(t *testing.T) {
	req := newRequest(t, "GET", "http://localhost:8000/users/A/certificates", nil)
	req.SetBasicAuth("userAEmail", "userApw")
	rec := executeRequest(req)

	a.GetUsersCertificates(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)

	res := rec.Result()
	if res == nil {
		t.Fatalf("Response is nil")
	}
	defer res.Body.Close()

	//Certificates in memory for user A
	a.InitInMemoryData()
	var usersCertificates []Certificate
	usersCertificates = append(usersCertificates, certificates[0], certificates[1])

	// Certificates from response body
	var respCertificates []Certificate
	if err := json.NewDecoder(res.Body).Decode(&respCertificates); err != nil {
		t.Fatalf("Could not decode body: %v", err)
	}
	rcb, err := json.Marshal(respCertificates)
	if err != nil {
		t.Fatalf("Could not marshal into JSON: %v", err)
	}
	ucb, err := json.Marshal(usersCertificates)
	if err != nil {
		t.Fatalf("Could not marshal into JSON: %v", err)
	}
	assert.JSONEq(t, string(ucb), string(rcb))

	//Clean up in memory data
	certificates = certificates[:0]
	users = users[:0]
}
