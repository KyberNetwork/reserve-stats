package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func newAssertGetEquation(expectedData []byte) assertFn {
	return func(t *testing.T, resp *httptest.ResponseRecorder) {
		t.Helper()

		if resp.Code != http.StatusOK {
			t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
		}

		type responseBody struct {
			Success bool
		}

		decoded := responseBody{}
		if aErr := json.NewDecoder(resp.Body).Decode(&decoded); aErr != nil {
			t.Fatal(aErr)
		}

		if decoded.Success != true {
			t.Errorf("wrong success status, expected: %t, got: %t", true, decoded.Success)
		}

	}
}

func testHTTPRequest(t *testing.T, tc testCase, handler http.Handler) {
	t.Helper()

	req, tErr := http.NewRequest(tc.method, tc.endpoint, nil)
	if tErr != nil {
		t.Fatal(tErr)
	}

	if len(tc.data) != 0 {
		form := url.Values{}
		for key, value := range tc.data {
			form.Add(key, value)
		}
		req.PostForm = form
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	tc.assert(t, resp)
}
