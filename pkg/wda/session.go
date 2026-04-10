package wda

import (
	"encoding/json"
	"fmt"
)

// StartSession starts a new W3C WebDriver session
func (c *Client) StartSession(bundleID string) error {
	// W3C Capabilities structure with alwaysMatch
	req := NewSessionRequest{
		Capabilities: Capabilities{
			AlwaysMatch: map[string]any{
				"bundleId": bundleID,
			},
			FirstMatch: []map[string]any{
				{"bundleId": bundleID},
			},
		},
	}

	respBody, err := c.doRequest("POST", EndpointSession, req)
	if err != nil {
		return fmt.Errorf("start session: %w", err)
	}

	// W3C response envelope: {"value": {...}, "sessionId": "..."}
	var resp struct {
		Value     json.RawMessage `json:"value"`
		SessionID string          `json:"sessionId"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("parse session response: %w", err)
	}

	// Extract session ID from value envelope if not at top level
	if resp.SessionID == "" {
		var valueResp struct {
			SessionID string `json:"sessionId"`
		}
		if err := json.Unmarshal(resp.Value, &valueResp); err == nil {
			resp.SessionID = valueResp.SessionID
		}
	}

	c.SessionID = resp.SessionID

	// Parse capabilities from value
	if len(resp.Value) > 0 && resp.Value[0] == '{' {
		if err := json.Unmarshal(resp.Value, &c.Capabilities); err != nil {
			// Non-fatal: capabilities parsing is best-effort
		}
	}

	return nil
}

// StopSession ends the current session
func (c *Client) StopSession() error {
	if c.SessionID == "" {
		return nil
	}

	path := fmt.Sprintf("%s/%s", EndpointSession, c.SessionID)
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		// Session might already be invalid, but we still clear the ID
		c.SessionID = ""
		return err
	}
	c.SessionID = ""
	c.Capabilities = nil
	return nil
}

// NewSessionRequest creates a new session request with custom capabilities
type NewSessionOptions struct {
	BundleID                string
	Arguments               []string
	Environment             map[string]string
	ShouldWaitForQuiescence *bool
	MaxTypingFrequency      int
	ScreenshotOrientation   string
}

func (c *Client) StartSessionWithOptions(opts NewSessionOptions) error {
	capabilities := map[string]any{
		"bundleId": opts.BundleID,
	}

	if len(opts.Arguments) > 0 {
		capabilities["arguments"] = opts.Arguments
	}
	if len(opts.Environment) > 0 {
		capabilities["environment"] = opts.Environment
	}
	if opts.ShouldWaitForQuiescence != nil {
		capabilities["shouldWaitForQuiescence"] = *opts.ShouldWaitForQuiescence
	}
	if opts.MaxTypingFrequency > 0 {
		capabilities["maxTypingFrequency"] = opts.MaxTypingFrequency
	}
	if opts.ScreenshotOrientation != "" {
		capabilities["screenshotOrientation"] = opts.ScreenshotOrientation
	}

	req := NewSessionRequest{
		Capabilities: Capabilities{
			AlwaysMatch: capabilities,
			FirstMatch:  []map[string]any{capabilities},
		},
	}

	respBody, err := c.doRequest("POST", EndpointSession, req)
	if err != nil {
		return fmt.Errorf("start session with options: %w", err)
	}

	var resp struct {
		Value     json.RawMessage `json:"value"`
		SessionID string          `json:"sessionId"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return fmt.Errorf("parse session response: %w", err)
	}

	if resp.SessionID == "" {
		var valueResp struct {
			SessionID string `json:"sessionId"`
		}
		if err := json.Unmarshal(resp.Value, &valueResp); err == nil {
			resp.SessionID = valueResp.SessionID
		}
	}

	c.SessionID = resp.SessionID
	return nil
}
