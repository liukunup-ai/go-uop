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

func (c *Client) StartSession(bundleID string) error {
	body := map[string]interface{}{
		"capabilities": map[string]interface{}{
			"bundleId": bundleID,
		},
	}

	respBody, err := c.doRequest("POST", EndpointSession, body)
	if err != nil {
		return fmt.Errorf("start session: %w", err)
	}

	var resp struct {
		SessionID string `json:"sessionId"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("parse session: %w", err)
	}

	c.SessionID = resp.SessionID
	return nil
}

func (c *Client) StopSession() error {
	if c.SessionID == "" {
		return nil
	}

	path := fmt.Sprintf("%s/%s", EndpointSession, c.SessionID)
	_, err := c.doRequest("DELETE", path, nil)
	c.SessionID = ""
	return err
}

func (c *Client) Screenshot() ([]byte, error) {
	respBody, err := c.doRequest("GET", EndpointScreenshot, nil)
	if err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}

	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parse screenshot: %w", err)
	}

	return []byte(resp.Value), nil
}

func (c *Client) Tap(x, y int) error {
	path := fmt.Sprintf(EndpointTap, x, y)
	_, err := c.doRequest("POST", path, nil)
	return err
}

func (c *Client) SendKeys(text string) error {
	body := map[string]interface{}{
		"value": []string{text},
	}
	_, err := c.doRequest("POST", EndpointKeys, body)
	return err
}

func (c *Client) GetSource() (string, error) {
	respBody, err := c.doRequest("GET", EndpointSource, nil)
	if err != nil {
		return "", fmt.Errorf("get source: %w", err)
	}

	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse source: %w", err)
	}

	return resp.Value, nil
}

func (c *Client) LaunchApp(bundleID string) error {
	body := map[string]interface{}{
		"bundleId": bundleID,
	}
	_, err := c.doRequest("POST", EndpointAppLaunch, body)
	return err
}

func (c *Client) TerminateApp(bundleID string) error {
	path := fmt.Sprintf(EndpointAppTerminate, bundleID)
	_, err := c.doRequest("POST", path, nil)
	return err
}

func (c *Client) IsHealthy() bool {
	_, err := c.doRequest("GET", EndpointStatus, nil)
	return err == nil
}
