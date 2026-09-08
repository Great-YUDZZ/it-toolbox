package formatters

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrEmptyRegexPattern = errors.New("regular expression pattern cannot be empty")
)

// TestRegex compiles a regular expression and finds all matching substrings in the input
func TestRegex(pattern, input string) ([]string, error) {
	if pattern == "" {
		return nil, ErrEmptyRegexPattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex syntax: %w", err)
	}

	matches := re.FindAllString(input, -1)
	if matches == nil {
		return []string{}, nil
	}

	return matches, nil
}
