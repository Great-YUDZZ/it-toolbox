package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewHeroHeader renders a top page header with icon badge, bold title, and description
func NewHeroHeader(title, description string, tagBadge fyne.CanvasObject) fyne.CanvasObject {
	titleText := canvas.NewText(title, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
	titleText.TextSize = 18
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	descText := canvas.NewText(description, color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF})
	descText.TextSize = 12

	titleBox := container.NewVBox(titleText, descText)

	headerBar := container.NewBorder(nil, nil, titleBox, tagBadge)
	sep := widget.NewSeparator()

	return container.NewVBox(
		container.NewPadded(headerBar),
		sep,
	)
}
