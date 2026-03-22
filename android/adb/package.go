package adb

import (
	"fmt"
	"regexp"
	"time"
)

func (c *Client) Install(apkPath string, grantPerms bool) error {
	args := []string{"install"}
	if grantPerms {
		args = append(args, "-g")
	}
	args = append(args, apkPath)
	_, err := c.exec(args...)
	return err
}

func (c *Client) Uninstall(packageName string) error {
	_, err := c.exec("uninstall", packageName)
	return err
}

func (c *Client) CurrentPackage() (string, error) {
	output, err := c.Shell("dumpsys activity activities | grep mResumedActivity")
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`([\w.]+)/`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find foreground package")
	}

	return matches[1], nil
}

func (c *Client) WaitForIdle(timeout time.Duration) error {
	end := time.Now().Add(timeout)
	for time.Now().Before(end) {
		_, err := c.Shell("sleep 0.1 && dumpsys activity activities")
		if err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for idle")
}
