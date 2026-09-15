package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/yudz/it-toolbox/ui/constants"
)

// Badge renders a theme-adaptive sticker badge (Blocky sticker for Neo-Brutalism, Pill for Neumorphism/Glassmorphism)
func Badge(text string, textColor, bgColor, borderColor color.Color) fyne.CanvasObject {
	bg := canvas.NewRectangle(bgColor)
	bg.StrokeColor = borderColor
	bg.StrokeWidth = constants.CurrentBorderWidth
	bg.CornerRadius = constants.CurrentBadgeRadius

	txt := canvas.NewText(text, textColor)
	txt.TextSize = constants.FontSizeLabel // 9.5px
	txt.TextStyle = fyne.TextStyle{Bold: true}
	txt.Alignment = fyne.TextAlignCenter

	padded := container.NewPadded(txt)
	return container.NewStack(bg, padded)
}

// BadgeYellow - Amber / Yellow Badge
func BadgeYellow(text string) fyne.CanvasObject {
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		return Badge(text,
			color.RGBA{R: 0x92, G: 0x40, B: 0x0E, A: 0xFF},
			color.RGBA{R: 0xFE, G: 0xF3, B: 0xC7, A: 0xF5},
			color.RGBA{R: 0xFC, G: 0xD3, B: 0x4D, A: 0xCC},
		)
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		return Badge(text,
			color.RGBA{R: 0xFD, G: 0xE6, B: 0x8A, A: 0xFF},
			color.RGBA{R: 0x78, G: 0x35, B: 0x0F, A: 0xFF},
			color.RGBA{R: 0x92, G: 0x40, B: 0x0E, A: 0xFF},
		)
	}
	// Neo-Brutalism Signature Sticker
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorAccentYellow,
		constants.ColorBorderSubtle,
	)
}

// BadgeCyan - Cyber Sky Cyan Badge
func BadgeCyan(text string) fyne.CanvasObject {
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		return Badge(text,
			color.RGBA{R: 0x03, G: 0x69, B: 0xA1, A: 0xFF},
			color.RGBA{R: 0xE0, G: 0xF2, B: 0xFE, A: 0xF5},
			color.RGBA{R: 0x7D, G: 0xD3, B: 0xFC, A: 0xCC},
		)
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		return Badge(text,
			color.RGBA{R: 0x7D, G: 0xD3, B: 0xFC, A: 0xFF},
			color.RGBA{R: 0x0C, G: 0x4A, B: 0x6E, A: 0xFF},
			color.RGBA{R: 0x02, G: 0x84, B: 0xC7, A: 0xFF},
		)
	}
	// Neo-Brutalism Signature Sticker
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorAccentCyan,
		constants.ColorBorderSubtle,
	)
}

// BadgeSuccess - Vivid Mint / Emerald Badge
func BadgeSuccess(text string) fyne.CanvasObject {
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		return Badge(text,
			color.RGBA{R: 0x06, G: 0x5F, B: 0x46, A: 0xFF},
			color.RGBA{R: 0xD1, G: 0xFA, B: 0xE5, A: 0xF5},
			color.RGBA{R: 0x6E, G: 0xE7, B: 0xB7, A: 0xCC},
		)
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		return Badge(text,
			color.RGBA{R: 0x6E, G: 0xE7, B: 0xB7, A: 0xFF},
			color.RGBA{R: 0x06, G: 0x4E, B: 0x3B, A: 0xFF},
			color.RGBA{R: 0x05, G: 0x96, B: 0x69, A: 0xFF},
		)
	}
	// Neo-Brutalism Signature Sticker
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorSuccess,
		constants.ColorBorderSubtle,
	)
}

// BadgeWarning - Tangerine / Orange Badge
func BadgeWarning(text string) fyne.CanvasObject {
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		return Badge(text,
			color.RGBA{R: 0x9A, G: 0x34, B: 0x12, A: 0xFF},
			color.RGBA{R: 0xFF, G: 0xED, B: 0xD5, A: 0xF5},
			color.RGBA{R: 0xFD, G: 0xBA, B: 0x74, A: 0xCC},
		)
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		return Badge(text,
			color.RGBA{R: 0xFD, G: 0xBA, B: 0x74, A: 0xFF},
			color.RGBA{R: 0x43, G: 0x14, B: 0x07, A: 0xFF},
			color.RGBA{R: 0x9A, G: 0x34, B: 0x12, A: 0xFF},
		)
	}
	// Neo-Brutalism Signature Sticker
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorWarning,
		constants.ColorBorderSubtle,
	)
}

// BadgeDanger - Neo Coral Red Badge
func BadgeDanger(text string) fyne.CanvasObject {
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		return Badge(text,
			color.RGBA{R: 0x9F, G: 0x12, B: 0x39, A: 0xFF},
			color.RGBA{R: 0xFF, G: 0xE4, B: 0xE6, A: 0xF5},
			color.RGBA{R: 0xFD, G: 0xA4, B: 0xAF, A: 0xCC},
		)
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		return Badge(text,
			color.RGBA{R: 0xFD, G: 0xA4, B: 0xAF, A: 0xFF},
			color.RGBA{R: 0x4C, G: 0x05, B: 0x19, A: 0xFF},
			color.RGBA{R: 0x9F, G: 0x12, B: 0x39, A: 0xFF},
		)
	}
	// Neo-Brutalism Signature Sticker
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorDanger,
		constants.ColorBorderSubtle,
	)
}

// BadgeError is an alias for BadgeDanger
func BadgeError(text string) fyne.CanvasObject {
	return BadgeDanger(text)
}

// BadgeMuted - Neutral Soft Slate / Kraft Paper Tag
func BadgeMuted(text string) fyne.CanvasObject {
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		return Badge(text,
			color.RGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xFF},
			color.RGBA{R: 0xF1, G: 0xF5, B: 0xF9, A: 0xF5},
			color.RGBA{R: 0xCB, G: 0xD5, B: 0xE1, A: 0xCC},
		)
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		return Badge(text,
			color.RGBA{R: 0xE2, G: 0xE8, B: 0xF0, A: 0xFF},
			color.RGBA{R: 0x1E, G: 0x29, B: 0x3B, A: 0xFF},
			color.RGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xFF},
		)
	}
	// Neo-Brutalism Warm Kraft Tag
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		color.RGBA{R: 0xEE, G: 0xE8, B: 0xDD, A: 0xFF},
		constants.ColorBorderSubtle,
	)
}

// BadgeIndigo - Tech Violet / Purple Badge
func BadgeIndigo(text string) fyne.CanvasObject {
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		return Badge(text,
			color.RGBA{R: 0x37, G: 0x30, B: 0xA3, A: 0xFF},
			color.RGBA{R: 0xEE, G: 0xF2, B: 0xFF, A: 0xF5},
			color.RGBA{R: 0xA5, G: 0xB4, B: 0xFC, A: 0xCC},
		)
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		return Badge(text,
			color.RGBA{R: 0xC4, G: 0xB5, B: 0xFD, A: 0xFF},
			color.RGBA{R: 0x31, G: 0x2E, B: 0x81, A: 0xFF},
			color.RGBA{R: 0x43, G: 0x38, B: 0xCA, A: 0xFF},
		)
	}
	// Neo-Brutalism Electric Violet
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		color.RGBA{R: 0xC0, G: 0x84, B: 0xFC, A: 0xFF},
		constants.ColorBorderSubtle,
	)
}
