package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/ui/constants"
)

// ItemCard is a reusable Neo-Brutalist card container displaying a title, sticker badge, description, and action
type ItemCard struct {
	Widget fyne.CanvasObject
}

func NewItemCard(title, subtitle, description string, onAction func(), actionIcon fyne.Resource) *ItemCard {
	titleText := canvas.NewText(title, constants.ColorTextPrimary)
	titleText.TextSize = 14
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
		tagBadge = BadgeYellow(subtitle)
	}

	headerLeft := container.NewHBox(titleText)
	if tagBadge != nil {
		headerLeft.Add(tagBadge)
	}

	header := container.NewBorder(nil, nil, headerLeft, headerRight)

	descText := widget.NewLabel(description)
	descText.Wrapping = fyne.TextWrapWord

	bg := canvas.NewRectangle(constants.ColorBgCard)
	bg.StrokeColor = constants.ColorBorderSubtle
	bg.StrokeWidth = constants.BorderWidthHeavy
	bg.CornerRadius = constants.CornerRadiusBrutal

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		descText,
	)

	stack := container.NewStack(bg, container.NewPadded(content))
	shadowed := WrapHardShadow(stack, constants.ShadowOffsetHeavy, constants.ShadowOffsetHeavy)
	return &ItemCard{Widget: shadowed}
}
