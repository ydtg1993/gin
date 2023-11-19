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

func Tag() {

}

func Search(c *gin.Context) {

}

func Video(c *gin.Context) {
	//id := c.Param("id")
	staticContent, _ := ioutil.ReadFile("resources/templates/static.html")
	data := gin.H{
		"Title":  "Hello, Gin!",
		"header": template.HTML(staticContent),
	}
	c.HTML(http.StatusOK, "video.html", data)
}
