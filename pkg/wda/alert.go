package wda

import (
	"encoding/json"
	"fmt"
)

// AcceptAlert accepts the current alert
func (c *Client) AcceptAlert() error {
	_, err := c.doSessionRequest("POST", EndpointAlert+"/"+string(AlertAccept), nil)
	return err
}

// DismissAlert dismisses the current alert
func (c *Client) DismissAlert() error {
	_, err := c.doSessionRequest("POST", EndpointAlert+"/"+string(AlertDismiss), nil)
	return err
}

// GetAlertText returns the text of the current alert
func (c *Client) GetAlertText() (string, error) {
	respBody, err := c.doSessionRequest("GET", EndpointAlert+"/"+string(AlertText), nil)
	if err != nil {
		return "", fmt.Errorf("get alert text: %w", err)
	}

	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse alert text: %w", err)
	}

	if val, ok := resp.Value.(string); ok {
		return val, nil
	}

	return "", nil
}

// PostAlertText sends text to the alert prompt
func (c *Client) PostAlertText(text string) error {
	path := EndpointAlert + "/" + string(AlertText)
	body := map[string]any{
		"text": text,
	}
	_, err := c.doSessionRequest("POST", path, body)
	return err
}

// HasAlert checks if there is an active alert
func (c *Client) HasAlert() (bool, error) {
	_, err := c.doSessionRequest("GET", EndpointAlert+"/text", nil)
	if err != nil {
		// WDA returns error when no alert is present
		if w3cErr, ok := err.(*W3CError); ok {
			if w3cErr.Code == ErrNoSuchAlert {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}
