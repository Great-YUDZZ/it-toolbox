package converter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ErrNoInputPDFs = errors.New("at least one input PDF file is required")
	ErrInvalidPageRange = errors.New("invalid page range specified for PDF split")
)

// MergePDFs combines multiple PDF files into a single document using Ghostscript
func MergePDFs(inputPaths []string, outputPath string) error {
	if len(inputPaths) == 0 {
		return ErrNoInputPDFs
	}

	for _, p := range inputPaths {
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("%w: %s", ErrFileNotFound, p)
		}
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	gsPath, err := exec.LookPath("gs")
	if err != nil {
		return errors.New("ghostscript (gs) is required for PDF merge but was not found on PATH")
	}

	args := []string{
		"-dBATCH",
		"-dNOPAUSE",
		"-q",
		"-sDEVICE=pdfwrite",
		fmt.Sprintf("-sOutputFile=%s", outputPath),
	}
	args = append(args, inputPaths...)

	cmd := exec.Command(gsPath, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to merge PDFs: %s: %w", string(out), err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("merged PDF was not generated at %s", outputPath)
	}

	return nil
}

// SplitPDF extracts a page range from a PDF into a new PDF document
func SplitPDF(inputPath string, pageFrom, pageTo int, outputPath string) error {
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("%w: %s", ErrFileNotFound, inputPath)
	}

	if pageFrom <= 0 || pageTo < pageFrom {
		return ErrInvalidPageRange
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	gsPath, err := exec.LookPath("gs")
	if err != nil {
		return errors.New("ghostscript (gs) is required for PDF split but was not found on PATH")
	}

	args := []string{
		"-dBATCH",
		"-dNOPAUSE",
		"-q",
		"-sDEVICE=pdfwrite",
		fmt.Sprintf("-dFirstPage=%d", pageFrom),
		fmt.Sprintf("-dLastPage=%d", pageTo),
		fmt.Sprintf("-sOutputFile=%s", outputPath),
		inputPath,
	}

	cmd := exec.Command(gsPath, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to split PDF: %s: %w", string(out), err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("split PDF was not generated at %s", outputPath)
	}

	return nil
}

// CompressPDF compresses an existing PDF file using Ghostscript PDFSETTINGS
// level can be: "screen" (low-res, smallest), "ebook" (medium), "printer" (high quality)
func CompressPDF(inputPath string, level string, outputPath string) error {
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("%w: %s", ErrFileNotFound, inputPath)
	}

	if level == "" {
		level = "ebook"
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	gsPath, err := exec.LookPath("gs")
	if err != nil {
		return errors.New("ghostscript (gs) is required for PDF compression but was not found on PATH")
	}

	args := []string{
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		fmt.Sprintf("-dPDFSETTINGS=/%s", level),
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		fmt.Sprintf("-sOutputFile=%s", outputPath),
		inputPath,
	}

	cmd := exec.Command(gsPath, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to compress PDF: %s: %w", string(out), err)
	}

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("compressed PDF was not generated at %s", outputPath)
	}

	return nil
}
