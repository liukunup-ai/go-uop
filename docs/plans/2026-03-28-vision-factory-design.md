# Vision Module Refactoring Design

**Date**: 2026-03-28
**Author**: Sisyphus Agent
**Status**: Approved

## Overview

Refactor `pkg/vision/template.go` to use factory pattern for visual finding. Support multiple matching algorithms with a unified interface.

## Architecture

```
pkg/vision/
├── matcher.go          # Interface + factory + options
├── result.go           # MatchResult + helper methods
├── sift/
│   └── matcher.go      # SIFT + RANSAC implementation
├── loftr/
│   └── matcher.go      # LoFTR ONNX inference implementation
└── matchers/
    ├── template.go     # OpenCV single-scale template matching
    └── multiscale.go   # Multi-scale template matching + NMS + dynamic threshold
```

## Dependencies

| Library | Purpose | Required |
|---------|---------|----------|
| `gocv.io/x/gocv` | Template matching, SIFT, image processing | Yes |
| `github.com/yalue/onnxruntime_go` | LoFTR ONNX inference | Yes |

## Core Types

### MatchResult

```go
type MatchResult struct {
    Found      bool
    X, Y       int      // Top-left corner
    Width      int
    Height     int
    Score      float64  // Confidence [0, 1]
}

func (r *MatchResult) Center() (int, int)
func (r *MatchResult) Rectangle() image.Rectangle
```

### Matcher Interface

```go
type Matcher interface {
    Find(screenshot, template []byte) ([]*MatchResult, error)
    Name() string
    DebugRender(screenshot []byte, results []*MatchResult) []byte
}
```

### Factory Function

```go
func NewMatcher(algo string, opts ...Option) (Matcher, error)
```

Supported algorithms:
- `"template"` - OpenCV single-scale template matching
- `"multiscale"` - Multi-scale template matching + NMS + dynamic threshold
- `"sift"` - SIFT + RANSAC geometric verification
- `"loftr"` - LoFTR ONNX inference

## Configuration Options

```go
func WithThreshold(t float64) Option           // Confidence threshold
func WithScaleRange(min, max float64) Option   // Multi-scale range
func WithScaleStep(step float64) Option        // Scale step
func WithNMSThreshold(t float64) Option        // NMS threshold
func WithDebug(outputDir string) Option         // Debug output directory
```

## Debug Rendering

When debug is enabled, renders:
- Rectangle around each matched location
- Text label with score and coordinates (top-left)
- Cross marker at center point

Output saved to specified debug directory.

## Usage Example

```go
matcher, err := vision.NewMatcher("multiscale",
    vision.WithThreshold(0.8),
    vision.WithScaleRange(0.8, 1.2),
    vision.WithNMSThreshold(0.5),
    vision.WithDebug("/tmp/vision-debug"),
)
if err != nil {
    log.Fatal(err)
}

results, err := matcher.Find(screenshot, template)
if err != nil {
    log.Fatal(err)
}

for i, r := range results {
    if r.Found {
        x, y := r.Center()
        fmt.Printf("[%d] Found at (%d,%d), confidence: %.2f%%\n",
            i, x, y, r.Score*100)
    }
}

// Debug image auto-saved if debug enabled
```

## Implementation Tasks

1. Create `pkg/vision/matcher.go` with interface and factory
2. Create `pkg/vision/result.go` with MatchResult and helpers
3. Create `pkg/vision/matchers/template.go` - OpenCV single-scale
4. Create `pkg/vision/matchers/multiscale.go` - Multi-scale + NMS
5. Create `pkg/vision/sift/matcher.go` - SIFT + RANSAC
6. Create `pkg/vision/loftr/matcher.go` - LoFTR ONNX
7. Update existing `pkg/vision/template.go` to use new structure
8. Add tests for each matcher
