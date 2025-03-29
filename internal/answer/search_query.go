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
	"github.com/nktauserum/aisearch/prompt"
)

type SearchInfo struct {
	Topic   string
	Queries []string
}

func GetSearchInfo(ctx context.Context, question string) (*SearchInfo, error) {
	if question == "" {
		return nil, fmt.Errorf("question equals nil")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conversation := client.NewConversation(prompt.SearchQuery())
	result, err := conversation.Continue(ctx, models.Message{Text: question})
	if err != nil {
		log.Printf("error getting search queries: %v", err)
		return nil, err
	}

	contentTemp := strings.Split(result, ".")
	if len(contentTemp) < 2 {
		log.Printf("invalid search_info format: %s\n", result)
		return nil, fmt.Errorf("invalid search_info format")
	}

	topic := contentTemp[0]
	log.Printf("topic: %s\n", topic)
	queries := strings.Split(contentTemp[1], ";")
	log.Printf("queries: %s\n", queries)

	content := SearchInfo{Topic: topic, Queries: queries}

	return &content, nil
}

func GenerateRefineQueries(ctx context.Context, old_conversation *models.Conversation, query string) ([]string, error) {
	conversation := client.NewConversation(prompt.RefineQuery())
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

	if strings.Contains(result, "not_needed") {
		return nil, nil
	}

	queries := strings.Split(result, ";")
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
