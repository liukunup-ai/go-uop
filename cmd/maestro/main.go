package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/liukunup/go-uop/core"
	"github.com/liukunup/go-uop/ios"
	"github.com/liukunup/go-uop/maestro"
)

const (
	version     = "0.1.0"
	exitSuccess = 0
	exitFailure = 1
)

var (
	deviceFlag string
	outputFlag string
	configFlag string
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flag.StringVar(&deviceFlag, "device", "", "Target device platform (ios|android)")
	flag.StringVar(&outputFlag, "output", "", "Directory for screenshots")
	flag.StringVar(&configFlag, "config", "", "Config file path")

	remainingArgs := parseArgs(args)

	if len(remainingArgs) > 0 && remainingArgs[0] == "--version" {
		fmt.Printf("maestro version %s\n", version)
		return exitSuccess
	}

	if len(remainingArgs) == 0 || remainingArgs[0] == "--help" || remainingArgs[0] == "-h" {
		printHelp()
		return exitSuccess
	}

	cmd := remainingArgs[0]
	file := ""
	if len(remainingArgs) > 1 {
		file = remainingArgs[1]
	}

	switch cmd {
	case "test":
		return testCmd(file)
	case "validate":
		return validateCmd(file)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printHelp()
		return exitFailure
	}
}

func testCmd(file string) int {
	if file == "" {
		fmt.Fprintln(os.Stderr, "Error: file path required for test command")
		fmt.Fprintln(os.Stderr, "Usage: maestro test <file> [--device ios|android] [--output dir]")
		return exitFailure
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to open file: %v\n", err)
		return exitFailure
	}
	defer f.Close()

	flow, err := maestro.ParseFlow(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return exitFailure
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to seek file: %v\n", err)
		return exitFailure
	}

	var device core.Device
	if deviceFlag != "" {
		device, err = createDevice(deviceFlag, flow.AppID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to create device: %v\n", err)
			return exitFailure
		}
		defer device.Close()
	}

	translator := maestro.NewTranslator()
	actions, err := translator.TranslateFlow(flow, device)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to translate flow: %v\n", err)
		return exitFailure
	}

	executor := maestro.NewExecutor(device, outputFlag)
	if err := executor.Execute(actions, flow.Name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return exitFailure
	}

	fmt.Println("Test completed successfully")
	return exitSuccess
}

func createDevice(platform string, appID string) (core.Device, error) {
	switch core.Platform(platform) {
	case core.IOS:
		if appID == "" {
			return nil, fmt.Errorf("appId required for iOS")
		}
		return ios.NewDevice(appID)
	case core.Android:
		return nil, fmt.Errorf("Android device creation not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

func validateCmd(file string) int {
	if file == "" {
		fmt.Fprintln(os.Stderr, "Error: file path required for validate command")
		fmt.Fprintln(os.Stderr, "Usage: maestro validate <file>")
		return exitFailure
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to open file: %v\n", err)
		return exitFailure
	}
	defer f.Close()

	if err := validateYAML(f); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return exitFailure
	}

	fmt.Printf("Valid: %s\n", file)
	return exitSuccess
}

func validateYAML(r io.Reader) error {
	_, err := maestro.ParseFlow(r)
	return err
}

func printHelp() {
	fmt.Printf(`Maestro CLI - Mobile automation flow runner %s

Usage:
  maestro <command> [arguments] [flags]

Commands:
  test <file>      Execute Maestro flow
  validate <file>  Validate YAML syntax

Flags:
  --device <platform>  Target device: ios or android
  --output <dir>       Directory for screenshots
  --config <file>      Config file path
  --help, -h           Show this help message
  --version            Show version information

Examples:
  maestro validate flow.yaml
  maestro test flow.yaml
  maestro test flow.yaml --device ios --output ./screenshots

Exit codes:
  0  Success
  1  Failure`, version)
}

func init() {
	flag.Usage = func() {}
}

func parseArgs(args []string) []string {
	var remaining []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			if arg == "--device" || arg == "-device" {
				if i+1 < len(args) {
					deviceFlag = args[i+1]
					i++
				}
			} else if arg == "--output" || arg == "-output" {
				if i+1 < len(args) {
					outputFlag = args[i+1]
					i++
				}
			} else if arg == "--config" || arg == "-config" {
				if i+1 < len(args) {
					configFlag = args[i+1]
					i++
				}
			}
		} else {
			remaining = append(remaining, arg)
		}
	}
	return remaining
}
