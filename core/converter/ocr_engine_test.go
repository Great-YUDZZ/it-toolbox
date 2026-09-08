package converter

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractTextFromImage(t *testing.T) {
	tmpDir := t.TempDir()
	testImg := filepath.Join(tmpDir, "test_ocr.png")

	// Create test image with text using ImageMagick convert
	cmd := exec.Command("convert", "-background", "white", "-fill", "black", "-pointsize", "30", "label:Hello OCR Test", testImg)
	if err := cmd.Run(); err != nil {
		t.Skip("convert not available to generate test image")
	}

	text, err := ExtractTextFromImage(testImg, "eng")
	if err != nil {
		t.Fatalf("ExtractTextFromImage failed: %v", err)
	}

	if !strings.Contains(strings.ToLower(text), "ocr") {
		t.Errorf("OCR text does not contain expected word 'ocr', got: %q", text)
	}
}
