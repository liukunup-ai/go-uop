package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/liukunup/go-uop/internal/report"
)

type SuiteResult struct {
	FlowResults []*FlowResult
	TotalSteps  int
	PassedSteps int
	FailedSteps int
}

type FlowResult struct {
	FlowName string
	Status   string
	Steps    int
	Error    error
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
	suite, err := ParseSuite(r)
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

func (s *SuiteRunner) RunSuite(suite *Suite) (*SuiteResult, error) {
	result := &SuiteResult{
		FlowResults: make([]*FlowResult, 0, len(suite.Flows)),
	}

	for _, flowSpec := range suite.Flows {
		flowResult := s.runFlow(flowSpec)
		result.FlowResults = append(result.FlowResults, flowResult)
		if flowResult.Status == "failed" {
			result.FailedSteps += flowResult.Steps
		} else {
			result.PassedSteps += flowResult.Steps
		}
		result.TotalSteps += flowResult.Steps
	}

	return result, nil
}

func (s *SuiteRunner) runFlow(flowSpec SuiteFlow) *FlowResult {
	s.reportGen.StartTest(flowSpec.Name)

	flow, err := ParseFlowFile(flowSpec.Path)
	if err != nil {
		s.reportGen.EndTest("failed", err)
		return &FlowResult{
			FlowName: flowSpec.Name,
			Status:   "failed",
			Error:    err,
		}
	}

	err = s.executor.ExecuteFlow(flow)
	if err != nil {
		s.reportGen.EndTest("failed", err)
		return &FlowResult{
			FlowName: flow.Name,
			Status:   "failed",
			Steps:    len(flow.Steps),
			Error:    err,
		}
	}

	s.reportGen.EndTest("passed", nil)
	return &FlowResult{
		FlowName: flow.Name,
		Status:   "passed",
		Steps:    len(flow.Steps),
	}
}
