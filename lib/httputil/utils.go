package httputil

import (
	"github.com/gin-gonic/gin"
)

//MiddlewareHandler handle middleware error
func MiddlewareHandler(c *gin.Context) {
	c.Next()
	defer func(c *gin.Context) {
		if len(c.Errors) > 0 {
			c.JSON(
				c.Writer.Status(),
				c.Errors,
			)
		}
	}(c)
}
