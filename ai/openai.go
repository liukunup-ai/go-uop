package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIProvider struct {
	config     Config
	client     openai.Client
	embedModel string
}

func NewOpenAI(config Config) (*OpenAIProvider, error) {
	if config.Model == "" {
		config.Model = "gpt-4o"
	}
	clientOpts := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
	}
	if config.BaseURL != "" {
		clientOpts = append(clientOpts, option.WithBaseURL(config.BaseURL))
	}
	client := openai.NewClient(clientOpts...)
	return &OpenAIProvider{
		config:     config,
		client:     client,
		embedModel: openai.EmbeddingModelTextEmbedding3Small,
	}, nil
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) Assert(ctx context.Context, actual string, expected string) (bool, float64, error) {
	systemPrompt := `You are a precise text comparison judge.
Given an "actual" text and an "expected" description, determine if the actual text satisfies the expected description.
Output a JSON object with the following fields:
- "match": boolean (true if actual satisfies expected, false otherwise)
- "confidence": number between 0 and 1 (how certain you are)

Example output:
{"match": true, "confidence": 0.95}
{"match": false, "confidence": 0.85}

Rules:
- Set match to true only if the actual text clearly satisfies the expected description
- Set match to false if the actual text contradicts or fails to satisfy the expected description
- Confidence reflects how certain you are (1.0 = completely certain)`

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(fmt.Sprintf("Expected: %s\n\nActual: %s", expected, actual)),
		},
		Model: openai.ChatModel(p.config.Model),
	}
	if p.config.Temperature > 0 {
		params.Temperature = openai.Float(p.config.Temperature)
	}
	if p.config.TopP > 0 {
		params.TopP = openai.Float(p.config.TopP)
	}
	if p.config.MaxTokens > 0 {
		params.MaxTokens = openai.Int(int64(p.config.MaxTokens))
	}

	resp, err := p.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return false, 0, fmt.Errorf("assert: %w", err)
	}

	if len(resp.Choices) == 0 {
		return false, 0, fmt.Errorf("assert: no choices in response")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	return parseAssertResponse(content)
}

func (p *OpenAIProvider) Rerank(ctx context.Context, texts []string, expected string) ([]RankedText, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	reqEmbed, err := p.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(expected),
		},
		Model: p.embedModel,
	})
	if err != nil {
		return nil, fmt.Errorf("rerank: embedding expected: %w", err)
	}

	textList := make([]string, len(texts))
	copy(textList, texts)

	docEmbeds, err := p.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: textList,
		},
		Model: p.embedModel,
	})
	if err != nil {
		return nil, fmt.Errorf("rerank: embedding documents: %w", err)
	}

	results := make([]RankedText, len(texts))
	for i, text := range texts {
		score := cosineSimilarity(reqEmbed.Data[0].Embedding, docEmbeds.Data[i].Embedding)
		results[i] = RankedText{
			Text:  text,
			Score: score,
		}
	}

	sortByScoreDescending(results)

	return results, nil
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func parseAssertResponse(content string) (bool, float64, error) {
	content = strings.TrimSpace(content)

	var result struct {
		Match      bool    `json:"match"`
		Confidence float64 `json:"confidence"`
	}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return false, 0, fmt.Errorf("parse assert response: %w", err)
	}

	if result.Confidence < 0 {
		result.Confidence = 0
	} else if result.Confidence > 1 {
		result.Confidence = 1
	}

	return result.Match, result.Confidence, nil
}

func sortByScoreDescending(results []RankedText) {
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}
