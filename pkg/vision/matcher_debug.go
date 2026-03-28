package vision

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"gocv.io/x/gocv"
)

// DebugRender renders match results on screenshot for debugging
func DebugRender(screenshot []byte, results []*MatchResult, config *Config) []byte {
	if screenshot == nil || len(results) == 0 || config == nil || config.DebugDir == "" {
		return nil
	}

	img, err := gocv.IMDecode(screenshot, gocv.IMReadColor)
	if err != nil {
		return nil
	}
	defer img.Close()

	for i, r := range results {
		if !r.Found {
			continue
		}

		// Draw rectangle
		rect := image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
		gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)

		// Draw center cross
		cx, cy := r.Center()
		gocv.Line(&img, image.Point{cx - 10, cy}, image.Point{cx + 10, cy}, color.RGBA{0, 255, 0, 0}, 2)
		gocv.Line(&img, image.Point{cx, cy - 10}, image.Point{cx, cy + 10}, color.RGBA{0, 255, 0, 0}, 2)

		// Draw label
		label := fmt.Sprintf("[%d] Score: %.2f Pos: (%d,%d)", i, r.Score, r.X, r.Y)
		gocv.PutText(&img, label, image.Point{r.X, r.Y - 10}, gocv.FontHersheySimplex, 0.5, color.RGBA{255, 255, 255, 0}, 2)
	}

	// Save to debug directory
	os.MkdirAll(config.DebugDir, 0755)
	filename := fmt.Sprintf("debug_%d.png", len(results))
	outputPath := filepath.Join(config.DebugDir, filename)
	gocv.IMWrite(outputPath, img)

	// Return encoded PNG
	buf, err := gocv.IMEncode(".png", img)
	if err != nil {
		return nil
	}
	defer buf.Close()
	return buf.GetBytes()
}
