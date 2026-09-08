package converter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedImgFormat = errors.New("unsupported image format, supported formats: jpg, jpeg, png, webp")
)

// SupportedImageFormats contains all valid image extension conversions
var SupportedImageFormats = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"webp": true,
}

// ConvertImage converts an image between JPG, PNG, and WEBP with customizable quality/compression
func ConvertImage(inputPath, targetFormat string, quality int, outputDir string) (string, error) {
	if _, err := os.Stat(inputPath); err != nil {
		return "", fmt.Errorf("%w: %s", ErrFileNotFound, inputPath)
	}

	targetFormat = strings.ToLower(strings.TrimPrefix(targetFormat, "."))
	if targetFormat == "jpeg" {
		targetFormat = "jpg"
	}

	if !SupportedImageFormats[targetFormat] {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedImgFormat, targetFormat)
	}

	if quality <= 0 || quality > 100 {
		quality = 85
	}

	if outputDir == "" {
		outputDir = filepath.Dir(inputPath)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outFile := filepath.Join(outputDir, fmt.Sprintf("%s.%s", baseName, targetFormat))

	// If source and target are same path, append _converted
	if outFile == inputPath {
		outFile = filepath.Join(outputDir, fmt.Sprintf("%s_compressed.%s", baseName, targetFormat))
	}

	convertBin, err := exec.LookPath("convert")
	if err != nil {
		return "", errors.New("ImageMagick convert binary is required for image conversion but was not found")
	}

	// Use ImageMagick with quality parameter for compression
	args := []string{
		inputPath,
		"-quality", fmt.Sprintf("%d", quality),
		outFile,
	}

	cmd := exec.Command(convertBin, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("image conversion failed: %s: %w", string(out), err)
	}

	if _, err := os.Stat(outFile); err != nil {
		return "", fmt.Errorf("converted image was not created at %s", outFile)
	}

	return outFile, nil
}
