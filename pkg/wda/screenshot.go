package wda

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Screenshot returns a screenshot as PNG bytes
func (c *Client) Screenshot() ([]byte, error) {
	respBody, err := c.doSessionRequest("GET", EndpointScreenshot, nil)
	if err != nil {
		return nil, fmt.Errorf("screenshot: %w", err)
	}

	// W3C response: {"value": "base64-encoded-png...", "sessionId": "..."}
	var resp ScreenshotResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parse screenshot response: %w", err)
	}

	// Decode base64 value
	imgData, err := base64.StdEncoding.DecodeString(resp.Value)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	return imgData, nil
}

// ScreenshotBase64 returns the raw base64-encoded screenshot
func (c *Client) ScreenshotBase64() (string, error) {
	respBody, err := c.doSessionRequest("GET", EndpointScreenshot, nil)
	if err != nil {
		return "", fmt.Errorf("screenshot: %w", err)
	}

	var resp ScreenshotResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse screenshot response: %w", err)
	}

	return resp.Value, nil
}
