package maestro

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Configuration file names
const (
	ConfigFileName = "config.yaml"
	MaestroDirName = ".maestro"
)

// ExecutionOrder defines the order of flow execution
type ExecutionOrder string

const (
	ExecutionOrderNatural ExecutionOrder = "natural"
	ExecutionOrderRandom  ExecutionOrder = "random"
)

// PlatformConfig holds platform-specific configuration
type PlatformConfig struct {
	IOS     IOSConfig     `yaml:"ios,omitempty"`
	Android AndroidConfig `yaml:"android,omitempty"`
}

// IOSConfig holds iOS-specific configuration
type IOSConfig struct {
	BundleID   string `yaml:"bundleId,omitempty"`
	DeviceType string `yaml:"deviceType,omitempty"`
	UDID       string `yaml:"udid,omitempty"`
	LaunchArgs string `yaml:"launchArgs,omitempty"`
}

// AndroidConfig holds Android-specific configuration
type AndroidConfig struct {
	Package     string `yaml:"package,omitempty"`
	Activity    string `yaml:"activity,omitempty"`
	DeviceID    string `yaml:"deviceId,omitempty"`
	LaunchFlags string `yaml:"launchFlags,omitempty"`
}

// Config represents the Maestro workspace configuration
type Config struct {
	// Flows specifies which flow files to include (glob patterns)
	Flows []string `yaml:"flows,omitempty"`

	// TestOutputDir specifies the directory for test outputs
	TestOutputDir string `yaml:"testOutputDir,omitempty"`

	// IncludeTags specifies tags to include in test run
	IncludeTags []string `yaml:"includeTags,omitempty"`

	// ExcludeTags specifies tags to exclude from test run
	ExcludeTags []string `yaml:"excludeTags,omitempty"`

	// ExecutionOrder specifies the order of flow execution
	ExecutionOrder ExecutionOrder `yaml:"executionOrder,omitempty"`

	// Platform holds platform-specific configuration
	Platform PlatformConfig `yaml:"platform,omitempty"`
}

// ErrConfigNotFound is returned when no config file is found
var ErrConfigNotFound = errors.New("config file not found")

// LoadConfig loads the Maestro configuration from the workspace
// It searches in the following order:
// 1. ./config.yaml (project root)
// 2. ./.maestro/config.yaml
// 3. ./maestro.yaml
func LoadConfig() (*Config, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get current directory: %w", err)
	}

	return LoadConfigFromDir(cwd)
}

// LoadConfigFromDir loads the configuration from a specific directory
func LoadConfigFromDir(dir string) (*Config, error) {
	configPaths := []string{
		filepath.Join(dir, ConfigFileName),
		filepath.Join(dir, MaestroDirName, ConfigFileName),
		filepath.Join(dir, "maestro.yaml"),
	}

	for _, path := range configPaths {
		config, err := LoadConfigFromPath(path)
		if err == nil {
			return config, nil
		}
		if !errors.Is(err, ErrConfigNotFound) {
			return nil, err
		}
	}

	return nil, ErrConfigNotFound
}

// LoadConfigFromPath loads the configuration from a specific file path
func LoadConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrConfigNotFound
		}
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsingFailed, err)
	}

	// Set defaults
	if config.ExecutionOrder == "" {
		config.ExecutionOrder = ExecutionOrderNatural
	}

	return &config, nil
}

// HasFlows returns true if the config specifies any flow patterns
func (c *Config) HasFlows() bool {
	return len(c.Flows) > 0
}

// GetFlows returns the flow patterns, with a default pattern if none specified
func (c *Config) GetFlows() []string {
	if !c.HasFlows() {
		return []string{"**/*.maestro.yaml"}
	}
	return c.Flows
}

// HasTags returns true if any tag filtering is configured
func (c *Config) HasTags() bool {
	return len(c.IncludeTags) > 0 || len(c.ExcludeTags) > 0
}

// HasPlatformConfig returns true if platform-specific config is set
func (c *Config) HasPlatformConfig() bool {
	return c.Platform.IOS.BundleID != "" ||
		c.Platform.IOS.UDID != "" ||
		c.Platform.Android.Package != "" ||
		c.Platform.Android.DeviceID != ""
}
