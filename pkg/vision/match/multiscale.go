package match

import (
	"gocv.io/x/gocv"
	"image"
	"math"
	"sort"
)

type multiscaleMatcher struct {
	config *Config
}

func newMultiscaleMatcher(cfg *Config) *multiscaleMatcher {
	return &multiscaleMatcher{config: cfg}
}

func (m *multiscaleMatcher) Name() string {
	return "multiscale"
}

func (m *multiscaleMatcher) Find(screenshot, templateImg []byte) ([]*MatchResult, error) {
	if screenshot == nil || templateImg == nil {
		return []*MatchResult{}, nil
	}

	cfg := defaultConfig()
	if m.config != nil {
		cfg = m.config
	}

	screen, err := gocv.IMDecode(screenshot, gocv.IMReadGrayScale)
	if err != nil {
		return nil, err
	}
	defer screen.Close()

	tmpl, err := gocv.IMDecode(templateImg, gocv.IMReadGrayScale)
	if err != nil {
		return nil, err
	}
	defer tmpl.Close()

	var allMatches []matchCandidate

	for scale := cfg.ScaleMin; scale <= cfg.ScaleMax; scale += cfg.ScaleStep {
		newWidth := int(float64(tmpl.Cols()) * scale)
		newHeight := int(float64(tmpl.Rows()) * scale)

		if newWidth < 1 || newHeight < 1 {
			continue
		}
		if screen.Cols() < newWidth || screen.Rows() < newHeight {
			continue
		}

		resized := gocv.NewMat()
		err = gocv.Resize(tmpl, &resized, image.Point{X: newWidth, Y: newHeight}, 0, 0, gocv.InterpolationLinear)
		if err != nil {
			continue
		}
		defer resized.Close()

		result := gocv.NewMat()
		defer result.Close()

		var mask gocv.Mat
		err = gocv.MatchTemplate(screen, resized, &result, gocv.TmCcoeffNormed, mask)
		if err != nil {
			continue
		}

		for y := 0; y < result.Rows(); y++ {
			for x := 0; x < result.Cols(); x++ {
				val := float64(result.GetFloatAt(y, x))
				if val >= cfg.Threshold {
					allMatches = append(allMatches, matchCandidate{
						x: x, y: y,
						width: newWidth, height: newHeight,
						score: float32(val),
					})
				}
			}
		}
	}

	if len(allMatches) == 0 {
		return []*MatchResult{}, nil
	}

	sort.Slice(allMatches, func(i, j int) bool {
		return allMatches[i].score > allMatches[j].score
	})

	var results []*MatchResult
	used := make([]bool, len(allMatches))

	for i := 0; i < len(allMatches); i++ {
		if used[i] {
			continue
		}

		cand := allMatches[i]
		results = append(results, &MatchResult{
			Found:  true,
			X:      cand.x,
			Y:      cand.y,
			Width:  cand.width,
			Height: cand.height,
			Score:  float64(cand.score),
		})

		for j := i + 1; j < len(allMatches); j++ {
			if used[j] {
				continue
			}
			if overlapRatio(cand, allMatches[j]) > cfg.NMSThreshold {
				used[j] = true
			}
		}
	}

	return results, nil
}

func (m *multiscaleMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	return debugRender(screenshot, results, m.config)
}

type matchCandidate struct {
	x, y   int
	width  int
	height int
	score  float32
}

func overlapRatio(a, b matchCandidate) float64 {
	x1 := math.Max(float64(a.x), float64(b.x))
	y1 := math.Max(float64(a.y), float64(b.y))
	x2 := math.Min(float64(a.x+a.width), float64(b.x+b.width))
	y2 := math.Min(float64(a.y+a.height), float64(b.y+b.height))

	if x2 <= x1 || y2 <= y1 {
		return 0
	}

	inter := (x2 - x1) * (y2 - y1)
	union := float64(a.width*a.height+b.width*b.height) - inter
	if union <= 0 {
		return 0
	}
	return inter / union
}
