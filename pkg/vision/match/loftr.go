package match

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"runtime"

	ort "github.com/yalue/onnxruntime_go"
	"gocv.io/x/gocv"
)

const (
	LoFTRInputWidth  = 512
	LoFTRInputHeight = 512
)

type loftrMatcher struct {
	config      *Config
	modelPath   string
	session     *onnxSession
	initialized bool
}

type onnxSession struct {
	session *ort.DynamicAdvancedSession
}

func newLoFTRMatcher(cfg *Config) *loftrMatcher {
	return &loftrMatcher{
		config:    cfg,
		modelPath: getLoFTRModelPath(),
	}
}

func getLoFTRModelPath() string {
	if path := os.Getenv("LOFTR_MODEL_PATH"); path != "" {
		return path
	}
	return "./models/loftr.onnx"
}

func (m *loftrMatcher) Name() string {
	return "loftr"
}

func (m *loftrMatcher) initialize() error {
	if m.initialized {
		return nil
	}

	if _, err := os.Stat(m.modelPath); os.IsNotExist(err) {
		return fmt.Errorf("LoFTR model not found at %s", m.modelPath)
	}

	sharedLibPath := getONNXRuntimeLibPath()
	ort.SetSharedLibraryPath(sharedLibPath)

	if err := ort.InitializeEnvironment(); err != nil {
		return fmt.Errorf("failed to initialize ONNX Runtime: %w", err)
	}

	session, err := ort.NewDynamicAdvancedSession(
		m.modelPath,
		[]string{"image0", "image1"},
		[]string{"keypoints0", "keypoints1", "confidence"},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create ONNX session: %w", err)
	}

	m.session = &onnxSession{session: session}
	m.initialized = true
	return nil
}

func getONNXRuntimeLibPath() string {
	if path := os.Getenv("ONNXRUNTIME_LIB_PATH"); path != "" {
		return path
	}

	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return "/opt/homebrew/lib/libonnxruntime.dylib"
		}
		return "/usr/local/lib/libonnxruntime.dylib"
	case "linux":
		return "/usr/local/lib/libonnxruntime.so"
	case "windows":
		return "C:\\onnxruntime\\onnxruntime.dll"
	}
	return "./libonnxruntime.so"
}

func (m *loftrMatcher) Find(screenshot, templateImg []byte) ([]*MatchResult, error) {
	if screenshot == nil || templateImg == nil {
		return []*MatchResult{}, nil
	}

	if err := m.initialize(); err != nil {
		return []*MatchResult{}, nil
	}

	cfg := defaultConfig()
	if m.config != nil {
		cfg = m.config
	}

	img0Data, origSize1, err := preprocessForLoFTR(screenshot)
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess screenshot: %w", err)
	}

	img1Data, _, err := preprocessForLoFTR(templateImg)
	if err != nil {
		return nil, fmt.Errorf("failed to preprocess template: %w", err)
	}

	matches0, matches1, confidences, err := runLoFTRInference(m.session, img0Data, img1Data)
	if err != nil {
		return nil, fmt.Errorf("LoFTR inference failed: %w", err)
	}

	var results []*MatchResult
	threshold := cfg.Threshold
	if threshold <= 0 {
		threshold = 0.5
	}

	scaleX := float64(origSize1.X) / LoFTRInputWidth
	scaleY := float64(origSize1.Y) / LoFTRInputHeight

	for i, conf := range confidences {
		if float64(conf) >= threshold {
			x0, y0 := matches0[i*2], matches0[i*2+1]
			x1, y1 := matches1[i*2], matches1[i*2+1]

			_ = x0
			_ = y0

			centerX := int(float64(x1) * scaleX)
			centerY := int(float64(y1) * scaleY)

			width, height := 50, 50

			results = append(results, &MatchResult{
				Found:  true,
				X:      centerX - width/2,
				Y:      centerY - height/2,
				Width:  width,
				Height: height,
				Score:  float64(conf),
			})
		}
	}

	if len(results) == 0 {
		return []*MatchResult{}, nil
	}

	return results, nil
}

func preprocessForLoFTR(imgData []byte) ([]float32, image.Point, error) {
	img, err := gocv.IMDecode(imgData, gocv.IMReadGrayScale)
	if err != nil {
		return nil, image.Point{}, err
	}
	defer img.Close()

	origWidth, origHeight := img.Cols(), img.Rows()

	targetSize := image.Point{X: LoFTRInputWidth, Y: LoFTRInputHeight}
	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(img, &resized, targetSize, 0, 0, gocv.InterpolationLinear)

	pixels := make([]float32, LoFTRInputWidth*LoFTRInputHeight)
	for y := 0; y < LoFTRInputHeight; y++ {
		for x := 0; x < LoFTRInputWidth; x++ {
			pixels[y*LoFTRInputWidth+x] = float32(resized.GetDoubleAt(y, x)) / 255.0
		}
	}

	return pixels, image.Point{X: origWidth, Y: origHeight}, nil
}

func runLoFTRInference(session *onnxSession, img0, img1 []float32) (matches0, matches1, confidences []float32, err error) {
	inputShape := ort.NewShape(1, 1, LoFTRInputHeight, LoFTRInputWidth)

	in0, err := ort.NewTensor(inputShape, img0)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create input tensor0: %w", err)
	}
	defer in0.Destroy()

	in1, err := ort.NewTensor(inputShape, img1)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create input tensor1: %w", err)
	}
	defer in1.Destroy()

	outputs := []ort.Value{nil, nil, nil}
	err = session.session.Run([]ort.Value{in0, in1}, outputs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("session.Run failed: %w", err)
	}

	out0 := outputs[0].(*ort.Tensor[float32])
	out1 := outputs[1].(*ort.Tensor[float32])
	out2 := outputs[2].(*ort.Tensor[float32])

	return out0.GetData(), out1.GetData(), out2.GetData(), nil
}

func (m *loftrMatcher) DebugRender(screenshot []byte, results []*MatchResult) []byte {
	return debugRender(screenshot, results, m.config)
}

func debugRender(screenshot []byte, results []*MatchResult, config *Config) []byte {
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

		rect := image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
		gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)

		cx, cy := r.Center()
		gocv.Line(&img, image.Point{cx - 10, cy}, image.Point{cx + 10, cy}, color.RGBA{0, 255, 0, 0}, 2)
		gocv.Line(&img, image.Point{cx, cy - 10}, image.Point{cx, cy + 10}, color.RGBA{0, 255, 0, 0}, 2)

		label := fmt.Sprintf("[%d] Score: %.2f Pos: (%d,%d)", i, r.Score, r.X, r.Y)
		gocv.PutText(&img, label, image.Point{r.X, r.Y - 10}, gocv.FontHersheySimplex, 0.5, color.RGBA{255, 255, 255, 0}, 2)
	}

	os.MkdirAll(config.DebugDir, 0755)
	filename := filepath.Join(config.DebugDir, fmt.Sprintf("debug_%d.png", len(results)))
	gocv.IMWrite(filename, img)

	buf, err := gocv.IMEncode(".png", img)
	if err != nil {
		return nil
	}
	return buf.GetBytes()
}
