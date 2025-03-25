package answer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/nktauserum/aisearch/internal/search"
	"github.com/nktauserum/aisearch/pkg/ai/client"
	"github.com/nktauserum/aisearch/pkg/ai/models"
)

var search_query_prompt = `Вы интеллектуальный помощник, часть сервиса по поиску в интернете.

В качестве ответа возвращайте три готовых поисковых запроса и развернутую тему (current topic). Рекомендуется убирать ненужные эпитеты. Отвечайте корректно и соблюдайте орфографические нормы. В начале следует тема запроса, далее поисковые запросы, отделенные точкой, как в примере. Разделяйте запросы одним символом ";". Не заканчивайте предложение точкой. Старайтесь охватить как можно больше информации о вопросе, составляя запросы.

Пример ответа на запрос "Посоветуй какие-нибудь книги о диком Западе": "Книги о Диком Западе.дикий запад книги;книги в жанре вестерн;литература о ковбоях".

Пример ответа на запрос "Расскажи об австрийской экономической школе": "Австрийская экономическая школа: история и принципы.австрийская экономическая школа;экономисты австрийской экономической школы;либертарианство принципы".

`

var refine_query_prompt = `Вы интеллектуальный помощник, часть сервиса по поиску в интернете.

В качестве ответа возвращайте три готовых поисковых запроса. Рекомендуется убирать ненужные эпитеты. Отвечайте корректно и соблюдайте орфографические нормы. Разделяйте запросы одним символом ";". Не заканчивайте предложение точкой. Старайтесь охватить как можно больше информации о вопросе, составляя запросы, но не затрагивайте лишние темы.

Пример ответа на запрос "Посоветуй какие-нибудь книги о диком Западе": "дикий запад книги;книги в жанре вестерн;литература о ковбоях".

Пример ответа на запрос "Расскажи об австрийской экономической школе": "австрийская экономическая школа;экономисты австрийской экономической школы;либертарианство принципы".`

func GetSearchInfo(question string) (string, error) {
	//resultChan := make(chan []string, 1)

	if question == "" {
		return "", fmt.Errorf("question equals nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conversation := client.NewConversation(search_query_prompt)
	result, err := conversation.Continue(ctx, models.Message{Text: question})
	if err != nil {
		log.Printf("error getting search queries: %v", err)
		return "", nil
	}

	return result, nil
}

func GenerateRefineQueries(ctx context.Context, old_conversation *models.Conversation, query string) ([]string, error) {
	conversation := client.NewConversation(refine_query_prompt)
	messages := old_conversation.GetMessages()

	// Проверяем, есть ли хотя бы два элемента, чтобы безопасно срезать массив.
	if messages == nil || len(*messages) <= 1 {
		return nil, fmt.Errorf("slice messages is too short")
	}
	conversation.Append((*messages)[1:]...)

	result, err := conversation.Continue(ctx, models.Message{Text: query})
	if err != nil {
		return nil, err
	}
	content_temp := strings.Split(result, ".")
	queries := strings.Split(content_temp[1], ";")
	log.Printf("refine queries: %s\n", queries)

	return queries, nil
}

func DoSearchQueries(queries []string) ([]string, error) {
	var wg sync.WaitGroup

	result := make(chan struct {
		res []string
		err error
	}, len(queries))

	for _, query := range queries {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()

			searchResult, err := search.SearchTavily(query)
			var urls []string

			// В начале необходимо проверить, не присутствует ли полученный
			// результат поиска в итоговом массиве
			for _, searchRes := range searchResult {
				urlExists := false
				select {
				case r := <-result:
					for _, existingURL := range r.res {
						if existingURL == searchRes.URL {
							urlExists = true
							break
						}
					}
					result <- r
				default:
				}

				// Если не присутствует, добавляем со спокойной совестью
				if !urlExists {
					urls = append(urls, searchRes.URL)
				}
			}

			result <- struct {
				res []string
				err error
			}{
				res: urls,
				err: err,
			}
		}(query)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	var allResults []string
	var firstErr error

	// Create a map to track unique URLs
	uniqueURLs := make(map[string]bool)

	for r := range result {
		if r.err != nil && firstErr == nil {
			firstErr = r.err
		}
		// Check each URL and only append if it's unique
		for _, url := range r.res {
			if !uniqueURLs[url] {
				uniqueURLs[url] = true
				allResults = append(allResults, url)
			}
		}
	}

	return allResults, firstErr
}
