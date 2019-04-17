package authenticator

import (
	"errors"

	"github.com/gin-contrib/httpsign"
	"github.com/gin-contrib/httpsign/crypto"
	"github.com/gin-contrib/httpsign/validator"
)

// KeyPair difine AccessKeyID and SecretAccessKey
type KeyPair struct {
	AccessKeyID     string
	SecretAccessKey string
}

// NewAuthenticator create a httpsign.Authenticator to check the message signing is valid or not
func NewAuthenticator(keyPairs ...KeyPair) (*httpsign.Authenticator, error) {
	var secrets = make(httpsign.Secrets)

	if len(keyPairs) == 0 {
		return nil, errors.New("keyPairs are required")
	}

	for _, keyPair := range keyPairs {
		signKeyID := httpsign.KeyID(keyPair.AccessKeyID)
		secrets[signKeyID] = &httpsign.Secret{
			Key:       keyPair.SecretAccessKey,
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
