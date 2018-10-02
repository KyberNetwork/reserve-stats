package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type responseBody struct {
	Success bool
}

type userResponse struct {
	Success bool `json:"success"`
	KYC     bool `json:"kyced"`
}

// expectStatus asserts that given response matches the expected status.
func expectStatus(t *testing.T, resp *httptest.ResponseRecorder, success bool) {
	t.Helper()

	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
	decoded := &responseBody{}
	if err := json.NewDecoder(resp.Body).Decode(decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Success != success {
		t.Errorf("wrong success status, expected: %t, got: %t", success, decoded.Success)
	}
}

// ExpectSuccess asserts that given response is a success response.
func ExpectSuccess(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	expectStatus(t, resp, true)
}

// ExpectFailure asserts that given response is a failure response.
func ExpectFailure(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	expectStatus(t, resp, false)
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
