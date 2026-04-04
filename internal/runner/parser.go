package runner

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Flow struct {
	Name          string   `yaml:"name"`
	Devices       []Device `yaml:"devices"`
	DefaultDevice string   `yaml:"defaultDevice"`
	Steps         []Step   `yaml:"steps"`
}

type Suite struct {
	Name          string      `yaml:"name"`
	Devices       []Device    `yaml:"devices"`
	DefaultDevice string      `yaml:"defaultDevice"`
	Flows         []SuiteFlow `yaml:"flows"`
}

type SuiteFlow struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Device struct {
	ID     string `yaml:"id"`
	Type   string `yaml:"type"`
	Serial string `yaml:"serial"`
}

type Step map[string]any

func ParseFlow(r io.Reader) (*Flow, error) {
	var flow Flow
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&flow); err != nil {
		return nil, fmt.Errorf("failed to parse flow: %w", err)
	}
	return &flow, nil
}

func ParseSuite(r io.Reader) (*Suite, error) {
	var suite Suite
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&suite); err != nil {
		return nil, fmt.Errorf("failed to parse suite: %w", err)
	}
	return &suite, nil
}

func ParseFlowFile(path string) (*Flow, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open flow file: %w", err)
	}
	defer file.Close()
	return ParseFlow(file)
}

func ParseFlowString(content string) (*Flow, error) {
	return ParseFlow(strings.NewReader(content))
}

func ParseSuiteFile(path string) (*Suite, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open suite file: %w", err)
	}
	defer file.Close()
	return ParseSuite(file)
}
