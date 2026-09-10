package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/ui/constants"
)

// NewHeroHeader renders an authentic Neo-Brutalist page header with bold typography, sticker tag, and stark divider
func NewHeroHeader(title, description string, tagBadge fyne.CanvasObject) fyne.CanvasObject {
	titleText := canvas.NewText(title, constants.ColorTextPrimary)
	titleText.TextSize = constants.FontSizeDisplay // 22px bold statement
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	descLabel := widget.NewLabel(description)
	descLabel.Wrapping = fyne.TextWrapWord

	titleBox := container.NewVBox(titleText, descLabel)
	var rightItem fyne.CanvasObject
	if tagBadge != nil {
		rightItem = container.NewCenter(tagBadge)
	}
	headerBar := container.NewBorder(nil, nil, nil, rightItem, titleBox)

	sep := widget.NewSeparator()

	return container.NewVBox(
		container.NewPadded(headerBar),
		sep,
	)
}
