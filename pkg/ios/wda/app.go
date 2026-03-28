package wda

import (
	"fmt"
)

func (c *Client) LaunchApp(bundleID string) error {
	body := map[string]interface{}{
		"bundleId": bundleID,
	}
	_, err := c.doRequest("POST", EndpointAppLaunch, body)
	return err
}

func (c *Client) TerminateApp(bundleID string) error {
	path := fmt.Sprintf(EndpointAppTerminate, bundleID)
	_, err := c.doRequest("POST", path, nil)
	return err
}
