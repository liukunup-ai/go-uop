package maestro

import (
	"testing"

	"github.com/liukunup/go-uop/internal/selector"
)

func TestTranslateTapOn(t *testing.T) {
	translator := NewTranslator()

	cmd := &MaestroCommand{
		TapOn: &TapOnCommand{
			Text: "Login",
		},
	}

	act, err := translator.TranslateCommand(cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrapper, ok := act.(*TapOnWrapper)
	if !ok {
		t.Fatalf("expected *tapOnWrapper, got %T", act)
	}

	if wrapper.element == nil {
		t.Fatal("expected element to be set")
	}

	if wrapper.element.Type != selector.SelectorTypeText {
		t.Errorf("expected SelectorTypeText, got %v", wrapper.element.Type)
	}

	if wrapper.element.Value != "Login" {
		t.Errorf("expected 'Login', got %v", wrapper.element.Value)
	}
}

func TestTranslateTapOnWithID(t *testing.T) {
	translator := NewTranslator()

	cmd := &MaestroCommand{
		TapOn: &TapOnCommand{
			ID: "com.example.app:id/button",
		},
	}

	act, err := translator.TranslateCommand(cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrapper, ok := act.(*TapOnWrapper)
	if !ok {
		t.Fatalf("expected *tapOnWrapper, got %T", act)
	}

	if wrapper.element == nil {
		t.Fatal("expected element to be set")
	}

	if wrapper.element.Type != selector.SelectorTypeID {
		t.Errorf("expected SelectorTypeID, got %v", wrapper.element.Type)
	}

	if wrapper.element.Value != "com.example.app:id/button" {
		t.Errorf("expected 'com.example.app:id/button', got %v", wrapper.element.Value)
	}
}

func TestTranslateUnknown(t *testing.T) {
	translator := NewTranslator()

	cmd := &MaestroCommand{
		Log: "test log",
	}

	_, err := translator.TranslateCommand(cmd, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !IsUnsupportedCommand(err) {
		t.Errorf("expected ErrUnsupportedCommand, got %v", err)
	}
}

func TestTranslateFlow(t *testing.T) {
	translator := NewTranslator()

	flow := &MaestroFlow{
		Steps: []MaestroCommand{
			{Launch: "com.example.app"},
			{TapOn: &TapOnCommand{Text: "Login"}},
		},
	}

	actions, err := translator.TranslateFlow(flow, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(actions) != 2 {
		t.Errorf("expected 2 actions, got %d", len(actions))
	}

	if _, ok := actions[0].(*LaunchWrapper); !ok {
		t.Errorf("expected first action to be *launchWrapper, got %T", actions[0])
	}

	if _, ok := actions[1].(*TapOnWrapper); !ok {
		t.Errorf("expected second action to be *tapOnWrapper, got %T", actions[1])
	}
}

func TestTranslateInputText(t *testing.T) {
	translator := NewTranslator()

	cmd := &MaestroCommand{
		InputText: &InputTextCommand{
			Text: "hello@example.com",
			ID:   "email_input",
		},
	}

	act, err := translator.TranslateCommand(cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrapper, ok := act.(*SendKeysWrapper)
	if !ok {
		t.Fatalf("expected *sendKeysWrapper, got %T", act)
	}

	if wrapper.text != "hello@example.com" {
		t.Errorf("expected 'hello@example.com', got %v", wrapper.text)
	}

	if wrapper.element == nil {
		t.Fatal("expected element to be set")
	}

	if wrapper.element.Type != selector.SelectorTypeID {
		t.Errorf("expected SelectorTypeID, got %v", wrapper.element.Type)
	}
}

func TestTranslateSwipe(t *testing.T) {
	translator := NewTranslator()

	cmd := &MaestroCommand{
		Swipe: &SwipeCommand{
			StartX:   100,
			StartY:   200,
			EndX:     100,
			EndY:     400,
			Duration: 500,
		},
	}

	act, err := translator.TranslateCommand(cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrapper, ok := act.(*SwipeWrapper)
	if !ok {
		t.Fatalf("expected *swipeWrapper, got %T", act)
	}

	if wrapper.startX != 100 {
		t.Errorf("expected startX 100, got %d", wrapper.startX)
	}

	if wrapper.startY != 200 {
		t.Errorf("expected startY 200, got %d", wrapper.startY)
	}

	if wrapper.endX != 100 {
		t.Errorf("expected endX 100, got %d", wrapper.endX)
	}

	if wrapper.endY != 400 {
		t.Errorf("expected endY 400, got %d", wrapper.endY)
	}
}

func TestTranslateLaunch(t *testing.T) {
	translator := NewTranslator()

	cmd := &MaestroCommand{
		Launch: "com.example.app",
	}

	act, err := translator.TranslateCommand(cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrapper, ok := act.(*LaunchWrapper)
	if !ok {
		t.Fatalf("expected *launchWrapper, got %T", act)
	}

	if wrapper.appID != "com.example.app" {
		t.Errorf("expected 'com.example.app', got %v", wrapper.appID)
	}

	if !wrapper.waitIdle {
		t.Error("expected waitIdle to be true")
	}
}

func TestTranslateTap(t *testing.T) {
	translator := NewTranslator()

	cmd := &MaestroCommand{
		Tap: &PointCommand{
			X: 150,
			Y: 200,
		},
	}

	act, err := translator.TranslateCommand(cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrapper, ok := act.(*TapWrapper)
	if !ok {
		t.Fatalf("expected *tapWrapper, got %T", act)
	}

	if wrapper.x != 150 {
		t.Errorf("expected x 150, got %d", wrapper.x)
	}

	if wrapper.y != 200 {
		t.Errorf("expected y 200, got %d", wrapper.y)
	}
}

func TestTranslateWait(t *testing.T) {
	translator := NewTranslator()

	cmd := &MaestroCommand{
		Wait: 1000,
	}

	act, err := translator.TranslateCommand(cmd, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wrapper, ok := act.(*WaitWrapper)
	if !ok {
		t.Fatalf("expected *waitWrapper, got %T", act)
	}

	if wrapper.duration.Milliseconds() != 1000 {
		t.Errorf("expected 1000ms, got %d", wrapper.duration.Milliseconds())
	}
}

func TestTranslateNilFlow(t *testing.T) {
	translator := NewTranslator()

	actions, err := translator.TranslateFlow(nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actions != nil {
		t.Errorf("expected nil actions, got %v", actions)
	}
}

func TestTranslateEmptyFlow(t *testing.T) {
	translator := NewTranslator()

	flow := &MaestroFlow{
		Steps: []MaestroCommand{},
	}

	actions, err := translator.TranslateFlow(flow, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(actions))
	}
}
