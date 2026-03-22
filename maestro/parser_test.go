package maestro

import (
	"strings"
	"testing"
)

func TestParseFlowValid(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Login Flow
tags:
  - login
  - smoke
steps:
  - launch: com.example.app
  - tapOn:
      text: "Login"
  - inputText:
      text: "user@example.com"
  - tapOn:
      id: "submit_btn"
      index: 0
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if flow.AppID != "com.example.app" {
		t.Errorf("expected AppID 'com.example.app', got '%s'", flow.AppID)
	}

	if flow.Name != "Login Flow" {
		t.Errorf("expected Name 'Login Flow', got '%s'", flow.Name)
	}

	if len(flow.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(flow.Tags))
	}

	if len(flow.Steps) != 4 {
		t.Errorf("expected 4 steps, got %d", len(flow.Steps))
	}
}

func TestParseFlowFromReader(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Test Flow
steps:
  - launch: com.example.app
`

	reader := strings.NewReader(yamlContent)
	flow, err := ParseFlow(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if flow.AppID != "com.example.app" {
		t.Errorf("expected AppID 'com.example.app', got '%s'", flow.AppID)
	}

	if len(flow.Steps) != 1 {
		t.Errorf("expected 1 step, got %d", len(flow.Steps))
	}
}

func TestParseFlowInvalid(t *testing.T) {
	invalidYAML := `
appId: com.example.app
name: Invalid Flow
steps:
  - invalidCommand:
      this is not valid yaml: [
`

	_, err := ParseFlowFromString(invalidYAML)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}

	if !strings.Contains(err.Error(), "parsing failed") {
		t.Errorf("expected error to contain 'parsing failed', got: %v", err)
	}
}

func TestParseFlowEmpty(t *testing.T) {
	_, err := ParseFlowFromString("")
	if err == nil {
		t.Error("expected error for empty YAML, got nil")
	}
}

func TestParseFlowDocumentSeparator(t *testing.T) {
	yamlContent := `---
appId: com.example.app
name: Multi Doc Flow
steps:
  - launch: com.example.app
---
appId: another.app
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if flow.AppID != "com.example.app" {
		t.Errorf("expected AppID from first document 'com.example.app', got '%s'", flow.AppID)
	}
}

func TestParseShorthandSelectors(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Shorthand Test
steps:
  - tapOn:
      text: "Login Extended"
  - tapOn:
      id: "btn"
      index: 1
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(flow.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(flow.Steps))
	}

	if flow.Steps[0].TapOn == nil || flow.Steps[0].TapOn.Text != "Login Extended" {
		t.Error("expected tapOn with text 'Login Extended'")
	}

	if flow.Steps[1].TapOn == nil || flow.Steps[1].TapOn.ID != "btn" || flow.Steps[1].TapOn.Index != 1 {
		t.Error("expected tapOn with id 'btn' and index 1")
	}
}

func TestParseFlowWithoutAppId(t *testing.T) {
	yamlContent := `
name: No AppId Flow
steps:
  - launch: com.example.app
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if flow.AppID != "" {
		t.Errorf("expected empty AppID, got '%s'", flow.AppID)
	}

	if flow.Name != "No AppId Flow" {
		t.Errorf("expected Name 'No AppId Flow', got '%s'", flow.Name)
	}
}

func TestParseFlowVariousCommands(t *testing.T) {
	yamlContent := `
appId: com.example.app
name: Various Commands Test
steps:
  - launch: com.example.app
  - terminate: com.example.app
  - tapOn:
      text: "Button"
  - inputText:
      text: "Hello World"
  - wait: 1000
  - screenshot:
      name: "test_screenshot"
  - back:
`

	flow, err := ParseFlowFromString(yamlContent)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(flow.Steps) != 7 {
		t.Errorf("expected 7 steps, got %d", len(flow.Steps))
	}

	if flow.Steps[0].Launch != "com.example.app" {
		t.Error("expected first step to be launch command")
	}

	if flow.Steps[2].TapOn == nil || flow.Steps[2].TapOn.Text != "Button" {
		t.Error("expected tapOn with text 'Button'")
	}

	if flow.Steps[4].Wait != 1000 {
		t.Error("expected wait 1000")
	}
}
