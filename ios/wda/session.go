package wda

import (
	"encoding/json"
	"fmt"
)

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
