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
	ErrTesseractNotFound = errors.New("tesseract OCR engine is not available")
)

// getTesseractCommand prepares an exec.Cmd configured for either system or portable tesseract
func getTesseractCommand(args ...string) (*exec.Cmd, error) {
	// 1. Check system tesseract
	if path, err := exec.LookPath("tesseract"); err == nil {
		return exec.Command(path, args...), nil
	}

	// 2. Check portable tesseract in ~/.it-toolbox/tools/tesseract
	home, err := os.UserHomeDir()
	if err == nil {
		localTess := filepath.Join(home, ".it-toolbox", "tools", "tesseract", "usr", "bin", "tesseract")
		if _, err := os.Stat(localTess); err == nil {
			cmd := exec.Command(localTess, args...)
			libDir := filepath.Join(home, ".it-toolbox", "tools", "tesseract", "usr", "lib", "x86_64-linux-gnu")
			tessdataDir := filepath.Join(home, ".it-toolbox", "tools", "tesseract", "usr", "share", "tesseract-ocr", "5", "tessdata")

			env := os.Environ()
			env = append(env, fmt.Sprintf("LD_LIBRARY_PATH=%s:%s", libDir, os.Getenv("LD_LIBRARY_PATH")))
			env = append(env, fmt.Sprintf("TESSDATA_PREFIX=%s", tessdataDir))
			cmd.Env = env
			return cmd, nil
		}
	}

	return nil, ErrTesseractNotFound
}

// ExtractTextFromImage extracts text from a JPG, PNG, or WEBP image using Tesseract OCR
func ExtractTextFromImage(imagePath string, lang string) (string, error) {
	if _, err := os.Stat(imagePath); err != nil {
		return "", fmt.Errorf("%w: %s", ErrFileNotFound, imagePath)
	}

	if lang == "" {
		lang = "eng"
	}

	tmpDir, err := os.MkdirTemp("", "it_ocr_*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	outBase := filepath.Join(tmpDir, "ocr_output")
	args := []string{
		imagePath,
		outBase,
		"-l", lang,
	}

	cmd, err := getTesseractCommand(args...)
	if err != nil {
		return "", err
	}

	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("OCR execution failed: %s: %w", string(out), err)
	}

	txtFile := outBase + ".txt"
	content, err := os.ReadFile(txtFile)
	if err != nil {
		return "", fmt.Errorf("failed to read OCR output file: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}

// ExtractTextFromPDF extracts text from a PDF document using pdftotext, with OCR fallback for scanned pages
func ExtractTextFromPDF(pdfPath string) (string, error) {
	if _, err := os.Stat(pdfPath); err != nil {
		return "", fmt.Errorf("%w: %s", ErrFileNotFound, pdfPath)
	}

	tmpDir, err := os.MkdirTemp("", "it_pdf_ocr_*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	// Step 1: Try pdftotext
	txtOut := filepath.Join(tmpDir, "extracted.txt")
	pdftotextPath, err := exec.LookPath("pdftotext")
	if err == nil {
		cmd := exec.Command(pdftotextPath, pdfPath, txtOut)
		if err := cmd.Run(); err == nil {
			content, _ := os.ReadFile(txtOut)
			trimmed := strings.TrimSpace(string(content))
			if len(trimmed) > 0 {
				return trimmed, nil
			}
		}
	}

	// Step 2: Fallback for scanned/image PDF -> render pages with pdftoppm then run OCR
	pdftoppmPath, err := exec.LookPath("pdftoppm")
	if err != nil {
		return "", errors.New("pdftoppm is required for scanned PDF OCR but was not found")
	}

	pagePrefix := filepath.Join(tmpDir, "page")
	renderCmd := exec.Command(pdftoppmPath, "-png", "-r", "150", pdfPath, pagePrefix)
	if out, err := renderCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to render PDF pages: %s: %w", string(out), err)
	}

	matches, err := filepath.Glob(filepath.Join(tmpDir, "page*.png"))
	if err != nil || len(matches) == 0 {
		return "", errors.New("no rendered pages found in scanned PDF")
	}

	var sb strings.Builder
	for i, pageImg := range matches {
		text, err := ExtractTextFromImage(pageImg, "eng")
		if err == nil && len(text) > 0 {
			sb.WriteString(fmt.Sprintf("--- Halaman %d ---\n%s\n\n", i+1, text))
		}
	}

	return strings.TrimSpace(sb.String()), nil
}
