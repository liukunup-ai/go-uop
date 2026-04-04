package runner

import (
	"strings"
	"testing"

	"github.com/liukunup/go-uop/internal/report"
)

// TestExecutor_Creation tests that an Executor can be created
func TestExecutor_Creation(t *testing.T) {
	pool := NewDevicePool()
	reportGen := report.NewGenerator("test-suite")

	executor := NewExecutor(pool, reportGen)

	if executor == nil {
		t.Fatal("NewExecutor returned nil")
	}

	if executor.pool != pool {
		t.Error("Executor pool not set correctly")
	}

	if executor.reportGen != reportGen {
		t.Error("Executor reportGen not set correctly")
	}

	if executor.executors == nil {
		t.Error("Executor executors map is nil")
	}

	// Verify command registry has expected commands
	expectedCommands := []string{"launch", "tapOn", "inputText", "swipe", "pressKey", "wait", "screenshot", "device"}
	for _, cmd := range expectedCommands {
		if _, exists := executor.executors[cmd]; !exists {
			t.Errorf("Expected command %q not found in registry", cmd)
		}
	}
}

// TestExecutor_ExecuteFlow_Basic tests basic flow execution
func TestExecutor_ExecuteFlow_Basic(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	// Add mock device to pool
	err := pool.AddDevice("test-device-1", "ios", "test-serial-1")
	if err != nil {
		t.Fatalf("Failed to add device to pool: %v", err)
	}

	reportGen := report.NewGenerator("test-suite")
	executor := NewExecutor(pool, reportGen)

	// Define a simple flow
	flow := &Flow{
		Name: "test-flow",
		Steps: []Step{
			{
				"tapOn": map[string]interface{}{
					"x": 100,
					"y": 200,
				},
			},
		},
	}

	// Execute the flow
	err = executor.ExecuteFlow(flow)
	if err != nil {
		t.Fatalf("ExecuteFlow failed: %v", err)
	}
}

// TestExecutor_ExecuteFlow_WithDeviceSwitch tests device switching
func TestExecutor_ExecuteFlow_WithDeviceSwitch(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	// Add two devices to pool
	err := pool.AddDevice("device1", "ios", "serial1")
	if err != nil {
		t.Fatalf("Failed to add device1: %v", err)
	}
	err = pool.AddDevice("device2", "android", "serial2")
	if err != nil {
		t.Fatalf("Failed to add device2: %v", err)
	}

	reportGen := report.NewGenerator("test-suite")
	executor := NewExecutor(pool, reportGen)

	// Flow that switches devices
	flow := &Flow{
		Name: "switch-device-flow",
		Steps: []Step{
			{
				"device": "device2",
			},
		},
	}

	err = executor.ExecuteFlow(flow)
	if err != nil {
		t.Fatalf("ExecuteFlow with device switch failed: %v", err)
	}

	// Verify device switched
	current := pool.CurrentDevice()
	if current.ID != "device2" {
		t.Errorf("Expected current device to be 'device2', got '%s'", current.ID)
	}
}

// TestExecutor_ExecuteFlow_ParsesYAML tests that YAML parsing works with executor
func TestExecutor_ExecuteFlow_ParsesYAML(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	// Add a device to pool
	err := pool.AddDevice("test-device", "ios", "test-serial")
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}

	reportGen := report.NewGenerator("test-suite")
	executor := NewExecutor(pool, reportGen)

	// Parse YAML flow
	yamlContent := `
name: yaml-test-flow
steps:
  - launch: com.example.app
  - tapOn:
      x: 100
      y: 200
  - inputText:
      text: "hello world"
`
	flow, err := ParseFlow(strings.NewReader(yamlContent))
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	if flow.Name != "yaml-test-flow" {
		t.Errorf("Expected flow name 'yaml-test-flow', got '%s'", flow.Name)
	}

	if len(flow.Steps) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(flow.Steps))
	}

	// Execute the flow
	err = executor.ExecuteFlow(flow)
	if err != nil {
		t.Fatalf("ExecuteFlow failed: %v", err)
	}
}

// TestExecutor_ReportIntegration tests that executor reports to Generator
func TestExecutor_ReportIntegration(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	// Add device
	err := pool.AddDevice("report-test-device", "ios", "report-serial")
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}

	reportGen := report.NewGenerator("report-test")
	executor := NewExecutor(pool, reportGen)

	flow := &Flow{
		Name: "report-test-flow",
		Steps: []Step{
			{
				"tapOn": map[string]interface{}{
					"x": 50,
					"y": 50,
				},
			},
		},
	}

	reportGen.StartTest(flow.Name)
	err = executor.ExecuteFlow(flow)
	reportGen.EndTest("passed", err)

	// Verify report has step results
	result := reportGen.Generate()
	if len(result.Results) != 1 {
		t.Fatalf("Expected 1 test result, got %d", len(result.Results))
	}

	if len(result.Results[0].Steps) == 0 {
		t.Error("Expected step results in report")
	}
}

// TestExecutor_CommandRegistry tests that all expected commands are registered
func TestExecutor_CommandRegistry(t *testing.T) {
	pool := NewDevicePool()
	reportGen := report.NewGenerator("registry-test")

	executor := NewExecutor(pool, reportGen)

	// Check all expected commands exist
	expectedCommands := map[string]bool{
		"launch":     true,
		"tapOn":      true,
		"inputText":  true,
		"swipe":      true,
		"pressKey":   true,
		"wait":       true,
		"screenshot": true,
		"device":     true,
	}

	for cmd := range expectedCommands {
		if _, exists := executor.executors[cmd]; !exists {
			t.Errorf("Command %q not registered", cmd)
		}
	}

	// Verify no extra commands
	if len(executor.executors) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(executor.executors))
	}
}

// TestExecutor_StepExecution tests step execution order
func TestExecutor_StepExecution(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	// Add device
	err := pool.AddDevice("order-test-device", "ios", "order-serial")
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}

	reportGen := report.NewGenerator("step-order-test")
	executor := NewExecutor(pool, reportGen)

	// Create flow with multiple steps
	flow := &Flow{
		Name: "step-order-flow",
		Steps: []Step{
			{"launch": nil},
			{"tapOn": map[string]interface{}{"x": 10, "y": 10}},
			{"inputText": map[string]interface{}{"text": "test"}},
		},
	}

	err = executor.ExecuteFlow(flow)
	if err != nil {
		t.Fatalf("ExecuteFlow failed: %v", err)
	}
}

// TestExecutor_DeviceKeyInStep tests that device key in step switches device
func TestExecutor_DeviceKeyInStep(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	// Add two devices
	err := pool.AddDevice("phone", "ios", "phone-serial")
	if err != nil {
		t.Fatalf("Failed to add phone: %v", err)
	}
	err = pool.AddDevice("tablet", "android", "tablet-serial")
	if err != nil {
		t.Fatalf("Failed to add tablet: %v", err)
	}

	reportGen := report.NewGenerator("device-switch-test")
	executor := NewExecutor(pool, reportGen)

	// Flow with device switch and subsequent commands
	flow := &Flow{
		Name: "device-switch-flow",
		Steps: []Step{
			{"device": "tablet"},
			{"tapOn": map[string]interface{}{"x": 100, "y": 200}},
		},
	}

	err = executor.ExecuteFlow(flow)
	if err != nil {
		t.Fatalf("ExecuteFlow failed: %v", err)
	}

	// Verify device switched to tablet
	current := pool.CurrentDevice()
	if current.ID != "tablet" {
		t.Errorf("Expected current device to be 'tablet', got '%s'", current.ID)
	}
}

// TestExecutor_MultipleFlows tests executing multiple flows
func TestExecutor_MultipleFlows(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("test-device", "ios", "test-serial")
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}

	reportGen := report.NewGenerator("multi-flow-test")
	executor := NewExecutor(pool, reportGen)

	// First flow
	flow1 := &Flow{
		Name: "flow-1",
		Steps: []Step{
			{"launch": nil},
		},
	}

	err = executor.ExecuteFlow(flow1)
	if err != nil {
		t.Fatalf("ExecuteFlow flow1 failed: %v", err)
	}

	// Second flow
	flow2 := &Flow{
		Name: "flow-2",
		Steps: []Step{
			{"tapOn": map[string]interface{}{"x": 50, "y": 50}},
		},
	}

	err = executor.ExecuteFlow(flow2)
	if err != nil {
		t.Fatalf("ExecuteFlow flow2 failed: %v", err)
	}
}
