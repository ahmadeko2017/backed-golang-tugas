package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodySizeMiddleware limits the size of the request body
func BodySizeMiddleware(limitBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limitBytes)
		c.Next()
	}
}
