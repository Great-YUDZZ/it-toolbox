package ui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

// Embed modern typography assets
//
//go:embed assets/fonts/GoogleSans-Regular.ttf
var fontGoogleSansBytes []byte

//go:embed assets/fonts/Gilroy-Bold.ttf
var fontGilroyBoldBytes []byte

//go:embed assets/fonts/NotoSansMono-Regular.ttf
var fontNotoMonoBytes []byte

var (
	// ResourceFontRegular provides clean, modern Google Sans for UI body and labels
	ResourceFontRegular = fyne.NewStaticResource("GoogleSans-Regular.ttf", fontGoogleSansBytes)

	// ResourceFontBold provides punchy, modern Gilroy-Bold for titles, branding, and emphasis
	ResourceFontBold = fyne.NewStaticResource("Gilroy-Bold.ttf", fontGilroyBoldBytes)

	// ResourceFontMono provides crisp Noto Sans Mono for code, hashes, and IP addresses
	ResourceFontMono = fyne.NewStaticResource("NotoSansMono-Regular.ttf", fontNotoMonoBytes)
)
