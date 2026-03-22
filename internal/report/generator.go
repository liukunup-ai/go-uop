package report

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type TestResult struct {
	Name        string           `json:"name"`
	StartTime   time.Time        `json:"startTime"`
	EndTime     time.Time        `json:"endTime"`
	Duration    time.Duration    `json:"duration"`
	Status      string           `json:"status"`
	Error       string           `json:"error,omitempty"`
	Steps       []StepResult     `json:"steps,omitempty"`
	Screenshots []ScreenshotInfo `json:"screenshots,omitempty"`
}

type StepResult struct {
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
	Status   string        `json:"status"`
	Error    string        `json:"error,omitempty"`
}

type ScreenshotInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Name      string    `json:"name"`
}

type SuiteResult struct {
	Name         string        `json:"name"`
	StartTime    time.Time     `json:"startTime"`
	EndTime      time.Time     `json:"endTime"`
	Duration     time.Duration `json:"duration"`
	TotalTests   int           `json:"totalTests"`
	PassedTests  int           `json:"passedTests"`
	FailedTests  int           `json:"failedTests"`
	SkippedTests int           `json:"skippedTests"`
	Results      []TestResult  `json:"results"`
	Metadata     Metadata      `json:"metadata,omitempty"`
}

type Metadata struct {
	Platform    string `json:"platform,omitempty"`
	DeviceID    string `json:"deviceId,omitempty"`
	AppVersion  string `json:"appVersion,omitempty"`
	BuildNumber string `json:"buildNumber,omitempty"`
}

type Generator struct {
	suiteName string
	results   []TestResult
	current   *TestResult
}

func NewGenerator(name string) *Generator {
	return &Generator{
		suiteName: name,
		results:   make([]TestResult, 0),
	}
}

func (g *Generator) StartTest(name string) {
	g.current = &TestResult{
		Name:      name,
		StartTime: time.Now(),
		Status:    "running",
	}
}

func (g *Generator) EndTest(status string, err error) {
	if g.current == nil {
		return
	}
	g.current.EndTime = time.Now()
	g.current.Duration = g.current.EndTime.Sub(g.current.StartTime)
	g.current.Status = status
	if err != nil {
		g.current.Error = err.Error()
	}
	g.results = append(g.results, *g.current)
	g.current = nil
}

func (g *Generator) AddStep(name string, duration time.Duration, status string, err error) {
	if g.current == nil {
		return
	}
	step := StepResult{
		Name:     name,
		Duration: duration,
		Status:   status,
	}
	if err != nil {
		step.Error = err.Error()
	}
	g.current.Steps = append(g.current.Steps, step)
}

func (g *Generator) AddScreenshot(path, name string) {
	if g.current == nil {
		return
	}
	g.current.Screenshots = append(g.current.Screenshots, ScreenshotInfo{
		Timestamp: time.Now(),
		Path:      path,
		Name:      name,
	})
}

func (g *Generator) Generate() *SuiteResult {
	endTime := time.Now()
	var startTime time.Time
	if len(g.results) > 0 {
		startTime = g.results[0].StartTime
	}

	passed := 0
	failed := 0
	skipped := 0
	for _, r := range g.results {
		switch r.Status {
		case "passed":
			passed++
		case "failed":
			failed++
		case "skipped":
			skipped++
		}
	}

	return &SuiteResult{
		Name:         g.suiteName,
		StartTime:    startTime,
		EndTime:      endTime,
		Duration:     endTime.Sub(startTime),
		TotalTests:   len(g.results),
		PassedTests:  passed,
		FailedTests:  failed,
		SkippedTests: skipped,
		Results:      g.results,
	}
}

func (g *Generator) ToJSON() ([]byte, error) {
	return json.MarshalIndent(g.Generate(), "", "  ")
}

func (g *Generator) ToJSONFile(path string) error {
	data, err := g.ToJSON()
	if err != nil {
		return fmt.Errorf("generate json: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func (g *Generator) Summary() string {
	s := g.Generate()
	return fmt.Sprintf("%s: %d/%d passed (%s)",
		s.Name, s.PassedTests, s.TotalTests, s.Duration)
}
