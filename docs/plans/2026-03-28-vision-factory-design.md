# Vision Module Refactoring Design

**Date**: 2026-03-28
**Author**: Sisor Agent
**Status**: Approved + Implemented

## Overview

Refactor `pkg/vision/template.go` to use factory pattern for visual finding. Support multiple matching algorithms with a unified interface.

## Architecture

```
pkg/vision/
├── matcher.go              # Interface + factory + options
├── result.go               # MatchResult + helper methods
├── matcher_template.go     # OpenCV single-scale template matching
├── matcher_multiscale.go   # Multi-scale template matching + NMS
├── matcher_sift.go        # SIFT + RANSAC implementation
├── matcher_loftr.go       # LoFTR ONNX inference
├── matcher_debug.go       # Debug rendering helper
├── integration_test.go    # Integration tests
└── *_test.go            # Unit tests
```

Note: All matcher implementations are in the same `pkg/vision` package to avoid import cycles.

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

## Algorithms

### Template (OpenCV)

- Uses `gocv.MatchTemplate` with `TmCcoeffNormed`
- Single-scale matching
- Returns best match above threshold

### Multiscale (OpenCV)

- Multi-scale search from `scaleMin` to `scaleMax` with `scaleStep`
- NMS (Non-Maximum Suppression) using IoU
- Returns multiple matches

### SIFT + RANSAC

- Uses `gocv.NewSIFT()` detector
- BFMatcher with kNN
- Lowe's ratio test (0.75)
- RANSAC homography verification
- Perspective transform for bounding box

### LoFTR (ONNX Runtime)

- Uses `github.com/yalue/onnxruntime_go` DynamicAdvancedSession
- Input: 512x512 grayscale, normalized to [0,1]
- Output: keypoints0, keypoints1, confidence
- Model: https://github.com/oooooha/loftr2onnx

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

## LoFTR Setup

To use LoFTR matcher:

1. Install ONNX Runtime shared library:
```bash
# macOS
brew install onnxruntime

# Linux/Windows: download from https://github.com/microsoft/onnxruntime/releases
```

2. Set library path (if not in default location):
```bash
export ONNXRUNTIME_LIB_PATH=/opt/homebrew/lib/libonnxruntime.dylib
```

3. Download LoFTR ONNX model:
```bash
export LOFTR_MODEL_PATH=./models/loftr_outdoor_ds.onnx
# Model from: https://github.com/oooooha/loftr2onnx
```

## Implementation Status

| Task | Status |
|------|--------|
| Create `pkg/vision/matcher.go` with interface and factory | ✅ |
| Create `pkg/vision/result.go` with MatchResult and helpers | ✅ |
| Create `pkg/vision/matcher_template.go` - OpenCV single-scale | ✅ |
| Create `pkg/vision/matcher_multiscale.go` - Multi-scale + NMS | ✅ |
| Create `pkg/vision/matcher_sift.go` - SIFT + RANSAC | ✅ |
| Create `pkg/vision/matcher_loftr.go` - LoFTR ONNX | ✅ |
| Create `pkg/vision/matcher_debug.go` - Debug rendering | ✅ |
| Remove old `pkg/vision/template.go` | ✅ |
| Add integration tests | ✅ |
