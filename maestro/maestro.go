package maestro

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const maestroFileSuffix = ".maestro.yaml"

func IsMaestroFile(path string) bool {
	return strings.HasSuffix(path, maestroFileSuffix)
}

func ParseMaestroFlow(path string) (*MaestroFlow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var flow MaestroFlow
	if err := yaml.Unmarshal(data, &flow); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsingFailed, err)
	}

	return &flow, nil
}
