package router

import (
	"github.com/gin-gonic/gin"
	"time"
	"xo/middleware"
)

func ApiRouter(g *gin.Engine) *gin.Engine {
	route := g.Group("/api")
	{
		route.Use(
			gin.Logger(),
			gin.Recovery(),
			middleware.ConcurrencyLimiterMiddleware(1024),
			middleware.RequestTimeoutMiddleware(30*time.Second),
			middleware.RequestDataSizeMiddleware(1024))
	}

	return g
}
