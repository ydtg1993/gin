package router

import (
	"github.com/gin-gonic/gin"
	"time"
	"xo/controller"
	"xo/middleware"
)

func WebRouter(g *gin.Engine) *gin.Engine {
	g.LoadHTMLGlob("templates/*")
	route := g.Group("/web")
	{
		route.Use(
			gin.Logger(),
			gin.Recovery(),
			middleware.ConcurrencyLimiterMiddleware(1024),
			middleware.RequestTimeoutMiddleware(30*time.Second),
			middleware.RequestDataSizeMiddleware(1024),
			middleware.CacheHTMLMiddleware())
		route.GET("/home", controller.Web)
	}

	return g
}
