package yaml

import (
	"testing"
)

func TestParseFlowFromString_Basic(t *testing.T) {
	yaml := `
name: login
steps:
  - launch: com.example.app
  - tapOn:
      text: "登录"
`
	flow, err := ParseFlowFromString(yaml)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flow.Name != "login" {
		t.Errorf("expected name 'login', got '%s'", flow.Name)
	}

	if len(flow.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(flow.Steps))
	}

	if flow.Steps[0].Launch != "com.example.app" {
		t.Errorf("expected launch 'com.example.app', got '%s'", flow.Steps[0].Launch)
	}
}

func TestParseFlowFromString_WithParams(t *testing.T) {
	yaml := `
name: login
params:
  username: string
  password: string
steps:
  - inputText:
      text: "${params.username}"
`
	flow, err := ParseFlowFromString(yaml)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if flow.Params["username"] != "string" {
		t.Errorf("expected param 'username', got '%s'", flow.Params["username"])
	}
}

func TestParseFlowFromString_WithControlFlow(t *testing.T) {
	yaml := `
name: test
steps:
  - foreach:
      variable: item
      in: "a,b,c"
      do:
        - tapOn:
            text: "${item}"
`
	flow, err := ParseFlowFromString(yaml)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(flow.Steps) != 1 {
		t.Errorf("expected 1 step, got %d", len(flow.Steps))
	}

	if flow.Steps[0].Foreach == nil {
		t.Fatal("expected foreach command")
	}

	if flow.Steps[0].Foreach.Variable != "item" {
		t.Errorf("expected variable 'item', got '%s'", flow.Steps[0].Foreach.Variable)
	}
}

func TestParseFlowFromString_WithWait(t *testing.T) {
	yaml := `
name: wait
steps:
  - wait: 2000
`
	flow, err := ParseFlowFromString(yaml)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(flow.Steps) != 1 {
		t.Errorf("expected 1 step, got %d", len(flow.Steps))
	}

	if flow.Steps[0].Wait != 2000 {
		t.Errorf("expected wait 2000, got %d", flow.Steps[0].Wait)
	}
}

func TestParseFlowFromString_InvalidYAML(t *testing.T) {
	yaml := `
name: test
  invalid: yaml
    structure: broken
`
	_, err := ParseFlowFromString(yaml)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
