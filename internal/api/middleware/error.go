package mw

import (
	"log"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, err error, statusCode int) {
	// Обрабатываем ошибку
	log.Printf("%d error: %v\n", statusCode, err)
	w.WriteHeader(statusCode)
}
