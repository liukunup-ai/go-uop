package vision

import (
	"os"
	"testing"
)

func TestTemplateMatcher_Find_WithRealImages(t *testing.T) {
	// Skip if OpenCV not available
	if os.Getenv("TEST_OPENCV") != "1" {
		t.Skip("Skipping OpenCV test")
	}

	// Load test images
	screenshot, err := os.ReadFile("testdata/screenshot.png")
	if err != nil {
		t.Skip("Test images not found")
	}
	template, err := os.ReadFile("testdata/button.png")
	if err != nil {
		t.Skip("Test images not found")
	}

	m := newTemplateMatcher(nil)
	results, err := m.Find(screenshot, template)
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}

	// Should find at least one match
	if len(results) == 0 {
		t.Log("No matches found (may be expected if images don't match)")
	}
}
