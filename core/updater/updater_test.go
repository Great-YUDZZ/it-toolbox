package updater

import (
	"testing"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		candidate string
		baseline  string
		expected  bool
	}{
		{"v1.3.2", "v1.3.1", true},
		{"1.3.2", "1.3.1", true},
		{"1.5.0", "1.4.1", true},
		{"1.4.1", "1.4.0", true},
		{"1.4.0", "1.3.1", true},
		{"2.0.0", "1.9.9", true},
		{"v1.3.1", "1.3.1", false},
		{"1.3.1", "v1.3.1", false},
		{"1.3.0", "1.3.1", false},
		{"1.2.9", "1.3.1", false},
		{"v1.3.2-beta", "v1.3.1", true},
	}

	for _, tt := range tests {
		got := IsNewerVersion(tt.candidate, tt.baseline)
		if got != tt.expected {
			t.Errorf("IsNewerVersion(%q, %q) = %v; want %v", tt.candidate, tt.baseline, got, tt.expected)
		}
	}
}

func TestParseSemVer(t *testing.T) {
	tests := []struct {
		input    string
		expected [3]int
	}{
		{"v1.3.1", [3]int{1, 3, 1}},
		{"1.0.0", [3]int{1, 0, 0}},
		{"v2.10.5", [3]int{2, 10, 5}},
		{"v0.9.1-beta+build123", [3]int{0, 9, 1}},
	}

	for _, tt := range tests {
		got := parseSemVer(tt.input)
		if got != tt.expected {
			t.Errorf("parseSemVer(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}
