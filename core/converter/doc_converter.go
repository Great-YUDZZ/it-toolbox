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
	ErrUnsupportedDocFormat = errors.New("unsupported document conversion format")
	ErrFileNotFound         = errors.New("input file does not exist")
)

// SupportedDocConversions lists supported target formats per source extension
var SupportedDocConversions = map[string][]string{
	".pdf":  {"docx", "txt"},
	".docx": {"pdf", "txt"},
	".xlsx": {"pdf", "csv", "txt"},
	".pptx": {"pdf"},
	".txt":  {"pdf", "docx"},
}

// ConvertDocument converts a document file into the target format using LibreOffice / pdftotext
func ConvertDocument(inputPath, targetFormat, outputDir string) (string, error) {
	if _, err := os.Stat(inputPath); err != nil {
		return "", fmt.Errorf("%w: %s", ErrFileNotFound, inputPath)
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	targetFormat = strings.ToLower(strings.TrimPrefix(targetFormat, "."))

	if outputDir == "" {
		outputDir = filepath.Dir(inputPath)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	expectedOutFile := filepath.Join(outputDir, fmt.Sprintf("%s.%s", baseName, targetFormat))

	// Special case: PDF to TXT -> use pdftotext for maximum speed and accuracy
	if ext == ".pdf" && targetFormat == "txt" {
		cmd := exec.Command("pdftotext", inputPath, expectedOutFile)
		if output, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("pdftotext failed: %s: %w", string(output), err)
		}
		if _, err := os.Stat(expectedOutFile); err == nil {
			return expectedOutFile, nil
		}
	}

	// For LibreOffice headless conversions
	sofficePath, err := exec.LookPath("libreoffice")
	if err != nil {
		sofficePath, err = exec.LookPath("soffice")
		if err != nil {
			return "", errors.New("libreoffice is required for document conversion but was not found on system PATH")
		}
	}

	var args []string
	if ext == ".pdf" && targetFormat == "docx" {
		// PDF to DOCX requires writer_pdf_import infilter
		args = []string{
			"--headless",
			"--infilter=writer_pdf_import",
			"--convert-to", "docx:MS Word 2007 XML",
			"--outdir", outputDir,
			inputPath,
		}
	} else {
		args = []string{
			"--headless",
			"--convert-to", targetFormat,
			"--outdir", outputDir,
			inputPath,
		}
	}

	cmd := exec.Command(sofficePath, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("conversion failed (%s): %s: %w", targetFormat, string(out), err)
	}

	// Verify output exists
	if _, err := os.Stat(expectedOutFile); err == nil {
		return expectedOutFile, nil
	}

	// Search output directory in case LibreOffice altered filename casing
	matches, _ := filepath.Glob(filepath.Join(outputDir, fmt.Sprintf("%s*.%s", baseName, targetFormat)))
	if len(matches) > 0 {
		return matches[0], nil
	}

	return "", fmt.Errorf("conversion finished but output file %s was not generated", expectedOutFile)
}
