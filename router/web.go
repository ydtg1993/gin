package router

import (
	"fmt"
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
		route.GET("/sitemap.xml", func(c *gin.Context) {
			sitemapContent, err := ioutil.ReadFile("resources/sitemap.xml")
			if err != nil {
				return
			}
			fmt.Println("Sitemap Content:", string(sitemapContent))
			c.Header("Content-Type", "application/xml")
			c.Data(http.StatusOK, "application/xml", sitemapContent)
		})
		route.GET("/BingSiteAuth.xml", func(c *gin.Context) {
			siteContent, err := ioutil.ReadFile("resources/BingSiteAuth.xml")
			if err != nil {
				return
			}
			c.Header("Content-Type", "application/xml")
			c.Data(http.StatusOK, "application/xml", siteContent)
		})
	}

	return g
}
