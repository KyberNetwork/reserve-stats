package httputil

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type assertFn func(t *testing.T, resp *httptest.ResponseRecorder)

// HTTPTestCase struct for http test case
type HTTPTestCase struct {
	Msg      string
	Endpoint string
	Method   string
	Data     map[string]string
	Assert   assertFn
}

// RunHTTPTestCase run http request test case
func RunHTTPTestCase(t *testing.T, tc HTTPTestCase, handler http.Handler) {
	t.Helper()

	req, err := http.NewRequest(tc.Method, tc.Endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(tc.Data) != 0 {
		form := url.Values{}
		for key, value := range tc.Data {
			form.Add(key, value)
		}
		req.Form = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	tc.Assert(t, resp)
}
