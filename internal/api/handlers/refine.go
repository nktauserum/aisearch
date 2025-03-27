package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nktauserum/aisearch/internal/answer"
	mw "github.com/nktauserum/aisearch/internal/api/middleware"
	"github.com/nktauserum/aisearch/pkg/ai/client"
	"github.com/nktauserum/aisearch/shared"
)

func RefineSearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Объявляем контекст
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Обрабатываем запрос пользователя
	request := new(shared.RefineRequest)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, request)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	log.Printf("request: %s\n", request.Query)

	// Получаем из памяти сессию поиска
	memory := client.GetMemory()
	conversation, err := memory.GetConversation(ctx, request.UUID)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	// Я думаю, здесь уместно было бы генерировать поисковые запросы, основываясь на предыдущем поиске,
	// в т.ч. на результатах.
	queries, err := answer.GenerateRefineQueries(ctx, conversation, request.Query)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Запрос: %s\n\n", request.Query))
	if queries != nil {
		content, err := answer.Search(ctx, queries...)
		if err != nil {
			mw.ErrorHandler(w, err, http.StatusInternalServerError)
			return
		}
		log.Println("analyzing is done")

		for _, site := range content {
			builder.WriteString(fmt.Sprintf("Title: %s\nWeb-ресурс: %s\nТекст: %s\n\n", site.Title, site.Sitename, site.Content))
		}
	}

	answer, err := answer.Research(ctx, conversation, builder.String())
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	fmt.Println(answer)

	//Переводим в JSON и возвращаем ответ
	response := shared.RefineResponse{Response: answer, Session: shared.SearchSession{UUID: request.UUID, Topic: conversation.Session.Topic}}
	responseBytes, err := json.Marshal(&response)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	// Успешно!
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

	err = memory.SaveConversation(ctx, request.UUID, conversation)
	if err != nil {
		log.Printf("error saving conversation: %v", err)
	}
}
