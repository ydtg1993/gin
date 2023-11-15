package router

import (
	"github.com/gin-gonic/gin"
	"xo/controller"
)

func ApiRouter(g *gin.Engine) *gin.Engine {
	g.Use(gin.Recovery())
	api := g.Group("/api")
	{
		api.GET("/home", controller.Index)
	}

	return g
}
