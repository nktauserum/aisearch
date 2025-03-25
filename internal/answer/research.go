package answer

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nktauserum/aisearch/pkg/ai/models"
	"github.com/nktauserum/aisearch/pkg/parser"
	"github.com/nktauserum/aisearch/shared"
)

func Search(ctx context.Context, queries ...string) ([]shared.Website, error) {
	// Execute search queries
	log.Println("ищем в интернете")
	urls, err := DoSearchQueries(queries)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Собираем контент с сайтов асинхронно
	ch := make(chan shared.Website)
	for _, url := range urls {
		go func(ctx context.Context, url string) {
			//defer log.Printf("done %s\n", url)
			ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()
			log.Printf("сайт %s\n", url)
			siteInfo, err := parser.GetContent(ctx, url)
			if err != nil {
				fmt.Println(err)
				ch <- shared.Website{}
				return
			}
			ch <- siteInfo
		}(ctx, url)
	}

	var content []shared.Website
	for range urls {
		if siteInfo := <-ch; siteInfo.URL != "" {
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
		filename := fmt.Sprintf("content/summary_%d.txt", time.Now().Unix())
		err = os.WriteFile(filename, []byte(context), 0644)
		if err != nil {
			log.Printf("error writing summary to file: %v", err)
		}
		log.Printf("summary written to file: %s", filename)
	}(content)

	return summary, nil
}
