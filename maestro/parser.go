package maestro

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseFlow(r io.Reader) (*MaestroFlow, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read reader: %w", err)
	}

	return parseMaestroYAML(data)
}

func ParseFlowFromString(content string) (*MaestroFlow, error) {
	return parseMaestroYAML([]byte(content))
}

func parseMaestroYAML(data []byte) (*MaestroFlow, error) {
	docs := splitDocuments(data)
	if len(docs) == 0 {
		return nil, fmt.Errorf("%w: empty YAML content", ErrParsingFailed)
	}

	content := strings.TrimSpace(string(docs[0]))
	if content == "" {
		return nil, fmt.Errorf("%w: empty YAML content", ErrParsingFailed)
	}

	var flow MaestroFlow
	if err := yaml.Unmarshal(docs[0], &flow); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParsingFailed, err)
	}

	return &flow, nil
}

func splitDocuments(data []byte) [][]byte {
	lines := bytes.Split(data, []byte("\n"))
	var docs [][]byte
	var currentDoc []string

	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if bytes.HasPrefix(trimmed, []byte("---")) {
			if len(currentDoc) > 0 {
				docs = append(docs, []byte(strings.Join(currentDoc, "\n")))
			}
			currentDoc = nil
		} else {
			currentDoc = append(currentDoc, string(line))
		}
	}

	if len(currentDoc) > 0 {
		docs = append(docs, []byte(strings.Join(currentDoc, "\n")))
	}

	return docs
}
