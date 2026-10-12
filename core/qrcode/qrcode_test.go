package qrcode

import (
	"bytes"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateQRCode(t *testing.T) {
	tests := []struct {
		name    string
		content string
		size    int
		level   RecoveryLevel
	}{
		{"URL Standard", "https://github.com/Great-YUDZZ/it-toolbox", 256, RecoveryMedium},
		{"Large Size", "https://example.com/stream.mp4", 512, RecoveryHigh},
		{"Small Text", "Hello IT Toolbox", 128, RecoveryLow},
		{"Maximum Recovery", "Critical security token 12345", 256, RecoveryHighest},
	}

	pngHeader := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := GenerateQRCode(tt.content, tt.size, tt.level)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(b) < 8 {
				t.Fatalf("expected png data, got %d bytes", len(b))
			}
			if !bytes.Equal(b[:8], pngHeader) {
				t.Errorf("invalid PNG signature: %v", b[:8])
			}
		})
	}
}

func TestGenerateCustomQRCode(t *testing.T) {
	fg := color.RGBA{R: 0x02, G: 0x84, B: 0xC7, A: 0xFF}
	bg := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	b, err := GenerateCustomQRCode("Custom Color Test", 256, RecoveryMedium, fg, bg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(b) < 100 {
		t.Fatalf("expected valid png bytes, got length %d", len(b))
	}
}

func TestSaveQRCodePNG(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "test_qr.png")

	err := SaveQRCodePNG("https://it-toolbox.local/download", outPath, 256, RecoveryMedium)
	if err != nil {
		t.Fatalf("unexpected error saving qr file: %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
	if info.Size() < 100 {
		t.Errorf("file size unexpectedly small: %d bytes", info.Size())
	}
}

func TestParseRecoveryLevel(t *testing.T) {
	cases := []struct {
		input    string
		expected RecoveryLevel
	}{
		{"Rendah (7%)", RecoveryLow},
		{"low", RecoveryLow},
		{"Sedang (15%)", RecoveryMedium},
		{"Tinggi (25%)", RecoveryHigh},
		{"high", RecoveryHigh},
		{"Maksimal (30%)", RecoveryHighest},
		{"highest", RecoveryHighest},
		{"unknown", RecoveryMedium},
	}

	for _, c := range cases {
		got := ParseRecoveryLevel(c.input)
		if got != c.expected {
			t.Errorf("ParseRecoveryLevel(%q) = %v; want %v", c.input, got, c.expected)
		}
	}
}

func TestGetAvailableColorSchemes(t *testing.T) {
	schemes := GetAvailableColorSchemes()
	if len(schemes) < 3 {
		t.Errorf("expected at least 3 color schemes, got %d", len(schemes))
	}
}
