package vision

import (
	"gocv.io/x/gocv"
)

type templateMatcher struct {
	config *Config
}

func newTemplateMatcher(cfg *Config) *templateMatcher {
	return &templateMatcher{config: cfg}
}

func (m *templateMatcher) Name() string {
	return "template"
}

func (m *templateMatcher) Find(screenshot, templateImg []byte) ([]*MatchResult, error) {
	if screenshot == nil || templateImg == nil {
		return nil, nil
	}

	cfg := defaultConfig()
	if m.config != nil {
		cfg = m.config
	}

	// Decode images using gocv
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

	// Ensure template is smaller than screen
	if screen.Cols() < tmpl.Cols() || screen.Rows() < tmpl.Rows() {
		return nil, nil
	}

	// Perform template matching
	result := gocv.NewMat()
	defer result.Close()

	var mask gocv.Mat
	err = gocv.MatchTemplate(screen, tmpl, &result, gocv.TmCcoeffNormed, mask)
	if err != nil {
		return nil, err
	}

	// Find best match
	_, maxVal, _, maxLoc := gocv.MinMaxLoc(result)

	if float64(maxVal) < cfg.Threshold {
		return nil, nil
	}

	return []*MatchResult{
		{
			Found:  true,
			X:      maxLoc.X,
			Y:      maxLoc.Y,
			Width:  tmpl.Cols(),
			Height: tmpl.Rows(),
			Score:  float64(maxVal),
		},
	}, nil
}

func (m *templateMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	return DebugRender(screenshot, results, m.config)
}
