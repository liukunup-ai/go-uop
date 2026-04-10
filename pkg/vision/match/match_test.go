package match

import (
	"testing"
)

func TestNew(t *testing.T) {
	m, err := New("template")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.Name() != "template" {
		t.Errorf("Name() = %s, want template", m.Name())
	}
}

func TestNew_Multiscale(t *testing.T) {
	m, err := New("multiscale")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.Name() != "multiscale" {
		t.Errorf("Name() = %s, want multiscale", m.Name())
	}
}

func TestNew_UnknownAlgorithm(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("New() expected error for unknown algorithm")
	}
}

func TestOptions(t *testing.T) {
	m, err := New("template",
		WithThreshold(0.9),
		WithScaleRange(0.5, 1.5),
		WithScaleStep(0.05),
		WithNMSThreshold(0.3),
		WithDebug("/tmp"),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_ = m
}

func TestNew_Sift(t *testing.T) {
	m, err := New("sift")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.Name() != "sift" {
		t.Errorf("Name() = %s, want sift", m.Name())
	}
}

func TestNew_LoFTR(t *testing.T) {
	m, err := New("loftr")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.Name() != "loftr" {
		t.Errorf("Name() = %s, want loftr", m.Name())
	}
}

func TestMatchResult_Center(t *testing.T) {
	r := &MatchResult{
		X:      10,
		Y:      20,
		Width:  100,
		Height: 50,
	}
	cx, cy := r.Center()
	if cx != 60 || cy != 45 {
		t.Errorf("Center() = (%d, %d), want (60, 45)", cx, cy)
	}
}

func TestMatchResult_Rectangle(t *testing.T) {
	r := &MatchResult{
		X:      10,
		Y:      20,
		Width:  100,
		Height: 50,
	}
	rect := r.Rectangle()
	if rect.Min.X != 10 || rect.Min.Y != 20 || rect.Max.X != 110 || rect.Max.Y != 70 {
		t.Errorf("Rectangle() = %v, want Rect(10, 20, 110, 70)", rect)
	}
}
