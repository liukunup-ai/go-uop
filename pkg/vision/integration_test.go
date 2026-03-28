package vision

import (
	"os"
	"testing"
)

func TestIntegration_Factory(t *testing.T) {
	// Test that all algorithms can be created
	algorithms := []string{"template", "multiscale", "sift", "loftr"}

	for _, algo := range algorithms {
		t.Run(algo, func(t *testing.T) {
			m, err := NewMatcher(algo)
			if err != nil {
				t.Fatalf("NewMatcher(%s) failed: %v", algo, err)
			}
			if m.Name() != algo {
				t.Errorf("Name() = %q, want %q", m.Name(), algo)
			}
		})
	}
}

func TestIntegration_UnknownAlgorithm(t *testing.T) {
	_, err := NewMatcher("unknown")
	if err == nil {
		t.Error("NewMatcher(unknown) should return error")
	}
}

func TestIntegration_TemplateMatcher(t *testing.T) {
	if os.Getenv("TEST_OPENCV") != "1" {
		t.Skip("Skipping OpenCV test")
	}

	m, err := NewMatcher("template")
	if err != nil {
		t.Fatalf("NewMatcher(template) failed: %v", err)
	}

	// Test with nil inputs - should return empty results
	results, err := m.Find(nil, nil)
	if err != nil {
		t.Fatalf("Find(nil, nil) error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Find(nil, nil) returned %d results, want 0", len(results))
	}

	// Test DebugRender with nil - should return nil
	debug := m.DebugRender(nil, nil)
	if debug != nil {
		t.Error("DebugRender(nil, nil) should return nil")
	}
}

func TestIntegration_MultiscaleMatcher(t *testing.T) {
	if os.Getenv("TEST_OPENCV") != "1" {
		t.Skip("Skipping OpenCV test")
	}

	m, err := NewMatcher("multiscale")
	if err != nil {
		t.Fatalf("NewMatcher(multiscale) failed: %v", err)
	}

	// Test with nil inputs
	results, err := m.Find(nil, nil)
	if err != nil {
		t.Fatalf("Find(nil, nil) error = %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Find(nil, nil) returned %d results, want 0", len(results))
	}
}

func TestIntegration_Options(t *testing.T) {
	m, err := NewMatcher("template",
		WithThreshold(0.9),
		WithScaleRange(0.5, 1.5),
		WithScaleStep(0.2),
		WithNMSThreshold(0.4),
		WithDebug("/tmp/vision-debug"),
	)
	if err != nil {
		t.Fatalf("NewMatcher with options failed: %v", err)
	}
	if m.Name() != "template" {
		t.Errorf("Name() = %q, want %q", m.Name(), "template")
	}
}
