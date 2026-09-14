package youtube

import (
	"testing"
	"time"
)

func TestExtractVideoID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://www.youtube.com/watch?v=aqz-KE-bpKQ", "aqz-KE-bpKQ"},
		{"https://youtu.be/aqz-KE-bpKQ", "aqz-KE-bpKQ"},
		{"https://www.youtube.com/embed/aqz-KE-bpKQ", "aqz-KE-bpKQ"},
		{"https://www.youtube.com/shorts/aqz-KE-bpKQ", "aqz-KE-bpKQ"},
		{"https://m.youtube.com/watch?v=aqz-KE-bpKQ&list=RD12345", "aqz-KE-bpKQ"},
		{"aqz-KE-bpKQ", "aqz-KE-bpKQ"},
		{"https://example.com/not-youtube", ""},
		{"", ""},
	}

	for _, tc := range tests {
		got := ExtractVideoID(tc.input)
		if got != tc.expected {
			t.Errorf("ExtractVideoID(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My: Awesome / Video | 1080p? *great*", "My_Awesome_Video_1080p_great"},
		{"Normal Video Title", "Normal Video Title"},
		{"", "youtube_video"},
		{"::::////", "youtube_video"},
	}

	for _, tc := range tests {
		got := SanitizeFilename(tc.input)
		if got != tc.expected {
			t.Errorf("SanitizeFilename(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d        time.Duration
		expected string
	}{
		{45 * time.Second, "00:45"},
		{10*time.Minute + 35*time.Second, "10:35"},
		{1*time.Hour + 23*time.Minute + 45*time.Second, "01:23:45"},
	}

	for _, tc := range tests {
		got := FormatDuration(tc.d)
		if got != tc.expected {
			t.Errorf("FormatDuration(%v) = %q; want %q", tc.d, got, tc.expected)
		}
	}
}

func TestParseResolutionInt(t *testing.T) {
	tests := []struct {
		label    string
		expected int
	}{
		{"1080p60", 1080},
		{"720p", 720},
		{"480p", 480},
		{"360p", 360},
		{"144p", 144},
		{"audio", 0},
	}

	for _, tc := range tests {
		got := parseResolutionInt(tc.label)
		if got != tc.expected {
			t.Errorf("parseResolutionInt(%q) = %d; want %d", tc.label, got, tc.expected)
		}
	}
}

func TestDownloaderCreation(t *testing.T) {
	dl := NewDownloader()
	if dl == nil || dl.client == nil {
		t.Fatal("expected non-nil Downloader and client")
	}
}

func TestFetchVideoDetailsLive(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live network test in short mode")
	}
	dl := NewDownloader()
	details, err := dl.FetchVideoDetails("aqz-KE-bpKQ")
	if err != nil {
		t.Fatalf("FetchVideoDetails failed: %v", err)
	}

	if details.Title == "" {
		t.Errorf("expected non-empty Title")
	}
	if len(details.Options) == 0 {
		t.Fatalf("expected at least 1 resolution option")
	}

	t.Logf("Video Title: %s", details.Title)
	t.Logf("Author: %s, Duration: %s", details.Author, details.DurationStr)
	t.Logf("Found %d resolution/audio options:", len(details.Options))
	for i, opt := range details.Options {
		t.Logf("  [%d] %s", i+1, opt.DisplayLabel)
	}
}
