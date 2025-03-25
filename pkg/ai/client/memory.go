package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/nktauserum/aisearch/pkg/ai/models"
	"github.com/nktauserum/aisearch/shared"
)

var (
	memory Memory
	once   sync.Once
)

type Memory struct {
	memory map[string]*models.Conversation
	m      sync.Mutex
}

func GetMemory() *Memory {
	once.Do(func() {
		memory = Memory{memory: make(map[string]*models.Conversation)}
	})

	return &memory
}

// Сохраняет в памяти диалог. Возвращает его ID
func (m *Memory) NewConversation(ctx context.Context, conv *models.Conversation) (string, error) {
	m.m.Lock()
	defer m.m.Unlock()

	id := uuid.New().String()
	if m.memory[id] != nil {
		return "", fmt.Errorf("id %s already exists", id)
	}
	m.memory[id] = conv

	return id, nil
}

func (m *Memory) SaveConversation(ctx context.Context, UUID string, conv *models.Conversation) error {
	m.m.Lock()
	defer m.m.Unlock()

	if m.memory[UUID] == nil {
		return fmt.Errorf("conversation with that UUID is not exists. Use memory.NewConversation() instead")
	}

	m.memory[UUID] = conv

	return nil
}

func (m *Memory) GetConversation(ctx context.Context, UUID string) (*models.Conversation, error) {
	m.m.Lock()
	defer m.m.Unlock()

	if m.memory[UUID] == nil {
		return nil, fmt.Errorf("conversation with UUID %s isn't exists", UUID)
	}

	return m.memory[UUID], nil
}

func (m *Memory) GetConversationList(ctx context.Context) []shared.SearchSession {
	var result []shared.SearchSession
	for i, conv := range m.memory {
		result = append(result, shared.SearchSession{UUID: i, Topic: conv.Session.Topic})
	}

	return result
}
