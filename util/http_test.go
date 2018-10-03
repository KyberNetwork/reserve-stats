package util

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

// HTTPTestCase struct for http test case
type HTTPTestCase struct {
	Msg      string
	Endpoint string
	Method   string
	Data     map[string]string
	Assert   assertFn
}

// CallTestHTTPRequest run http request test case
func CallTestHTTPRequest(t *testing.T, tc HTTPTestCase, handler http.Handler) {
	t.Helper()

	req, tErr := http.NewRequest(tc.Method, tc.Endpoint, nil)
	if tErr != nil {
		t.Fatal(tErr)
	}

	if len(tc.Data) != 0 {
		form := url.Values{}
		// log.Printf("request data: %+v", tc.data)
		for key, value := range tc.Data {
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
	tc.Assert(t, resp)
}
