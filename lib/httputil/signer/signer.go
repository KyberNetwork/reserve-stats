package signer

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"go.uber.org/zap"
	"net/http"
)

// GenerateNonce returns nonce header required to use Core API,
// which is current timestamp in milliseconds.
func GenerateNonce() string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(now, 10)
}

// sign will sign a message for authentication at server
func sign(key, msg string) (string, error) {
	mac := hmac.New(sha512.New, []byte(key))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

//NewRequest create a http request with params and header accordingly
func NewRequest(sugar *zap.SugaredLogger, url, method, signingKey, endpoint string, params map[string]string) (*http.Request, error) {
	logger := sugar.With(
		"func", "lib/http/signer/signer.go/NewRequest",
		"method", method,
		"endpoint", endpoint,
	)
	fullurl := fmt.Sprintf("%s/%s",
		strings.TrimRight(url, "/"),
		strings.Trim(endpoint, "/"),
	)

	req, err := http.NewRequest(method, fullurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	_, ok := params["nonce"]
	if ok {
		signed, err := sign(signingKey, q.Encode())
		if err != nil {
			return nil, err
		}
		req.Header.Add("signed", signed)
		logger = logger.With("signed", signed)
	}

	logger.Debugw("HTTP request created", "host", req.URL.Host, "raw query", req.URL.RawQuery)
	return req, nil
}
