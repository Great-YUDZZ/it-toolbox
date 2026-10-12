package qrcode

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFormatFileSize(t *testing.T) {
	cases := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
	}

	for _, c := range cases {
		got := FormatFileSize(c.bytes)
		if got != c.expected {
			t.Errorf("FormatFileSize(%d) = %q; want %q", c.bytes, got, c.expected)
		}
	}
}

func TestGetLocalIPv4Addresses(t *testing.T) {
	ips, err := GetLocalIPv4Addresses()
	if err != nil {
		t.Fatalf("unexpected error getting local IPs: %v", err)
	}
	if len(ips) == 0 {
		t.Fatal("expected at least one IP (e.g. 127.0.0.1 fallback), got empty")
	}
}

func TestFileShareServerLifecycle(t *testing.T) {
	tmpDir := t.TempDir()
	sampleContent := []byte("Hello from IT Toolbox QR Code Share!")
	testFilePath := filepath.Join(tmpDir, "test_share.txt")
	if err := os.WriteFile(testFilePath, sampleContent, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	server := NewFileShareServer()
	if server.IsRunning() {
		t.Fatal("new server should not be running")
	}

	// Start server on 127.0.0.1 and random free port
	shareURL, err := server.Start(testFilePath, 18765, "127.0.0.1")
	if err != nil {
		t.Fatalf("failed to start file server: %v", err)
	}
	defer server.Stop()

	if !server.IsRunning() {
		t.Fatal("server should be marked running")
	}
	if !strings.HasPrefix(shareURL, "http://127.0.0.1:") {
		t.Errorf("unexpected share URL: %s", shareURL)
	}

	// Give a moment for goroutine listener
	time.Sleep(50 * time.Millisecond)

	// 1. Test Landing Page GET /
	resp, err := http.Get(shareURL)
	if err != nil {
		t.Fatalf("failed to GET landing page: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("landing page status: %d; want 200", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "test_share.txt") {
		t.Errorf("expected landing page to contain filename, got: %s", string(body))
	}

	// 2. Test Download GET /download
	downloadURL := strings.TrimSuffix(shareURL, "/") + "/download"
	dlResp, err := http.Get(downloadURL)
	if err != nil {
		t.Fatalf("failed to GET /download: %v", err)
	}
	defer dlResp.Body.Close()
	if dlResp.StatusCode != http.StatusOK {
		t.Errorf("/download status: %d; want 200", dlResp.StatusCode)
	}
	dlData, err := io.ReadAll(dlResp.Body)
	if err != nil {
		t.Fatalf("failed to read download body: %v", err)
	}
	if !bytes.Equal(dlData, sampleContent) {
		t.Errorf("downloaded content mismatch: got %q, want %q", string(dlData), string(sampleContent))
	}

	// 3. Test View GET /view
	viewURL := strings.TrimSuffix(shareURL, "/") + "/view"
	viewResp, err := http.Get(viewURL)
	if err != nil {
		t.Fatalf("failed to GET /view: %v", err)
	}
	defer viewResp.Body.Close()
	if viewResp.StatusCode != http.StatusOK {
		t.Errorf("/view status: %d; want 200", viewResp.StatusCode)
	}

	// 4. Verify file metadata & download count
	name, size, mime, count := server.GetFileInfo()
	if name != "test_share.txt" {
		t.Errorf("GetFileInfo name = %q; want 'test_share.txt'", name)
	}
	if size != int64(len(sampleContent)) {
		t.Errorf("GetFileInfo size = %d; want %d", size, len(sampleContent))
	}
	if !strings.Contains(mime, "text/plain") {
		t.Logf("GetFileInfo mime: %s", mime)
	}
	if count < 1 {
		t.Errorf("expected download count >= 1, got %d", count)
	}

	// 5. Test Stop
	if err := server.Stop(); err != nil {
		t.Errorf("unexpected error stopping server: %v", err)
	}
	if server.IsRunning() {
		t.Error("server should be stopped")
	}
}
