package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/ui/constants"
)

// hardShadowLayout renders a signature Neo-Brutalism zero-blur hard offset shadow
type hardShadowLayout struct {
	offsetX float32
	offsetY float32
}

func (l *hardShadowLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) < 2 {
		for _, o := range objects {
			o.Resize(size)
			o.Move(fyne.NewPos(0, 0))
		}
		return
	}
	cardW := size.Width - l.offsetX
	cardH := size.Height - l.offsetY
	if cardW < 0 {
		cardW = 0
	}
	if cardH < 0 {
		cardH = 0
	}
	cardSize := fyne.NewSize(cardW, cardH)

	// Layer 0: Hard shadow block offset at (offsetX, offsetY)
	objects[0].Resize(cardSize)
	objects[0].Move(fyne.NewPos(l.offsetX, l.offsetY))

	// Layer 1: Foreground card at (0, 0)
	objects[1].Resize(cardSize)
	objects[1].Move(fyne.NewPos(0, 0))
}

func (l *hardShadowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	if len(objects) < 2 {
		if len(objects) == 1 {
			return objects[0].MinSize()
		}
		return fyne.NewSize(0, 0)
	}
	cardMin := objects[1].MinSize()
	return fyne.NewSize(cardMin.Width+l.offsetX, cardMin.Height+l.offsetY)
}

// WrapHardShadow places an object over a brutalist hard offset shadow
func WrapHardShadow(card fyne.CanvasObject, offsetX, offsetY float32) fyne.CanvasObject {
	shadow := canvas.NewRectangle(constants.ColorShadow)
	shadow.CornerRadius = constants.CornerRadiusBrutal
	if constants.IsDarkTheme {
		shadow.StrokeColor = constants.ColorBorderSubtle
		shadow.StrokeWidth = 1
	}
	return container.New(&hardShadowLayout{offsetX: offsetX, offsetY: offsetY}, shadow, card)
}

// ModernCard wraps canvas objects inside a Neo-Brutalist card with 2.5px solid border and hard offset shadow
type ModernCard struct {
	Widget fyne.CanvasObject
}

// NewModernCard creates a brutalist card with optional header, badge, and actions
func NewModernCard(title string, badge fyne.CanvasObject, content fyne.CanvasObject, action fyne.CanvasObject) *ModernCard {
	return NewModernCardWithAccent(title, "", constants.ColorBorderSubtle, badge, content, action)
}

// NewModernCardWithAccent creates a card with 4px left semantic accent and hard shadow
func NewModernCardWithAccent(title, subtitle string, accentColor color.Color, badge fyne.CanvasObject, content fyne.CanvasObject, action fyne.CanvasObject) *ModernCard {
	if accentColor == nil {
		accentColor = constants.ColorBorderSubtle
	}

	bg := canvas.NewRectangle(constants.ColorBgCard)
	bg.StrokeColor = constants.ColorBorderSubtle
	bg.StrokeWidth = constants.BorderWidthHeavy
	bg.CornerRadius = constants.CornerRadiusBrutal

	leftAccent := canvas.NewRectangle(accentColor)
	leftAccent.SetMinSize(fyne.NewSize(4, 0))

	var cardContent *fyne.Container

	if title != "" || badge != nil || action != nil || subtitle != "" {
		titleText := canvas.NewText(title, constants.ColorTextPrimary)
		titleText.TextSize = constants.FontSizeH2
		titleText.TextStyle = fyne.TextStyle{Bold: true}

		headerLeft := container.NewHBox(titleText)
		if badge != nil {
			headerLeft.Add(badge)
		}

		header := container.NewBorder(nil, nil, headerLeft, action)
		headerItems := []fyne.CanvasObject{header}

		if subtitle != "" {
			subText := canvas.NewText(subtitle, constants.ColorTextMuted)
			subText.TextSize = constants.FontSizeSmall
			headerItems = append(headerItems, subText)
		}

		sep := widget.NewSeparator()
		headerItems = append(headerItems, sep)

		cardContent = container.NewVBox(
			container.NewPadded(container.NewVBox(headerItems...)),
			container.NewPadded(content),
		)
	} else {
		cardContent = container.NewVBox(container.NewPadded(content))
	}

	cardWithBorder := container.NewBorder(nil, nil, leftAccent, nil, cardContent)
	cardFace := container.NewStack(bg, cardWithBorder)
	shadowedCard := WrapHardShadow(cardFace, constants.ShadowOffsetHeavy, constants.ShadowOffsetHeavy)

	return &ModernCard{Widget: shadowedCard}
}

// NewPlainCard creates a Neo-Brutalist card with 2.5px border and hard offset shadow
func NewPlainCard(content fyne.CanvasObject) fyne.CanvasObject {
	return NewPlainCardWithAccent(content, constants.ColorBorderSubtle)
}

// NewPlainCardWithAccent creates a card with a 4px semantic left accent border and hard offset shadow
func NewPlainCardWithAccent(content fyne.CanvasObject, accentColor color.Color) fyne.CanvasObject {
	if accentColor == nil {
		accentColor = constants.ColorBorderSubtle
	}

	bg := canvas.NewRectangle(constants.ColorBgCard)
	bg.StrokeColor = constants.ColorBorderSubtle
	bg.StrokeWidth = constants.BorderWidthHeavy
	bg.CornerRadius = constants.CornerRadiusBrutal

	leftAccent := canvas.NewRectangle(accentColor)
	leftAccent.SetMinSize(fyne.NewSize(4, 0))

	paddedContent := container.NewPadded(content)
	cardWithBorder := container.NewBorder(nil, nil, leftAccent, nil, paddedContent)
	cardFace := container.NewStack(bg, cardWithBorder)

	return WrapHardShadow(cardFace, constants.ShadowOffsetHeavy, constants.ShadowOffsetHeavy)
}
