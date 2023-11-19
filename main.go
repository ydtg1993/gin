package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
	"xo/core"
	_ "xo/core"
	"xo/middleware"
	"xo/router"
)

func main() {
	gin.SetMode("debug")
	g := gin.New()
	g.NoRoute(func(c *gin.Context) {
		staticContent, _ := ioutil.ReadFile("resources/templates/static.html")
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"title":  core.Config.GetString("app.name"),
			"header": template.HTML(staticContent),
			"error":  "page cannot be found"})
	})
	g.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.ConcurrencyLimiterMiddleware(1024),
		middleware.RequestTimeoutMiddleware(30*time.Second),
		middleware.RequestDataSizeMiddleware(1024))
	g = router.ApiRouter(g)
	g = router.WebRouter(g)
	g.Run(":8081")
}
