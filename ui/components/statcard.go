package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

// StatCard displays a bold metric stat box for KPIs (Subnet Sizer, VLSM, File Converter, etc.)
type StatCard struct {
	Widget     fyne.CanvasObject
	valueText  *canvas.Text
	labelTxt   *canvas.Text
	sublabel   *canvas.Text
	borderRect *canvas.Rectangle
}

// NewStatCard builds a modern KPI widget
func NewStatCard(label, initialValue string, valColor color.Color) *StatCard {
	if valColor == nil {
		valColor = color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0xFF} // Default Electric Cyan
	}

	bg := canvas.NewRectangle(color.RGBA{R: 0x11, G: 0x16, B: 0x22, A: 0xFF})
	bg.StrokeColor = color.RGBA{R: 0x21, G: 0x2B, B: 0x3C, A: 0xFF}
	bg.StrokeWidth = 1
	bg.CornerRadius = 8

	lbl := canvas.NewText(label, color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF})
	lbl.TextSize = 10.5
	lbl.TextStyle = fyne.TextStyle{Bold: true}

	val := canvas.NewText(initialValue, valColor)
	val.TextSize = 16
	val.TextStyle = fyne.TextStyle{Bold: true}

	sub := canvas.NewText("", color.RGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF})
	sub.TextSize = 9.5

	content := container.NewVBox(lbl, val, sub)
	stack := container.NewStack(bg, container.NewPadded(content))

	return &StatCard{
		Widget:     stack,
		valueText:  val,
		labelTxt:   lbl,
		sublabel:   sub,
		borderRect: bg,
	}
}

func (s *StatCard) SetValue(val string) {
	s.valueText.Text = val
	s.valueText.Refresh()
}

func (s *StatCard) SetSubtext(txt string) {
	s.sublabel.Text = txt
	s.sublabel.Refresh()
}

func (s *StatCard) SetColor(c color.Color) {
	s.valueText.Color = c
	s.valueText.Refresh()
}
