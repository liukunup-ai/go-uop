//go:build integration

package maestro

import (
	"strings"
	"testing"

	"github.com/liukunup/go-uop/core"
)

// mockDevice implements core.Device for integration testing
type mockDevice struct {
	platform_      core.Platform
	info_          map[string]interface{}
	tapCalls       []tapCall
	sendKeysCalls  []string
	launchCalled   bool
	screenshotData []byte
	shouldFail     bool
	failError      error
}

type tapCall struct {
	x, y int
}

func newMockDevice(platform core.Platform) *mockDevice {
	return &mockDevice{
		platform_: platform,
		info_: map[string]interface{}{
			"platform":      platform,
			"model":         "Mock Device",
			"serial":        "mock-serial-001",
			"os_version":    "14.0",
			"screen_width":  1080,
			"screen_height": 1920,
		},
		screenshotData: []byte("mock-screenshot-png-data"),
	}
}

func (m *mockDevice) Platform() core.Platform {
	return m.platform_
}

func (m *mockDevice) Info() (map[string]interface{}, error) {
	if m.shouldFail {
		return nil, m.failError
	}
	return m.info_, nil
}

func (m *mockDevice) Screenshot() ([]byte, error) {
	if m.shouldFail {
		return nil, m.failError
	}
	return m.screenshotData, nil
}

func (m *mockDevice) Tap(x, y int) error {
	if m.shouldFail {
		return m.failError
	}
	m.tapCalls = append(m.tapCalls, tapCall{x: x, y: y})
	return nil
}

func (m *mockDevice) SendKeys(text string) error {
	if m.shouldFail {
		return m.failError
	}
	m.sendKeysCalls = append(m.sendKeysCalls, text)
	return nil
}

func (m *mockDevice) Launch() error {
	if m.shouldFail {
		return m.failError
	}
	m.launchCalled = true
	return nil
}

func (m *mockDevice) Close() error {
	return nil
}

// TestIntegrationParseTranslateExecute tests the complete flow
func TestIntegrationParseTranslateExecute(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Login Flow
tags:
  - login
  - smoke
steps:
  - launch: com.example.app
  - tap:
      x: 100
      y: 200
  - wait: 100
  - inputText:
      text: "user@example.com"
      id: "email_input"
`

	// Parse
	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if flow.Name != "Login Flow" {
		t.Errorf("expected flow name 'Login Flow', got '%s'", flow.Name)
	}

	if len(flow.Steps) != 4 {
		t.Fatalf("expected 4 steps, got %d", len(flow.Steps))
	}

	// Translate
	device := newMockDevice(core.Android)
	translator := NewTranslator()
	actions, err := translator.TranslateFlow(flow, device)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	if len(actions) != 4 {
		t.Fatalf("expected 4 actions, got %d", len(actions))
	}

	// Execute
	executor := NewExecutor(device, "")
	err = executor.Execute(actions, flow.Name)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify mock device was called correctly
	if !device.launchCalled {
		t.Error("expected Launch to be called")
	}

	if len(device.tapCalls) != 1 {
		t.Errorf("expected 1 tap call, got %d", len(device.tapCalls))
	}

	if device.tapCalls[0].x != 100 || device.tapCalls[0].y != 200 {
		t.Errorf("expected tap at (100, 200), got (%d, %d)", device.tapCalls[0].x, device.tapCalls[0].y)
	}

	if len(device.sendKeysCalls) != 1 {
		t.Errorf("expected 1 sendKeys call, got %d", len(device.sendKeysCalls))
	}

	if device.sendKeysCalls[0] != "user@example.com" {
		t.Errorf("expected sendKeys 'user@example.com', got '%s'", device.sendKeysCalls[0])
	}
}

// TestIntegrationIOSFlow tests iOS-specific flow
func TestIntegrationIOSFlow(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: iOS Flow
steps:
  - launch: com.example.app
  - tap:
      x: 50
      y: 100
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	device := newMockDevice(core.IOS)
	translator := NewTranslator()
	actions, err := translator.TranslateFlow(flow, device)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	executor := NewExecutor(device, "")
	err = executor.Execute(actions, flow.Name)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if device.platform_ != core.IOS {
		t.Errorf("expected iOS platform, got %s", device.platform_)
	}
}

// TestIntegrationTapOnCommand tests tapOn with selectors
func TestIntegrationTapOnCommand(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: TapOn Flow
steps:
  - tapOn:
      text: "Login"
  - tapOn:
      id: "submit_button"
      index: 0
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	device := newMockDevice(core.Android)
	translator := NewTranslator()

	// Translate each command individually to test selectors
	for i, cmd := range flow.Steps {
		action, err := translator.TranslateCommand(&cmd, device)
		if err != nil {
			t.Fatalf("Step %d: Translate failed: %v", i, err)
		}

		if action == nil {
			t.Fatalf("Step %d: expected action, got nil", i)
		}
	}
}

// TestIntegrationEmptyFlow tests handling of empty flow
func TestIntegrationEmptyFlow(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Empty Flow
steps: []
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	device := newMockDevice(core.Android)
	translator := NewTranslator()
	actions, err := translator.TranslateFlow(flow, device)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}

	executor := NewExecutor(device, "")
	err = executor.Execute(actions, flow.Name)
	if err != nil {
		t.Fatalf("Execute should not fail for empty flow: %v", err)
	}
}

// TestIntegrationSwipeCommand tests swipe command translation
func TestIntegrationSwipeCommand(t *testing.T) {
	translator := NewTranslator()
	device := newMockDevice(core.Android)

	cmd := &MaestroCommand{
		Swipe: &SwipeCommand{
			StartX:   100,
			StartY:   500,
			EndX:     100,
			EndY:     200,
			Duration: 300,
		},
	}

	action, err := translator.TranslateCommand(cmd, device)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	wrapper, ok := action.(*SwipeWrapper)
	if !ok {
		t.Fatalf("expected *SwipeWrapper, got %T", action)
	}

	if wrapper.startX != 100 || wrapper.startY != 500 {
		t.Errorf("unexpected start position: (%d, %d)", wrapper.startX, wrapper.startY)
	}

	if wrapper.endX != 100 || wrapper.endY != 200 {
		t.Errorf("unexpected end position: (%d, %d)", wrapper.endX, wrapper.endY)
	}
}

// TestIntegrationTerminateCommand tests terminate command
func TestIntegrationTerminateCommand(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Terminate Flow
steps:
  - terminate: com.example.app
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	device := newMockDevice(core.Android)
	translator := NewTranslator()
	actions, err := translator.TranslateFlow(flow, device)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}

	wrapper, ok := actions[0].(*LaunchWrapper)
	if !ok {
		t.Fatalf("expected *LaunchWrapper, got %T", actions[0])
	}

	if wrapper.appID != "com.example.app" {
		t.Errorf("expected appID 'com.example.app', got '%s'", wrapper.appID)
	}

	if wrapper.waitIdle {
		t.Error("expected waitIdle to be false for terminate")
	}
}

// TestIntegrationParseFromFile tests parsing from file path
func TestIntegrationParseFromFile(t *testing.T) {
	flow, err := ParseMaestroFlow("../examples/sample.maestro.yaml")
	if err != nil {
		t.Fatalf("ParseMaestroFlow failed: %v", err)
	}

	if flow.Name == "" {
		t.Error("expected flow name to be set")
	}

	if flow.AppID == "" {
		t.Error("expected appId to be set")
	}

	if len(flow.Steps) == 0 {
		t.Error("expected at least one step")
	}
}

// TestIntegrationErrorRecoveryInvalidSelector tests error handling
func TestIntegrationErrorRecoveryInvalidSelector(t *testing.T) {
	// Test with unsupported command - should return ErrUnsupportedCommand
	yamlContent := `
appId: com.example.app
name: Invalid Selector Flow
steps:
  - unsupportedCommand: true
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse should not fail: %v", err)
	}

	device := newMockDevice(core.Android)
	translator := NewTranslator()

	_, err = translator.TranslateFlow(flow, device)
	if err == nil {
		t.Fatal("expected error for unsupported command")
	}

	if !IsUnsupportedCommand(err) {
		t.Errorf("expected ErrUnsupportedCommand, got: %v", err)
	}
}

// TestIntegrationErrorRecoveryDeviceFailure tests device failure handling
func TestIntegrationErrorRecoveryDeviceFailure(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Device Failure Flow
steps:
  - launch: com.example.app
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Setup device to fail
	device := newMockDevice(core.Android)
	device.shouldFail = true
	device.failError = &mockError{msg: "device connection lost"}

	translator := NewTranslator()
	actions, err := translator.TranslateFlow(flow, device)
	if err != nil {
		t.Fatalf("Translate should not fail: %v", err)
	}

	executor := NewExecutor(device, "")
	err = executor.Execute(actions, flow.Name)
	if err == nil {
		t.Fatal("expected error when device fails")
	}

	// Verify error contains step information
	if !strings.Contains(err.Error(), "step 1 failed") {
		t.Errorf("expected error to mention step 1 failure, got: %v", err)
	}
}

// TestIntegrationWaitCommand tests wait command
func TestIntegrationWaitCommand(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Wait Flow
steps:
  - wait: 50
  - launch: com.example.app
  - wait: 100
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	device := newMockDevice(core.Android)
	translator := NewTranslator()
	actions, err := translator.TranslateFlow(flow, device)
	if err != nil {
		t.Fatalf("Translate failed: %v", err)
	}

	if len(actions) != 3 {
		t.Fatalf("expected 3 actions, got %d", len(actions))
	}

	// Execute with timing
	executor := NewExecutor(device, "")
	err = executor.Execute(actions, flow.Name)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
}

// TestIntegrationMultipleFlows tests multiple flows in sequence
func TestIntegrationMultipleFlows(t *testing.T) {
	flows := []string{
		`
appId: com.example.app
name: Flow 1
steps:
  - launch: com.example.app
`,
		`
appId: com.example.app
name: Flow 2
steps:
  - tap:
      x: 100
      y: 100
`,
	}

	device := newMockDevice(core.Android)
	translator := NewTranslator()
	executor := NewExecutor(device, "")

	for i, yamlContent := range flows {
		flow, err := ParseFlowFromString(yamlContent)
		if err != nil {
			t.Fatalf("Flow %d: Parse failed: %v", i+1, err)
		}

		actions, err := translator.TranslateFlow(flow, device)
		if err != nil {
			t.Fatalf("Flow %d: Translate failed: %v", i+1, err)
		}

		err = executor.Execute(actions, flow.Name)
		if err != nil {
			t.Fatalf("Flow %d: Execute failed: %v", i+1, err)
		}
	}

	// Verify all actions were executed
	if !device.launchCalled {
		t.Error("expected Launch to be called for Flow 1")
	}

	if len(device.tapCalls) != 1 {
		t.Errorf("expected 1 tap call from Flow 2, got %d", len(device.tapCalls))
	}
}

// TestIntegrationTagsParsing tests tag parsing
func TestIntegrationTagsParsing(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Tagged Flow
tags:
  - smoke
  - regression
  - login
steps:
  - launch: com.example.app
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(flow.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(flow.Tags))
	}

	expectedTags := []string{"smoke", "regression", "login"}
	for i, tag := range expectedTags {
		if flow.Tags[i] != tag {
			t.Errorf("expected tag '%s' at index %d, got '%s'", tag, i, flow.Tags[i])
		}
	}
}

// mockError for testing error propagation
type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}
