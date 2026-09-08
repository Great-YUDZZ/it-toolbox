package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CodePanel represents a dark terminal/editor panel with a top action bar
type CodePanel struct {
	Widget fyne.CanvasObject
	Entry  *widget.Entry
}

// NewCodePanel creates a sleek code viewer / terminal output box
func NewCodePanel(title string, initialContent string, win fyne.Window) *CodePanel {
	entry := widget.NewMultiLineEntry()
	entry.SetText(initialContent)
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.Wrapping = fyne.TextWrapWord

	titleLabel := canvas.NewText(title, color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF})
	titleLabel.TextSize = 11
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}

	copyBtn := widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() {
		if win != nil && entry.Text != "" {
			win.Clipboard().SetContent(entry.Text)
		}
	})
	copyBtn.Importance = widget.LowImportance

	topBar := container.NewBorder(nil, nil, titleLabel, copyBtn)

	bg := canvas.NewRectangle(color.RGBA{R: 0x09, G: 0x0D, B: 0x15, A: 0xFF})
	bg.StrokeColor = color.RGBA{R: 0x1E, G: 0x27, B: 0x38, A: 0xFF}
	bg.StrokeWidth = 1
	bg.CornerRadius = 8

	panelContent := container.NewVBox(
		container.NewPadded(topBar),
		widget.NewSeparator(),
		container.NewPadded(entry),
	)

	stack := container.NewStack(bg, panelContent)

	return &CodePanel{
		Widget: stack,
		Entry:  entry,
	}
}
