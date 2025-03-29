package mw

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Error string `json:"error"`
}

func ErrorHandler(c *gin.Context, err error, statusCode int) {
	// Обрабатываем ошибку
	log.Printf("%d error: %v\n", statusCode, err)
	c.JSON(statusCode, Error{Error: err.Error()})
}

func ErrorStreamHandler(c *gin.Context, err error, statusCode int) {
	// Обрабатываем ошибку
	log.Printf("%d error: %v\n", statusCode, err)
	c.Status(statusCode)
	c.SSEvent("error", err.Error())
	if c.Writer.Status() != http.StatusOK {
		return
	}
}
