package httputil

import (
	"fmt"
	"net/http"
	"strings"
)

//NewRequest create a httpRequest
//if there is no params, pass a nil map to this func
func NewRequest(method, endpoint, host string, params map[string]string) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s",
		strings.TrimRight(host, "/"),
		strings.Trim(endpoint, "/"),
	)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}
