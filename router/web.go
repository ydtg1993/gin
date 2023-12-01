package router

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
	"xo/controller"
	"xo/middleware"
)

func WebRouter(g *gin.Engine) *gin.Engine {
	g.LoadHTMLGlob("resources/templates/*")
	g.Static("/static", "resources/static")
	g.Static("/img", "resources/img")
	route := g.Group("/")
	{
		route.GET("/main.html", middleware.CacheHTMLMiddleware(8*time.Hour, "page"), controller.Home)
		route.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/main.html")
		})
		route.GET("/video/:id/detail.html", middleware.CacheHTMLMiddleware(8*time.Hour), controller.Video)
		route.GET("/tag/:id/index.html", middleware.CacheHTMLMiddleware(8*time.Hour, "page"), controller.Tag)
		route.GET("/search", controller.Search)

		route.GET("/robots.txt", func(c *gin.Context) {
			robotsTxtContent := `
User-agent: *
Disallow:

Sitemap: https://www.apebt.com/sitemap.xml
		`
			c.String(200, robotsTxtContent)
		})
		route.GET("/:file", func(c *gin.Context) {
			file := c.Param("file")
			fileContent, err := ioutil.ReadFile("resources/" + file)
			if err != nil {
				return
			}
			extension := filepath.Ext(file)

			switch extension {
			case ".xml":
				c.Data(http.StatusOK, "application/xml", fileContent)
			default:
				c.Data(http.StatusOK, "text/plain; charset=utf-8", fileContent)
			}
		})
	}

	return g
}
