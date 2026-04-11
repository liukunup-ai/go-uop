package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/liukunup/go-uop/internal/report"
	"github.com/liukunup/go-uop/internal/runner"
)

const (
	version     = "0.2.0"
	exitSuccess = 0
	exitFailure = 1
)

var (
	platformFlag  string
	serialFlag    string
	addressFlag   string
	reportFormats string
	reportOutput  string
	debugFlag     bool
	configFile    string
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fset := flag.NewFlagSet("uop", flag.ContinueOnError)
	fset.StringVar(&platformFlag, "platform", "", "Target platform (ios|android)")
	fset.StringVar(&serialFlag, "serial", "", "Device serial number (Android)")
	fset.StringVar(&addressFlag, "address", "", "Device address (iOS)")
	fset.StringVar(&reportFormats, "report", "json,html", "Report formats (json,html,junit)")
	fset.StringVar(&reportOutput, "output", "reports", "Output directory for reports")
	fset.BoolVar(&debugFlag, "debug", false, "Enable debug mode")
	fset.StringVar(&configFile, "config", "", "Config file path")

	fset.Usage = func() { printHelp() }

	if err := fset.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return exitSuccess
		}
		return exitFailure
	}

	remaining := fset.Args()
	if len(remaining) == 0 {
		fset.Usage()
		return exitFailure
	}

	cmd := remaining[0]
	remaining = remaining[1:]

	switch cmd {
	case "run":
		return runCmd(remaining)
	case "debug":
		return debugCmd(remaining)
	case "test":
		return testCmd(remaining)
	case "devices":
		return devicesCmd()
	case "version":
		fmt.Printf("uop version %s\n", version)
		return exitSuccess
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		fset.Usage()
		return exitFailure
	}
}

func runCmd(args []string) int {
	fset := flag.NewFlagSet("uop run", flag.ContinueOnError)
	var deviceID string
	fset.StringVar(&deviceID, "device", "", "Device ID to use")

	if err := fset.Parse(args); err != nil {
		return exitFailure
	}

	targets := fset.Args()
	if len(targets) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no YAML files specified")
		fset.Usage()
		return exitFailure
	}

	pool := runner.NewDevicePool()
	if err := setupDevices(pool); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up devices: %v\n", err)
		return exitFailure
	}

	if deviceID != "" {
		if err := pool.SwitchDevice(deviceID); err != nil {
			fmt.Fprintf(os.Stderr, "Error switching to device: %v\n", err)
			return exitFailure
		}
	}

	reportGen := report.NewGenerator("uop-run")
	executor := runner.NewExecutor(pool, reportGen)

	for _, target := range targets {
		flow, err := runner.ParseFlowFile(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", target, err)
			continue
		}

		if err := executor.ExecuteSuite(flow); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing %s: %v\n", flow.Name, err)
		} else {
			fmt.Printf("✓ %s completed\n", flow.Name)
		}
	}

	exportReports(reportGen)
	return exitSuccess
}

func debugCmd(args []string) int {
	fset := flag.NewFlagSet("uop debug", flag.ContinueOnError)
	var deviceID string
	fset.StringVar(&deviceID, "device", "", "Device ID to use")
	fset.BoolVar(&debugFlag, "debug", true, "Enable debug mode")

	if err := fset.Parse(args); err != nil {
		return exitFailure
	}

	targets := fset.Args()
	if len(targets) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no YAML files specified")
		fset.Usage()
		return exitFailure
	}

	pool := runner.NewDevicePool()
	if err := setupDevices(pool); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up devices: %v\n", err)
		return exitFailure
	}

	if deviceID != "" {
		if err := pool.SwitchDevice(deviceID); err != nil {
			fmt.Fprintf(os.Stderr, "Error switching to device: %v\n", err)
			return exitFailure
		}
	}

	reportGen := report.NewGenerator("uop-debug")
	executor := runner.NewExecutor(pool, reportGen)

	fmt.Println("Debug mode - step through commands with interactive prompts")

	for _, target := range targets {
		flow, err := runner.ParseFlowFile(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", target, err)
			continue
		}

		fmt.Printf("Flow: %s (%d testcases)\n", flow.Name, len(flow.TestCases))
		debugger := runner.NewDebugger(flow.TestCases)

		if err := debugger.ExecuteWithDebug(flow, func(tc, step int, s runner.Step) error {
			fmt.Printf("[%d:%d] %v\n", tc, step, s)
			return executor.ExecuteStep(fmt.Sprintf("%s-%d", s.Type, step), s)
		}); err != nil {
			fmt.Fprintf(os.Stderr, "Debug error: %v\n", err)
		}
	}

	return exitSuccess
}

func testCmd(args []string) int {
	fset := flag.NewFlagSet("uop test", flag.ContinueOnError)
	var suitePath string
	fset.StringVar(&suitePath, "suite", "", "Suite file path")

	if err := fset.Parse(args); err != nil {
		return exitFailure
	}

	targets := fset.Args()
	if len(targets) == 0 && suitePath == "" {
		fmt.Fprintln(os.Stderr, "Error: no YAML files or suite specified")
		fset.Usage()
		return exitFailure
	}

	pool := runner.NewDevicePool()
	if err := setupDevices(pool); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up devices: %v\n", err)
		return exitFailure
	}

	reportGen := report.NewGenerator("uop-test")

	if suitePath != "" {
		result, err := runner.ParseAndRunSuiteFile(suitePath, pool, reportGen)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running suite: %v\n", err)
			return exitFailure
		}
		fmt.Printf("Suite: %d/%d passed\n", result.PassedSteps, result.TotalSteps)
	} else {
		executor := runner.NewExecutor(pool, reportGen)
		for _, target := range targets {
			flow, err := runner.ParseFlowFile(target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", target, err)
				continue
			}

			reportGen.StartTest(flow.Name)
			if err := executor.ExecuteSuite(flow); err != nil {
				reportGen.EndTest("failed", err)
				fmt.Printf("✗ %s failed: %v\n", flow.Name, err)
			} else {
				reportGen.EndTest("passed", nil)
				fmt.Printf("✓ %s passed\n", flow.Name)
			}
		}
	}

	exportReports(reportGen)
	return exitSuccess
}

func devicesCmd() int {
	pool := runner.NewDevicePool()
	if err := setupDevices(pool); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up devices: %v\n", err)
		return exitFailure
	}

	fmt.Println("Available devices:")
	for id, dev := range pool.ListDevices() {
		fmt.Printf("  [%s] %s (%s)\n", id, dev.Type, dev.Serial)
	}
	return exitSuccess
}

func setupDevices(pool *runner.DevicePool) error {
	if platformFlag != "" {
		deviceType := platformFlag
		serial := serialFlag
		if deviceType == "ios" {
			serial = addressFlag
		}
		id := platformFlag + "-default"
		return pool.AddDevice(id, deviceType, serial)
	}
	return nil
}

func exportReports(reportGen *report.Generator) {
	formats := strings.Split(reportFormats, ",")
	for _, format := range formats {
		format = strings.TrimSpace(format)
		path := fmt.Sprintf("%s/report.%s", reportOutput, format)
		if err := reportGen.WriteFormat(format, path); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s report: %v\n", format, err)
		} else {
			fmt.Printf("Report written to %s\n", path)
		}
	}
}

func printHelp() {
	fmt.Printf(`uop CLI - Unified mobile automation %s

Usage:
  uop <command> [arguments] [flags]

Commands:
  run <file...>       Run YAML flow(s)
  debug <file...>      Debug YAML flow(s) with step-through
  test <file...>      Run test(s) with reporting
  devices             List available devices
  version             Show version information

Flags:
  --platform <p>      Target platform: ios or android
  --serial <s>        Device serial number (Android)
  --address <a>       Device address (iOS)
  --report <f>        Report formats: json,html,junit (comma-separated)
  --output <dir>      Output directory for reports (default: reports)
  --device <id>       Device ID to use
  --debug             Enable debug mode
  --config <path>     Config file path
  --help, -h          Show this help message

Examples:
  uop run flow.yaml
  uop run flow1.yaml flow2.yaml --platform ios --address http://localhost:8100
  uop debug flow.yaml --device my-phone
  uop test --suite suite.yaml
  uop test flow.yaml --report json,html,junit

Exit codes:
  0  Success
  1  Failure`, version)
}
