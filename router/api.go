package router

import (
	"read/controller"

	"github.com/gin-gonic/gin"
)

func ApiRouter(g *gin.Engine) *gin.Engine{
	g.Use(gin.Recovery())
	api := g.Group("/api")
	{
		api.GET("/home", controller.Index)
	}

	return g
}
