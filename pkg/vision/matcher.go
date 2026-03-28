package vision

import (
	"fmt"
)

type Matcher interface {
	Find(screenshot, template []byte) ([]*MatchResult, error)
	Name() string
	DebugRender(screenshot []byte, results []*MatchResult) []byte
}

type Config struct {
	Threshold    float64
	ScaleMin     float64
	ScaleMax     float64
	ScaleStep    float64
	NMSThreshold float64
	DebugDir     string
}

var defaultConfig = func() *Config {
	return &Config{
		Threshold:    0.8,
		ScaleMin:     0.8,
		ScaleMax:     1.2,
		ScaleStep:    0.1,
		NMSThreshold: 0.5,
		DebugDir:     "",
	}
}

type Option func(*Config)

func WithThreshold(t float64) Option {
	return func(c *Config) { c.Threshold = t }
}

func WithScaleRange(min, max float64) Option {
	return func(c *Config) {
		c.ScaleMin = min
		c.ScaleMax = max
	}
}

func WithScaleStep(step float64) Option {
	return func(c *Config) { c.ScaleStep = step }
}

func WithNMSThreshold(t float64) Option {
	return func(c *Config) { c.NMSThreshold = t }
}

func WithDebug(outputDir string) Option {
	return func(c *Config) { c.DebugDir = outputDir }
}

func NewMatcher(algo string, opts ...Option) (Matcher, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	switch algo {
	case "template":
		return newTemplateMatcher(cfg), nil
	case "multiscale":
		return newMultiscaleMatcher(cfg), nil
	case "sift":
		return newSIFTMatcher(cfg), nil
	case "loftr":
		return newLoFTRMatcher(cfg), nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
}
