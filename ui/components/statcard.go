package components

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/yudz/it-toolbox/ui/constants"
)

type boundedLayout struct {
	minWidth float32
}

func (b *boundedLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, o := range objects {
		o.Resize(size)
		o.Move(fyne.NewPos(0, 0))
	}
}

func (b *boundedLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var h float32 = 68
	for _, o := range objects {
		if mh := o.MinSize().Height; mh > h {
			h = mh
		}
	}
	return fyne.NewSize(b.minWidth, h)
}

// ResolveCardBg maps a semantic accent color to its matching vivid Neo-Brutalist background
func ResolveCardBg(c color.Color) color.Color {
	if c == nil {
		return constants.ColorCardBgCyan
	}
	if c == constants.ColorCardBgCyan || c == constants.ColorCardBgGreen ||
		c == constants.ColorCardBgYellow || c == constants.ColorCardBgOrange ||
		c == constants.ColorCardBgPurple || c == constants.ColorCardBgCoral ||
		c == constants.ColorCardBgBlue {
		return c
	}
	r, g, b, _ := c.RGBA()
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)

	// Red / Danger / Coral (e.g. 0xFF385C, 0xFF4D4D)
	if r8 > 0xD0 && g8 < 0x70 && b8 < 0x80 {
		return constants.ColorCardBgCoral
	}
	// Green / Success (e.g. 0x00F076, 0x22E565, 0x10B981)
	if g8 > 0xA0 && r8 < 0x80 {
		return constants.ColorCardBgGreen
	}
	// Yellow / Primary (e.g. 0xFFE600, 0xFDE047)
	if r8 > 0xD0 && g8 > 0xB0 && b8 < 0x60 {
		return constants.ColorCardBgYellow
	}
	// Orange / Warning (e.g. 0xFF9000, 0xFF8800, 0xFF9F1C)
	if r8 > 0xD0 && g8 > 0x60 && b8 < 0x60 {
		return constants.ColorCardBgOrange
	}
	// Purple / Indigo (e.g. 0xA855F7, 0xB87CF8)
	if r8 > 0x70 && b8 > 0xC0 {
		return constants.ColorCardBgPurple
	}
	// Cyan / Info / Blue (e.g. 0x00E0FF, 0x00E5FF, 0x3882F6)
	if b8 > 0xB0 {
		return constants.ColorCardBgCyan
	}
	return constants.ColorCardBgCyan
}

// StatCard displays a bold metric stat box for KPIs with Neo-Brutalist styling
type StatCard struct {
	Widget     fyne.CanvasObject
	valueText  *canvas.Text
	labelTxt   *canvas.Text
	sublabel   *canvas.Text
	accentBar  *canvas.Rectangle
	borderRect *canvas.Rectangle
}

// NewStatCard builds a Neo-Brutalist KPI widget with vivid background, 2.5px border, and solid black fonts
func NewStatCard(label, initialValue string, accentColor color.Color) *StatCard {
	if accentColor == nil {
		accentColor = constants.ColorCardBgCyan
	}

	cardBg := ResolveCardBg(accentColor)
	bg := canvas.NewRectangle(cardBg)
	bg.StrokeColor = constants.ColorBorderSubtle
	bg.StrokeWidth = constants.BorderWidthHeavy
	bg.CornerRadius = constants.CornerRadiusBrutal

	// 4px Left Accent Bar
	accentBar := canvas.NewRectangle(color.Black)
	accentBar.SetMinSize(fyne.NewSize(4, 38))

	textColor := color.Black
	if constants.IsDarkTheme {
		textColor = color.White
	}

	lbl := canvas.NewText(strings.ToUpper(label), textColor)
	lbl.TextSize = constants.FontSizeLabel // 9.5px
	lbl.TextStyle = fyne.TextStyle{Bold: true}

	val := canvas.NewText(initialValue, textColor)
	val.TextSize = 18
	val.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	sub := canvas.NewText("", textColor)
	sub.TextSize = constants.FontSizeSmall // 11px
	sub.TextStyle = fyne.TextStyle{Bold: true}

	content := container.NewVBox(lbl, val, sub)
	paddedContent := container.NewPadded(content)

	// Combine left accent bar and padded content
	cardBody := container.NewBorder(nil, nil, accentBar, nil, paddedContent)
	cardFace := container.NewStack(bg, cardBody)

	// Wrap in Neo-Brutalist hard offset shadow (+3px, +3px)
	shadow := canvas.NewRectangle(constants.ColorShadow)
	shadow.CornerRadius = constants.CornerRadiusBrutal
	if constants.IsDarkTheme {
		shadow.StrokeColor = constants.ColorBorderSubtle
		shadow.StrokeWidth = 1
	}
	shadowedWidget := container.New(&hardShadowLayout{
		offsetX: constants.ShadowOffsetMedium,
		offsetY: constants.ShadowOffsetMedium,
	}, shadow, cardFace)

	// Wrap in boundedLayout so MinSize().Width never forces window to balloon beyond screen
	boundedWidget := container.New(&boundedLayout{minWidth: 140}, shadowedWidget)

	return &StatCard{
		Widget:     boundedWidget,
		valueText:  val,
		labelTxt:   lbl,
		sublabel:   sub,
		accentBar:  accentBar,
		borderRect: bg,
	}
}

func (s *StatCard) SetValue(val string) {
	s.valueText.Text = val
	textColor := color.Black
	if constants.IsDarkTheme {
		textColor = color.White
	}
	s.valueText.Color = textColor

	// Dynamically adjust font size so long values fit comfortably
	if len(val) > 24 {
		s.valueText.TextSize = 13.5
	} else if len(val) > 16 {
		s.valueText.TextSize = 15.5
	} else {
		s.valueText.TextSize = 18
	}
	s.valueText.Refresh()
}

func (s *StatCard) SetSubtext(txt string) {
	s.sublabel.Text = txt
	textColor := color.Black
	if constants.IsDarkTheme {
		textColor = color.White
	}
	s.sublabel.Color = textColor
	s.sublabel.Refresh()
}

func (s *StatCard) SetLabel(lbl string) {
	s.labelTxt.Text = strings.ToUpper(lbl)
	textColor := color.Black
	if constants.IsDarkTheme {
		textColor = color.White
	}
	s.labelTxt.Color = textColor
	s.labelTxt.Refresh()
}

func (s *StatCard) SetColor(c color.Color) {
	textColor := color.Black
	if constants.IsDarkTheme {
		textColor = color.White
	}
	s.valueText.Color = textColor
	s.labelTxt.Color = textColor
	s.sublabel.Color = textColor
	s.valueText.Refresh()
	s.labelTxt.Refresh()
	s.sublabel.Refresh()

	// Update card background with matching vivid color
	s.borderRect.FillColor = ResolveCardBg(c)
	s.borderRect.Refresh()

	if s.accentBar != nil {
		if constants.IsDarkTheme {
			s.accentBar.FillColor = c
		} else {
			s.accentBar.FillColor = color.Black
		}
		s.accentBar.Refresh()
	}
}

func (s *StatCard) SetAccentColor(c color.Color) {
	s.SetColor(c)
}
