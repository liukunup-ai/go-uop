package report

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"
)

type JUnitTestSuites struct {
	XMLName   xml.Name         `xml:"testsuites"`
	Name      string           `xml:"name,attr"`
	Tests     int              `xml:"tests,attr"`
	Failures  int              `xml:"failures,attr"`
	Skipped   int              `xml:"skipped,attr"`
	Errors    int              `xml:"errors,attr"`
	Time      string           `xml:"time,attr"`
	TestSuite []JUnitTestSuite `xml:"testsuite"`
}

type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Skipped   int             `xml:"skipped,attr"`
	Errors    int             `xml:"errors,attr"`
	Time      string          `xml:"time,attr"`
	TestCases []JUnitTestCase `xml:"testcase,omitempty"`
}

type JUnitFailure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
	Type    string   `xml:"type,attr"`
	Content string   `xml:",chardata"`
}

type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Name      string        `xml:"name,attr"`
	Classname string        `xml:"classname,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
	Skipped   string        `xml:"skipped,omitempty"`
}

func (g *Generator) ToJUnitXML() ([]byte, error) {
	result := g.Generate()
	suites := convertToJUnit(result)

	return xml.MarshalIndent(suites, "", "  ")
}

func (g *Generator) ToJUnitXMLFile(path string) error {
	data, err := g.ToJUnitXML()
	if err != nil {
		return fmt.Errorf("generate JUnit XML: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func convertToJUnit(result *SuiteResult) *JUnitTestSuites {
	suites := &JUnitTestSuites{
		Name:      result.Name,
		Tests:     result.TotalTests,
		Failures:  result.FailedTests,
		Skipped:   result.SkippedTests,
		Time:      formatDuration(result.Duration),
		TestSuite: make([]JUnitTestSuite, 0, len(result.Results)),
	}

	for _, test := range result.Results {
		suite := JUnitTestSuite{
			Name:      test.Name,
			Tests:     len(test.Steps),
			Time:      formatDuration(test.Duration),
			TestCases: make([]JUnitTestCase, 0, len(test.Steps)),
		}

		if test.Status == "failed" && test.Error != "" {
			suite.Failures = 1
			suites.Failures++
		}

		suites.TestSuite = append(suites.TestSuite, suite)
	}

	return suites
}

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3f", d.Seconds())
}

func RenderJUnitXML(result *SuiteResult) (string, error) {
	data, err := xml.MarshalIndent(convertToJUnit(result), "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal JUnit XML: %w", err)
	}
	return string(data), nil
}

func init() {
	RegisterFormat("junit", FormatHandler{
		Generate: func(g *Generator) ([]byte, error) {
			return g.ToJUnitXML()
		},
		WriteFile: func(g *Generator, path string) error {
			return g.ToJUnitXMLFile(path)
		},
	})

	RegisterFormat("xml", FormatHandler{
		Generate: func(g *Generator) ([]byte, error) {
			return g.ToJUnitXML()
		},
		WriteFile: func(g *Generator, path string) error {
			return g.ToJUnitXMLFile(path)
		},
	})
}

type MultiFormatExporter struct {
	Generator *Generator
	Configs   map[string]OutputConfig
}

type OutputConfig struct {
	Enabled bool
	Path    string
}

func NewMultiFormatExporter(g *Generator) *MultiFormatExporter {
	return &MultiFormatExporter{
		Generator: g,
		Configs: map[string]OutputConfig{
			"json":  {Enabled: true, Path: "report.json"},
			"html":  {Enabled: true, Path: "report.html"},
			"junit": {Enabled: true, Path: "report.xml"},
		},
	}
}

func (m *MultiFormatExporter) ExportAll() error {
	for format, config := range m.Configs {
		if !config.Enabled {
			continue
		}
		if err := m.Generator.WriteFormat(format, config.Path); err != nil {
			return fmt.Errorf("export %s: %w", format, err)
		}
	}
	return nil
}

func (m *MultiFormatExporter) ExportIfEnabled(format, path string) error {
	if !m.isFormatEnabled(format) {
		return nil
	}
	return m.Generator.WriteFormat(format, path)
}

func (m *MultiFormatExporter) isFormatEnabled(format string) bool {
	if config, ok := m.Configs[format]; ok {
		return config.Enabled
	}
	return strings.HasSuffix(format, ".json") ||
		strings.HasSuffix(format, ".html") ||
		strings.HasSuffix(format, ".xml")
}
