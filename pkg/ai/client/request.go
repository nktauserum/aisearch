package client

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/nktauserum/aisearch/config"
	"github.com/nktauserum/aisearch/pkg/ai/models"
)

// возвращает указатель на новый объект Conversation.
// TODO: Добавить возможность создания объекта вместе с переданным пользователем
// ^ что это за херня?
func NewConversation(system_prompt ...string) *models.Conversation {
	config := config.GetConfig()

	client := openai.NewClient(
		option.WithBaseURL(config.OpenAI.URL),
		option.WithAPIKey(config.OpenAI.Key),
		option.WithMaxRetries(0),
	)

	msg := []openai.ChatCompletionMessageParamUnion{}

	if len(system_prompt) > 0 && system_prompt[0] != "" {
		msg = []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(system_prompt[0])}
	}

	params := openai.ChatCompletionNewParams{
		Messages:    openai.F(msg),
		Model:       openai.F(config.OpenAI.Model),
		Temperature: openai.Float(config.OpenAI.Temperature),
	}

	return &models.Conversation{
		Client: client,
		Params: params,
	}
}
