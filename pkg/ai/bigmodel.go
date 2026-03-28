package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type BigModelProvider struct {
	config  Config
	baseURL string
	client  *http.Client
}

func NewBigModel(config Config) (*BigModelProvider, error) {
	if config.Model == "" {
		config.Model = "glm-4"
	}
	baseURL := "https://open.bigmodel.cn/api/paas/v4"
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}
	return &BigModelProvider{
		config:  config,
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

func (p *BigModelProvider) Name() string {
	return "bigmodel"
}

func (p *BigModelProvider) Assert(ctx context.Context, actual string, expected string) (bool, float64, error) {
	systemPrompt := `你是一个精确的文本对比判断专家。
给定一个"actual"文本和一个"expected"描述，判断actual文本是否满足expected描述。
请输出一个JSON对象，包含以下字段：
- "match": boolean (true表示actual满足expected，false表示不满足)
- "confidence": 0到1之间的数字（表示你的确定程度）

示例输出：
{"match": true, "confidence": 0.95}
{"match": false, "confidence": 0.85}

规则：
- 只有当actual文本明确满足expected描述时match才为true
- 如果actual文本与expected描述矛盾或不满足expected描述，match为false
- confidence反映你的确定程度（1.0表示完全确定）`

	reqBody := map[string]any{
		"model": p.config.Model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": fmt.Sprintf("Expected: %s\n\nActual: %s", expected, actual)},
		},
	}
	if p.config.Temperature > 0 {
		reqBody["temperature"] = p.config.Temperature
	}
	if p.config.TopP > 0 {
		reqBody["top_p"] = p.config.TopP
	}
	if p.config.MaxTokens > 0 {
		reqBody["max_tokens"] = p.config.MaxTokens
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return false, 0, fmt.Errorf("assert: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return false, 0, fmt.Errorf("assert: create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return false, 0, fmt.Errorf("assert: send request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, 0, fmt.Errorf("assert: decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return false, 0, fmt.Errorf("assert: no choices in response")
	}

	content := strings.TrimSpace(result.Choices[0].Message.Content)
	return parseAssertResponse(content)
}

func (p *BigModelProvider) Rerank(ctx context.Context, texts []string, expected string) ([]RankedText, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	reqBody := map[string]any{
		"model":             "rerank",
		"query":             expected,
		"documents":         texts,
		"top_n":             len(texts),
		"return_documents":  true,
		"return_raw_scores": true,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("rerank: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL+"/rerank", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("rerank: create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rerank: send request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Document       string  `json:"document"`
			Index          int     `json:"index"`
			RelevanceScore float64 `json:"relevance_score"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("rerank: decode response: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, nil
	}

	results := make([]RankedText, len(result.Results))
	for i, r := range result.Results {
		results[i] = RankedText{
			Text:  r.Document,
			Score: r.RelevanceScore,
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}
