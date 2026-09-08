package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConvertDocument(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Create sample text file
	sampleText := filepath.Join(tmpDir, "sample.txt")
	err := os.WriteFile(sampleText, []byte("Unit test sample text document for conversion"), 0644)
	if err != nil {
		t.Fatalf("Failed to create sample.txt: %v", err)
	}

	// 2. Convert TXT to PDF
	pdfPath, err := ConvertDocument(sampleText, "pdf", tmpDir)
	if err != nil {
		t.Fatalf("ConvertDocument(txt -> pdf) error: %v", err)
	}
	if _, err := os.Stat(pdfPath); err != nil {
		t.Fatalf("Output PDF does not exist: %v", err)
	}

	// 3. Convert PDF to TXT
	txtOutPath, err := ConvertDocument(pdfPath, "txt", tmpDir)
	if err != nil {
		t.Fatalf("ConvertDocument(pdf -> txt) error: %v", err)
	}
	content, err := os.ReadFile(txtOutPath)
	if err != nil || len(content) == 0 {
		t.Fatalf("Extracted TXT from PDF is empty or missing: %v", err)
	}

	// 4. Non-existent file should error
	_, errNotFound := ConvertDocument("non_existent.docx", "pdf", tmpDir)
	if errNotFound == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}
