package formatters

import (
	"testing"
)

func TestTestRegex(t *testing.T) {
	tests := []struct {
		pattern  string
		input    string
		wantLen  int
		wantErr  bool
	}{
		{`\d+`, "User 123 has 45 apples and 99 oranges", 3, false},
		{`[a-zA-Z]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`, "Contact test@example.com or info@domain.org", 2, false},
		{`xyz`, "No match here", 0, false},
		{`[unclosed`, "sample", 0, true},
		{"", "sample", 0, true},
	}

	for _, tc := range tests {
		got, err := TestRegex(tc.pattern, tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("TestRegex(%q, %q) error = %v; wantErr %v", tc.pattern, tc.input, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && len(got) != tc.wantLen {
			t.Errorf("TestRegex(%q, %q) returned %d matches; want %d", tc.pattern, tc.input, len(got), tc.wantLen)
		}
	}
}
