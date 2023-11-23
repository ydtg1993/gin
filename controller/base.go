package controller

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Encrypt(plaintext string, key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return ""
	}

	blockSize := len(key)
	padding := blockSize - len(plaintext)%blockSize
	if padding == 0 {
		padding = blockSize
	}

	// 填充 padding 个 byte(padding) 到 plaintext
	plaintext += string(bytes.Repeat([]byte{byte(padding)}, padding))
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = rand.Read(iv); err != nil {
		return ""
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext)
}
