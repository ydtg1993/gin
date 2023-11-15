package router

import (
	"github.com/gin-gonic/gin"
	"xo/controller"
)

func WebRouter(g *gin.Engine) *gin.Engine {
	g.LoadHTMLGlob("templates/*")
	g.Use(gin.Recovery())
	api := g.Group("/web")
	{
		api.GET("/home", controller.Web)
	}

	return g
}
