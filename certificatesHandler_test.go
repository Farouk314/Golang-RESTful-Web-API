package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
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

}

func TestGetCertificate(t *testing.T) {
	req := newRequest(t, "GET", "http://localhost:8000/certificates/1", nil)
	rec := executeRequest(req)

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
}

func TestGetUsersCertificates(t *testing.T) {
	req := newRequest(t, "GET", "http://localhost:8000/users/A/certificates", nil)
	// req.SetBasicAuth("userAEmail", "userApw")
	rec := executeRequest(req)

	a.GetUsersCertificates(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)

	res := rec.Result()
	if res == nil {
		t.Fatalf("Response is nil")
	}
	defer res.Body.Close()

	//Certificates in memory for user A
	var usersCertificates []Certificate
	usersCertificates = append(usersCertificates, certificates[0], certificates[1])
	//fmt.Printf("usersCertificates: %+v", usersCertificates)

	//Certificates from response body
	var rc []byte
	fmt.Println("Body:")
	_, err := res.Body.Read(rc)
	if err != nil {
		t.Fatalf("Could not read resp body: %v", err)
	}
	fmt.Println(rec.Body) //rec.Body
	fmt.Println("Attempting to decode body")
	if err := json.NewDecoder(res.Body).Decode(&rc); err != nil {
		t.Fatalf("Could not decode body: %v", err)
	}
}

// func newRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
// 	req, err := http.NewRequest(method, url, body)
// 	if err != nil {
// 		t.Fatalf("Could not create %s request to url: %v, error: %v", method, url, err)
// 	}
// 	return req
// }

// func executeRequest(req *http.Request) *httptest.ResponseRecorder {
// 	rec := httptest.NewRecorder()
// 	a.Initialise()
// 	a.Handler.ServeHTTP(rec, req)

// 	return rec
// }

// func checkResponseCode(t *testing.T, expected, actual int) {
// 	fmt.Println("Checking repsonse code...")
// 	if expected != actual {
// 		t.Errorf("Expected response code %d. Got %d", expected, actual)
// 	}
// }

// func TestGetCertificates(t *testing.T) {
// 	req := newRequest(t, "GET", "http://localhost:8000/api/certificates", nil)
// 	rec := executeRequest(req)

// 	a.getCertificate(rec, req)

// 	checkResponseCode(t, http.StatusOK, rec.Code)
// 	res := rec.Result()
// 	defer res.Body.Close()

// 	var certificates []Certificate
// 	json.NewDecoder(res.Body).Decode(&certificates)
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatalf("Expected to receive a valid response in the body")
// 	}

// 	json.Unmarshal(body, &certificates)
// 	t.Fatalf("%v", certificates)
// 	// jsonData, err := json.Marshal(certificates)
// 	// if err != nil {
// 	// 	t.Fatalf("Could not Marshal JSON")
// 	// }

// 	expected := "[{1 A certificate title do something for dates guid1 2019 First Certificate 0xc000004ea0} {2 Another certificate title do something for dates guid2 2015 Second Certificate 0xc000004ec0}]"
// 	assert.JSONEq(t,
// 		expected,
// 		string(body))

// 	// var m map[string]gruckful.Certificate
// 	// json.Unmarshal(rec.Body.Bytes(), &m)
// 	// t.Fatalf("m: %v", m)
// }
