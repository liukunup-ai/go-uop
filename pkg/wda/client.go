package wda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is the W3C WebDriver client for iOS WebDriverAgent
type Client struct {
	BaseURL      *url.URL
	SessionID    string
	Capabilities map[string]any
	HTTPClient   *http.Client
}

// NewClient creates a new WDA client
func NewClient(baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		BaseURL:      u,
		SessionID:    "",
		Capabilities: nil,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// NewClientWithSession creates a client and starts a session
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

// doRequest performs an HTTP request with W3C compliance
func (c *Client) doRequest(method, path string, body any) ([]byte, error) {
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

	// Add W3C Session-Id header if session exists
	if c.SessionID != "" {
		req.Header.Set("X-Session-Id", c.SessionID)
	}

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
		// Try to parse as W3C error
		var errResp W3CErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Value.Error != "" {
			return nil, NewW3CError(&errResp)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// doSessionRequest performs a request within an active session
func (c *Client) doSessionRequest(method, path string, body any) ([]byte, error) {
	if c.SessionID == "" {
		return nil, fmt.Errorf("no active session")
	}
	sessionPath := fmt.Sprintf("%s/%s%s", strings.TrimSuffix(EndpointSession, "/"), c.SessionID, path)
	return c.doRequest(method, sessionPath, body)
}

// IsHealthy checks if the WDA server is responsive
func (c *Client) IsHealthy() bool {
	_, err := c.doRequest("GET", EndpointStatus, nil)
	return err == nil
}

// Close closes the session and releases resources
func (c *Client) Close() error {
	return c.StopSession()
}

// GetSessionID returns the current session ID
func (c *Client) GetSessionID() string {
	return c.SessionID
}

// GetCapabilities returns the session capabilities
func (c *Client) GetCapabilities() map[string]any {
	return c.Capabilities
}
