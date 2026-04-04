package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testFlowYAML = `name: 登录测试
devices:
  - id: iphone
    type: ios
    serial: 00001234-00123456789
  - id: android-tablet
    type: android
    serial: emulator-5554
defaultDevice: iphone
steps:
  - launch: com.example.app
  - tapOn: { text: "用户名" }
  - device: android-tablet
    tapOn: { text: "跳过" }
`

const testSuiteYAML = `name: 登录测试套件
devices:
  - id: iphone
    type: ios
    serial: 00001234-00123456789
defaultDevice: iphone
flows:
  - name: 登录流程
    path: ./flows/login.yaml
  - name: 登出流程
    path: ./flows/logout.yaml
`

func TestParseFlow(t *testing.T) {
	r := strings.NewReader(testFlowYAML)
	flow, err := ParseFlow(r)
	if err != nil {
		t.Fatalf("ParseFlow() error = %v", err)
	}

	// Verify Flow.Name
	if flow.Name != "登录测试" {
		t.Errorf("Flow.Name = %q, want %q", flow.Name, "登录测试")
	}

	// Verify Flow.Devices length
	if len(flow.Devices) != 2 {
		t.Errorf("len(Flow.Devices) = %d, want %d", len(flow.Devices), 2)
	}

	// Verify Flow.DefaultDevice
	if flow.DefaultDevice != "iphone" {
		t.Errorf("Flow.DefaultDevice = %q, want %q", flow.DefaultDevice, "iphone")
	}

	// Verify Flow.Steps length
	if len(flow.Steps) != 3 {
		t.Errorf("len(Flow.Steps) = %d, want %d", len(flow.Steps), 3)
	}
}

func TestParseFlowDevices(t *testing.T) {
	r := strings.NewReader(testFlowYAML)
	flow, err := ParseFlow(r)
	if err != nil {
		t.Fatalf("ParseFlow() error = %v", err)
	}

	// Verify first device
	if flow.Devices[0].ID != "iphone" {
		t.Errorf("Devices[0].ID = %q, want %q", flow.Devices[0].ID, "iphone")
	}
	if flow.Devices[0].Type != "ios" {
		t.Errorf("Devices[0].Type = %q, want %q", flow.Devices[0].Type, "ios")
	}
	if flow.Devices[0].Serial != "00001234-00123456789" {
		t.Errorf("Devices[0].Serial = %q, want %q", flow.Devices[0].Serial, "00001234-00123456789")
	}

	// Verify second device
	if flow.Devices[1].ID != "android-tablet" {
		t.Errorf("Devices[1].ID = %q, want %q", flow.Devices[1].ID, "android-tablet")
	}
	if flow.Devices[1].Type != "android" {
		t.Errorf("Devices[1].Type = %q, want %q", flow.Devices[1].Type, "android")
	}
	if flow.Devices[1].Serial != "emulator-5554" {
		t.Errorf("Devices[1].Serial = %q, want %q", flow.Devices[1].Serial, "emulator-5554")
	}
}

func TestParseFlowSteps(t *testing.T) {
	r := strings.NewReader(testFlowYAML)
	flow, err := ParseFlow(r)
	if err != nil {
		t.Fatalf("ParseFlow() error = %v", err)
	}

	// Verify step 1: launch command
	if _, ok := flow.Steps[0]["launch"]; !ok {
		t.Error("Step[0] should have 'launch' key")
	}
	if flow.Steps[0]["launch"] != "com.example.app" {
		t.Errorf("Step[0]['launch'] = %q, want %q", flow.Steps[0]["launch"], "com.example.app")
	}

	// Verify step 2: tapOn command with map value
	if _, ok := flow.Steps[1]["tapOn"]; !ok {
		t.Error("Step[1] should have 'tapOn' key")
	}
}

func TestParseFlowFile(t *testing.T) {
	// Create temp file with test content
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_flow.yaml")
	if err := os.WriteFile(tmpFile, []byte(testFlowYAML), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	flow, err := ParseFlowFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFlowFile() error = %v", err)
	}

	if flow.Name != "登录测试" {
		t.Errorf("Flow.Name = %q, want %q", flow.Name, "登录测试")
	}
	if len(flow.Devices) != 2 {
		t.Errorf("len(Flow.Devices) = %d, want %d", len(flow.Devices), 2)
	}
	if flow.DefaultDevice != "iphone" {
		t.Errorf("Flow.DefaultDevice = %q, want %q", flow.DefaultDevice, "iphone")
	}
	if len(flow.Steps) != 3 {
		t.Errorf("len(Flow.Steps) = %d, want %d", len(flow.Steps), 3)
	}
}

func TestParseFlowFileNotFound(t *testing.T) {
	_, err := ParseFlowFile("/nonexistent/path/flow.yaml")
	if err == nil {
		t.Error("ParseFlowFile() should return error for nonexistent file")
	}
}

func TestParseSuite(t *testing.T) {
	r := strings.NewReader(testSuiteYAML)
	suite, err := ParseSuite(r)
	if err != nil {
		t.Fatalf("ParseSuite() error = %v", err)
	}

	// Verify Suite.Name
	if suite.Name != "登录测试套件" {
		t.Errorf("Suite.Name = %q, want %q", suite.Name, "登录测试套件")
	}

	// Verify Suite.Devices length
	if len(suite.Devices) != 1 {
		t.Errorf("len(Suite.Devices) = %d, want %d", len(suite.Devices), 1)
	}

	// Verify Suite.DefaultDevice
	if suite.DefaultDevice != "iphone" {
		t.Errorf("Suite.DefaultDevice = %q, want %q", suite.DefaultDevice, "iphone")
	}

	// Verify Suite.Flows length
	if len(suite.Flows) != 2 {
		t.Errorf("len(Suite.Flows) = %d, want %d", len(suite.Flows), 2)
	}

	// Verify SuiteFlow entries
	if suite.Flows[0].Name != "登录流程" {
		t.Errorf("Suite.Flows[0].Name = %q, want %q", suite.Flows[0].Name, "登录流程")
	}
	if suite.Flows[0].Path != "./flows/login.yaml" {
		t.Errorf("Suite.Flows[0].Path = %q, want %q", suite.Flows[0].Path, "./flows/login.yaml")
	}
	if suite.Flows[1].Name != "登出流程" {
		t.Errorf("Suite.Flows[1].Name = %q, want %q", suite.Flows[1].Name, "登出流程")
	}
	if suite.Flows[1].Path != "./flows/logout.yaml" {
		t.Errorf("Suite.Flows[1].Path = %q, want %q", suite.Flows[1].Path, "./flows/logout.yaml")
	}
}

func TestParseSuiteFile(t *testing.T) {
	// Create temp file with test content
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_suite.yaml")
	if err := os.WriteFile(tmpFile, []byte(testSuiteYAML), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	suite, err := ParseSuiteFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseSuiteFile() error = %v", err)
	}

	if suite.Name != "登录测试套件" {
		t.Errorf("Suite.Name = %q, want %q", suite.Name, "登录测试套件")
	}
	if len(suite.Devices) != 1 {
		t.Errorf("len(Suite.Devices) = %d, want %d", len(suite.Devices), 1)
	}
	if len(suite.Flows) != 2 {
		t.Errorf("len(Suite.Flows) = %d, want %d", len(suite.Flows), 2)
	}
}

func TestParseSuiteFileNotFound(t *testing.T) {
	_, err := ParseSuiteFile("/nonexistent/path/suite.yaml")
	if err == nil {
		t.Error("ParseSuiteFile() should return error for nonexistent file")
	}
}

func TestParseInvalidYAML(t *testing.T) {
	invalidYAML := `name: test
devices:
  - id: iphone
    type: [invalid
steps:
  - launch: app`
	r := strings.NewReader(invalidYAML)
	_, err := ParseFlow(r)
	if err == nil {
		t.Error("ParseFlow() should return error for invalid YAML")
	}
}

func TestParseFlowEmptyYAML(t *testing.T) {
	r := strings.NewReader("")
	_, err := ParseFlow(r)
	if err == nil {
		t.Error("ParseFlow() should return error for empty YAML")
	}
}

func TestParseSuiteEmptyYAML(t *testing.T) {
	r := strings.NewReader("")
	_, err := ParseSuite(r)
	if err == nil {
		t.Error("ParseSuite() should return error for empty YAML")
	}
}
