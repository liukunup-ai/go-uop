package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ClaudeProvider struct {
	config Config
	client *http.Client
}

func NewClaude(config Config) (*ClaudeProvider, error) {
	if config.Model == "" {
		config.Model = "claude-3-5-sonnet-20241022"
	}
	return &ClaudeProvider{
		config: config,
		client: &http.Client{},
	}, nil
}

func (p *ClaudeProvider) Name() string {
	return "claude"
}

func (p *ClaudeProvider) Chat(ctx context.Context, messages []Message) (string, error) {
	reqBody := map[string]interface{}{
		"model":      p.config.Model,
		"messages":   messages,
		"max_tokens": 1024,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", p.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return result.Content[0].Text, nil
}

func (p *ClaudeProvider) ChatWithImage(ctx context.Context, messages []Message, imageData []byte) (string, error) {
	return p.Chat(ctx, messages)
}
