package wda

import (
	"encoding/json"
	"fmt"
)

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
