package authorization

import (
	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	readID                 = KeyID("read")
	writeID                = KeyID("write")
	invalidSignatureHeader = "Signature signature=\"Hello world\""
	readSignatureHeader    = "Signature keyId=\"read\",signature=\"Hello world\""
	writeSignatureHeader   = "Signature keyId=\"write\",signature=\"Hello world\""
)

var (
	writePermission = Permission{writeID}
	readPermission  = Permission{readID, writeID}
	readAuthorizer  = &Authorizer{readPermission}
	writeAuthorizer = &Authorizer{writePermission}
)

func middlewareTest(a *Authorizer, r *http.Request) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = r
	a.CheckPermission(c)
	return c
}

func TestNoAuthHeader(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	c := middlewareTest(readAuthorizer, r)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrCouldNotGetKeyID)
}

func TestInvalidSignHeader(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set(authorizationHeader, invalidSignatureHeader)
	c := middlewareTest(readAuthorizer, r)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrCouldNotGetKeyID)
}

func TestNotAuthorize(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set(authorizationHeader, readSignatureHeader)
	c := middlewareTest(writeAuthorizer, r)
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
	assert.Equal(t, c.Errors[0].Err, ErrNotEnoughPermission)
}

func TestAuthorize(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set(authorizationHeader, readSignatureHeader)
	c := middlewareTest(readAuthorizer, r)
	assert.Equal(t, http.StatusOK, c.Writer.Status())
}
