package ios

import (
	"github.com/liukunup/go-uop/pkg/wda"
)

type wdaClient struct {
	client   *wda.Client
	bundleID string
}

func newWDAClient(addr, bundleID string) (*wdaClient, error) {
	client, err := wda.NewClient(addr)
	if err != nil {
		return nil, err
	}

	return &wdaClient{
		client:   client,
		bundleID: bundleID,
	}, nil
}

func (c *wdaClient) Tap(x, y int) error {
	return c.client.Tap(x, y)
}

func (c *wdaClient) SendKeys(text string) error {
	return c.client.SendKeys(text)
}

func (c *wdaClient) GetSource() (string, error) {
	return c.client.GetSource()
}

func (c *wdaClient) GetAlertText() (string, error) {
	return c.client.GetAlertText()
}

func (c *wdaClient) AcceptAlert() error {
	return c.client.AcceptAlert()
}

func (c *wdaClient) DismissAlert() error {
	return c.client.DismissAlert()
}

func (c *wdaClient) Close() error {
	return c.client.Close()
}
