package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RequestTimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		timer := time.NewTimer(timeout)
		done := make(chan bool)
		go func() {
			c.Next()
			done <- true
		}()

		select {
		case <-done:
			return
		case <-timer.C:
			c.String(http.StatusGatewayTimeout, "Request took too long")
			c.Abort()
		}
	}
}
