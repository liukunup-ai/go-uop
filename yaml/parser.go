package yaml

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseFlow(path string) (*Flow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var flow Flow
	if err := yaml.Unmarshal(data, &flow); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	return &flow, nil
}

func ParseSuite(path string) (*TestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var suite TestSuite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	return &suite, nil
}

func ParseFlowFromString(content string) (*Flow, error) {
	var flow Flow
	if err := yaml.Unmarshal([]byte(content), &flow); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	return &flow, nil
}
