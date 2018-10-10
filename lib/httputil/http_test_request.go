package httputil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type assertFn func(t *testing.T, resp *httptest.ResponseRecorder)

// HTTPTestCase struct for http test case
type HTTPTestCase struct {
	Msg      string
	Endpoint string
	Method   string
	Data     map[string]interface{}
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
		reqBody, _ := json.Marshal(tc.Data)
		if err != nil {
			t.Fatal(err)
		}
		req.Body = ioutil.NopCloser(strings.NewReader(string(reqBody)))
		req.Header.Add("Content-Type", "application/json")
	}

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	tc.Assert(t, resp)
}
