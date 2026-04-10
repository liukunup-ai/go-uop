package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/liukunup/go-uop/pkg/uiautomator2"
)

type Manager struct {
	client *uiautomator2.Client
	serial string
}

func NewManager(client *uiautomator2.Client, serial string) *Manager {
	return &Manager{client: client, serial: serial}
}

func (m *Manager) Install(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download APK: %w", err)
	}
	defer resp.Body.Close()

	tmpFile := filepath.Join(os.TempDir(), "uiautomator2_install.apk")
	f, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile)

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	f.Close()

	devicePath := "/sdcard/tmp_install.apk"
	if err := m.pushToDevice(tmpFile, devicePath); err != nil {
		return fmt.Errorf("push to device: %w", err)
	}

	return m.client.InstallApk(devicePath)
}

func (m *Manager) InstallLocal(apkPath string) error {
	devicePath := "/sdcard/tmp_install.apk"
	if err := m.pushToDevice(apkPath, devicePath); err != nil {
		return fmt.Errorf("push to device: %w", err)
	}
	return m.client.InstallApk(devicePath)
}

func (m *Manager) Start(pkg string, activity string, useMonkey bool) error {
	if activity != "" {
		return m.client.StartActivity(pkg, activity)
	}
	return m.client.StartApp(pkg)
}

func (m *Manager) Stop(pkg string) error {
	return m.client.ForceStop(pkg)
}

func (m *Manager) Clear(pkg string) error {
	return m.client.Clear(pkg)
}

func (m *Manager) StopAll(excludes []string) error {
	running, err := m.client.ListRunningApps()
	if err != nil {
		return err
	}
	for _, p := range running {
		if contains(excludes, p) {
			continue
		}
		m.client.ForceStop(p)
	}
	return nil
}

func (m *Manager) Info(pkg string) (*uiautomator2.AppInfo, error) {
	return m.client.AppInfo(pkg)
}

func (m *Manager) Icon(pkg string) ([]byte, error) {
	return m.client.AppIcon(pkg)
}

func (m *Manager) ListRunning() ([]string, error) {
	return m.client.ListRunningApps()
}

func (m *Manager) Wait(pkg string, front bool, timeout float64) (int, error) {
	return m.client.AppWait(pkg, front, timeout)
}

func (m *Manager) Push(localPath, remotePath string, mode int) error {
	if err := m.pushToDevice(localPath, remotePath); err != nil {
		return err
	}
	if mode != 0 {
		return m.setFileMode(remotePath, mode)
	}
	return nil
}

func (m *Manager) Pull(remotePath, localPath string) error {
	return m.pullFromDevice(remotePath, localPath)
}

func (m *Manager) AutoGrantPermissions(pkg string) error {
	return m.client.GrantPermissions(pkg)
}

func (m *Manager) OpenUrl(url string) error {
	return m.client.OpenUrl(url)
}

func (m *Manager) pushToDevice(localPath, remotePath string) error {
	args := []string{"push", localPath, remotePath}
	if m.serial != "" {
		args = append([]string{"-s", m.serial}, args...)
	}
	cmd := exec.Command("adb", args...)
	_, err := cmd.CombinedOutput()
	return err
}

func (m *Manager) pullFromDevice(remotePath, localPath string) error {
	args := []string{"pull", remotePath, localPath}
	if m.serial != "" {
		args = append([]string{"-s", m.serial}, args...)
	}
	cmd := exec.Command("adb", args...)
	_, err := cmd.CombinedOutput()
	return err
}

func (m *Manager) setFileMode(path string, mode int) error {
	args := []string{"shell", "chmod", fmt.Sprintf("%o", mode), path}
	if m.serial != "" {
		args = append([]string{"-s", m.serial}, args...)
	}
	cmd := exec.Command("adb", args...)
	_, err := cmd.CombinedOutput()
	return err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

type Session struct {
	client *uiautomator2.Client
	pkg    string
}

func NewSession(client *uiautomator2.Client, pkg string) *Session {
	return &Session{client: client, pkg: pkg}
}

func (s *Session) Start() error {
	return s.client.StartApp(s.pkg)
}

func (s *Session) Close() error {
	return s.client.ForceStop(s.pkg)
}

func (s *Session) Restart() error {
	if err := s.client.ForceStop(s.pkg); err != nil {
		return err
	}
	return s.client.StartApp(s.pkg)
}

func (s *Session) Running() bool {
	apps, err := s.client.ListRunningApps()
	if err != nil {
		return false
	}
	return contains(apps, s.pkg)
}

func (s *Session) Attach() error {
	return s.client.StartApp(s.pkg)
}
