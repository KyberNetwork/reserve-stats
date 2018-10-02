package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type assertFn func(t *testing.T, resp *httptest.ResponseRecorder)

type testCase struct {
	msg      string
	endpoint string
	method   string
	data     map[string]string
	assert   assertFn
}

func testHTTPRequest(t *testing.T, tc testCase, handler http.Handler) {
	t.Helper()

	req, tErr := http.NewRequest(tc.method, tc.endpoint, nil)
	if tErr != nil {
		t.Fatal(tErr)
	}

	if len(tc.data) != 0 {
		form := url.Values{}
		// log.Printf("request data: %+v", tc.data)
		for key, value := range tc.data {
			form.Add(key, value)
		}
		req.Form = form
		req.PostForm = form
		req.URL.RawQuery = form.Encode()
		req.Body = ioutil.NopCloser(strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Close = true
	log.Printf("request: %+v", req)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	tc.assert(t, resp)
}
