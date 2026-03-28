package vision

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"gocv.io/x/gocv"
)

type loftrMatcher struct {
	config      *Config
	modelPath   string
	session     *onnxSession
	initialized bool
}

type onnxSession struct {
	session     interface{} // *onnxruntime_go.AdvancedSession or similar
	inputNames  []string
	outputNames []string
}

// loftrMatcher implements LoFTR (LoFTR: Detector-Free Local Feature Matching with Transformers)
// using ONNX Runtime for inference. LoFTR provides robust feature matching even in
// challenging conditions where traditional feature detectors struggle.
//
// Model source: https://github.com/oooooha/loftr2onnx
//
// Input format:
//   - Two images: screenshot (reference) and template (to find)
//   - Preprocessed to 640x480 grayscale, normalized to [0,1]
//
// Output format:
//   - Match confidence scores
//   - Keypoint coordinates

func newLoFTRMatcher(cfg *Config) *loftrMatcher {
	return &loftrMatcher{
		config:    cfg,
		modelPath: getLoFTRModelPath(),
	}
}

func getLoFTRModelPath() string {
	// Check environment variable for model path
	if path := os.Getenv("LOFTR_MODEL_PATH"); path != "" {
		return path
	}
	// Default path
	return "./models/loftr.onnx"
}

func (m *loftrMatcher) Name() string {
	return "loftr"
}

// Initialize sets up the ONNX Runtime session with the LoFTR model.
// This is called lazily on first use.
func (m *loftrMatcher) initialize() error {
	if m.initialized {
		return nil
	}

	// Check if model file exists
	if _, err := os.Stat(m.modelPath); os.IsNotExist(err) {
		return fmt.Errorf("LoFTR model not found at %s. Set LOFTR_MODEL_PATH or download from https://github.com/oooooha/loftr2onnx", m.modelPath)
	}

	// Initialize ONNX Runtime (lazy import)
	// Note: We use dynamic import to avoid hard dependency when not using LoFTR
	if err := initONNXRuntime(); err != nil {
		return fmt.Errorf("failed to initialize ONNX Runtime: %w", err)
	}

	session, err := createONNXSession(m.modelPath, []string{"image0", "image1"}, []string{"matches0", "matches1", "scores"})
	if err != nil {
		return fmt.Errorf("failed to create ONNX session: %w", err)
	}

	m.session = session
	m.initialized = true
	return nil
}

func (m *loftrMatcher) Find(screenshot, templateImg []byte) ([]*MatchResult, error) {
	if screenshot == nil || templateImg == nil {
		return []*MatchResult{}, nil
	}

	// Initialize ONNX session if not already done
	if err := m.initialize(); err != nil {
		// Return empty results instead of error - LoFTR might not be available
		return []*MatchResult{}, nil
	}

	cfg := defaultConfig()
	if m.config != nil {
		cfg = m.config
	}

	// Preprocess images for LoFTR
	img0, _, err := preprocessForLoFTR(screenshot)
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess screenshot: %w", err)
	}

	img1, _, err := preprocessForLoFTR(templateImg)
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess template: %w", err)
	}

	// Run inference
	matches, scores, err := runLoFTRInference(m.session, img0, img1)
	if err != nil {
		return nil, fmt.Errorf("LoFTR inference failed: %w", err)
	}

	// Parse matches and apply threshold
	var results []*MatchResult
	threshold := cfg.Threshold
	if threshold <= 0 {
		threshold = 0.5 // Default threshold for LoFTR
	}

	for i, score := range scores {
		if float64(score) >= threshold {
			m0 := matches[i*2]
			m1 := matches[i*2+1]

			// Calculate bounding box centered on match
			x, y := int(m0), int(m1)
			width := 50 // Approximate feature size
			height := 50

			results = append(results, &MatchResult{
				Found:  true,
				X:      x - width/2,
				Y:      y - height/2,
				Width:  width,
				Height: height,
				Score:  float64(score),
			})
		}
	}

	if len(results) == 0 {
		return []*MatchResult{}, nil
	}

	return results, nil
}

// preprocessForLoFTR prepares an image for LoFTR inference:
// - Converts to grayscale
// - Resizes to 640x480 (LoFTR default input size)
// - Normalizes pixel values to [0, 1]
func preprocessForLoFTR(imgData []byte) ([]float32, image.Point, error) {
	// Decode image
	img, err := gocv.IMDecode(imgData, gocv.IMReadGrayScale)
	if err != nil {
		return nil, image.Point{}, err
	}
	defer img.Close()

	originalSize := image.Point{X: img.Cols(), Y: img.Rows()}

	// Resize to LoFTR input size (640x480)
	targetSize := image.Point{X: 640, Y: 480}
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, targetSize, 0, 0, gocv.InterpolationLinear)

	// Convert to float32 and normalize to [0, 1]
	pixels := make([]float32, 640*480)
	idx := 0
	for y := 0; y < 480; y++ {
		for x := 0; x < 640; x++ {
			pixels[idx] = float32(resized.GetDoubleAt(y, x)) / 255.0
			idx++
		}
	}

	return pixels, originalSize, nil
}

// runLoFTRInference executes the LoFTR model on two preprocessed images
func runLoFTRInference(session *onnxSession, img0, img1 []float32) (matches []float32, scores []float32, err error) {
	// This is a placeholder - actual implementation depends on onnxruntime_go API
	// The session would run inference like:
	//
	// inputTensor0, _ := ort.NewTensor(ort.NewShape(1, 1, 480, 640), img0)
	// inputTensor1, _ := ort.NewTensor(ort.NewShape(1, 1, 480, 640), img1)
	// outputMatches, _ := ort.NewEmptyTensor[float32](ort.NewShape(1, N, 2))
	// outputScores, _ := ort.NewEmptyTensor[float32](ort.NewShape(1, N))
	//
	// err = session.Run()
	// matches = outputMatches.GetData()
	// scores = outputScores.GetData()

	return nil, nil, fmt.Errorf("ONNX Runtime not yet integrated - requires onnxruntime_go dependency")
}

func (m *loftrMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	if screenshot == nil || len(results) == 0 || m.config == nil || m.config.DebugDir == "" {
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

	buf, err := gocv.IMEncode(".png", img)
	if err != nil {
		return nil
	}
	return buf.GetBytes()
}

// Placeholder functions for ONNX Runtime integration
// These will be implemented when onnxruntime_go is added as a dependency

var onnxInitialized = false

func initONNXRuntime() error {
	// TODO: Initialize onnxruntime_go
	// ort.SetSharedLibraryPath("/path/to/onnxruntime.so")
	// return ort.InitializeEnvironment()
	onnxInitialized = true
	return nil
}

func createONNXSession(modelPath string, inputNames, outputNames []string) (*onnxSession, error) {
	// TODO: Create ONNX session
	// return &onnxSession{
	//     session: session,
	//     inputNames: inputNames,
	//     outputNames: outputNames,
	// }, nil
	return &onnxSession{
		inputNames:  inputNames,
		outputNames: outputNames,
	}, nil
}

// Calculate IoU for overlap computation (utility function)
func calculateIoU(x1, y1, w1, h1, x2, y2, w2, h2 int) float64 {
	x_left := math.Max(float64(x1), float64(x2))
	y_top := math.Max(float64(y1), float64(y2))
	x_right := math.Min(float64(x1+w1), float64(x2+w2))
	y_bottom := math.Min(float64(y1+h1), float64(y2+h2))

	if x_right < x_left || y_bottom < y_top {
		return 0
	}

	inter_area := (x_right - x_left) * (y_bottom - y_top)
	area1 := float64(w1 * h1)
	area2 := float64(w2 * h2)
	union_area := area1 + area2 - inter_area

	if union_area <= 0 {
		return 0
	}
	return inter_area / union_area
}
