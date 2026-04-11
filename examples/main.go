package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/liukunup/go-uop/internal/runner"
	"github.com/liukunup/go-uop/internal/selector"
)

const (
	version     = "0.1.0"
	exitSuccess = 0
	exitFailure = 1
)

var (
	fileFlag   string
	dryRunFlag bool
	debugFlag  bool
)

func init() {
	flag.Usage = func() {}
}

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if argsContains(args, "--version", "-v", "version") {
		fmt.Printf("uop version %s\n", version)
		return exitSuccess
	}

	if argsContains(args, "--help", "-h", "help") || len(args) == 0 {
		printHelp()
		return exitSuccess
	}

	flag.StringVar(&fileFlag, "file", "", "YAML flow file to execute")
	flag.BoolVar(&dryRunFlag, "dry-run", false, "Parse flow file without execution")
	flag.BoolVar(&debugFlag, "debug", false, "Enable debug output")

	flag.Parse()

	if fileFlag == "" {
		fmt.Fprintf(os.Stderr, "Error: --file is required\n\n")
		printHelp()
		return exitFailure
	}

	return runFileCmd()
}

func argsContains(args []string, targets ...string) bool {
	for _, arg := range args {
		for _, target := range targets {
			if arg == target {
				return true
			}
		}
	}
	return false
}

func runFileCmd() int {
	if debugFlag {
		fmt.Printf("[DEBUG] Loading file: %s\n", fileFlag)
		fmt.Printf("[DEBUG] Dry-run mode: %v\n", dryRunFlag)
	}

	suite, err := runner.ParseFlowFile(fileFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse %s: %v\n", fileFlag, err)
		return exitFailure
	}

	if debugFlag {
		fmt.Printf("[DEBUG] Parsed suite: %s\n", suite.Name)
		fmt.Printf("[DEBUG] Devices: %d, TestCases: %d\n",
			len(suite.Devices), len(suite.TestCases))
	}

	if dryRunFlag {
		fmt.Println("=== Dry Run: TestSuite ===")
		fmt.Printf("Name: %s\n", suite.Name)
		fmt.Printf("Description: %s\n", suite.Description)
		fmt.Printf("AppID: %s\n", suite.AppID)
		fmt.Printf("TestOutputDir: %s\n", suite.TestOutputDir)
		fmt.Printf("Devices: %d\n", len(suite.Devices))
		for i, d := range suite.Devices {
			defaultMark := ""
			if d.Default {
				defaultMark = " (default)"
			}
			fmt.Printf("  [%d] ID: %s%s, Type: %s, Serial: %s, UDID: %s, Dev: %s\n",
				i+1, d.ID, defaultMark, d.Type, d.Serial, d.UDID, d.Dev)
		}
		fmt.Printf("TestCases: %d\n", len(suite.TestCases))

		for i, tc := range suite.TestCases {
			device := tc.Device
			if device == "" {
				device = "-"
			}
			fmt.Printf("  [%d] Name: %s, Device: %s, Steps: %d\n",
				i+1, tc.Name, device, len(tc.Steps))
			for j, step := range tc.Steps {
				selectorStr := formatSelector(step.Selector)
				paramsStr := formatStepParams(step.Params)
				fmt.Printf("      [%d] %s%s%s\n", j+1, step.Type, selectorStr, paramsStr)
			}
		}
		fmt.Println("\n[Dry Run] Execution skipped.")
		return exitSuccess
	}

	fmt.Printf("Executing: %s\n", suite.Name)
	fmt.Printf("AppID: %s\n", suite.AppID)
	fmt.Printf("Devices: %d, TestCases: %d\n", len(suite.Devices), len(suite.TestCases))

	for i, tc := range suite.TestCases {
		device := tc.Device
		if device == "" {
			device = "-"
		}
		fmt.Printf("[%d] %s @ %s\n", i+1, tc.Name, device)
		for j, step := range tc.Steps {
			selectorStr := formatSelector(step.Selector)
			paramsStr := formatStepParams(step.Params)
			fmt.Printf("    [%d] %s%s%s\n", j+1, step.Type, selectorStr, paramsStr)
		}
	}

	fmt.Println("\n[Execution would happen here]")
	return exitSuccess
}

func formatSelector(s *selector.Selector) string {
	if s == nil || s.IsEmpty() {
		return ""
	}
	return "(" + s.String() + ")"
}

func formatStepParams(params map[string]any) string {
	if params == nil || len(params) == 0 {
		return ""
	}
	var parts []string
	for k, v := range params {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return "(" + strings.Join(parts, ", ") + ")"
}

func printHelp() {
	fmt.Printf(`go-uop CLI - Unified mobile terminal automation %s

Usage:
  uop --file <yaml> [flags]

Flags:
  --file <path>             YAML flow file to execute (required)
  --dry-run                 Parse flow file without execution
  --debug                   Enable debug output
  --help, -h, help          Show this help message
  --version, -v, version    Show version information

Examples:
  uop --file demo.yaml
  uop --file demo.yaml --dry-run --debug
  uop --help

Exit codes:
  0  Success
  1  Failure
`, version)
}
