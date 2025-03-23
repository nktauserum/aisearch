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

var research_prompt = `Вы интеллектуальный сервис-помощник. Вам предоставлен запрос и контекст к нему. После тщательного анализа контекста, необходимо дать ответ на поставленный запрос, используя информацию оттуда.
`

// Handler, обрабатывающий запросы на поиск.
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Объявляем контекст
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Обрабатываем запрос пользователя
	request := new(shared.SearchRequest)
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

	search_info, err := answer.GetSearchInfo(request.Query)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	content_temp := strings.Split(search_info, ".")

	// Это у нас будет тема запроса
	topic := content_temp[0]
	log.Printf("topic: %s\n", topic)
	// А это поисковые запросы для нас
	queries := strings.Split(content_temp[1], ";")
	log.Printf("queries: %s\n", queries)

	// сайты, прочёсанные нашим сервисом
	content, err := answer.Search(ctx, queries...)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	log.Println("analyzing is done")

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Запрос: %s\n\n", request.Query))
	for _, site := range content {
		builder.WriteString(fmt.Sprintf("Title: %s\nWeb-ресурс: %s\nТекст: %s\n\n", site.Title, site.Sitename, site.Content))
	}

	// Делаем запрос к нейросети
	conversation := client.NewConversation(research_prompt)
	answer, err := answer.Research(ctx, conversation, builder.String())
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	//log.Printf("response: %s", answer)

	// Сохраняем результат запроса, получая uuid
	memory := client.GetMemory()
	uuid, err := memory.NewConversation(ctx, conversation)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}
	conversation.Session = shared.SearchSession{UUID: uuid, Topic: topic}

	// Сохраняем в структуру ответ и источники
	researchResponse := new(shared.Research)
	researchResponse.Answer = answer
	for _, site := range content {
		researchResponse.Sources = append(researchResponse.Sources, site.URL)
	}

	// Переводим в JSON и возвращаем ответ
	response := shared.SearchResponse{Response: *researchResponse, Session: conversation.Session}
	responseBytes, err := json.Marshal(&response)
	if err != nil {
		mw.ErrorHandler(w, err, http.StatusInternalServerError)
		return
	}

	// Успешно!
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
