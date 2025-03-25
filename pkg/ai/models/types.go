package models

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"github.com/nktauserum/aisearch/shared"
)

type Storage interface {
	Continue(ctx context.Context, message string) (string, error)
	SaveConversation(ctx context.Context, conv *Conversation) (int64, error)
	GetConversation(ctx context.Context, ID int64) (*Conversation, error)
}

type Conversation struct {
	User    *User
	Session shared.SearchSession
	Client  *openai.Client
	Params  openai.ChatCompletionNewParams
}

func (c *Conversation) GetMessages() *[]openai.ChatCompletionMessageParamUnion {
	return &c.Params.Messages.Value
}

type Message struct {
	Text     string
	ImageURL *Image
}

type Image struct {
	URL string
}

func (c *Conversation) Append(messages ...openai.ChatCompletionMessageParamUnion) error {
	to := c.GetMessages()
	if messages == nil {
		return fmt.Errorf("messages equals nil")
	}

	*to = append(*to, messages...)

	return nil
}

func (c *Conversation) Continue(ctx context.Context, message Message) (string, error) {
	messages := c.GetMessages()
	if message.ImageURL != nil {
		*messages = append(*messages, openai.UserMessageParts(openai.ImagePart(message.ImageURL.URL)))
	} else {
		*messages = append(*messages, openai.UserMessageParts(openai.TextPart(message.Text)))
	}

	chatCompletion, err := c.Client.Chat.Completions.New(ctx, c.Params,
		option.WithMaxRetries(2),
	)
	if err != nil {
		return "", err
	}

	if len(chatCompletion.Choices) == 0 {
		return "", fmt.Errorf("chatCompletion.Choices == 0")
	}

	c.Params.Messages.Value = append(c.Params.Messages.Value, chatCompletion.Choices[0].Message)
	return chatCompletion.Choices[0].Message.Content, nil
}

type User struct {
	ID       int64
	Username string
	Email    string
	password string
	Userpic  string
}
