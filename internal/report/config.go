package report

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Reports   ReportSettings    `json:"reports"`
	Output    OutputSettings    `json:"output"`
	Retention RetentionSettings `json:"retention"`
}

type ReportSettings struct {
	HTML  HTMLSettings  `json:"html"`
	JSON  JSONSettings  `json:"json"`
	JUnit JUnitSettings `json:"junit"`
}

type HTMLSettings struct {
	Enabled  bool   `json:"enabled"`
	Path     string `json:"path"`
	Template string `json:"template,omitempty"`
	Compact  bool   `json:"compact"`
}

type JSONSettings struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
	Indent  bool   `json:"indent"`
}

type JUnitSettings struct {
	Enabled   bool   `json:"enabled"`
	Path      string `json:"path"`
	Package   string `json:"package,omitempty"`
	SkipEmpty bool   `json:"skipEmpty"`
}

type OutputSettings struct {
	Directory string `json:"directory"`
	Filename  string `json:"filename"`
	Overwrite bool   `json:"overwrite"`
}

type RetentionSettings struct {
	MaxReports int `json:"maxReports"`
	MaxAgeDays int `json:"maxAgeDays"`
}

func DefaultConfig() *Config {
	return &Config{
		Reports: ReportSettings{
			HTML: HTMLSettings{
				Enabled: true,
				Path:    "report.html",
				Compact: false,
			},
			JSON: JSONSettings{
				Enabled: true,
				Path:    "report.json",
				Indent:  true,
			},
			JUnit: JUnitSettings{
				Enabled: true,
				Path:    "report.xml",
			},
		},
		Output: OutputSettings{
			Directory: "reports",
			Filename:  "test-report",
			Overwrite: true,
		},
		Retention: RetentionSettings{
			MaxReports: 10,
			MaxAgeDays: 30,
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func (c *Config) IsFormatEnabled(format string) bool {
	switch format {
	case "html":
		return c.Reports.HTML.Enabled
	case "json":
		return c.Reports.JSON.Enabled
	case "junit", "xml":
		return c.Reports.JUnit.Enabled
	default:
		return false
	}
}

func (c *Config) GetFormatPath(format string) string {
	switch format {
	case "html":
		return c.Reports.HTML.Path
	case "json":
		return c.Reports.JSON.Path
	case "junit", "xml":
		return c.Reports.JUnit.Path
	default:
		return "report." + format
	}
}

func (c *Config) GetOutputPath(format string) string {
	filename := c.GetFormatPath(format)
	if c.Output.Directory == "" {
		return filename
	}
	return c.Output.Directory + "/" + filename
}

func (c *Config) ApplyToGenerator(g *Generator) error {
	for _, format := range []string{"html", "json", "junit"} {
		if !c.IsFormatEnabled(format) {
			continue
		}
		path := c.GetOutputPath(format)
		if err := g.WriteFormat(format, path); err != nil {
			return fmt.Errorf("write %s: %w", format, err)
		}
	}
	return nil
}
