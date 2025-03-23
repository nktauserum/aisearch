package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nktauserum/aisearch/internal/api/handlers"
	mw "github.com/nktauserum/aisearch/internal/api/middleware"
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

	mux := mux.NewRouter()

	// index.html ))
	mux.HandleFunc("/", mw.LogHandler(handlers.IndexHandler))

	// Базовый поисковый запрос
	mux.HandleFunc("/api/v1/search", mw.LogHandler(handlers.SearchHandler)).Methods("POST")

	// Список текущих сохранённых сессий поиска
	mux.HandleFunc("/internal/sessions", mw.LogHandler(handlers.SessionListHandler)).Methods("GET")

	// Уточняющие запросы в рамках текущей сессии поиска
	mux.HandleFunc("/api/v1/refine", mw.LogHandler(handlers.RefineSearchHandler)).Methods("POST")

	// Получение истории поиска для конкретной сессии
	//mux.HandleFunc("/api/v1/search/{sessionId}/history", mw.LogHandler(handlers.SearchHistoryHandler)).Methods("GET")

	// Получение контекста текущей сессии поиска
	//mux.HandleFunc("/api/v1/search/{sessionId}", mw.LogHandler(handlers.SearchSessionHandler)).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%d", a.port), mux)
}
