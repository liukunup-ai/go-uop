package wda

import (
	"encoding/json"
	"fmt"
)

// LaunchApp launches the specified app
func (c *Client) LaunchApp(bundleID string) error {
	// Try W3C app/launch endpoint
	body := map[string]any{
		"bundleId": bundleID,
	}

	_, err := c.doSessionRequest("POST", EndpointAppLaunch, body)
	if err != nil {
		// Fallback to WDA legacy endpoint
		_, err = c.doRequest("POST", EndpointWDAAppLaunch, body)
		if err != nil {
			return fmt.Errorf("launch app: %w", err)
		}
	}
	return nil
}

// TerminateApp terminates the specified app
func (c *Client) TerminateApp(bundleID string) error {
	// Try W3C app/terminate endpoint
	body := map[string]any{
		"bundleId": bundleID,
	}

	respBody, err := c.doSessionRequest("POST", EndpointAppTerminate, body)
	if err != nil {
		return fmt.Errorf("terminate app: %w", err)
	}

	// Parse boolean response
	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil // Non-fatal
	}

	return nil
}

// ActivateApp brings the app to foreground
func (c *Client) ActivateApp(bundleID string) error {
	body := map[string]any{
		"bundleId": bundleID,
	}

	_, err := c.doSessionRequest("POST", EndpointAppActivate, body)
	if err != nil {
		return fmt.Errorf("activate app: %w", err)
	}

	return nil
}

// GetCurrentApp returns the bundle ID of the current app
func (c *Client) GetCurrentApp() (string, error) {
	respBody, err := c.doSessionRequest("GET", "/app/active", nil)
	if err != nil {
		return "", fmt.Errorf("get current app: %w", err)
	}

	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse current app: %w", err)
	}

	if val, ok := resp.Value.(string); ok {
		return val, nil
	}

	return "", nil
}

// GetRunningApps returns list of running app bundle IDs
func (c *Client) GetRunningApps() ([]string, error) {
	respBody, err := c.doSessionRequest("GET", "/app/device/available_count", nil)
	if err != nil {
		return nil, fmt.Errorf("get running apps: %w", err)
	}

	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parse running apps: %w", err)
	}

	if val, ok := resp.Value.([]any); ok {
		apps := make([]string, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				apps = append(apps, s)
			}
		}
		return apps, nil
	}

	return nil, nil
}
