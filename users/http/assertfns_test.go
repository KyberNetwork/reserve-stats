package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// expectSuccess asserts that given response is a success response.
func expectSuccess(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
}

//expectBadRequest assert that given response is bad request
func expectBadRequest(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("wrong return code, expected: %d, got %d", http.StatusBadRequest, resp.Code)
	}
}
