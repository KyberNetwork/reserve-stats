package permission

import (
	"errors"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader  = "Authorization"
	signatureHeader      = "Signature"
	keyIDHeaderSeperator = "keyId=\""
)

var (
	//ErrCouldNotGetKeyID error when could not get key id in header
	ErrCouldNotGetKeyID = errors.New("Could not get key id in header")
	//ErrNotEnoughPermission error when keyid do not have enough permission
	ErrNotEnoughPermission = errors.New("KeyID do not have permission")
)

//KeyID is the abstract key needed for authentication
type KeyID string

// Permissioner is struct for control permission
type Permissioner struct {
	enforcer *casbin.Enforcer
}

// NewPermissioner return a gin HandleFunc middleware
func NewPermissioner(e *casbin.Enforcer) gin.HandlerFunc {
	p := &Permissioner{enforcer: e}
	return func(c *gin.Context) {
		if !p.CheckPermission(c.Request) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

// CheckPermission return a gin middleware which check if a request is authorize to continue or not
func (p *Permissioner) CheckPermission(r *http.Request) bool {
	keyID, err := getKeyID(r)
	if err != nil {
		return false
	}
	method := r.Method
	path := r.URL.Path
	return p.enforcer.Enforce(string(keyID), path, method)
}

func extractKeyID(s []string) (KeyID, error) {
	s1 := strings.Split(s[0], keyIDHeaderSeperator)
	if len(s1) < 2 {
		return KeyID(""), ErrCouldNotGetKeyID
	}
	keyIDStr := strings.Split(s1[1], "\"")[0]
	return KeyID(keyIDStr), nil
}

func getKeyID(r *http.Request) (KeyID, error) {
	if s, ok := r.Header[authorizationHeader]; ok {
		return extractKeyID(s)
	}
	if s, ok := r.Header[signatureHeader]; ok {
		return extractKeyID(s)
	}
	return "", ErrCouldNotGetKeyID
}
