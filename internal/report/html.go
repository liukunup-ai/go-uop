package report

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

func (g *Generator) ToHTML() (string, error) {
	result := g.Generate()
	return RenderHTML(result)
}

func (g *Generator) ToHTMLFile(path string) error {
	html, err := g.ToHTML()
	if err != nil {
		return fmt.Errorf("generate HTML: %w", err)
	}
	return os.WriteFile(path, []byte(html), 0644)
}

func RenderHTML(result *SuiteResult) (string, error) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Report - {{.Name}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .header {
            background: white;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 20px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .header h1 {
            margin: 0 0 10px 0;
            color: #333;
        }
        .stats {
            display: flex;
            gap: 20px;
        }
        .stat {
            padding: 10px 20px;
            border-radius: 4px;
            background: #f0f0f0;
        }
        .stat.passed { background: #d4edda; color: #155724; }
        .stat.failed { background: #f8d7da; color: #721c24; }
        .stat.skipped { background: #fff3cd; color: #856404; }
        .test-card {
            background: white;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 10px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .test-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            cursor: pointer;
        }
        .test-name {
            font-weight: 600;
            color: #333;
        }
        .test-status {
            padding: 4px 12px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: 600;
        }
        .test-status.passed { background: #d4edda; color: #155724; }
        .test-status.failed { background: #f8d7da; color: #721c24; }
        .test-status.skipped { background: #fff3cd; color: #856404; }
        .test-duration {
            color: #666;
            font-size: 12px;
        }
        .test-error {
            margin-top: 10px;
            padding: 10px;
            background: #fff3f3;
            border-left: 3px solid #dc3545;
            color: #721c24;
            font-family: monospace;
            font-size: 13px;
            white-space: pre-wrap;
        }
        .steps {
            margin-top: 10px;
            padding-left: 20px;
            border-left: 2px solid #e0e0e0;
        }
        .step {
            padding: 5px 0;
            font-size: 13px;
            color: #555;
        }
        .step.failed { color: #dc3545; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Name}}</h1>
            <p>Executed: {{.StartTime.Format "2006-01-02 15:04:05"}}</p>
            <div class="stats">
                <div class="stat">Total: {{.TotalTests}}</div>
                <div class="stat passed">Passed: {{.PassedTests}}</div>
                <div class="stat failed">Failed: {{.FailedTests}}</div>
                <div class="stat skipped">Skipped: {{.SkippedTests}}</div>
                <div class="stat">Duration: {{.Duration.Round 1}}</div>
            </div>
        </div>

        {{range .Results}}
        <div class="test-card">
            <div class="test-header">
                <span class="test-name">{{.Name}}</span>
                <div>
                    <span class="test-status {{.Status}}">{{.Status}}</span>
                    <span class="test-duration">{{.Duration.Round 1}}</span>
                </div>
            </div>
            {{if .Error}}
            <div class="test-error">{{.Error}}</div>
            {{end}}
            {{if .Steps}}
            <div class="steps">
                {{range .Steps}}
                <div class="step {{.Status}}">{{.Name}}: {{.Status}} ({{.Duration.Round 1}}){{if .Error}} - {{.Error}}{{end}}</div>
                {{end}}
            </div>
            {{end}}
        </div>
        {{end}}
    </div>
</body>
</html>`

	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var sb strings.Builder
	if err := t.Execute(&sb, result); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return sb.String(), nil
}

func init() {
	RegisterFormat("html", FormatHandler{
		Generate: func(g *Generator) ([]byte, error) {
			html, err := g.ToHTML()
			return []byte(html), err
		},
		WriteFile: func(g *Generator, path string) error {
			return g.ToHTMLFile(path)
		},
	})
}

type FormatHandler struct {
	Generate  func(*Generator) ([]byte, error)
	WriteFile func(*Generator, string) error
}

var formatHandlers = make(map[string]FormatHandler)

func RegisterFormat(name string, handler FormatHandler) {
	formatHandlers[name] = handler
}

func GetFormatHandler(name string) (FormatHandler, bool) {
	h, ok := formatHandlers[name]
	return h, ok
}

func (g *Generator) WriteFormat(format, path string) error {
	h, ok := GetFormatHandler(format)
	if !ok {
		return fmt.Errorf("unknown format: %s", format)
	}
	return h.WriteFile(g, path)
}

type ReportConfig struct {
	Formats   []string
	OutputDir string
	HTML      HTMLConfig
	JUnit     JUnitConfig
}

type HTMLConfig struct {
	Template string
	Compact  bool
}

type JUnitConfig struct {
	Package   string
	SkipEmpty bool
}

type ReportFormat int

const (
	FormatJSON ReportFormat = iota
	FormatHTML
	FormatJUnit
)

func (f ReportFormat) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatHTML:
		return "html"
	case FormatJUnit:
		return "junit"
	default:
		return "unknown"
	}
}

func FormatFromString(s string) (ReportFormat, error) {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON, nil
	case "html":
		return FormatHTML, nil
	case "junit", "xml":
		return FormatJUnit, nil
	default:
		return 0, fmt.Errorf("unknown format: %s", s)
	}
}
