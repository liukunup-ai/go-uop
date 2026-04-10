package diff

import "fmt"

type Differ interface {
	Compare(img1, img2 []byte, cfg *Config) (*DiffResult, error)
	Name() string
}

type Config struct {
	Threshold float64
	Region    *Rect
	OutputDir string
}

type Rect struct {
	X, Y, Width, Height int
}

func defaultConfig() *Config {
	return &Config{
		Threshold: 0.1,
		Region:    nil,
		OutputDir: "",
	}
}

type Option func(*Config)

func WithThreshold(t float64) Option {
	return func(c *Config) { c.Threshold = t }
}

func WithRegion(x, y, w, h int) Option {
	return func(c *Config) { c.Region = &Rect{X: x, Y: y, Width: w, Height: h} }
}

func WithOutputDir(dir string) Option {
	return func(c *Config) { c.OutputDir = dir }
}

func New(algo string, opts ...Option) (Differ, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	switch algo {
	case "pixel":
		return newPixelDiffer(cfg), nil
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algo)
	}
}
