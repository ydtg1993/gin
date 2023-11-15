package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	SendResponse(c, 0, "success", nil)
}

func Web(c *gin.Context) {
	data := gin.H{
		"Title": "Hello, Gin!",
	}
	c.HTML(http.StatusOK, "example.html", data)
}
