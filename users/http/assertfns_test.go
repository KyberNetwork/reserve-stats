package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type userResponse struct {
	KYC  bool `json:"kyced"`
	Rich bool `json:"rich"`
}

// expectSuccess asserts that given response is a success response.
func expectSuccess(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
}

func expectKYCStatus(t *testing.T, resp *httptest.ResponseRecorder, kyc bool) {
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

func expectRichStatus(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
	decoded := &userResponse{}
	if err := json.NewDecoder(resp.Body).Decode(decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Rich {
		t.Errorf("wrong kyc statuc, expected %t, got: %t", false, decoded.Rich)
	}
}

//expectKYCed assert that givien response is a kyced user
func expectKYCed(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	expectKYCStatus(t, resp, true)
}

//expectNonKYCed assert that given response is a non kyced user
func expectNonKYCed(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	expectKYCStatus(t, resp, false)
}

//expectBadRequest assert that given response is bad request
func expectBadRequest(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("wrong return code, expected: %d, got %d", http.StatusBadRequest, resp.Code)
	}
}

//expectInternalServerError assert that given response is error from server
func expectInternalServerError(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	log.Printf("response: %+v", resp)
	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("wrong return code, expected: %d, got %d", http.StatusInternalServerError, resp.Code)
	}
}
