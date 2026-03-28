package adb

func (c *Client) StartActivity(component string) error {
	_, err := c.exec("shell", "am", "start", "-n", component)
	return err
}

func (c *Client) StopPackage(packageName string) error {
	_, err := c.exec("shell", "am", "force-stop", packageName)
	return err
}
