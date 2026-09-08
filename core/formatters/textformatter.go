package formatters

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrEmptyInput = errors.New("input text is empty")
)

// FormatJSON formats and indents a JSON string with 2 spaces
func FormatJSON(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", ErrEmptyInput
	}

	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(trimmed), "", "  "); err != nil {
		return "", fmt.Errorf("invalid JSON syntax: %w", err)
	}
	return buf.String(), nil
}

// FormatYAML validates and pretty-prints YAML content
func FormatYAML(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", ErrEmptyInput
	}

	var parsed any
	decoder := yaml.NewDecoder(strings.NewReader(trimmed))
	if err := decoder.Decode(&parsed); err != nil {
		return "", fmt.Errorf("invalid YAML syntax: %w", err)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(parsed); err != nil {
		return "", fmt.Errorf("failed to re-encode YAML: %w", err)
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

// FormatXML parses and pretty-formats an XML string with indentation
func FormatXML(input string) (string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", ErrEmptyInput
	}

	decoder := xml.NewDecoder(strings.NewReader(trimmed))
	var buf bytes.Buffer
	encoder := xml.NewEncoder(&buf)
	encoder.Indent("", "  ")

	for {
		token, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", fmt.Errorf("invalid XML syntax: %w", err)
		}
		if err := encoder.EncodeToken(token); err != nil {
			return "", fmt.Errorf("failed to re-encode XML: %w", err)
		}
	}

	if err := encoder.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush formatted XML: %w", err)
	}

	return buf.String(), nil
}
