package models_test

import (
	"context"
	"testing"

	"github.com/nktauserum/aisearch/pkg/ai/client"
	"github.com/nktauserum/aisearch/pkg/ai/models"
)

func TestStream(t *testing.T) {
	ctx := context.Background()
	defer ctx.Done()
	conversation := client.NewConversation("подробно отвечай на русском")

	result := make(chan string)
	user_message := models.Message{Text: "Расскажи про народовластие в Новгороде. Чем режим в то время отличался от либеральной демократии?"}
	err := conversation.Stream(ctx, user_message, result)
	if err != nil {
		t.Fatal(err)
	}

	for msg := range result {
		t.Log(msg)
	}
}
