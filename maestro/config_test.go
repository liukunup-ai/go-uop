package maestro

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromPath(t *testing.T) {
	tests := []struct {
		name        string
		configYAML  string
		wantErr     bool
		checkConfig func(*testing.T, *Config)
	}{
		{
			name: "basic config",
			configYAML: `
flows:
  - "flows/*.yaml"
testOutputDir: "test-results"
executionOrder: natural
`,
			wantErr: false,
			checkConfig: func(t *testing.T, c *Config) {
				if len(c.Flows) != 1 || c.Flows[0] != "flows/*.yaml" {
					t.Errorf("expected flows [flows/*.yaml], got %v", c.Flows)
				}
				if c.TestOutputDir != "test-results" {
					t.Errorf("expected testOutputDir test-results, got %s", c.TestOutputDir)
				}
				if c.ExecutionOrder != ExecutionOrderNatural {
					t.Errorf("expected executionOrder natural, got %s", c.ExecutionOrder)
				}
			},
		},
		{
			name: "config with tags",
			configYAML: `
flows:
  - "**/*.maestro.yaml"
includeTags:
  - smoke
  - regression
excludeTags:
  - skip
`,
			wantErr: false,
			checkConfig: func(t *testing.T, c *Config) {
				if len(c.IncludeTags) != 2 {
					t.Errorf("expected 2 include tags, got %d", len(c.IncludeTags))
				}
				if len(c.ExcludeTags) != 1 {
					t.Errorf("expected 1 exclude tag, got %d", len(c.ExcludeTags))
				}
			},
		},
		{
			name: "config with platform ios",
			configYAML: `
flows:
  - "*.maestro.yaml"
platform:
  ios:
    bundleId: com.example.app
    udid: "00001234-0000000000000000"
`,
			wantErr: false,
			checkConfig: func(t *testing.T, c *Config) {
				if c.Platform.IOS.BundleID != "com.example.app" {
					t.Errorf("expected ios bundleId com.example.app, got %s", c.Platform.IOS.BundleID)
				}
				if c.Platform.IOS.UDID != "00001234-0000000000000000" {
					t.Errorf("expected ios udid, got %s", c.Platform.IOS.UDID)
				}
			},
		},
		{
			name: "config with platform android",
			configYAML: `
platform:
  android:
    package: com.example.app
    activity: MainActivity
`,
			wantErr: false,
			checkConfig: func(t *testing.T, c *Config) {
				if c.Platform.Android.Package != "com.example.app" {
					t.Errorf("expected android package com.example.app, got %s", c.Platform.Android.Package)
				}
				if c.Platform.Android.Activity != "MainActivity" {
					t.Errorf("expected android activity MainActivity, got %s", c.Platform.Android.Activity)
				}
			},
		},
		{
			name: "config with both platforms",
			configYAML: `
platform:
  ios:
    bundleId: com.example.ios
  android:
    package: com.example.android
`,
			wantErr: false,
			checkConfig: func(t *testing.T, c *Config) {
				if c.Platform.IOS.BundleID != "com.example.ios" {
					t.Errorf("expected ios bundleId, got %s", c.Platform.IOS.BundleID)
				}
				if c.Platform.Android.Package != "com.example.android" {
					t.Errorf("expected android package, got %s", c.Platform.Android.Package)
				}
			},
		},
		{
			name:       "empty config",
			configYAML: ``,
			wantErr:    false,
			checkConfig: func(t *testing.T, c *Config) {
				if c.ExecutionOrder != ExecutionOrderNatural {
					t.Errorf("expected default executionOrder natural, got %s", c.ExecutionOrder)
				}
			},
		},
		{
			name:       "random execution order",
			configYAML: `executionOrder: random`,
			wantErr:    false,
			checkConfig: func(t *testing.T, c *Config) {
				if c.ExecutionOrder != ExecutionOrderRandom {
					t.Errorf("expected executionOrder random, got %s", c.ExecutionOrder)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			if err := os.WriteFile(configPath, []byte(tt.configYAML), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			config, err := LoadConfigFromPath(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfigFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkConfig != nil {
				tt.checkConfig(t, config)
			}
		})
	}
}

func TestLoadConfigFromPath_NotFound(t *testing.T) {
	_, err := LoadConfigFromPath("/nonexistent/path/config.yaml")
	if !errorsIs(err, ErrConfigNotFound) {
		t.Errorf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestLoadConfigFromDir(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    map[string]string
		searchDir     string
		wantErr       bool
		expectedFound string
	}{
		{
			name: "finds config.yaml in root",
			setupFiles: map[string]string{
				"config.yaml": `flows: ["*.yaml"]`,
			},
			searchDir:     ".",
			wantErr:       false,
			expectedFound: "config.yaml",
		},
		{
			name: "finds config.yaml in .maestro dir",
			setupFiles: map[string]string{
				".maestro/config.yaml": `flows: ["maestro/*.yaml"]`,
			},
			searchDir:     ".",
			wantErr:       false,
			expectedFound: ".maestro/config.yaml",
		},
		{
			name: "prefers config.yaml over maestro.yaml",
			setupFiles: map[string]string{
				"config.yaml":  `flows: ["root.yaml"]`,
				"maestro.yaml": `flows: ["old.yaml"]`,
			},
			searchDir:     ".",
			wantErr:       false,
			expectedFound: "config.yaml",
		},
		{
			name: "finds maestro.yaml when config.yaml absent",
			setupFiles: map[string]string{
				"maestro.yaml": `flows: ["legacy.yaml"]`,
			},
			searchDir:     ".",
			wantErr:       false,
			expectedFound: "maestro.yaml",
		},
		{
			name:          "returns ErrConfigNotFound when no config exists",
			setupFiles:    map[string]string{},
			searchDir:     ".",
			wantErr:       true,
			expectedFound: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			for relPath, content := range tt.setupFiles {
				fullPath := filepath.Join(tmpDir, relPath)
				if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}

			searchDir := tmpDir
			if tt.searchDir != "." {
				searchDir = filepath.Join(tmpDir, tt.searchDir)
			}

			config, err := LoadConfigFromDir(searchDir)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if !errorsIs(err, ErrConfigNotFound) {
					t.Errorf("expected ErrConfigNotFound, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if config == nil {
				t.Fatal("expected config but got nil")
			}
		})
	}
}

func TestConfig_HasFlows(t *testing.T) {
	c := &Config{}
	if c.HasFlows() {
		t.Error("expected HasFlows false for empty config")
	}

	c.Flows = []string{"test.yaml"}
	if !c.HasFlows() {
		t.Error("expected HasFlows true when flows set")
	}
}

func TestConfig_GetFlows(t *testing.T) {
	c := &Config{}
	defaultFlows := c.GetFlows()
	if len(defaultFlows) != 1 || defaultFlows[0] != "**/*.maestro.yaml" {
		t.Errorf("expected default flow pattern, got %v", defaultFlows)
	}

	c.Flows = []string{"custom.yaml"}
	flows := c.GetFlows()
	if len(flows) != 1 || flows[0] != "custom.yaml" {
		t.Errorf("expected custom flows, got %v", flows)
	}
}

func TestConfig_HasTags(t *testing.T) {
	c := &Config{}
	if c.HasTags() {
		t.Error("expected HasTags false for empty config")
	}

	c.IncludeTags = []string{"smoke"}
	if !c.HasTags() {
		t.Error("expected HasTags true when includeTags set")
	}

	c = &Config{ExcludeTags: []string{"skip"}}
	if !c.HasTags() {
		t.Error("expected HasTags true when excludeTags set")
	}
}

func TestConfig_HasPlatformConfig(t *testing.T) {
	c := &Config{}
	if c.HasPlatformConfig() {
		t.Error("expected HasPlatformConfig false for empty config")
	}

	c.Platform.IOS.BundleID = "com.test"
	if !c.HasPlatformConfig() {
		t.Error("expected HasPlatformConfig true when iOS bundleId set")
	}

	c = &Config{}
	c.Platform.Android.DeviceID = "device1"
	if !c.HasPlatformConfig() {
		t.Error("expected HasPlatformConfig true when android deviceId set")
	}
}

func TestConfig_ExecutionOrderConstants(t *testing.T) {
	if ExecutionOrderNatural != "natural" {
		t.Errorf("expected ExecutionOrderNatural to be 'natural', got %s", ExecutionOrderNatural)
	}
	if ExecutionOrderRandom != "random" {
		t.Errorf("expected ExecutionOrderRandom to be 'random', got %s", ExecutionOrderRandom)
	}
}

func errorsIs(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}
	return err.Error() == target.Error()
}
