package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/liukunup/go-uop/pkg/android"
	"github.com/liukunup/go-uop/pkg/ios"
)

const (
	version     = "0.1.0"
	exitSuccess = 0
	exitFailure = 1
)

var (
	platformFlag string
	serialFlag   string
	addressFlag  string
	appIDFlag    string
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flag.StringVar(&platformFlag, "platform", "", "Target platform (ios|android)")
	flag.StringVar(&serialFlag, "serial", "", "Device serial number (Android)")
	flag.StringVar(&addressFlag, "address", "", "Device address (iOS)")
	flag.StringVar(&appIDFlag, "app", "", "Application ID")

	remainingArgs := parseArgs(args)

	if len(remainingArgs) > 0 && remainingArgs[0] == "--version" {
		fmt.Printf("uop version %s\n", version)
		return exitSuccess
	}

	if len(remainingArgs) == 0 || remainingArgs[0] == "--help" || remainingArgs[0] == "-h" {
		printHelp()
		return exitSuccess
	}

	cmd := remainingArgs[0]

	switch cmd {
	case "devices":
		return devicesCmd()
	case "connect":
		fmt.Println("Connect command - TODO")
		return exitSuccess
	case "screenshot":
		fmt.Println("Screenshot command - TODO")
		return exitSuccess
	case "shell":
		fmt.Println("Shell command - TODO")
		return exitSuccess
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printHelp()
		return exitFailure
	}
}

func devicesCmd() int {
	switch platformFlag {
	case "ios":
		return listIOSDevices()
	case "android":
		return listAndroidDevices()
	default:
		fmt.Println("=== iOS Devices ===")
		listIOSDevices()
		fmt.Println("\n=== Android Devices ===")
		listAndroidDevices()
		return exitSuccess
	}
}

func listIOSDevices() int {
	fmt.Println("iOS devices - use 'uop connect ios --address <addr>' to connect")
	return exitSuccess
}

func listAndroidDevices() int {
	fmt.Println("Android devices - use 'uop connect android --serial <serial>' to connect")
	return exitSuccess
}

func createIOSDevice() (*ios.Device, error) {
	if addressFlag == "" {
		return nil, fmt.Errorf("address is required for iOS (--address)")
	}
	return ios.NewDevice(appIDFlag, ios.WithAddress(addressFlag))
}

func createAndroidDevice() (*android.Device, error) {
	if serialFlag == "" {
		return nil, fmt.Errorf("serial is required for Android (--serial)")
	}
	return android.NewDevice(android.WithUDID(serialFlag), android.WithPackage(appIDFlag))
}

func printHelp() {
	fmt.Printf(`uop CLI - Unified mobile automation %s

Usage:
  uop <command> [arguments] [flags]

Commands:
  devices              List available devices
  connect              Connect to a device
  screenshot           Take a screenshot
  shell                Open interactive shell

Flags:
  --platform <p>       Target platform: ios or android
  --serial <s>         Device serial number (Android)
  --address <a>        Device address (iOS, WebDriverAgent URL)
  --app <id>           Application ID / bundle ID
  --help, -h           Show this help message
  --version            Show version information

Examples:
  uop devices
  uop devices --platform ios
  uop connect --platform ios --address http://localhost:8100 --app com.example.app

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
		if arg == "-platform" || arg == "--platform" {
			if i+1 < len(args) {
				platformFlag = args[i+1]
				i++
			}
		} else if arg == "-serial" || arg == "--serial" {
			if i+1 < len(args) {
				serialFlag = args[i+1]
				i++
			}
		} else if arg == "-address" || arg == "--address" {
			if i+1 < len(args) {
				addressFlag = args[i+1]
				i++
			}
		} else if arg == "-app" || arg == "--app" {
			if i+1 < len(args) {
				appIDFlag = args[i+1]
				i++
			}
		} else if !flag.Parsed() {
			remaining = append(remaining, arg)
		} else {
			remaining = append(remaining, arg)
		}
	}
	flag.Parse()
	return remaining
}
