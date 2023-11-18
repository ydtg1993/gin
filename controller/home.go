package controller

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"net/http"
)

func Home(c *gin.Context) {
	staticContent, _ := ioutil.ReadFile("resources/templates/static.html")
	data := gin.H{
		"Title":  "Hello, Gin!",
		"header": template.HTML(staticContent),
	}
	c.HTML(http.StatusOK, "home.html", data)
}
