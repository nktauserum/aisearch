package search

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/openai/openai-go"

	"github.com/nktauserum/aisearch/pkg/ai/client"
	"github.com/nktauserum/aisearch/pkg/ai/models"
)

func errorHandler(err error) error {
	var apiErr *openai.Error
	if errors.As(err, &apiErr) {
		// Структурированная обработка ошибок API
		return fmt.Errorf("openai api error: %w", apiErr)
	}
	return err
}

func Dialog(messages chan string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conversation := client.NewConversation()

	for message := range messages {
		answer, err := conversation.Continue(ctx, models.Message{Text: message})
		if err != nil {
			log.Printf("error adding message: %v", err)
			return
		}
		fmt.Println(answer)
	}
}
