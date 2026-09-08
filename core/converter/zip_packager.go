package converter

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrNoFilesToZip = errors.New("no files provided to create zip archive")
)

// CreateZipArchive bundles a list of files into a single zip archive
func CreateZipArchive(files []string, zipOutputPath string) error {
	if len(files) == 0 {
		return ErrNoFilesToZip
	}

	if err := os.MkdirAll(filepath.Dir(zipOutputPath), 0755); err != nil {
		return err
	}

	zipFile, err := os.Create(zipOutputPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			continue // skip missing files
		}

		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", file, err)
		}

		info, err := f.Stat()
		if err != nil {
			f.Close()
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			f.Close()
			return err
		}
		header.Name = filepath.Base(file)
		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			f.Close()
			return fmt.Errorf("failed to create entry in zip for %s: %w", file, err)
		}

		if _, err := io.Copy(writer, f); err != nil {
			f.Close()
			return fmt.Errorf("failed to write content to zip for %s: %w", file, err)
		}
		f.Close()
	}

	return nil
}
