package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func SendResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Message: message,
		Data: data,
	})
}
