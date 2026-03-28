package vision

import "image"

type MatchResult struct {
	Found  bool
	X, Y   int
	Width  int
	Height int
	Score  float64
}

func (r *MatchResult) Center() (int, int) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

func (r *MatchResult) Rectangle() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
}
