package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// Badge renders a refined pill tag with colored background, border, and text
func Badge(text string, textColor, bgColor, borderColor color.Color) fyne.CanvasObject {
	bg := canvas.NewRectangle(bgColor)
	bg.StrokeColor = borderColor
	bg.StrokeWidth = 1
	bg.CornerRadius = 6

	txt := canvas.NewText(text, textColor)
	txt.TextSize = 10.5
	txt.TextStyle = fyne.TextStyle{Bold: true}
	txt.Alignment = fyne.TextAlignCenter

	// Subtle padding for the pill
	padded := container.NewPadded(txt)
	return container.NewStack(bg, padded)
}

// Preset modern badges
func BadgeCyan(text string) fyne.CanvasObject {
	return Badge(text,
		color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0xFF}, // Text #00D4FF
		color.RGBA{R: 0x00, G: 0x48, B: 0x5C, A: 0x50}, // Bg Tint
		color.RGBA{R: 0x00, G: 0x8C, B: 0xAA, A: 0x80}, // Border
	)
}

func BadgeIndigo(text string) fyne.CanvasObject {
	return Badge(text,
		color.RGBA{R: 0x81, G: 0x8C, B: 0xF8, A: 0xFF}, // Text #818CF8
		color.RGBA{R: 0x31, G: 0x2E, B: 0x81, A: 0x50}, // Bg Tint
		color.RGBA{R: 0x4F, G: 0x46, B: 0xE5, A: 0x80}, // Border
	)
}

func BadgeSuccess(text string) fyne.CanvasObject {
	return Badge(text,
		color.RGBA{R: 0x34, G: 0xD3, B: 0x99, A: 0xFF}, // Text #34D399
		color.RGBA{R: 0x06, G: 0x4E, B: 0x3B, A: 0x50}, // Bg Tint
		color.RGBA{R: 0x05, G: 0x96, B: 0x69, A: 0x80}, // Border
	)
}

func BadgeWarning(text string) fyne.CanvasObject {
	return Badge(text,
		color.RGBA{R: 0xFB, G: 0xBF, B: 0x24, A: 0xFF}, // Text #FBBF24
		color.RGBA{R: 0x78, G: 0x35, B: 0x0F, A: 0x50}, // Bg Tint
		color.RGBA{R: 0xD9, G: 0x77, B: 0x06, A: 0x80}, // Border
	)
}

func BadgeError(text string) fyne.CanvasObject {
	return Badge(text,
		color.RGBA{R: 0xFB, G: 0x71, B: 0x85, A: 0xFF}, // Text #FB7185
		color.RGBA{R: 0x88, G: 0x13, B: 0x37, A: 0x50}, // Bg Tint
		color.RGBA{R: 0xE1, G: 0x1D, B: 0x48, A: 0x80}, // Border
	)
}

func BadgeMuted(text string) fyne.CanvasObject {
	return Badge(text,
		color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF}, // Text #94A3B8
		color.RGBA{R: 0x1E, G: 0x29, B: 0x3B, A: 0x60}, // Bg Tint
		color.RGBA{R: 0x33, G: 0x41, B: 0x55, A: 0x80}, // Border
	)
}
