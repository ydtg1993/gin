package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(c *gin.Context) {
	data := gin.H{
		"Title": "Hello, Gin!",
	}
	c.HTML(http.StatusOK, "home.html", data)
}
