package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/liukunup/go-uop/internal/report"
)

type SuiteResult struct {
	TCResults   []*TCResult
	TotalSteps  int
	PassedSteps int
	FailedSteps int
}

type TCResult struct {
	TCName string
	Status string
	Steps  int
	Error  error
}

type SuiteRunner struct {
	pool      *DevicePool
	executor  *Executor
	reportGen *report.Generator
}

func NewSuiteRunner(pool *DevicePool, reportGen *report.Generator) *SuiteRunner {
	return &SuiteRunner{
		pool:      pool,
		executor:  NewExecutor(pool, reportGen),
		reportGen: reportGen,
	}
}

func ParseAndRunSuite(r io.Reader, pool *DevicePool, reportGen *report.Generator) (*SuiteResult, error) {
	suite, err := ParseFlow(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse suite: %w", err)
	}

	runner := NewSuiteRunner(pool, reportGen)
	return runner.RunSuite(suite)
}

func ParseAndRunSuiteFile(path string, pool *DevicePool, reportGen *report.Generator) (*SuiteResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open suite file: %w", err)
	}
	defer file.Close()

	return ParseAndRunSuite(file, pool, reportGen)
}

func (s *SuiteRunner) RunSuite(suite *TestSuite) (*SuiteResult, error) {
	result := &SuiteResult{
		TCResults: make([]*TCResult, 0, len(suite.TestCases)),
	}

	for _, tc := range suite.TestCases {
		tcResult := s.runTestCase(tc)
		result.TCResults = append(result.TCResults, tcResult)
		if tcResult.Status == "failed" {
			result.FailedSteps += tcResult.Steps
		} else {
			result.PassedSteps += tcResult.Steps
		}
		result.TotalSteps += tcResult.Steps
	}

	return result, nil
}

func (s *SuiteRunner) runTestCase(tc TestCase) *TCResult {
	s.reportGen.StartTest(tc.Name)

	err := s.executor.ExecuteTestCase(0, tc)
	if err != nil {
		s.reportGen.EndTest("failed", err)
		return &TCResult{
			TCName: tc.Name,
			Status: "failed",
			Steps:  len(tc.Steps),
			Error:  err,
		}
	}

	s.reportGen.EndTest("passed", nil)
	return &TCResult{
		TCName: tc.Name,
		Status: "passed",
		Steps:  len(tc.Steps),
	}
}
