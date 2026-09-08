package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPDFTools(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Create a base PDF using doc_converter
	sampleTxt := filepath.Join(tmpDir, "base.txt")
	_ = os.WriteFile(sampleTxt, []byte("Page 1 content\n\nPage 2 content\n\nPage 3 content"), 0644)

	pdf1, err := ConvertDocument(sampleTxt, "pdf", tmpDir)
	if err != nil {
		t.Fatalf("Failed to generate test PDF: %v", err)
	}

	// 2. Test Merge
	mergedPDF := filepath.Join(tmpDir, "merged.pdf")
	err = MergePDFs([]string{pdf1, pdf1}, mergedPDF)
	if err != nil {
		t.Fatalf("MergePDFs failed: %v", err)
	}
	if info, err := os.Stat(mergedPDF); err != nil || info.Size() == 0 {
		t.Fatalf("Merged PDF is missing or empty: %v", err)
	}

	// 3. Test Split
	splitPDF := filepath.Join(tmpDir, "split_page1.pdf")
	err = SplitPDF(mergedPDF, 1, 1, splitPDF)
	if err != nil {
		t.Fatalf("SplitPDF failed: %v", err)
	}
	if info, err := os.Stat(splitPDF); err != nil || info.Size() == 0 {
		t.Fatalf("Split PDF is missing or empty: %v", err)
	}

	// 4. Test Compress
	compressedPDF := filepath.Join(tmpDir, "compressed.pdf")
	err = CompressPDF(mergedPDF, "screen", compressedPDF)
	if err != nil {
		t.Fatalf("CompressPDF failed: %v", err)
	}
	if info, err := os.Stat(compressedPDF); err != nil || info.Size() == 0 {
		t.Fatalf("Compressed PDF is missing or empty: %v", err)
	}

	// 5. Test invalid page range
	errInv := SplitPDF(mergedPDF, 5, 2, filepath.Join(tmpDir, "bad.pdf"))
	if errInv == nil {
		t.Errorf("Expected error for invalid page range, got nil")
	}
}
