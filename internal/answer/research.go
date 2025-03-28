package answer

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nktauserum/aisearch/pkg/ai/models"
	"github.com/nktauserum/aisearch/pkg/parser"
	"github.com/nktauserum/aisearch/shared"
)

func ExtractInfo(ctx context.Context, queries ...string) ([]shared.Website, error) {
	// Execute search queries
	log.Println("ищем в интернете")
	urls, err := DoSearchQueries(queries)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Create a timeout for the entire search
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	ch := make(chan shared.Website, len(urls))

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(ctx context.Context, url string) {
			defer wg.Done()
			defer cancel()

			select {
			case <-ctx.Done():
				return
			default:
			}
			// start := time.Now()
			log.Printf("сайт %s\n", url)

			siteInfo, err := parser.GetContent(ctx, url)
			if err != nil {
				fmt.Println(err)
				ch <- shared.Website{}
				//log.Printf("done %s in %.3f\n", url, time.Since(start).Seconds())
				return
			}

			ch <- siteInfo
			// log.Printf("done %s in %.3f\n", url, time.Since(start).Seconds())
		}(ctxWithTimeout, url)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var content []shared.Website
	for siteInfo := range ch {
		if siteInfo.URL != "" {
			content = append(content, siteInfo)
		}
	}

	return content, nil
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
	go func(context string) {
		filename := fmt.Sprintf("content/summary_%d.md", time.Now().Unix())
		err = os.WriteFile(filename, []byte(context), 0644)
		if err != nil {
			log.Printf("error writing summary to file: %v", err)
		}
		log.Printf("summary written to file: %s", filename)
	}(content)

	return summary, nil
}
