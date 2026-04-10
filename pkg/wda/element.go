package wda

import (
	"encoding/json"
	"fmt"
)

// FindElement finds a single element using the given strategy
func (c *Client) FindElement(strategy, selector string) (ElementID, error) {
	req := FindElementRequest{
		Using: strategy,
		Value: selector,
	}

	respBody, err := c.doSessionRequest("POST", EndpointElement, req)
	if err != nil {
		return "", fmt.Errorf("find element: %w", err)
	}

	// W3C response: {"value": {"element-60611-11-0-1": {...}}, "sessionId": "..."}
	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse element response: %w", err)
	}

	// Extract element ID from value map
	if resp.Value == nil {
		return "", fmt.Errorf("element not found")
	}

	valueMap, ok := resp.Value.(map[string]any)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	// Find the element key (format: "element-UUID-index")
	for k, v := range valueMap {
		if k != "ELEMENT" && k != "element-60611-11-0-1" {
			continue
		}
		if elemID, ok := v.(string); ok {
			return ElementID(elemID), nil
		}
	}

	// Try ELEMENT key (JSON-Wire compatibility)
	if elemID, ok := valueMap["ELEMENT"].(string); ok {
		return ElementID(elemID), nil
	}

	return "", fmt.Errorf("element ID not found in response")
}

// FindElements finds multiple elements using the given strategy
func (c *Client) FindElements(strategy, selector string) ([]ElementID, error) {
	req := FindElementRequest{
		Using: strategy,
		Value: selector,
	}

	respBody, err := c.doSessionRequest("POST", EndpointElements, req)
	if err != nil {
		return nil, fmt.Errorf("find elements: %w", err)
	}

	// W3C response: {"value": [{"element-60611-11-0-1": {...}}, ...]}
	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("parse elements response: %w", err)
	}

	valueList, ok := resp.Value.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	elements := make([]ElementID, 0, len(valueList))
	for _, item := range valueList {
		if elemMap, ok := item.(map[string]any); ok {
			// Try ELEMENT key first
			if elemID, ok := elemMap["ELEMENT"].(string); ok {
				elements = append(elements, ElementID(elemID))
				continue
			}
			// Find any element key
			for k, v := range elemMap {
				if k == "ELEMENT" || (len(k) > 7 && k[:7] == "element") {
					if elemID, ok := v.(string); ok {
						elements = append(elements, ElementID(elemID))
						break
					}
				}
			}
		}
	}

	return elements, nil
}

// GetSource returns the page source
func (c *Client) GetSource() (string, error) {
	respBody, err := c.doSessionRequest("GET", EndpointSource, nil)
	if err != nil {
		// Try WDA legacy endpoint
		respBody, err = c.doRequest("GET", EndpointWDASource, nil)
		if err != nil {
			return "", fmt.Errorf("get source: %w", err)
		}
	}

	var resp SourceResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse source: %w", err)
	}

	return resp.Value, nil
}

// Click clicks on an element
func (c *Client) Click(elemID ElementID) error {
	path := fmt.Sprintf(EndpointElementID, elemID) + "/click"
	_, err := c.doSessionRequest("POST", path, nil)
	return err
}

// SendKeys sends text to an element or the active element
func (c *Client) SendKeys(text string) error {
	_, err := c.doSessionRequest("POST", "/element/active/value", map[string]any{
		"text": text,
	})
	if err != nil {
		// Fallback to WDA legacy endpoint
		bodyLegacy := map[string]any{
			"value": []string{text},
		}
		_, err = c.doRequest("POST", EndpointWDAKeys, bodyLegacy)
		if err != nil {
			return fmt.Errorf("send keys: %w", err)
		}
	}
	return nil
}

// GetElementRect returns the element's bounding rectangle
func (c *Client) GetElementRect(elemID ElementID) (ElementRect, error) {
	path := fmt.Sprintf(EndpointElementID, elemID) + "/rect"
	respBody, err := c.doSessionRequest("GET", path, nil)
	if err != nil {
		return ElementRect{}, fmt.Errorf("get element rect: %w", err)
	}

	var rect ElementRect
	if err := json.Unmarshal(respBody, &rect); err != nil {
		// Try extracting from value envelope
		var resp W3CResponse
		if err := json.Unmarshal(respBody, &resp); err != nil {
			return ElementRect{}, fmt.Errorf("parse rect: %w", err)
		}
		if rectMap, ok := resp.Value.(map[string]any); ok {
			if x, ok := rectMap["x"].(float64); ok {
				rect.X = x
			}
			if y, ok := rectMap["y"].(float64); ok {
				rect.Y = y
			}
			if w, ok := rectMap["width"].(float64); ok {
				rect.Width = w
			}
			if h, ok := rectMap["height"].(float64); ok {
				rect.Height = h
			}
		}
	}

	return rect, nil
}

// Tap performs a tap at the given coordinates using W3C Actions
func (c *Client) Tap(x, y int) error {
	// Try W3C Actions API first
	actions := []map[string]any{
		{
			"type": "pointer",
			"id":   "finger1",
			"parameters": map[string]any{
				"pointerType": "touch",
			},
			"actions": []map[string]any{
				{"type": "pointerMove", "x": x, "y": y, "duration": 0},
				{"type": "pointerDown", "button": 0},
				{"type": "pointerUp", "button": 0},
			},
		},
	}

	_, err := c.doSessionRequest("POST", EndpointActions, map[string]any{
		"actions": actions,
	})
	if err != nil {
		// Fallback to WDA legacy endpoint
		path := fmt.Sprintf(EndpointWDATap, x, y)
		_, err = c.doRequest("POST", path, nil)
		if err != nil {
			return fmt.Errorf("tap: %w", err)
		}
	}
	return nil
}

// Swipe performs a swipe gesture
func (c *Client) Swipe(startX, startY, endX, endY int, durationMs int) error {
	actions := []map[string]any{
		{
			"type": "pointer",
			"id":   "finger1",
			"parameters": map[string]any{
				"pointerType": "touch",
			},
			"actions": []map[string]any{
				{"type": "pointerMove", "x": startX, "y": startY, "duration": 0},
				{"type": "pointerDown", "button": 0},
				{"type": "pause", "duration": 100},
				{"type": "pointerMove", "x": endX, "y": endY, "duration": durationMs},
				{"type": "pointerUp", "button": 0},
			},
		},
	}

	_, err := c.doSessionRequest("POST", EndpointActions, map[string]any{
		"actions": actions,
	})
	return err
}

// GetAttribute returns an element's attribute value
func (c *Client) GetAttribute(elemID ElementID, name string) (string, error) {
	path := fmt.Sprintf(EndpointElementID, elemID) + "/attribute/" + name
	respBody, err := c.doSessionRequest("GET", path, nil)
	if err != nil {
		return "", fmt.Errorf("get attribute: %w", err)
	}

	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse attribute: %w", err)
	}

	if resp.Value == nil {
		return "", nil
	}
	if val, ok := resp.Value.(string); ok {
		return val, nil
	}
	return "", nil
}

// GetText returns an element's text content
func (c *Client) GetText(elemID ElementID) (string, error) {
	path := fmt.Sprintf(EndpointElementID, elemID) + "/text"
	respBody, err := c.doSessionRequest("GET", path, nil)
	if err != nil {
		return "", fmt.Errorf("get text: %w", err)
	}

	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", fmt.Errorf("parse text: %w", err)
	}

	if val, ok := resp.Value.(string); ok {
		return val, nil
	}
	return "", nil
}

// IsDisplayed checks if an element is displayed
func (c *Client) IsDisplayed(elemID ElementID) (bool, error) {
	path := fmt.Sprintf(EndpointElementID, elemID) + "/displayed"
	respBody, err := c.doSessionRequest("GET", path, nil)
	if err != nil {
		return false, fmt.Errorf("check displayed: %w", err)
	}

	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return false, fmt.Errorf("parse displayed: %w", err)
	}

	if val, ok := resp.Value.(bool); ok {
		return val, nil
	}
	return false, nil
}

// IsEnabled checks if an element is enabled
func (c *Client) IsEnabled(elemID ElementID) (bool, error) {
	path := fmt.Sprintf(EndpointElementID, elemID) + "/enabled"
	respBody, err := c.doSessionRequest("GET", path, nil)
	if err != nil {
		return false, fmt.Errorf("check enabled: %w", err)
	}

	var resp W3CResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return false, fmt.Errorf("parse enabled: %w", err)
	}

	if val, ok := resp.Value.(bool); ok {
		return val, nil
	}
	return false, nil
}

// Clear clears an input element
func (c *Client) Clear(elemID ElementID) error {
	path := fmt.Sprintf(EndpointElementID, elemID) + "/clear"
	_, err := c.doSessionRequest("POST", path, nil)
	return err
}

// ReleaseActions releases all pressed keys/pointers
func (c *Client) ReleaseActions() error {
	_, err := c.doSessionRequest("DELETE", EndpointActions, nil)
	return err
}
