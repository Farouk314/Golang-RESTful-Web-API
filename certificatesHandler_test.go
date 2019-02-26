package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestCreateCertificate(t *testing.T) {
	var jsonStrCertificate = []byte(`{        
		"title": "Test Certificate",
		"createdAt": "2019-02-26T00:17:14.1805401Z",
		"note" : "Certificate sent in request body"
	}`)

	req := newRequest(t, "POST", "http://localhost:8000/certificates/1", bytes.NewBuffer(jsonStrCertificate))
	req.SetBasicAuth("userA", "")
	rec := executeRequest(req)

	a.CreateCertificate(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)

	createdAt, err := time.Parse(time.RFC3339, "2019-02-26T00:17:14.1805401Z")
	if err != nil {
		t.Fatalf("Could not parse time: %v", err)
	}
	expectedCertificate := Certificate{
		ID:        "1",
		Title:     "Test Certificate",
		CreatedAt: createdAt,
		OwnerID:   "A",
		Year:      2019,
		Note:      "Certificate sent in request body",
		Transfer:  &Transfer{To: "", Status: ""},
	}

	res := rec.Result()
	if res == nil {
		t.Fatalf("Response is nil")
	}

	//Response body certificate
	var resCertificate Certificate

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&resCertificate); err != nil {
		t.Fatalf("Could not decode body: %v", err)
	}

	rcb, err := json.Marshal(resCertificate)
	if err != nil {
		t.Fatalf("Could not Marshal to JSON: %v", err.Error())
	}

	ecb, err := json.Marshal(expectedCertificate)
	if err != nil {
		t.Fatalf("Could not Marshal to JSON: %v", err.Error())
	}

	assert.JSONEq(t, string(ecb), string(rcb))

	//Clean up in memory data
	certificates = certificates[:0]
	users = users[:0]
}

func TestGetUsersCertificates(t *testing.T) {
	req := newRequest(t, "GET", "http://localhost:8000/users/A/certificates", nil)
	req.SetBasicAuth("userA", "")
	rec := executeRequest(req)

	a.GetUsersCertificates(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)

	res := rec.Result()
	if res == nil {
		t.Fatalf("Response is nil")
	}

	//Certificates in memory for user A
	var usersCertificates []Certificate
	usersCertificates = append(usersCertificates, certificates[0], certificates[1])

	// Certificates from response body
	var respCertificates []Certificate

	defer res.Body.Close()
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

func TestUpdateCertificate(t *testing.T) {
	var jsonStrCertificate = []byte(`{        
		"title": "Updated Title",
		"note" : "Updated Note"
	}`)
	req := newRequest(t, "PATCH", "http://localhost:8000/certificates/1", bytes.NewBuffer(jsonStrCertificate))
	req.SetBasicAuth("userA", "")
	rec := executeRequest(req)

	a.UpdateCertificate(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)

	res := rec.Result()
	if res == nil {
		t.Fatalf("Response is nil")
	}

	//Update certificate 1
	certificates[0].Title = "Updated Title"
	certificates[0].Note = "Updated Note"

	//Response certificate
	var rUpdatedCertificate Certificate

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&rUpdatedCertificate); err != nil {
		t.Fatalf("Could not decode JSON: %v", err.Error())
	}

	rcb, err := json.Marshal(rUpdatedCertificate)
	if err != nil {
		t.Fatalf("Could not marshal into JSON: %v", err.Error())
	}

	ecb, err := json.Marshal(certificates[0])
	if err != nil {
		t.Fatalf("Could not marshal into JSON: %v", err.Error())
	}

	assert.JSONEq(t, string(ecb), string(rcb))

	//Clean up in memory data
	certificates = certificates[:0]
	users = users[:0]
}

func TestDeleteCertificate(t *testing.T) {
	req := newRequest(t, "DELETE", "http://localhost:8000/certificates/1", nil)
	req.SetBasicAuth("userA", "")
	rec := executeRequest(req)

	a.DeleteCertificate(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)

	res := rec.Result()
	if res == nil {
		t.Fatalf("Response is nil")
	}

	for _, item := range certificates {
		if item.ID == "1" {
			t.Fatalf("Certificate 1 was not deleted")
		}
	}

	//Clean up in memory data
	certificates = certificates[:0]
	users = users[:0]
}

func TestCreateTransfer(t *testing.T) {
	var jsonStrTransferTo = []byte(`{        
		"transfer": {
			"to": "userB"
		}
	}`)
	req := newRequest(t, "PATCH", "http://localhost:8000/certificates/1/transfers", bytes.NewBuffer(jsonStrTransferTo))
	req.SetBasicAuth("userA", "")
	rec := executeRequest(req)

	a.CreateTransfer(rec, req)
	checkResponseCode(t, http.StatusOK, rec.Code)

	res := rec.Result()
	if res == nil {
		t.Fatalf("Response was nil")
	}

	//Response certificate
	var rCertificate Certificate

	//Update the transfer to and status of certificate 1
	certificates[0].Transfer.To = "userB"
	certificates[0].Transfer.Status = "Pending transfer"

	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&rCertificate); err != nil {
		t.Fatalf("Could not decode JSON: %v", err.Error())
	}

	rcb, err := json.Marshal(rCertificate)
	if err != nil {
		t.Fatalf("Could not marshal into JSON: %v", err.Error())
	}

	ecb, err := json.Marshal(certificates[0])
	if err != nil {
		t.Fatalf("Could not marshal into JSON: %v", err.Error())
	}

	assert.JSONEq(t, string(ecb), string(rcb))

	//Clean up in memory data
	certificates = certificates[:0]
	users = users[:0]
}

// func TestAcceptTransfer(t *testing.T) {
// 	req := newRequest(t, "PUT", "http://localhost:8000/certificates/1/transfers", nil)
// 	req.SetBasicAuth("userB", "")
// 	rec := executeRequest(req)

// 	//Change the transfer to and status of certificate 1, which userB will be accepting
// 	certificates[0].Transfer.To = "userB"
// 	certificates[0].Transfer.Status = "Pending transfer"

// 	a.AcceptTransfer(rec, req)
// 	checkResponseCode(t, http.StatusOK, rec.Code)

// 	res := rec.Result()
// 	if res == nil {
// 		t.Fatalf("Response was nil")
// 	}

// 	//Response certificate
// 	var rCertificate Certificate

// 	defer res.Body.Close()
// 	if err := json.NewDecoder(res.Body).Decode(&rCertificate); err != nil {
// 		t.Fatalf("Could not decode HERE JSON: %v", err.Error())
// 	}

// 	if rCertificate.Transfer.To != "" {
// 		t.Fatalf("Did not reset transfer field")
// 	}
// 	if rCertificate.Transfer.Status != "" {
// 		t.Fatalf("Did not reset status field")
// 	}
// 	userName, err := LookUpUserIDByName("userB")
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
// 	if rCertificate.OwnerID != userName {
// 		t.Fatalf("Did not modify ownerID field, expected B, got %v", rCertificate.OwnerID)
// 	}

// 	//Clean up in memory data
// 	certificates = certificates[:0]
// 	users = users[:0]
// }
