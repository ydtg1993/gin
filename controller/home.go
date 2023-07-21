package controller

import (
	"github.com/gin-gonic/gin"
	"time"
)


func Index(c *gin.Context) {
	time.Sleep(10*time.Second)
	SendResponse(c, 0, "success", nil)
}
