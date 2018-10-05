package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type responseBody struct {
	Success bool
}

type userResponse struct {
	KYC bool `json:"kyced"`
}

// ExpectSuccess asserts that given response is a success response.
func ExpectSuccess(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	log.Printf("response: %+v", resp)
	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
}

// ExpectInternalServerError assert that given response return internal server error
func ExpectInternalServerError(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
}

func expectKYC(t *testing.T, resp *httptest.ResponseRecorder, kyc bool) {
	t.Helper()

	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
	decoded := &userResponse{}
	if err := json.NewDecoder(resp.Body).Decode(decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.KYC != kyc {
		t.Errorf("wrong kyc status, expected: %t, got: %t", kyc, decoded.KYC)
	}
}

//ExpectKYC assert that givien response is a kyced user
func ExpectKYC(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	expectKYC(t, resp, true)
}

//ExpectNonKYC assert that given response is a non kyced user
func ExpectNonKYC(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	expectKYC(t, resp, false)
}

//ExpectBadRequest assert that given response is bad request
func ExpectBadRequest(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("wrong return code, expected: %d, got %d", http.StatusBadRequest, resp.Code)
	}
}
