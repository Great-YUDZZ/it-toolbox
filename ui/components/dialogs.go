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

	// Use a scrollable container with comfortable default dimensions for expansive form editing
	scrollContent := container.NewVScroll(container.NewPadded(formContent))
	scrollContent.SetMinSize(fyne.NewSize(620, 250))

	// Link any ScrollableMultiLineEntry inside the form to the dialog scrollContent
	var linkScrollers func(co fyne.CanvasObject)
	linkScrollers = func(co fyne.CanvasObject) {
		if co == nil {
			return
		}
		if se, ok := co.(*ScrollableMultiLineEntry); ok {
			se.SetParentScroller(scrollContent)
		}
		if c, ok := co.(*fyne.Container); ok {
			for _, child := range c.Objects {
				linkScrollers(child)
			}
		}
	}
	linkScrollers(formContent)

	dialogBody := container.NewBorder(
		headerBox,
		container.NewVBox(widget.NewSeparator(), actionBar),
		nil,
		nil,
		scrollContent,
	)

	// Wrap in signature card with theme-aware accent and shadow
	card := NewPlainCardWithAccent(dialogBody, accentColor)

	d = dialog.NewCustomWithoutButtons("", card, win)
	d.Resize(fyne.NewSize(660, 430))
	d.Show()
	return d
}

// ShowStyledConfirmDialog presents an authentic, theme-aware confirmation modal
// featuring a distinct accent border, badge, title, formatted content, and BOTH Cancel and Confirm buttons.
func ShowStyledConfirmDialog(
	win fyne.Window,
	badgeText string,
	accentColor color.Color,
	title string,
	subtitle string,
	messageContent fyne.CanvasObject,
	cancelText string,
	confirmText string,
	onConfirm func(),
) dialog.Dialog {
	var d dialog.Dialog

	if accentColor == nil {
		accentColor = constants.ColorAccentCobalt
	}

	titleTxt := canvas.NewText(title, constants.ColorTextPrimary)
	titleTxt.TextSize = constants.FontSizeH2
	titleTxt.TextStyle = fyne.TextStyle{Bold: true}

	headerLeft := container.NewHBox(titleTxt)
	if badgeText != "" {
		headerLeft.Add(BadgeCyan(badgeText))
	}

	headerBox := container.NewVBox(headerLeft)
	if subtitle != "" {
		subTxt := canvas.NewText(subtitle, constants.ColorTextMuted)
		subTxt.TextSize = constants.FontSizeSmall
		headerBox.Add(subTxt)
	}
	headerBox.Add(widget.NewSeparator())

	if cancelText == "" {
		cancelText = "Batal"
	}
	if confirmText == "" {
		confirmText = "Konfirmasi"
	}

	btnCancel := widget.NewButtonWithIcon(cancelText, theme.CancelIcon(), func() {
		if d != nil {
			d.Hide()
		}
	})
	btnCancel.Importance = widget.LowImportance

	btnConfirm := widget.NewButtonWithIcon(confirmText, theme.ConfirmIcon(), func() {
		if d != nil {
			d.Hide()
		}
		if onConfirm != nil {
			onConfirm()
		}
	})
	btnConfirm.Importance = widget.HighImportance

	actionBar := container.NewBorder(nil, nil, nil, container.NewHBox(btnCancel, btnConfirm))

	dialogBody := container.NewBorder(
		headerBox,
		container.NewVBox(widget.NewSeparator(), actionBar),
		nil,
		nil,
		container.NewPadded(messageContent),
	)

	card := NewPlainCardWithAccent(dialogBody, accentColor)

	d = dialog.NewCustomWithoutButtons("", card, win)
	d.Resize(fyne.NewSize(580, 290))
	d.Show()
	return d
}

// ShowStyledInformationDialog presents an informative theme-aware modal
// featuring a rich card with accent border, badge, title, message, and a Tutup/Batal dismiss button.
func ShowStyledInformationDialog(
	win fyne.Window,
	badgeText string,
	accentColor color.Color,
	title string,
	subtitle string,
	messageContent fyne.CanvasObject,
	dismissText string,
	onDismiss func(),
) dialog.Dialog {
	var d dialog.Dialog

	if accentColor == nil {
		accentColor = constants.ColorSuccess
	}

	titleTxt := canvas.NewText(title, constants.ColorTextPrimary)
	titleTxt.TextSize = constants.FontSizeH2
	titleTxt.TextStyle = fyne.TextStyle{Bold: true}

	headerLeft := container.NewHBox(titleTxt)
	if badgeText != "" {
		headerLeft.Add(BadgeSuccess(badgeText))
	}

	headerBox := container.NewVBox(headerLeft)
	if subtitle != "" {
		subTxt := canvas.NewText(subtitle, constants.ColorTextMuted)
		subTxt.TextSize = constants.FontSizeSmall
		headerBox.Add(subTxt)
	}
	headerBox.Add(widget.NewSeparator())

	if dismissText == "" {
		dismissText = "Tutup"
	}

	btnDismiss := widget.NewButtonWithIcon(dismissText, theme.CancelIcon(), func() {
		if d != nil {
			d.Hide()
		}
		if onDismiss != nil {
			onDismiss()
		}
	})
	btnDismiss.Importance = widget.MediumImportance

	actionBar := container.NewBorder(nil, nil, nil, container.NewHBox(btnDismiss))

	dialogBody := container.NewBorder(
		headerBox,
		container.NewVBox(widget.NewSeparator(), actionBar),
		nil,
		nil,
		container.NewPadded(messageContent),
	)

	card := NewPlainCardWithAccent(dialogBody, accentColor)

	d = dialog.NewCustomWithoutButtons("", card, win)
	d.Resize(fyne.NewSize(560, 260))
	d.Show()
	return d
}
