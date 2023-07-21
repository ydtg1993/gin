package main

import (
	"github.com/gin-gonic/gin"
	"read/core"
	"read/middleware"
	"read/router"
	"time"
)

func init()  {
	core.Env()
}

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

	g.Run(":8080")
}
