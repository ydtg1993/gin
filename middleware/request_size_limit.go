package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequestDataSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.String(http.StatusRequestEntityTooLarge, "Request data size exceeds the limit")
			c.Abort()
			return
		}

		c.Next()
	}
}
