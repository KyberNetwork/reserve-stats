package permission

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	signatureHeader     = "Signature"
	keyIDHeader         = "keyId"
)

var (
	//ErrCouldNotGetKeyID error when could not get key id in header
	ErrCouldNotGetKeyID = errors.New("could not get key id in header")
	//kvRegex is the regex to find key-value in a string
	kvRegex = regexp.MustCompile(`(\w+)="([^"]*)"`)
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
		if !p.checkPermission(c.Request) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

// checkPermission return a gin middleware which check if a request is authorize to continue or not
func (p *Permissioner) checkPermission(r *http.Request) bool {
	keyID, err := getKeyID(r)
	if err != nil {
		return false
	}
	method := r.Method
	path := r.URL.Path
	return p.enforcer.Enforce(string(keyID), path, method)
}

func extractKeyID(s string) (KeyID, error) {
	for _, match := range kvRegex.FindAllStringSubmatch(s, -1) {
		if len(match) < 3 {
			return KeyID(""), errors.New("malformed header")
		}
		k := match[1]
		v := match[2]
		if k == keyIDHeader {
			return KeyID(v), nil
		}
	}
	return KeyID(""), ErrCouldNotGetKeyID
}

func getKeyID(r *http.Request) (KeyID, error) {
	if s := r.Header.Get(authorizationHeader); len(s) > 0 {
		return extractKeyID(s)
	}
	if s := r.Header.Get(signatureHeader); len(s) > 0 {
		return extractKeyID(s)
	}
	return "", ErrCouldNotGetKeyID
}
