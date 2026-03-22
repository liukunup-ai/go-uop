package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/liukunup/go-uop/maestro"
)

func TestTranslateRunFlowWithPath(t *testing.T) {
	translator := NewFlowCommandTranslator()

	tmpDir := t.TempDir()
	subflowPath := filepath.Join(tmpDir, "subflow.yaml")
	if err := os.WriteFile(subflowPath, []byte("name: test\nsteps: []"), 0644); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	cmd := &maestro.RunFlowCommand{
		Path: subflowPath,
	}

	act, err := translator.TranslateRunFlow(cmd, nil, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if act == nil {
		t.Fatal("expected RunFlowAction, got nil")
	}

	if act.SubflowPath != subflowPath {
		t.Errorf("expected SubflowPath %s, got %s", subflowPath, act.SubflowPath)
	}

	if act.Depth != 1 {
		t.Errorf("expected Depth 1, got %d", act.Depth)
	}
}

func TestTranslateRunFlowWithEnvVars(t *testing.T) {
	translator := NewFlowCommandTranslator()

	tmpDir := t.TempDir()
	subflowPath := filepath.Join(tmpDir, "login.yaml")
	if err := os.WriteFile(subflowPath, []byte("name: login\nsteps: []"), 0644); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	parentVars := map[string]string{
		"base_url": "https://example.com",
	}

	cmd := &maestro.RunFlowCommand{
		Path: subflowPath,
		Params: map[string]string{
			"user":     "testuser",
			"password": "secret",
		},
	}

	act, err := translator.TranslateRunFlow(cmd, parentVars, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if act == nil {
		t.Fatal("expected RunFlowAction, got nil")
	}

	if act.EnvVars["base_url"] != "https://example.com" {
		t.Errorf("expected parent var base_url inherited, got %s", act.EnvVars["base_url"])
	}

	if act.EnvVars["user"] != "testuser" {
		t.Errorf("expected EnvVars user=testuser, got %s", act.EnvVars["user"])
	}

	if act.EnvVars["password"] != "secret" {
		t.Errorf("expected EnvVars password=secret, got %s", act.EnvVars["password"])
	}
}

func TestTranslateRunFlowWithName(t *testing.T) {
	translator := NewFlowCommandTranslator()

	tmpDir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	os.Chdir(tmpDir)

	subflowPath := filepath.Join(tmpDir, "login.yaml")
	if err := os.WriteFile(subflowPath, []byte("name: login\nsteps: []"), 0644); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	cmd := &maestro.RunFlowCommand{
		Name: "login",
	}

	act, err := translator.TranslateRunFlow(cmd, nil, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if act == nil {
		t.Fatal("expected RunFlowAction, got nil")
	}

	if !strings.HasSuffix(act.SubflowPath, "login.yaml") {
		t.Errorf("expected SubflowPath ending with login.yaml, got %s", act.SubflowPath)
	}
}

func TestTranslateRunFlowNilCommand(t *testing.T) {
	translator := NewFlowCommandTranslator()

	_, err := translator.TranslateRunFlow(nil, nil, 0)
	if err == nil {
		t.Error("expected error for nil command")
	}
}

func TestTranslateRunFlowNoPathOrName(t *testing.T) {
	translator := NewFlowCommandTranslator()

	cmd := &maestro.RunFlowCommand{}

	_, err := translator.TranslateRunFlow(cmd, nil, 0)
	if err == nil {
		t.Error("expected error for missing path and name")
	}
}

func TestTranslateRunFlowMaxDepth(t *testing.T) {
	translator := NewFlowCommandTranslator()

	cmd := &maestro.RunFlowCommand{
		Path: "subflow.yaml",
	}

	_, err := translator.TranslateRunFlow(cmd, nil, maxFlowDepth)
	if err == nil {
		t.Error("expected error when max depth exceeded")
	}
}

func TestExecuteRunFlowSubflowInheritance(t *testing.T) {
	translator := NewFlowCommandTranslator()

	tmpDir := t.TempDir()
	subflowPath := filepath.Join(tmpDir, "subflow.yaml")
	subflowContent := `
name: subflow
steps:
  - wait: 100
`
	if err := os.WriteFile(subflowPath, []byte(subflowContent), 0644); err != nil {
		t.Fatalf("create test file: %v", err)
	}

	cmd := &maestro.RunFlowCommand{
		Path: subflowPath,
		Params: map[string]string{
			"user": "testuser",
		},
	}

	act, err := translator.TranslateRunFlow(cmd, nil, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := translator.ExecuteRunFlow(act, nil, nil); err != nil {
		t.Fatalf("execute runFlow failed: %v", err)
	}
}

func TestResolveSubflowPath(t *testing.T) {
	translator := NewFlowCommandTranslator()

	absPath := "/absolute/path.yaml"
	result := translator.ResolveSubflowPath(absPath, "/base")
	if result != absPath {
		t.Errorf("expected absolute path unchanged, got %s", result)
	}

	relPath := "relative.yaml"
	baseDir := "/base"
	expected := "/base/relative.yaml"
	result = translator.ResolveSubflowPath(relPath, baseDir)
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
