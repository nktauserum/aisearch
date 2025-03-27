package mw

import (
	"log"

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
