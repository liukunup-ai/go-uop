package console

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"

	goios "github.com/danielpaulus/go-ios/ios"
	log "github.com/sirupsen/logrus"
)

type forwardProcess struct {
	cmd   *exec.Cmd
	udid  string
	ports []int
}

type IOSManager struct {
	mu           sync.RWMutex
	forwardProcs map[string]*forwardProcess
	wdaProcs     map[string]*exec.Cmd
	usedPorts    map[int]bool
	portToUDID   map[int]string
}

func NewIOSManager() *IOSManager {
	return &IOSManager{
		forwardProcs: make(map[string]*forwardProcess),
		wdaProcs:     make(map[string]*exec.Cmd),
		usedPorts:    make(map[int]bool),
		portToUDID:   make(map[int]string),
	}
}

func (m *IOSManager) ListDevices() ([]IOSDeviceInfo, error) {
	cmd := exec.Command("ios", "list", "--details")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list iOS devices: %w", err)
	}

	type iosListOutput struct {
		DeviceList []struct {
			DeviceID       int    `json:"deviceId"`
			Udid           string `json:"udid"`
			ProductName    string `json:"productName,omitempty"`
			ProductType    string `json:"productType,omitempty"`
			ProductVersion string `json:"productVersion,omitempty"`
		} `json:"deviceList"`
	}
	var listOutput iosListOutput
	if err := json.Unmarshal(output, &listOutput); err != nil {
		return nil, fmt.Errorf("failed to parse ios list output: %w", err)
	}

	var result []IOSDeviceInfo
	for _, dev := range listOutput.DeviceList {
		info := IOSDeviceInfo{
			UDID:       dev.Udid,
			Name:       dev.ProductName,
			Model:      dev.ProductType,
			IOSVersion: dev.ProductVersion,
			Status:     "available",
		}

		m.mu.RLock()
		if m.wdaProcs[dev.Udid] != nil {
			info.Status = "wda_running"
		} else if m.forwardProcs[dev.Udid] != nil {
			info.Status = "forwarding"
		}
		m.mu.RUnlock()

		result = append(result, info)
	}

	return result, nil
}

type IOSDeviceInfo struct {
	UDID       string `json:"udid"`
	Name       string `json:"name,omitempty"`
	Model      string `json:"model,omitempty"`
	IOSVersion string `json:"iosVersion,omitempty"`
	Status     string `json:"status"`
}

func (m *IOSManager) SetupPortForwarding(udid string, ports []int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if proc, ok := m.forwardProcs[udid]; ok {
		proc.cmd.Process.Kill()
		proc.cmd.Wait()
		delete(m.forwardProcs, udid)
	}

	for port := range m.usedPorts {
		if m.portToUDID[port] == udid {
			delete(m.usedPorts, port)
			delete(m.portToUDID, port)
		}
	}

	for _, port := range ports {
		if port == 9100 {
			cmd := exec.Command("ios", "screenshot", "--stream", "--udid="+udid)
			cmd.Stdout = nil
			cmd.Stderr = nil
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("failed to start screenshot stream: %w", err)
			}
			m.forwardProcs[udid] = &forwardProcess{
				cmd:   cmd,
				udid:  udid,
				ports: []int{9100},
			}
			m.usedPorts[9100] = true
			m.portToUDID[9100] = udid
			log.Infof("Screenshot stream started for device %s (MJPEG at localhost:3333)", udid)
			continue
		}

		args := []string{"forward", "--udid=" + udid, "--port=" + strconv.Itoa(port) + ":" + strconv.Itoa(port)}
		cmd := exec.Command("ios", args...)
		cmd.Stdout = nil
		cmd.Stderr = nil
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to forward port %d: %w", port, err)
		}
		m.usedPorts[port] = true
		m.portToUDID[port] = udid
		m.forwardProcs[udid] = &forwardProcess{
			cmd:   cmd,
			udid:  udid,
			ports: []int{port},
		}
		log.Infof("Port forwarding started: port %d for device %s", port, udid)
	}

	return nil
}

func (m *IOSManager) RemovePortForwarding(udid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	proc, ok := m.forwardProcs[udid]
	if !ok || proc == nil {
		return nil
	}

	proc.cmd.Process.Signal(syscall.SIGTERM)
	proc.cmd.Wait()

	delete(m.forwardProcs, udid)

	for port := range m.usedPorts {
		if m.portToUDID[port] == udid {
			delete(m.usedPorts, port)
			delete(m.portToUDID, port)
		}
	}

	log.Infof("Port forwarding stopped for device %s", udid)
	return nil
}

func (m *IOSManager) LaunchWDA(udid string, bundleID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.wdaProcs[udid] != nil {
		log.Infof("WDA already running for device %s", udid)
		return nil
	}

	if bundleID == "" {
		bundleID = "com.facebook.WebDriverAgentRunner.xctrunner"
	}

	args := []string{
		"runwda",
		"--udid=" + udid,
		"--bundleid=" + bundleID,
		"--testrunnerbundleid=" + bundleID,
		"--xctestconfig=WebDriverAgentRunner.xctest",
	}

	cmd := exec.Command("ios", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start WDA: %w", err)
	}

	m.wdaProcs[udid] = cmd

	log.Infof("WDA started: bundleID=%s for device %s", bundleID, udid)
	return nil
}

func (m *IOSManager) StopWDA(udid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmd, ok := m.wdaProcs[udid]
	if !ok || cmd == nil {
		return nil
	}

	cmd.Process.Signal(syscall.SIGTERM)
	cmd.Wait()

	delete(m.wdaProcs, udid)

	log.Infof("WDA stopped for device %s", udid)
	return nil
}

func (m *IOSManager) IsWDARunning(udid string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.wdaProcs[udid] != nil
}

func (m *IOSManager) IsForwarding(udid string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.forwardProcs[udid] != nil
}

func (m *IOSManager) GetForwardedPorts(udid string) []int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var ports []int
	for port, d := range m.portToUDID {
		if d == udid {
			ports = append(ports, port)
		}
	}
	return ports
}

func (m *IOSManager) GetIOSForwardCommand(udid string, ports []int) string {
	var portArgs []string
	for _, port := range ports {
		portArgs = append(portArgs, "--port="+strconv.Itoa(port)+":"+strconv.Itoa(port))
	}
	return "ios forward --udid=" + udid + " " + strings.Join(portArgs, " ")
}

func (m *IOSManager) GetPreviewURL(udid string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.forwardProcs[udid] != nil {
		for _, port := range m.forwardProcs[udid].ports {
			if port == 9100 {
				return "http://localhost:3333"
			}
		}
	}
	return ""
}

func (m *IOSManager) GetIOSRunWDACCommand(udid string, bundleID string) string {
	if bundleID == "" {
		bundleID = "com.facebook.WebDriverAgentRunner.xctrunner"
	}
	return "ios runwda --udid=" + udid + " --bundleid=" + bundleID + " --testrunnerbundleid=" + bundleID + " --xctestconfig=WebDriverAgentRunner.xctest"
}

func (m *IOSManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for udid, cmd := range m.wdaProcs {
		if cmd != nil {
			cmd.Process.Signal(syscall.SIGTERM)
			cmd.Wait()
		}
		delete(m.wdaProcs, udid)
	}

	for udid, proc := range m.forwardProcs {
		if proc != nil && proc.cmd != nil && proc.cmd.Process != nil {
			proc.cmd.Process.Signal(syscall.SIGTERM)
			proc.cmd.Wait()
		}
		delete(m.forwardProcs, udid)
	}

	return nil
}

func (m *IOSManager) getDevice(udid string) (goios.DeviceEntry, error) {
	deviceList, err := goios.ListDevices()
	if err != nil {
		return goios.DeviceEntry{}, err
	}

	for _, dev := range deviceList.DeviceList {
		if dev.Properties.SerialNumber == udid {
			return dev, nil
		}
	}
	return goios.DeviceEntry{}, fmt.Errorf("device not found: %s", udid)
}
