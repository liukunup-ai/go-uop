package diff

import (
	"testing"
)

func TestNew(t *testing.T) {
	d, err := New("pixel")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if d == nil {
		t.Fatal("New() returned nil")
	}
	if d.Name() != "pixel" {
		t.Errorf("Name() = %s, want pixel", d.Name())
	}
}

func TestNew_UnknownAlgorithm(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("New() expected error for unknown algorithm")
	}
}

func TestOptions(t *testing.T) {
	d, err := New("pixel",
		WithThreshold(0.05),
		WithRegion(0, 0, 100, 100),
		WithOutputDir("/tmp"),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	_ = d
}

func TestConfig_Defaults(t *testing.T) {
	cfg := defaultConfig()
	if cfg.Threshold != 0.1 {
		t.Errorf("defaultConfig().Threshold = %f, want 0.1", cfg.Threshold)
	}
	if cfg.Region != nil {
		t.Error("defaultConfig().Region = nil")
	}
	if cfg.OutputDir != "" {
		t.Errorf("defaultConfig().OutputDir = %q, want empty", cfg.OutputDir)
	}
}

func TestRect(t *testing.T) {
	r := &Rect{X: 10, Y: 20, Width: 100, Height: 50}
	if r.X != 10 || r.Y != 20 || r.Width != 100 || r.Height != 50 {
		t.Errorf("Rect = %+v, want {X:10, Y:20, Width:100, Height:50}", r)
	}
}
