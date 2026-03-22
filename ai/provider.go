package ai

import (
	"context"
)

type Provider interface {
	Name() string
	Chat(ctx context.Context, messages []Message) (string, error)
	ChatWithImage(ctx context.Context, messages []Message, imageData []byte) (string, error)
}

type Message struct {
	Role    string
	Content string
}

type Config struct {
	APIKey string
	Model  string
}

func NewProvider(providerType string, config Config) (Provider, error) {
	switch providerType {
	case "openai":
		return NewOpenAI(config)
	default:
		return nil, nil
	}
}
