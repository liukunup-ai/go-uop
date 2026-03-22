package adb

func (c *Client) Shell(command string) (string, error) {
	return c.exec("shell", command)
}
