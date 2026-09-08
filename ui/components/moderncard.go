package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ModernSurfaceColor defines the card surface tone
var (
	CardBgColor     = color.RGBA{R: 0x13, G: 0x18, B: 0x24, A: 0xFF} // #131824 Deep Slate Card
	CardBorderColor = color.RGBA{R: 0x21, G: 0x2B, B: 0x3C, A: 0xFF} // #212B3C 1px crisp border
	CardHoverBorder = color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0x60} // Cyan tint on focus
)

// ModernCard wraps canvas objects inside a modern glass-slate card with rounded corners and subtle border
type ModernCard struct {
	Widget fyne.CanvasObject
}

// NewModernCard creates a sleek container with optional header, badge, and actions
func NewModernCard(title string, badge fyne.CanvasObject, content fyne.CanvasObject, action fyne.CanvasObject) *ModernCard {
	bg := canvas.NewRectangle(CardBgColor)
	bg.StrokeColor = CardBorderColor
	bg.StrokeWidth = 1
	bg.CornerRadius = 10

	var cardContent *fyne.Container

	if title != "" || badge != nil || action != nil {
		titleText := canvas.NewText(title, color.RGBA{R: 0xF8, G: 0xFA, B: 0xFC, A: 0xFF})
		titleText.TextSize = 14
		titleText.TextStyle = fyne.TextStyle{Bold: true}

		headerLeft := container.NewHBox(titleText)
		if badge != nil {
			headerLeft.Add(badge)
		}

		header := container.NewBorder(nil, nil, headerLeft, action)
		sep := widget.NewSeparator()

		cardContent = container.NewVBox(
			container.NewPadded(header),
			sep,
			container.NewPadded(content),
		)
	} else {
		cardContent = container.NewVBox(container.NewPadded(content))
	}

	stack := container.NewStack(bg, cardContent)
	return &ModernCard{Widget: stack}
}

// NewPlainCard creates a card wrapper with no header line
func NewPlainCard(content fyne.CanvasObject) fyne.CanvasObject {
	bg := canvas.NewRectangle(CardBgColor)
	bg.StrokeColor = CardBorderColor
	bg.StrokeWidth = 1
	bg.CornerRadius = 10

	return container.NewStack(bg, container.NewPadded(content))
}
