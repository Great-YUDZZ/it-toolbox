package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/ui/constants"
)

// CodePanel represents a Neo-Brutalist terminal/editor panel with hard shadow and action bar
type CodePanel struct {
	Widget fyne.CanvasObject
	Entry  *widget.Entry
}

// NewCodePanel creates a brutalist code viewer / terminal output box with hard shadow
func NewCodePanel(title string, initialContent string, win fyne.Window) *CodePanel {
	entry := widget.NewMultiLineEntry()
	entry.SetText(initialContent)
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.Wrapping = fyne.TextWrapWord

	titleLabel := canvas.NewText("⌨ "+title, constants.ColorTextPrimary)
	titleLabel.TextSize = constants.FontSizeSmall
	titleLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	copyBtn := widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() {
		if win != nil && entry.Text != "" {
			win.Clipboard().SetContent(entry.Text)
		}
	})
	copyBtn.Importance = widget.LowImportance

	topBar := container.NewBorder(nil, nil, titleLabel, copyBtn)

	bg := canvas.NewRectangle(constants.ColorBgCardInner)
	bg.StrokeColor = constants.ColorBorderSubtle
	bg.StrokeWidth = constants.BorderWidthHeavy
	bg.CornerRadius = constants.CornerRadiusBrutal

	panelContent := container.NewVBox(
		container.NewPadded(topBar),
		widget.NewSeparator(),
		container.NewPadded(entry),
	)

	stack := container.NewStack(bg, panelContent)
	shadowed := WrapHardShadow(stack, constants.ShadowOffsetHeavy, constants.ShadowOffsetHeavy)

	return &CodePanel{
		Widget: shadowed,
		Entry:  entry,
	}
}
