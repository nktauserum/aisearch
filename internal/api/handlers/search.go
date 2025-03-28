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

	"github.com/gin-gonic/gin"
	"github.com/nktauserum/aisearch/internal/answer"
	mw "github.com/nktauserum/aisearch/internal/api/middleware"
	"github.com/nktauserum/aisearch/pkg/ai/client"
	"github.com/nktauserum/aisearch/prompt"
	"github.com/nktauserum/aisearch/shared"
)

// Handler, обрабатывающий запросы на поиск.
func SearchHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	log.Println("SearchHandler started")

	// Объявляем контекст
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

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

	parsemode := shared.NewFormatHTML()

	aiStart := time.Now()
	conversation := client.NewConversation(prompt.Research(parsemode))
	answer, err := answer.Research(ctx, conversation, builder.String())
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}
	log.Printf("AI research completed: %v ms", time.Since(aiStart).Milliseconds())

	memory := client.GetMemory()
	uuid, err := memory.NewConversation(ctx, conversation)
	if err != nil {
		mw.ErrorHandler(c, err, http.StatusInternalServerError)
		return
	}
	conversation.Session = shared.SearchSession{UUID: uuid, Topic: search_info.Topic}

	researchResponse := new(shared.Research)
	researchResponse.Answer = answer
	for _, site := range content {
		researchResponse.Sources = append(researchResponse.Sources, site.URL)
	}

	response := shared.SearchResponse{Response: *researchResponse, Session: conversation.Session}

	c.JSON(http.StatusOK, response)
	log.Println("SearchHandler completed")
}
