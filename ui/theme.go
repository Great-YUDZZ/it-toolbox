package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomCyanTheme provides a world-class, human-friendly dark theme with Deep Obsidian surfaces and Electric Cyan accents
type CustomCyanTheme struct{}

var _ fyne.Theme = (*CustomCyanTheme)(nil)

func NewCustomCyanTheme() fyne.Theme {
	return &CustomCyanTheme{}
}

func (m *CustomCyanTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	// ------------------------------------------------------------------------
	// Peak Modern Cyber-Slate Palette (Raycast & Linear Inspired)
	// ------------------------------------------------------------------------
	cyanPrimary := color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0xFF}  // #00D4FF Electric Cyan
	cyanHover := color.RGBA{R: 0x38, G: 0xDF, B: 0xFF, A: 0xFF}    // #38DFFF Bright Cyan Hover
	cyanFocus := color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0xFF}    // #00D4FF Focus Ring
	selection := color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0x30}    // Translucent glowing cyan

	bgCanvas := color.RGBA{R: 0x0B, G: 0x0F, B: 0x17, A: 0xFF}     // #0B0F17 Deep Obsidian Canvas
	bgSidebar := color.RGBA{R: 0x0F, G: 0x14, B: 0x20, A: 0xFF}    // #0F1420 Elevated Dark Midnight
	bgCard := color.RGBA{R: 0x15, G: 0x1B, B: 0x28, A: 0xFF}       // #151B28 Crisp Surface Card
	bgInput := color.RGBA{R: 0x09, G: 0x0D, B: 0x15, A: 0xFF}      // #090D15 Recessed Input / Terminal
	bgButton := color.RGBA{R: 0x18, G: 0x20, B: 0x30, A: 0xFF}     // #182030 Elevated Pill / Button

	textPrimary := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}   // #FFFFFF Crisp Pure White
	textMuted := color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF}     // #94A3B8 Slate 400 Readable Subtitle
	textDisabled := color.RGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF}  // #64748B Slate 500 Placeholder
	separator := color.RGBA{R: 0x1F, G: 0x29, B: 0x3A, A: 0xFF}     // #1F293A Thin Border Divider

	success := color.RGBA{R: 0x10, G: 0xB9, B: 0x81, A: 0xFF}       // #10B981 Emerald Green
	warning := color.RGBA{R: 0xF5, G: 0x9E, B: 0x0B, A: 0xFF}       // #F59E0B Amber
	errorRed := color.RGBA{R: 0xF4, G: 0x3F, B: 0x5E, A: 0xFF}      // #F43F5E Rose Crimson

	switch name {
	case theme.ColorNamePrimary:
		return cyanPrimary
	case theme.ColorNameHover:
		return cyanHover
	case theme.ColorNameFocus:
		return cyanFocus
	case theme.ColorNameSelection:
		return selection
	case theme.ColorNameHyperlink:
		return cyanPrimary
	case theme.ColorNameBackground:
		return bgCanvas
	case theme.ColorNameMenuBackground:
		return bgSidebar
	case theme.ColorNameOverlayBackground:
		return bgCard
	case theme.ColorNameInputBackground:
		return bgInput
	case theme.ColorNameButton:
		return bgButton
	case theme.ColorNameForeground:
		return textPrimary
	case theme.ColorNamePlaceHolder:
		return textMuted
	case theme.ColorNameDisabled:
		return textDisabled
	case theme.ColorNameSeparator:
		return separator
	case theme.ColorNameSuccess:
		return success
	case theme.ColorNameWarning:
		return warning
	case theme.ColorNameError:
		return errorRed
	case theme.ColorNameShadow:
		return color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x60}
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

func (m *CustomCyanTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		if len(fontNotoMonoBytes) > 0 {
			return ResourceFontMono
		}
	}
	if style.Bold {
		if len(fontGilroyBoldBytes) > 0 {
			return ResourceFontBold
		}
	}
	if len(fontGoogleSansBytes) > 0 {
		return ResourceFontRegular
	}
	return theme.DefaultTheme().Font(style)
}

func (m *CustomCyanTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *CustomCyanTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 13.5 // Modern crisp UI typography size (not too chunky, not too tiny)
	case theme.SizeNameHeadingText:
		return 20 // Bold prominent titles
	case theme.SizeNameSubHeadingText:
		return 15 // Clean subtitles
	case theme.SizeNameCaptionText:
		return 11 // Refined small labels
	case theme.SizeNameButtonRadius, theme.SizeNameInputRadius:
		return 8 // Rounded smooth corners
	case theme.SizeNameCardRadius:
		return 12 // Modern spacious card radius
	case theme.SizeNameSelectionRadius:
		return 6
	case theme.SizeNamePadding:
		return 8 // Comfortable breathing space
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameScrollBar:
		return 6 // Slim, non-intrusive scrollbar
	case theme.SizeNameScrollBarSmall:
		return 4
	default:
		return theme.DefaultTheme().Size(name)
	}
}
