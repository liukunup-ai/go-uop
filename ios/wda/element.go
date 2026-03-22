package wda

import (
	"encoding/json"
	"fmt"
)

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
