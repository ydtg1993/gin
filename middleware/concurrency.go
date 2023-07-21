package middleware

import "github.com/gin-gonic/gin"

func ConcurrencyLimiterMiddleware(limit int) gin.HandlerFunc  {
	semaphore := make(chan struct{}, limit)

	return func(c *gin.Context) {
		semaphore <- struct{}{}
		defer func() {
			<-semaphore
		}()

		c.Next()
	}
}
