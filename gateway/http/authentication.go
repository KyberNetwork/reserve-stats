package http

import (
	"github.com/gin-contrib/httpsign"
	"github.com/gin-contrib/httpsign/crypto"
	"github.com/gin-contrib/httpsign/validator"
)

func newAuthenticator(readKeyID, readKeySecret, writeKeyID, writeKeySecret string) (*httpsign.Authenticator, error) {

	var (
		secrets      = make(httpsign.Secrets)
		mapKeySecret = make(map[string]string)
	)

	mapKeySecret[readKeyID] = readKeySecret
	mapKeySecret[writeKeyID] = writeKeySecret

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
