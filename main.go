package main

import (
	"github.com/gin-gonic/gin"
	"time"
	_ "xo/core"
	"xo/middleware"
	"xo/router"
)

func main() {
	gin.SetMode("debug")
	g := gin.New()
	g.Use(gin.Logger(),
		gin.Recovery(),
		middleware.ConcurrencyLimiterMiddleware(1024),
		middleware.RequestTimeoutMiddleware(30*time.Second),
		middleware.RequestDataSizeMiddleware(1024),
	)
	g = router.ApiRouter(g)
	g = router.WebRouter(g)
	g.Run(":8080")
}
