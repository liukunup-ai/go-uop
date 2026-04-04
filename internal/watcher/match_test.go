package watcher

import (
	"context"
	"github.com/liukunup/go-uop/core"
	"testing"
)

// Mock device for testing
type mockDevice struct{}

func (m *mockDevice) Platform() core.Platform               { return core.IOS }
func (m *mockDevice) Info() (map[string]interface{}, error) { return nil, nil }
func (m *mockDevice) Screenshot() ([]byte, error)           { return []byte("screenshot"), nil }
func (m *mockDevice) Tap(x, y int) error                    { return nil }
func (m *mockDevice) SendKeys(text string) error            { return nil }
func (m *mockDevice) Launch() error                         { return nil }
func (m *mockDevice) PressKey(code int) error               { return nil }
func (m *mockDevice) GetAlertText() (string, error)         { return "", nil }
func (m *mockDevice) AcceptAlert() error                    { return nil }
func (m *mockDevice) DismissAlert() error                   { return nil }
func (m *mockDevice) Close() error                          { return nil }

var _ core.Device = (*mockDevice)(nil)

func TestImageMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	m := NewImageMatch("test_template.png", 0.8)
	matched, err := m.Match(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = matched
}

func TestTextMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	m := NewTextMatch("确定")
	matched, err := m.Match(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if matched {
		t.Log("TextMatch: popup detected")
	}
}

func TestRegexMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	m := NewRegexMatch("版本.*更新")
	matched, err := m.Match(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = matched
}

func TestCompoundMatch(t *testing.T) {
	ctx := context.Background()
	device := &mockDevice{}
	m := NewCompoundMatch("or", []MatchCondition{
		NewTextMatch("确定"),
		NewTextMatch("取消"),
	})
	matched, err := m.Match(ctx, device)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = matched
}
