package converter

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertImage(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Create dummy PNG image
	samplePNG := filepath.Join(tmpDir, "test.png")
	f, err := os.Create(samplePNG)
	if err != nil {
		t.Fatalf("Failed to create sample PNG: %v", err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for x := 0; x < 50; x++ {
		for y := 0; y < 50; y++ {
			img.Set(x, y, color.RGBA{R: 0, G: 212, B: 255, A: 255})
		}
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatalf("Failed to encode PNG: %v", err)
	}
	f.Close()

	// 2. Convert PNG to JPG
	jpgOut, err := ConvertImage(samplePNG, "jpg", 80, tmpDir)
	if err != nil {
		t.Fatalf("ConvertImage(png -> jpg) error: %v", err)
	}
	if _, err := os.Stat(jpgOut); err != nil {
		t.Fatalf("JPG output file does not exist: %v", err)
	}

	// 3. Convert JPG to WEBP
	webpOut, err := ConvertImage(jpgOut, "webp", 75, tmpDir)
	if err != nil {
		t.Fatalf("ConvertImage(jpg -> webp) error: %v", err)
	}
	if _, err := os.Stat(webpOut); err != nil {
		t.Fatalf("WEBP output file does not exist: %v", err)
	}

	// 4. Invalid format should fail
	_, errInv := ConvertImage(samplePNG, "invalid_fmt", 80, tmpDir)
	if errInv == nil {
		t.Errorf("Expected error for invalid format, got nil")
	}
}
