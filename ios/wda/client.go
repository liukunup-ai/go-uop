package wda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL    *url.URL
	SessionID  string
	HTTPClient *http.Client
}

func NewClient(baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		BaseURL:   u,
		SessionID: "",
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

func NewClientWithSession(baseURL string, bundleID string) (*Client, error) {
	client, err := NewClient(baseURL)
	if err != nil {
		return nil, err
	}

	if err := client.StartSession(bundleID); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	u := c.BaseURL.ResolveReference(&url.URL{Path: path})

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) IsHealthy() bool {
	_, err := c.doRequest("GET", EndpointStatus, nil)
	return err == nil
}
