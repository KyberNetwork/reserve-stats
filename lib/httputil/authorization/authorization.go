package authorization

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	keyIDHeader         = "keyId=\""
)

var (
	//ErrCouldNotGetKeyID error when could not get key id in header
	ErrCouldNotGetKeyID = errors.New("Could not get key id in header")

	//ErrNotEnoughPermission error when keyid do not have enough permission
	ErrNotEnoughPermission = errors.New("KeyID do not have permission")
)

// Authorizer is struct for control permission
type Authorizer struct {
	p Permission
}

// CheckPermission return a gin middleware which check if a request is authorize to continue or not
func (a *Authorizer) CheckPermission(c *gin.Context) {
	keyID, err := getKeyID(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if !a.p.isPermit(keyID) {
		c.AbortWithError(http.StatusUnauthorized, ErrNotEnoughPermission)
		return
	}
	c.Next()
}

func getKeyID(r *http.Request) (KeyID, error) {
	if s, ok := r.Header[authorizationHeader]; ok {
		s1 := strings.Split(s[0], keyIDHeader)
		if len(s1) < 2 {
			return KeyID(""), ErrCouldNotGetKeyID
		}
		keyID := strings.Split(s1[1], "\"")[0]
		return KeyID(keyID), nil
	}
	return KeyID(""), ErrCouldNotGetKeyID
}
