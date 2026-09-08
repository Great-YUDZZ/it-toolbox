package converter

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateZipArchive(t *testing.T) {
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	_ = os.WriteFile(file1, []byte("Content 1"), 0644)
	_ = os.WriteFile(file2, []byte("Content 2"), 0644)

	zipOut := filepath.Join(tmpDir, "output.zip")
	err := CreateZipArchive([]string{file1, file2}, zipOut)
	if err != nil {
		t.Fatalf("CreateZipArchive failed: %v", err)
	}

	// Verify zip contents
	reader, err := zip.OpenReader(zipOut)
	if err != nil {
		t.Fatalf("Failed to open generated zip: %v", err)
	}
	defer reader.Close()

	if len(reader.File) != 2 {
		t.Errorf("Expected 2 files in zip, got %d", len(reader.File))
	}

	// Test empty files list
	errEmpty := CreateZipArchive([]string{}, filepath.Join(tmpDir, "empty.zip"))
	if errEmpty == nil {
		t.Errorf("Expected error for empty files slice, got nil")
	}
}
