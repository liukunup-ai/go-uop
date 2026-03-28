package ai

import (
	"context"
)

type Provider interface {
	Name() string
	// Assert judges whether actual text matches the description in expected text.
	// Returns: (isMatch, confidence [0-1], error)
	Assert(ctx context.Context, actual string, expected string) (bool, float64, error)
	// Rerank reorders texts based on how well each matches the given requirement.
	// Returns: ranked texts sorted by score descending.
	Rerank(ctx context.Context, texts []string, expected string) ([]RankedText, error)
}

// RankedText represents a text with its relevance score.
type RankedText struct {
	Text  string
	Score float64
}

type Config struct {
	APIKey      string
	BaseURL     string
	Model       string
	TopP        float64
	Temperature float64
	MaxTokens   int
}

func NewProvider(providerType string, config Config) (Provider, error) {
	switch providerType {
	case "openai":
		return NewOpenAI(config)
	case "bigmodel":
		return NewBigModel(config)
	default:
		return nil, nil
	}
}
