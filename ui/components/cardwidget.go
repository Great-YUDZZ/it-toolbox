package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ItemCard is a reusable modern card container displaying a title, badge/tag, description, and optional action
type ItemCard struct {
	Widget fyne.CanvasObject
}

func NewItemCard(title, subtitle, description string, onAction func(), actionIcon fyne.Resource) *ItemCard {
	titleText := canvas.NewText(title, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
	titleText.TextSize = 13.5
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	var headerRight fyne.CanvasObject
	if onAction != nil {
		if actionIcon == nil {
			actionIcon = theme.ContentCopyIcon()
		}
		btn := widget.NewButtonWithIcon("", actionIcon, onAction)
		btn.Importance = widget.LowImportance
		headerRight = btn
	}

	var tagBadge fyne.CanvasObject
	if subtitle != "" {
		tagBadge = BadgeIndigo(subtitle)
	}

	headerLeft := container.NewHBox(titleText)
	if tagBadge != nil {
		headerLeft.Add(tagBadge)
	}

	header := container.NewBorder(nil, nil, headerLeft, headerRight)

	descText := widget.NewLabel(description)
	descText.Wrapping = fyne.TextWrapWord

	bg := canvas.NewRectangle(color.RGBA{R: 0x13, G: 0x18, B: 0x24, A: 0xFF})
	bg.StrokeColor = color.RGBA{R: 0x21, G: 0x2B, B: 0x3C, A: 0xFF}
	bg.StrokeWidth = 1
	bg.CornerRadius = 10

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		descText,
	)

	stack := container.NewStack(bg, container.NewPadded(content))
	return &ItemCard{Widget: stack}
}
