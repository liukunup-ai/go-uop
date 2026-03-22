package action

import (
	"testing"
	"time"
)

func TestTapAction_Struct(t *testing.T) {
	ta := TapAction{
		X: 100,
		Y: 200,
	}

	if ta.X != 100 {
		t.Errorf("expected X 100, got %d", ta.X)
	}
	if ta.Y != 200 {
		t.Errorf("expected Y 200, got %d", ta.Y)
	}
}

func TestTapAction_WithElement(t *testing.T) {
	ta := TapAction{
		X:       100,
		Y:       200,
		Element: nil,
	}

	if ta.Element != nil {
		t.Error("expected Element to be nil")
	}
}

func TestSwipeAction_Struct(t *testing.T) {
	sa := SwipeAction{
		StartX:   100,
		StartY:   200,
		EndX:     300,
		EndY:     400,
		Duration: 500 * time.Millisecond,
	}

	if sa.StartX != 100 {
		t.Errorf("expected StartX 100, got %d", sa.StartX)
	}
	if sa.EndY != 400 {
		t.Errorf("expected EndY 400, got %d", sa.EndY)
	}
	if sa.Duration != 500*time.Millisecond {
		t.Errorf("expected Duration 500ms, got %v", sa.Duration)
	}
}

func TestSendKeysAction_Struct(t *testing.T) {
	sk := SendKeysAction{
		Text:   "hello world",
		Secure: false,
	}

	if sk.Text != "hello world" {
		t.Errorf("expected Text 'hello world', got '%s'", sk.Text)
	}
	if sk.Secure != false {
		t.Error("expected Secure false")
	}
}

func TestSendKeysAction_Secure(t *testing.T) {
	sk := SendKeysAction{
		Text:   "password123",
		Secure: true,
	}

	if !sk.Secure {
		t.Error("expected Secure true")
	}
}

func TestLaunchAction_Struct(t *testing.T) {
	la := LaunchAction{
		AppID:     "com.example.app",
		Arguments: []string{"--debug"},
		WaitIdle:  true,
	}

	if la.AppID != "com.example.app" {
		t.Errorf("expected AppID 'com.example.app', got '%s'", la.AppID)
	}
	if len(la.Arguments) != 1 {
		t.Errorf("expected 1 argument, got %d", len(la.Arguments))
	}
	if !la.WaitIdle {
		t.Error("expected WaitIdle true")
	}
}

func TestPressKeyAction_Struct(t *testing.T) {
	pk := PressKeyAction{
		KeyCode: 3,
	}

	if pk.KeyCode != 3 {
		t.Errorf("expected KeyCode 3, got %d", pk.KeyCode)
	}
}

func TestWaitAction_Struct(t *testing.T) {
	wa := WaitAction{
		Duration: 2 * time.Second,
		Optional: false,
	}

	if wa.Duration != 2*time.Second {
		t.Errorf("expected Duration 2s, got %v", wa.Duration)
	}
	if wa.Optional != false {
		t.Error("expected Optional false")
	}
}

func TestWaitAction_WithElement(t *testing.T) {
	wa := WaitAction{
		Duration: 5 * time.Second,
		Element:  nil,
		Optional: true,
	}

	if wa.Element != nil {
		t.Error("expected Element to be nil in this test")
	}
	if !wa.Optional {
		t.Error("expected Optional true")
	}
}
