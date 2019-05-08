package httputil

import (
	"github.com/gin-gonic/gin"
)

//MiddlewareHandler handle middleware error
func MiddlewareHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		c.JSON(
			c.Writer.Status(),
			c.Errors,
		)
	}
}

//ResponseFailure sets response code and error to the given one in parameter.
func ResponseFailure(c *gin.Context, code int, err error) {
	_ = c.Error(err)
	c.JSON(
		code,
		gin.H{"error": err.Error()},
	)
}
