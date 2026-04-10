package server

import (
	"fmt"
	"os/exec"
	"strings"
)

const (
	atxAgentApk          = "apks/atx-agent.apk"
	uiautomatorServerApk = "apks/uiautomator-server.apk"
	atxAgentPkg          = "com.github.uiautomator"
	atxAgentActivity     = "com.github.uiautomator/.MainActivity"
)

func EnsureInstalled(serial string) error {
	if serial == "" {
		serial = getEnvSerial()
	}

	if isRunning(serial) {
		return nil
	}

	if err := installAtxAgent(serial); err != nil {
		return fmt.Errorf("install atx-agent: %w", err)
	}

	if err := installUiautomatorServer(serial); err != nil {
		return fmt.Errorf("install uiautomator-server: %w", err)
	}

	if err := startServices(serial); err != nil {
		return fmt.Errorf("start services: %w", err)
	}

	return nil
}

func isRunning(serial string) bool {
	output, err := exec.Command("adb", "-s", serial, "shell", "curl", "-s", "-m", "2", "http://localhost:7912/status").CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "ok")
}

func installAtxAgent(serial string) error {
	_, err := exec.Command("adb", "-s", serial, "install", "-r", "-t", atxAgentApk).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func installUiautomatorServer(serial string) error {
	_, err := exec.Command("adb", "-s", serial, "install", "-r", "-t", uiautomatorServerApk).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func startServices(serial string) error {
	_, err := exec.Command("adb", "-s", serial, "shell", "am", "start", "-n", atxAgentActivity).CombinedOutput()
	if err != nil {
		return err
	}
	return waitForService(serial)
}

func waitForService(serial string) error {
	for i := 0; i < 30; i++ {
		if isRunning(serial) {
			return nil
		}
	}
	return fmt.Errorf("service did not start in time")
}

func getEnvSerial() string {
	for _, s := range []string{"ANDROID_SERIAL", "ANDROID_DEVICE_ID"} {
		if val := strings.TrimSpace(getEnv(s)); val != "" {
			return val
		}
	}

	devices, err := getDevices()
	if err != nil || len(devices) == 0 {
		return ""
	}
	return devices[0]
}

func getEnv(key string) string {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo $%s", key))
	output, _ := cmd.CombinedOutput()
	return strings.TrimSpace(string(output))
}

func getDevices() ([]string, error) {
	cmd := exec.Command("adb", "devices")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var devices []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "device" {
			devices = append(devices, fields[0])
		}
	}
	return devices, nil
}
