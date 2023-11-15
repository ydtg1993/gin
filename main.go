package main

import (
	"github.com/gin-gonic/gin"
	_ "xo/core"
	"xo/router"
)

func main() {
	gin.SetMode("debug")
	g := gin.New()
	g = router.ApiRouter(g)
	g = router.WebRouter(g)
	g.Run(":8080")
}
