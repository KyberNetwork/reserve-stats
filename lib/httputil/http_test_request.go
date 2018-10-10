package httputil

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type assertFn func(t *testing.T, resp *httptest.ResponseRecorder)

// HTTPTestCase struct for http test case
type HTTPTestCase struct {
	Msg      string
	Endpoint string
	Method   string
	Body     []byte
	Assert   assertFn
}

// RunHTTPTestCase run http request test case
func RunHTTPTestCase(t *testing.T, tc HTTPTestCase, handler http.Handler) {
	t.Helper()
	req, err := http.NewRequest(tc.Method, tc.Endpoint, bytes.NewBuffer(tc.Body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	tc.Assert(t, resp)
}
