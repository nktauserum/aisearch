package answer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nktauserum/aisearch/pkg/ai/models"
	"github.com/nktauserum/aisearch/pkg/parser"
	"github.com/nktauserum/aisearch/shared"
)

func ExtractInfo(ctx context.Context, queries ...string) (chan shared.Website, error) {
	// Execute search queries
	log.Println("ищем в интернете")
	urls, err := DoSearchQueries(queries)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Create a timeout for the entire search
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	ch := make(chan shared.Website, len(urls))
	done := make(chan struct{})

	var wg sync.WaitGroup
	log.Printf("Got %d urls", len(urls))
	for _, url := range urls {
		wg.Add(1)
		go func(ctx context.Context, url string) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			siteInfo, err := parser.GetContent(ctx, url)
			if err != nil {
				fmt.Println(err)
				return
			}

			select {
			case <-ctx.Done():
				return
			case ch <- siteInfo:
			}
			log.Printf("сайт %s\n", url)
		}(ctxWithTimeout, url)
	}

	go func() {
		wg.Wait()
		close(ch)
		done <- struct{}{}
	}()

	select {
	case <-ctxWithTimeout.Done():
		close(ch) // Закрываем канал, чтобы не было утечки ресурсов
		return ch, nil
	case <-done:
		return ch, nil
	}
}

// Даёт ответ на запрос по переданному контенту
func Research(ctx context.Context, conversation *models.Conversation, content string) (string, error) {
	// Get summary from AI
	summary, err := conversation.Continue(ctx, models.Message{Text: content})
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Write summary to file

	return summary, nil
}
