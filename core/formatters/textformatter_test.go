package formatters

import (
	"strings"
	"testing"
)

func TestFormatJSON(t *testing.T) {
	input := `{"name":"alice","age":20,"skills":["go","linux"]}`
	got, err := FormatJSON(input)
	if err != nil {
		t.Fatalf("FormatJSON unexpected error: %v", err)
	}

	if !strings.Contains(got, "  \"name\": \"alice\"") {
		t.Errorf("FormatJSON did not indent properly:\n%s", got)
	}

	// Test invalid JSON
	_, errInv := FormatJSON(`{name: "invalid"}`)
	if errInv == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}
}

func TestFormatYAML(t *testing.T) {
	input := "name: alice\nage: 20\nitems:\n- one\n- two"
	got, err := FormatYAML(input)
	if err != nil {
		t.Fatalf("FormatYAML unexpected error: %v", err)
	}

	if !strings.Contains(got, "name: alice") {
		t.Errorf("FormatYAML output missing fields:\n%s", got)
	}

	// Test invalid YAML
	_, errInv := FormatYAML("foo: [unclosed")
	if errInv == nil {
		t.Errorf("Expected error for invalid YAML, got nil")
	}
}

func TestFormatXML(t *testing.T) {
	input := `<root><user id="1"><name>Bob</name></user></root>`
	got, err := FormatXML(input)
	if err != nil {
		t.Fatalf("FormatXML unexpected error: %v", err)
	}

	if !strings.Contains(got, "<name>Bob</name>") {
		t.Errorf("FormatXML output incorrect:\n%s", got)
	}

	// Test invalid XML
	_, errInv := FormatXML("<root><unclosed></root>")
	if errInv == nil {
		t.Errorf("Expected error for invalid XML, got nil")
	}
}
