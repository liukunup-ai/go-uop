package wda

import (
	"encoding/json"
	"fmt"
)

func (c *Client) AcceptAlert() error {
	path := fmt.Sprintf(EndpointAlert, string(AlertAccept))
	_, err := c.doRequest("POST", path, nil)
	return err
}

func (c *Client) DismissAlert() error {
	path := fmt.Sprintf(EndpointAlert, string(AlertDismiss))
	_, err := c.doRequest("POST", path, nil)
	return err
}

func (c *Client) GetAlertText() (string, error) {
	path := fmt.Sprintf(EndpointAlert, string(AlertText))
	respBody, err := c.doRequest("GET", path, nil)
	if err != nil {
		return "", fmt.Errorf("get alert text: %w", err)
	}

	var resp struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse alert text: %w", err)
	}

	return resp.Value, nil
}

func (c *Client) PostAlertText(text string) error {
	path := fmt.Sprintf(EndpointAlert, string(AlertText))
	body := map[string]interface{}{
		"value": text,
	}
	_, err := c.doRequest("POST", path, body)
	return err
}
