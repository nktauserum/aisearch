package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/nktauserum/aisearch/internal/api/handlers"
)

type Application struct {
	port int
}

func NewApplication(port int) *Application {
	return &Application{
		port: port,
	}
}

func (a *Application) Run() error {
	log.Println("Приложение запущено")

	r := gin.Default()

	// index.html ))
	r.GET("/", handlers.IndexHandler)

	r.POST("/api/v1/search", handlers.SearchHandler)

	// Список текущих сохранённых сессий поиска
	r.GET("/internal/sessions", handlers.SessionListHandler)

	// Уточняющие запросы в рамках текущей сессии поиска
	r.POST("/api/v1/refine", handlers.RefineSearchHandler)

	r.POST("/api/v1/stream", handlers.StreamHandler)
	r.POST("/api/v1/fish", handlers.FishHandler)

	// Получение истории поиска для конкретной сессии
	// r.GET("/api/v1/search/:sessionId/history", mw.GinLogMiddleware(), handlers.SearchHistoryHandler)

	// Получение контекста текущей сессии поиска
	// r.GET("/api/v1/search/:sessionId", mw.GinLogMiddleware(), handlers.SearchSessionHandler)

	return r.Run(fmt.Sprintf(":%d", a.port))
}
