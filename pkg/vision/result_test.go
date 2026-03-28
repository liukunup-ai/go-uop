package vision

import (
	"image"
	"testing"
)

func TestMatchResult_Center(t *testing.T) {
	r := &MatchResult{
		X: 100, Y: 200, Width: 50, Height: 60,
	}
	cx, cy := r.Center()
	if cx != 125 || cy != 230 {
		t.Errorf("Center() = (%d, %d), want (125, 230)", cx, cy)
	}
}

func TestMatchResult_Rectangle(t *testing.T) {
	r := &MatchResult{
		X: 100, Y: 200, Width: 50, Height: 60,
	}
	rect := r.Rectangle()
	if rect != image.Rect(100, 200, 150, 260) {
		t.Errorf("Rectangle() = %v, want Rect(100, 200, 150, 260)", rect)
	}
}
