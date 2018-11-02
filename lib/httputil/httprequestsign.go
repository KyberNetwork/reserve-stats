package httputil

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"io"
	"net/http"
	"strings"
)

const (
	requestTarget = "(request-target)"
	nonce         = "nonce"
	digest        = "digest"
)

var headers = []string{requestTarget, nonce, digest}

//Sign sign a http request with selected header and return a signed request
//Use with http client only
func Sign(r *http.Request, keyID, secret string) (*http.Request, error) {
	// Set digest to body
	digestBody, err := calculateDigest(r)
	if err != nil {
		return nil, err
	}
	r.Header.Set(digest, digestBody)
	// Set nonce
	currentNonce := timeutil.UnixMilliSecond()
	r.Header.Set(nonce, fmt.Sprintf("%d", currentNonce))
	// Create sign string
	var signBuffer bytes.Buffer
	for i, h := range headers {
		var value string
		switch h {
		case requestTarget:
			value = fmt.Sprintf("%s %s", strings.ToLower(r.Method), r.URL.RequestURI())
		default:
			value = r.Header.Get(h)
		}
		signString := fmt.Sprintf("%s: %s", h, value)
		signBuffer.WriteString(signString)
		if i < len(headers)-1 {
			signBuffer.WriteString("\n")
		}
	}
	signString := signBuffer.String()
	// Create signature header
	signature, err := sign(signString, secret)
	if err != nil {
		return nil, err
	}
	signatureHeader := constructHeader(headers, keyID, signature)
	r.Header.Set("Signature", signatureHeader)
	return r, nil
}

func calculateDigest(r *http.Request) (string, error) {
	if r.Body == nil || r.ContentLength == 0 {
		return "", nil
	}
	body, err := r.GetBody()
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, err = io.Copy(h, body)
	if err != nil {
		return "", err
	}
	digest := fmt.Sprintf("SHA-256=%s", base64.StdEncoding.EncodeToString(h.Sum(nil)))
	return digest, nil
}

func sign(msg, secret string) (string, error) {
	mac := hmac.New(sha512.New, []byte(secret))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature, nil
}

func constructHeader(headers []string, keyID, signature string) string {
	var signBuffer bytes.Buffer
	signBuffer.WriteString(fmt.Sprintf(`keyId="%s",`, keyID))
	signBuffer.WriteString(`algorithm="hmac-sha512",`)
	signBuffer.WriteString(fmt.Sprintf(`headers="%s",`, strings.Join(headers, " ")))
	signBuffer.WriteString(fmt.Sprintf(`signature="%s"`, signature))
	return signBuffer.String()
}
