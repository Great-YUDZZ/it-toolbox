package reference

import (
	"testing"
)

func TestLoadCommands(t *testing.T) {
	data, err := LoadCommands()
	if err != nil {
		t.Fatalf("LoadCommands failed: %v", err)
	}
	if len(data.Categories) == 0 {
		t.Fatalf("LoadCommands returned 0 categories")
	}

	// Search command test
	results := SearchCommands("git")
	if len(results) == 0 {
		t.Errorf("SearchCommands(\"git\") returned 0 results")
	}
}

func TestLoadPorts(t *testing.T) {
	data, err := LoadPorts()
	if err != nil {
		t.Fatalf("LoadPorts failed: %v", err)
	}
	if len(data.Ports) == 0 {
		t.Fatalf("LoadPorts returned 0 ports")
	}

	// Search port test
	results := SearchPorts("80")
	if len(results) == 0 {
		t.Errorf("SearchPorts(\"80\") returned 0 results")
	}
}

func TestLoadHTTPCodes(t *testing.T) {
	data, err := LoadHTTPCodes()
	if err != nil {
		t.Fatalf("LoadHTTPCodes failed: %v", err)
	}
	if len(data.Codes) == 0 {
		t.Fatalf("LoadHTTPCodes returned 0 codes")
	}

	// Search http test
	results := SearchHTTPCodes("404")
	if len(results) == 0 {
		t.Errorf("SearchHTTPCodes(\"404\") returned 0 results")
	}
}
