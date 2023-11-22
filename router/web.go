package router

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
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
		route.GET("/main.html", controller.Home)
		route.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, "/main.html")
		})
		//route.GET("/video/:id/detail.html", middleware.CacheHTMLMiddleware(5*time.Minute) ,controller.Video)
		route.GET("/video/:id/detail.html", controller.Video)
		route.GET("/tag/:id/index.html", middleware.CacheHTMLMiddleware(5*time.Minute, "page"), controller.Tag)
		route.GET("/search", controller.Search)

		route.GET("/robots.txt", func(c *gin.Context) {
			robotsTxtContent := `
User-agent: *
Disallow:

Sitemap: https://www.apebt.com/sitemap.xml
		`
			c.String(200, robotsTxtContent)
		})
		route.GET("/sitemap.xml", func(c *gin.Context) {
			sitemapContent, err := ioutil.ReadFile("resources/site.xml")
			if err != nil {
				return
			}
			c.Header("Content-Type", "application/xml")
			c.Data(http.StatusOK, "application/xml", sitemapContent)
		})
	}

	return g
}
