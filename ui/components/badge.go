package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/yudz/it-toolbox/ui/constants"
)

// Badge renders an authentic Neo-Brutalist sticker badge with solid 2.0px border and bold typography
func Badge(text string, textColor, bgColor, borderColor color.Color) fyne.CanvasObject {
	bg := canvas.NewRectangle(bgColor)
	bg.StrokeColor = borderColor
	bg.StrokeWidth = constants.BorderWidthMedium // 2.0px bold outline
	bg.CornerRadius = constants.CornerRadiusBrutal

	txt := canvas.NewText(text, textColor)
	txt.TextSize = constants.FontSizeLabel // 9.5px
	txt.TextStyle = fyne.TextStyle{Bold: true}
	txt.Alignment = fyne.TextAlignCenter

	padded := container.NewPadded(txt)
	return container.NewStack(bg, padded)
}

// BadgeYellow - Iconic Gumroad/Acid Yellow Sticker
func BadgeYellow(text string) fyne.CanvasObject {
	if constants.IsDarkTheme {
		return Badge(text,
			color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // Pitch Black text
			constants.ColorAccentYellow,                     // #FFE600
			constants.ColorBorderSubtle,                     // #FFFFFF
		)
	}
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // Pitch Black text
		constants.ColorAccentYellow,                     // #FFE600
		constants.ColorBorderSubtle,                     // #000000
	)
}

// BadgeCyan - Vivid Cyber Cyan Sticker
func BadgeCyan(text string) fyne.CanvasObject {
	if constants.IsDarkTheme {
		return Badge(text,
			color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
			constants.ColorAccentCyan, // #00E5FF
			constants.ColorBorderSubtle,
		)
	}
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorAccentCyan, // #00E5FF
		constants.ColorBorderSubtle,
	)
}

// BadgeSuccess - Vivid Neo Mint Sticker
func BadgeSuccess(text string) fyne.CanvasObject {
	if constants.IsDarkTheme {
		return Badge(text,
			color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
			constants.ColorSuccess, // #00F0A0
			constants.ColorBorderSubtle,
		)
	}
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorSuccess, // #00F0A0
		constants.ColorBorderSubtle,
	)
}

// BadgeWarning - Vivid Amber Tangerine Sticker
func BadgeWarning(text string) fyne.CanvasObject {
	if constants.IsDarkTheme {
		return Badge(text,
			color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
			constants.ColorWarning, // #FFB800
			constants.ColorBorderSubtle,
		)
	}
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorWarning, // #FF9F1C
		constants.ColorBorderSubtle,
	)
}

// BadgeDanger - Vivid Neo Red Sticker
func BadgeDanger(text string) fyne.CanvasObject {
	if constants.IsDarkTheme {
		return Badge(text,
			color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
			constants.ColorDanger, // #FF5353
			constants.ColorBorderSubtle,
		)
	}
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		constants.ColorDanger, // #FF4D4D
		constants.ColorBorderSubtle,
	)
}

// BadgeError is an alias for BadgeDanger
func BadgeError(text string) fyne.CanvasObject {
	return BadgeDanger(text)
}

// BadgeMuted - Neutral Retro Box Tag
func BadgeMuted(text string) fyne.CanvasObject {
	if constants.IsDarkTheme {
		return Badge(text,
			constants.ColorTextPrimary,                      // #FFFFFF
			color.RGBA{R: 0x27, G: 0x27, B: 0x30, A: 0xFF}, // Charcoal
			constants.ColorBorderSubtle,
		)
	}
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // #000000
		color.RGBA{R: 0xEE, G: 0xE8, B: 0xDD, A: 0xFF}, // Warm slate/kraft paper
		constants.ColorBorderSubtle,
	)
}

// BadgeIndigo - Tech Violet / Purple Sticker
func BadgeIndigo(text string) fyne.CanvasObject {
	if constants.IsDarkTheme {
		return Badge(text,
			color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
			constants.ColorTechIndigo, // #C084FC
			constants.ColorBorderSubtle,
		)
	}
	return Badge(text,
		color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF},
		color.RGBA{R: 0xC0, G: 0x84, B: 0xFC, A: 0xFF}, // Electric Violet
		constants.ColorBorderSubtle,
	)
}
