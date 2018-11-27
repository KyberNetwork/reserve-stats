package http

import (
	"github.com/gin-contrib/httpsign"
	"github.com/gin-contrib/httpsign/crypto"
	"github.com/gin-contrib/httpsign/validator"
)

// NewAuthenticator create a httpsign.Authenticator to check the message signing is valid or not
func NewAuthenticator(readKeyID, readKeySecret, writeKeyID, writeKeySecret string) (*httpsign.Authenticator, error) {
	var (
		secrets      = make(httpsign.Secrets)
		mapKeySecret = map[string]string{
			readKeyID:  readKeySecret,
			writeKeyID: writeKeySecret,
		}
	)

	for key, secret := range mapKeySecret {
		signKeyID := httpsign.KeyID(key)
		secrets[signKeyID] = &httpsign.Secret{
			Key:       secret,
			Algorithm: &crypto.HmacSha512{},
		}
	}

	auth := httpsign.NewAuthenticator(
		secrets,
		httpsign.WithValidator(
			NewNonceValidator(),
			validator.NewDigestValidator(),
		),
		httpsign.WithRequiredHeaders(
			[]string{"(request-target)", "nonce", "digest"},
		),
	)
	return auth, nil
}
