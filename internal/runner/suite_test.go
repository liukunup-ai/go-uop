package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/liukunup/go-uop/internal/report"
)

func TestSuiteRunner_Creation(t *testing.T) {
	pool := NewDevicePool()
	reportGen := report.NewGenerator("test-suite")
	runner := NewSuiteRunner(pool, reportGen)

	if runner == nil {
		t.Fatal("NewSuiteRunner returned nil")
	}

	if runner.pool != pool {
		t.Error("SuiteRunner pool not set correctly")
	}

	if runner.executor == nil {
		t.Error("SuiteRunner executor is nil")
	}

	if runner.reportGen != reportGen {
		t.Error("SuiteRunner reportGen not set correctly")
	}
}

func TestSuiteRunner_RunSuite(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("test-device", "ios", "test-serial")
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}

	reportGen := report.NewGenerator("test-suite")
	runner := NewSuiteRunner(pool, reportGen)

	tmpDir := t.TempDir()
	flowPath := filepath.Join(tmpDir, "test-flow.yaml")
	if err := os.WriteFile(flowPath, []byte(`
name: test-flow
steps:
  - tapOn:
      x: 100
      y: 200
`), 0644); err != nil {
		t.Fatalf("Failed to write flow file: %v", err)
	}

	suite := &Suite{
		Name: "test-suite",
		Flows: []SuiteFlow{
			{Name: "flow-1", Path: flowPath},
		},
	}

	result, err := runner.RunSuite(suite)
	if err != nil {
		t.Fatalf("RunSuite failed: %v", err)
	}

	if len(result.FlowResults) != 1 {
		t.Errorf("Expected 1 flow result, got %d", len(result.FlowResults))
	}

	if result.TotalSteps != 1 {
		t.Errorf("Expected 1 total step, got %d", result.TotalSteps)
	}
}

func TestSuiteRunner_MultipleFlows(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("test-device", "ios", "test-serial")
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}

	reportGen := report.NewGenerator("multi-flow-suite")
	runner := NewSuiteRunner(pool, reportGen)

	tmpDir := t.TempDir()
	flow1Path := filepath.Join(tmpDir, "flow1.yaml")
	flow2Path := filepath.Join(tmpDir, "flow2.yaml")

	if err := os.WriteFile(flow1Path, []byte(`
name: flow-1
steps:
  - tapOn:
      x: 100
      y: 200
`), 0644); err != nil {
		t.Fatalf("Failed to write flow1: %v", err)
	}

	if err := os.WriteFile(flow2Path, []byte(`
name: flow-2
steps:
  - launch: com.example.app
  - inputText:
      text: hello
`), 0644); err != nil {
		t.Fatalf("Failed to write flow2: %v", err)
	}

	suite := &Suite{
		Name: "multi-flow-suite",
		Flows: []SuiteFlow{
			{Name: "flow-1", Path: flow1Path},
			{Name: "flow-2", Path: flow2Path},
		},
	}

	result, err := runner.RunSuite(suite)
	if err != nil {
		t.Fatalf("RunSuite failed: %v", err)
	}

	if len(result.FlowResults) != 2 {
		t.Errorf("Expected 2 flow results, got %d", len(result.FlowResults))
	}

	if result.TotalSteps != 3 {
		t.Errorf("Expected 3 total steps, got %d", result.TotalSteps)
	}

	if result.PassedSteps != 3 {
		t.Errorf("Expected 3 passed steps, got %d", result.PassedSteps)
	}
}

func TestParseAndRunSuite(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("test-device", "ios", "test-serial")
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}

	reportGen := report.NewGenerator("parse-suite")
	tmpDir := t.TempDir()
	flowPath := filepath.Join(tmpDir, "test-flow.yaml")

	if err := os.WriteFile(flowPath, []byte(`
name: parse-test-flow
steps:
  - tapOn:
      x: 50
      y: 50
`), 0644); err != nil {
		t.Fatalf("Failed to write flow: %v", err)
	}

	yamlContent := `
name: parse-test-suite
flows:
  - name: flow-1
    path: ` + flowPath

	result, err := ParseAndRunSuite(strings.NewReader(yamlContent), pool, reportGen)
	if err != nil {
		t.Fatalf("ParseAndRunSuite failed: %v", err)
	}

	if len(result.FlowResults) != 1 {
		t.Errorf("Expected 1 flow result, got %d", len(result.FlowResults))
	}
}

func TestParseAndRunSuiteFile(t *testing.T) {
	factory := newMockDeviceFactory()
	pool := newDevicePoolWithFactory(factory)

	err := pool.AddDevice("test-device", "ios", "test-serial")
	if err != nil {
		t.Fatalf("Failed to add device: %v", err)
	}

	reportGen := report.NewGenerator("file-suite")
	tmpDir := t.TempDir()
	flowPath := filepath.Join(tmpDir, "test-flow.yaml")
	suitePath := filepath.Join(tmpDir, "suite.yaml")

	if err := os.WriteFile(flowPath, []byte(`
name: file-test-flow
steps:
  - launch: com.example.app
`), 0644); err != nil {
		t.Fatalf("Failed to write flow: %v", err)
	}

	if err := os.WriteFile(suitePath, []byte(`
name: file-test-suite
flows:
  - name: flow-1
    path: `+flowPath), 0644); err != nil {
		t.Fatalf("Failed to write suite: %v", err)
	}

	result, err := ParseAndRunSuiteFile(suitePath, pool, reportGen)
	if err != nil {
		t.Fatalf("ParseAndRunSuiteFile failed: %v", err)
	}

	if len(result.FlowResults) != 1 {
		t.Errorf("Expected 1 flow result, got %d", len(result.FlowResults))
	}
}

func TestSuiteResult_Aggregation(t *testing.T) {
	result := &SuiteResult{
		FlowResults: []*FlowResult{
			{FlowName: "flow-1", Status: "passed", Steps: 2},
			{FlowName: "flow-2", Status: "failed", Steps: 3},
			{FlowName: "flow-3", Status: "passed", Steps: 1},
		},
	}

	for _, fr := range result.FlowResults {
		result.TotalSteps += fr.Steps
		if fr.Status == "passed" {
			result.PassedSteps += fr.Steps
		} else {
			result.FailedSteps += fr.Steps
		}
	}

	if result.TotalSteps != 6 {
		t.Errorf("Expected 6 total steps, got %d", result.TotalSteps)
	}

	if result.PassedSteps != 3 {
		t.Errorf("Expected 3 passed steps, got %d", result.PassedSteps)
	}

	if result.FailedSteps != 3 {
		t.Errorf("Expected 3 failed steps, got %d", result.FailedSteps)
	}
}
