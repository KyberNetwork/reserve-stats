package middleware

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	errInvalidNonce            = errors.New("invalid nonce")
	errInvalidSign             = errors.New("invalid signed token")
	errAuthenticatedPermission = errors.New("server error, unknown permission")
)

// Permission defines the permission type.
type Permission uint32

type requiredFields struct {
	Nonce int64 `form:"nonce" binding:"required"`
}

// HMACAuthenticator is the gin authenticator middleware
// that uses HMAC with SHA512.
type HMACAuthenticator struct {
	secrets map[Permission]string
	v       NonceValidator
}

// NewHMACAuthenticator creates a new HMACAuthenticator instance with
// given allowed permissions and secret keys.
func NewHMACAuthenticator(secretKeys map[Permission]string, v NonceValidator) *HMACAuthenticator {
	return &HMACAuthenticator{secrets: secretKeys, v: v}
}

// Authenticated returns a gin middleware which permits given permissions in parameter.
func (ha *HMACAuthenticator) Authenticated(permissions []Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r requiredFields

		for _, p := range permissions {
			if _, ok := ha.secrets[p]; !ok {
				c.AbortWithError(http.StatusInternalServerError, errAuthenticatedPermission)
				return
			}
		}

		if err := c.Request.ParseForm(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if err := c.ShouldBind(&r); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if !ha.v.IsValid(r.Nonce) {
			c.AbortWithError(http.StatusBadRequest, errInvalidNonce)
			return
		}

		signed := c.GetHeader("signed")
		message := c.Request.Form.Encode()
		auth := ha.checkAuth(signed, message, permissions)
		if !auth {
			c.AbortWithError(http.StatusUnauthorized, errInvalidSign)
			return
		}
		c.Next()
	}
}

func (ha *HMACAuthenticator) checkAuth(signed string, message string, permissions []Permission) bool {
	for _, p := range permissions {
		secret, ok := ha.secrets[p]
		if !ok {
			return false
		}
		serverSigned, err := serverSign(message, secret)
		if err != nil {
			return false
		}
		if serverSigned == signed {
			return true
		}
	}
	return false
}

func serverSign(msg string, secret string) (string, error) {
	mac := hmac.New(sha512.New, []byte(secret))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}
