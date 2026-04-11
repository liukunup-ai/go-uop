package runner

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/liukunup/go-uop/internal/selector"
	"gopkg.in/yaml.v3"
)

type TestSuite struct {
	Name          string     `yaml:"name"`
	Description   string     `yaml:"description,omitempty"`
	AppID         string     `yaml:"appId,omitempty"`
	TestOutputDir string     `yaml:"testOutputDir,omitempty"`
	Devices       []Device   `yaml:"devices"`
	DefaultDevice string     `yaml:"defaultDevice,omitempty"`
	TestCases     []TestCase `yaml:"testCases,omitempty"`
}

type TestCase struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Device      string `yaml:"device,omitempty"`
	Steps       []Step `yaml:"steps"`
}

type Step struct {
	Type     string
	Selector *selector.Selector
	Params   map[string]any
}

type Device struct {
	ID      string `yaml:"name"`
	Type    string `yaml:"type,omitempty"`
	Serial  string `yaml:"serial,omitempty"`
	UDID    string `yaml:"udid,omitempty"`
	Dev     string `yaml:"dev,omitempty"`
	Default bool   `yaml:"default,omitempty"`
}

func ParseFlow(r io.Reader) (*TestSuite, error) {
	decoder := yaml.NewDecoder(r)

	var suite TestSuite
	if err := decoder.Decode(&suite); err != nil {
		return nil, fmt.Errorf("failed to parse flow: %w", err)
	}

	docIndex := 0
	for {
		var rawDoc any
		if err := decoder.Decode(&rawDoc); err != nil {
			if strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "overflow") {
				break
			}
			return nil, fmt.Errorf("failed to parse document %d: %w", docIndex+1, err)
		}
		docIndex++

		switch v := rawDoc.(type) {
		case map[string]any:
			if len(v) == 0 {
				continue
			}
			tc, err := parseTestCase(v)
			if err != nil {
				return nil, fmt.Errorf("failed to parse test case in document %d: %w", docIndex+1, err)
			}
			if tc != nil {
				suite.TestCases = append(suite.TestCases, *tc)
			}
		case []any:
			tc := &TestCase{}
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					mergeTestCase(tc, m)
				}
			}
			if tc.Name != "" {
				suite.TestCases = append(suite.TestCases, *tc)
			}
		}
	}

	return &suite, nil
}

func parseTestCase(doc map[string]any) (*TestCase, error) {
	if _, hasName := doc["name"]; !hasName {
		return nil, nil
	}
	tc := &TestCase{Steps: []Step{}}
	mergeTestCase(tc, doc)
	return tc, nil
}

func mergeTestCase(tc *TestCase, doc map[string]any) {
	if name, ok := doc["name"].(string); ok && tc.Name == "" {
		tc.Name = name
	}
	if desc, ok := doc["description"].(string); ok && tc.Description == "" {
		tc.Description = desc
	}
	if dev, ok := doc["device"].(string); ok && tc.Device == "" {
		tc.Device = dev
	}
	for key, value := range doc {
		if key == "name" || key == "description" || key == "device" {
			continue
		}
		if key == "steps" {
			if steps, ok := value.([]any); ok {
				for _, s := range steps {
					if m, ok := s.(map[string]any); ok {
						for stepType, params := range m {
							step := Step{Type: stepType}
							if p, ok := params.(map[string]any); ok {
								step.Selector = parseSelector(p)
								step.Params = filterSelectorParams(p)
							}
							tc.Steps = append(tc.Steps, step)
						}
					}
				}
			}
			continue
		}
		step := Step{Type: key}
		if m, ok := value.(map[string]any); ok {
			step.Selector = parseSelector(m)
			step.Params = filterSelectorParams(m)
		}
		tc.Steps = append(tc.Steps, step)
	}
}

func parseSelector(params map[string]any) *selector.Selector {
	s := &selector.Selector{}

	if v, ok := params["text"].(string); ok {
		s.Text = v
	}
	if v, ok := params["id"].(string); ok {
		s.ID = v
	}
	if v, ok := params["xpath"].(string); ok {
		s.XPath = v
	}
	if v, ok := params["index"].(int); ok {
		s.Index = v
	}
	if v, ok := params["index"].(float64); ok {
		s.Index = int(v)
	}
	if v, ok := params["point"].(string); ok {
		s.Point = v
	}
	if v, ok := params["css"].(string); ok {
		s.CSS = v
	}
	if v, ok := params["traits"].([]any); ok {
		var traits []string
		for _, t := range v {
			if str, ok := t.(string); ok {
				traits = append(traits, str)
			}
		}
		s.Traits = traits
	}
	if v, ok := params["traits"].(string); ok {
		s.Traits = []string{v}
	}
	if v, ok := params["enabled"].(bool); ok {
		s.Enabled = &v
	}
	if v, ok := params["checked"].(bool); ok {
		s.Checked = &v
	}
	if v, ok := params["focused"].(bool); ok {
		s.Focused = &v
	}
	if v, ok := params["selected"].(bool); ok {
		s.Selected = &v
	}
	if v, ok := params["width"].(int); ok {
		s.Width = v
	}
	if v, ok := params["height"].(int); ok {
		s.Height = v
	}
	if v, ok := params["tolerance"].(int); ok {
		s.Tolerance = v
	}
	if v, ok := params["image"].(string); ok {
		s.Image = v
	}
	if v, ok := params["algorithm"].(string); ok {
		s.Algorithm = v
	}
	if v, ok := params["threshold"].(float64); ok {
		s.Threshold = v
	}
	if v, ok := params["above"].(map[string]any); ok {
		s.Type = selector.SelectorTypeAbove
		s.Nested = parseSelector(v)
	}
	if v, ok := params["below"].(map[string]any); ok {
		s.Type = selector.SelectorTypeBelow
		s.Nested = parseSelector(v)
	}
	if v, ok := params["leftOf"].(map[string]any); ok {
		s.Type = selector.SelectorTypeLeftOf
		s.Nested = parseSelector(v)
	}
	if v, ok := params["rightOf"].(map[string]any); ok {
		s.Type = selector.SelectorTypeRightOf
		s.Nested = parseSelector(v)
	}
	if v, ok := params["containsChild"].(map[string]any); ok {
		s.Type = selector.SelectorTypeContainsChild
		s.Nested = parseSelector(v)
	}
	if v, ok := params["childOf"].(map[string]any); ok {
		s.Type = selector.SelectorTypeChildOf
		s.Nested = parseSelector(v)
	}
	if v, ok := params["containsDescendants"].(map[string]any); ok {
		s.Type = selector.SelectorTypeContainsDescendants
		s.Nested = parseSelector(v)
	}

	return s
}

var selectorKeys = map[string]bool{
	"text": true, "id": true, "xpath": true, "index": true, "point": true, "css": true,
	"traits": true, "enabled": true, "checked": true, "focused": true, "selected": true,
	"width": true, "height": true, "tolerance": true,
	"image": true, "algorithm": true, "threshold": true,
	"above": true, "below": true, "leftOf": true, "rightOf": true,
	"containsChild": true, "childOf": true, "containsDescendants": true,
}

func filterSelectorParams(params map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range params {
		if !selectorKeys[k] {
			result[k] = v
		}
	}
	return result
}

func ParseFlowFile(path string) (*TestSuite, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open flow file: %w", err)
	}
	defer file.Close()
	return ParseFlow(file)
}

func ParseFlowString(content string) (*TestSuite, error) {
	return ParseFlow(strings.NewReader(content))
}
