package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xo/controller"
)

func WebRouter(g *gin.Engine) *gin.Engine {
	g.LoadHTMLGlob("resources/templates/*")
	g.Static("/static", "resources/static")
	g.Static("/img", "resources/img")
	route := g.Group("/")
	{
		route.GET("/main.html", controller.Home)
		route.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/main.html")
		})
		route.GET("/main.html:string|/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/main.html")
		})
		route.GET("/video/:id/detail.html", controller.Video)
	}

	return g
}
