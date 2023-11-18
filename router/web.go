package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"xo/controller"
	"xo/middleware"
)

func WebRouter(g *gin.Engine) *gin.Engine {
	g.LoadHTMLGlob("resources/templates/*")
	g.Static("/static", "resources/static")
	route := g.Group("/web")
	{
		route.Use(
			gin.Logger(),
			gin.Recovery(),
			middleware.ConcurrencyLimiterMiddleware(1024),
			middleware.RequestTimeoutMiddleware(30*time.Second),
			middleware.RequestDataSizeMiddleware(1024))
		route.GET("/home", middleware.CacheHTMLMiddleware(1*time.Minute), controller.Home)
		route.GET("/home:string", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/web/home")
		})
	}

	return g
}
