package qrcode

import (
	"image/color"
	"os"
	"strings"

	skipqrcode "github.com/skip2/go-qrcode"
)

// RecoveryLevel aliases for error correction
type RecoveryLevel = skipqrcode.RecoveryLevel

const (
	RecoveryLow     RecoveryLevel = skipqrcode.Low     // ~7% recovery
	RecoveryMedium  RecoveryLevel = skipqrcode.Medium  // ~15% recovery
	RecoveryHigh    RecoveryLevel = skipqrcode.High    // ~25% recovery
	RecoveryHighest RecoveryLevel = skipqrcode.Highest // ~30% recovery
)

// ParseRecoveryLevel converts human-friendly label to skip2/go-qrcode RecoveryLevel
func ParseRecoveryLevel(levelStr string) RecoveryLevel {
	clean := strings.ToLower(strings.TrimSpace(levelStr))
	switch {
	case strings.Contains(clean, "rendah") || strings.Contains(clean, "low") || strings.Contains(clean, "7"):
		return RecoveryLow
	case strings.Contains(clean, "maksimal") || strings.Contains(clean, "highest") || strings.Contains(clean, "30"):
		return RecoveryHighest
	case strings.Contains(clean, "tinggi") || strings.Contains(clean, "high") || strings.Contains(clean, "25"):
		return RecoveryHigh
	default:
		// Default to medium (~15%)
		return RecoveryMedium
	}
}

// GenerateQRCode generates a standard black-on-white PNG QR code with default medium error recovery
func GenerateQRCode(content string, size int, level RecoveryLevel) ([]byte, error) {
	if size <= 0 {
		size = 256
	}
	return skipqrcode.Encode(content, level, size)
}

// GenerateCustomQRCode generates a QR code with custom foreground and background colors
func GenerateCustomQRCode(content string, size int, level RecoveryLevel, fg, bg color.Color) ([]byte, error) {
	if size <= 0 {
		size = 256
	}
	qr, err := skipqrcode.New(content, level)
	if err != nil {
		return nil, err
	}
	if fg != nil {
		qr.ForegroundColor = fg
	}
	if bg != nil {
		qr.BackgroundColor = bg
	}
	return qr.PNG(size)
}

// SaveQRCodePNG writes the QR code directly to a file on disk
func SaveQRCodePNG(content string, filePath string, size int, level RecoveryLevel) error {
	pngBytes, err := GenerateQRCode(content, size, level)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, pngBytes, 0644)
}

// ColorScheme represents named foreground and background colors for styling QR codes
type ColorScheme struct {
	Name string
	FG   color.Color
	BG   color.Color
}

// GetAvailableColorSchemes returns predefined palettes matching the application design
func GetAvailableColorSchemes() []ColorScheme {
	return []ColorScheme{
		{
			Name: "Standar (Hitam - Putih)",
			FG:   color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
			BG:   color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		},
		{
			Name: "Cyber Sky (Cyan Gelap - Putih)",
			FG:   color.RGBA{R: 0x02, G: 0x84, B: 0xC7, A: 0xFF},
			BG:   color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		},
		{
			Name: "Tech Violet (Ungu Gelap - Putih)",
			FG:   color.RGBA{R: 0x6D, G: 0x28, B: 0xD9, A: 0xFF},
			BG:   color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		},
		{
			Name: "Emerald (Hijau Gelap - Putih)",
			FG:   color.RGBA{R: 0x04, G: 0x78, B: 0x57, A: 0xFF},
			BG:   color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		},
		{
			Name: "Dark Inverted (Putih - Hitam Gelap)",
			FG:   color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
			BG:   color.RGBA{R: 0x18, G: 0x18, B: 0x1B, A: 0xFF},
		},
	}
}
