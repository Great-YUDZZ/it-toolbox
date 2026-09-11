package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/ui/constants"
)

// ShowBrutalistFormDialog presents an authentic Neo-Brutalist modal dialog
// featuring a 2.5px solid border, colored left accent, bold sticker badge, and styled buttons.
func ShowBrutalistFormDialog(
	win fyne.Window,
	badgeText string,
	accentColor color.Color,
	title string,
	subtitle string,
	formContent fyne.CanvasObject,
	submitText string,
	onSubmit func(),
) dialog.Dialog {
	var d dialog.Dialog

	if accentColor == nil {
		accentColor = constants.ColorAccentCyan
	}

	titleTxt := canvas.NewText(title, constants.ColorTextPrimary)
	titleTxt.TextSize = constants.FontSizeH2
	titleTxt.TextStyle = fyne.TextStyle{Bold: true}

	var badge fyne.CanvasObject
	if badgeText != "" {
		badge = BadgeCyan(badgeText)
	}

	headerLeft := container.NewHBox(titleTxt)
	if badge != nil {
		headerLeft.Add(badge)
	}

	headerBox := container.NewVBox(headerLeft)
	if subtitle != "" {
		subTxt := canvas.NewText(subtitle, constants.ColorTextMuted)
		subTxt.TextSize = constants.FontSizeSmall
		headerBox.Add(subTxt)
	}
	headerBox.Add(widget.NewSeparator())

	if submitText == "" {
		submitText = "Simpan"
	}

	btnCancel := widget.NewButtonWithIcon("Batal", theme.CancelIcon(), func() {
		if d != nil {
			d.Hide()
		}
	})
	btnCancel.Importance = widget.LowImportance

	btnSubmit := widget.NewButtonWithIcon(submitText, theme.ConfirmIcon(), func() {
		if onSubmit != nil {
			onSubmit()
		}
		if d != nil {
			d.Hide()
		}
	})
	btnSubmit.Importance = widget.HighImportance

	actionBar := container.NewBorder(nil, nil, nil, container.NewHBox(btnCancel, btnSubmit))

	dialogBody := container.NewVBox(
		headerBox,
		formContent,
		widget.NewSeparator(),
		actionBar,
	)

	// Wrap in signature Neo-Brutalist card with accent
	card := NewPlainCardWithAccent(dialogBody, accentColor)

	d = dialog.NewCustomWithoutButtons("", card, win)
	d.Show()
	return d
}
