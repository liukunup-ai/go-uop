package commands

import (
	"testing"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/yaml"
)

type mockDevice struct{}

func (m *mockDevice) Platform() core.Platform               { return core.IOS }
func (m *mockDevice) Info() (map[string]interface{}, error) { return nil, nil }
func (m *mockDevice) Screenshot() ([]byte, error)           { return nil, nil }
func (m *mockDevice) Tap(x, y int) error                    { return nil }
func (m *mockDevice) SendKeys(text string) error            { return nil }
func (m *mockDevice) Launch() error                         { return nil }
func (m *mockDevice) Close() error                          { return nil }

func TestExecutor_ExecuteIf_True(t *testing.T) {
	exec := NewExecutor(&mockDevice{}, map[string]interface{}{})

	cmd := yaml.IfCommand{
		Condition: "true",
		Then: []yaml.Command{
			{Log: "then branch"},
		},
		Else: []yaml.Command{
			{Log: "else branch"},
		},
	}

	err := exec.executeIf(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExecutor_ExecuteIf_False(t *testing.T) {
	exec := NewExecutor(&mockDevice{}, map[string]interface{}{})

	cmd := yaml.IfCommand{
		Condition: "false",
		Then: []yaml.Command{
			{Log: "then branch"},
		},
		Else: []yaml.Command{
			{Log: "else branch"},
		},
	}

	err := exec.executeIf(cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExecutor_isTrue(t *testing.T) {
	exec := NewExecutor(&mockDevice{}, nil)

	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"1", true},
		{"yes", true},
		{"y", true},
		{"false", false},
		{"0", false},
		{"no", false},
		{"", false},
		{"maybe", false},
	}

	for _, tc := range tests {
		result := exec.isTrue(tc.input)
		if result != tc.expected {
			t.Errorf("isTrue(%q) = %v, want %v", tc.input, result, tc.expected)
		}
	}
}
