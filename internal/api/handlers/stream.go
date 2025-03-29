package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nktauserum/aisearch/internal/answer"
	mw "github.com/nktauserum/aisearch/internal/api/middleware"
	"github.com/nktauserum/aisearch/pkg/ai/client"
	"github.com/nktauserum/aisearch/pkg/ai/models"
	"github.com/nktauserum/aisearch/prompt"
	"github.com/nktauserum/aisearch/shared"
)

func StreamHandler(c *gin.Context) {
	// Устанавливаем заголовки для потока
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		c.Status(http.StatusRequestTimeout)
		return
	default:
	}

	// Обрабатываем запрос пользователя
	request := new(shared.SearchRequest)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, request)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}
	log.Printf("request: %s\n", request.Query)

	searchInfoStart := time.Now()
	search_info, err := answer.GetSearchInfo(ctx, request.Query)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched search info: %v ms", time.Since(searchInfoStart).Milliseconds())

	searchStart := time.Now()
	content, err := answer.Search(ctx, *search_info)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}
	log.Printf("Search completed: %v ms", time.Since(searchStart).Milliseconds())

	if content == nil {
		c.Status(http.StatusNoContent)
		return
	}

	for _, query := range content {
		c.SSEvent("source", query.URL)
		if c.Writer.Status() != http.StatusOK {
			ctx.Done()
			return
		}
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# Запрос: %s\n# Topic:%s\n# Search queries: [", request.Query, search_info.Topic))
	for _, query := range search_info.Queries {
		builder.WriteString(fmt.Sprintf("%s; ", query))
	}
	builder.WriteString("]")
	builder.WriteString("\n")
	for _, site := range content {
		builder.WriteString(fmt.Sprintf("\n\n### Title: %s\n#### URL: %s\n#### Название ресурса:%s\n#### Текст: %s\n", site.Title, site.URL, site.Sitename, site.Content))
	}

	parsemode := shared.NewFormatXML()

	result := make(chan string)
	conversation := client.NewConversation(prompt.Research(parsemode))
	user_message := models.Message{Text: request.Query}
	err = conversation.Stream(ctx, user_message, result)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}

	DEBUG := true

	for msg := range result {
		if DEBUG {
			fmt.Println(msg)
		}

		// Отправляем данные в формате Server-Sent Events
		c.SSEvent("message", msg)
		// Проверяем, не закрыт ли контекст
		if c.Writer.Status() != http.StatusOK {
			ctx.Done()
			return
		}
	}

	memory := client.GetMemory()
	uuid, err := memory.NewConversation(ctx, conversation)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}

	session := shared.SearchSession{UUID: uuid, Topic: search_info.Topic}
	session_json, err := json.Marshal(session)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}

	c.SSEvent("info", string(session_json))

	filename := fmt.Sprintf("content/summary_%d.md", time.Now().Unix())
	err = os.WriteFile(filename, []byte(builder.String()), 0644)
	if err != nil {
		log.Printf("error writing summary to file: %v", err)
	}
	log.Printf("summary written to file: %s", filename)
}

func FishHandler(c *gin.Context) {

	// Устанавливаем заголовки для потока
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Создаем канал для передачи данных
	result := make(chan string)

	// Запускаем функцию Fish в отдельной горутине
	go func() {
		defer close(result)
		answer.Fish(result)
	}()

	// Читаем данные из канала и отправляем их клиенту
	for msg := range result {
		// Отправляем данные в формате Server-Sent Events
		c.SSEvent("message", msg)
		// Проверяем, не закрыт ли контекст
		if c.Writer.Status() != http.StatusOK {
			ctx.Done()
			return
		}
	}

	// Отправляем завершающее событие
	session := shared.SearchSession{UUID: "some uuid", Topic: "To be or not to be - that's a question"}
	session_json, err := json.Marshal(session)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}

	c.SSEvent("info", string(session_json))
}
