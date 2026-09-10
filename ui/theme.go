package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/yudz/it-toolbox/ui/constants"
)

// ============================================================================
// DESIGN TOKEN SYSTEM — Dynamic Light / Dark Theme Support
// ============================================================================

// Local mirror variables initialized from constants
var (
	ColorBgBase           = constants.ColorBgBase
	ColorBgSidebar        = constants.ColorBgSidebar
	ColorBgCard           = constants.ColorBgCard
	ColorBgCardInner      = constants.ColorBgCardInner
	ColorBgHover          = constants.ColorBgHover
	ColorBorderSubtle     = constants.ColorBorderSubtle
	ColorBorderActive     = constants.ColorBorderActive
	ColorShadow           = constants.ColorShadow
	ColorAccentYellow     = constants.ColorAccentYellow
	ColorAccentCobalt     = constants.ColorAccentCobalt
	ColorAccentCobaltDark = constants.ColorAccentCobaltDark
	ColorAccentCobaltDim  = constants.ColorAccentCobaltDim
	ColorAccentCyan       = constants.ColorAccentCyan
	ColorAccentCyanDim    = constants.ColorAccentCyanDim
	ColorSuccess          = constants.ColorSuccess
	ColorWarning          = constants.ColorWarning
	ColorWarningTangerine = constants.ColorWarningTangerine
	ColorDanger           = constants.ColorDanger
	ColorTechIndigo       = constants.ColorTechIndigo
	ColorInfo             = constants.ColorInfo
	ColorTextPrimary      = constants.ColorTextPrimary
	ColorTextSecondary    = constants.ColorTextSecondary
	ColorTextMuted        = constants.ColorTextMuted
	ColorTextDisabled     = constants.ColorTextDisabled

	ColorCardBgCyan   = constants.ColorCardBgCyan
	ColorCardBgGreen  = constants.ColorCardBgGreen
	ColorCardBgYellow = constants.ColorCardBgYellow
	ColorCardBgOrange = constants.ColorCardBgOrange
	ColorCardBgPurple = constants.ColorCardBgPurple
	ColorCardBgCoral  = constants.ColorCardBgCoral
	ColorCardBgBlue   = constants.ColorCardBgBlue
)

// ApplyTheme synchronizes both constants and ui packages with the chosen mode
func ApplyTheme(isDark bool) {
	constants.SetTheme(isDark)

	ColorBgBase = constants.ColorBgBase
	ColorBgSidebar = constants.ColorBgSidebar
	ColorBgCard = constants.ColorBgCard
	ColorBgCardInner = constants.ColorBgCardInner
	ColorBgHover = constants.ColorBgHover
	ColorBorderSubtle = constants.ColorBorderSubtle
	ColorBorderActive = constants.ColorBorderActive
	ColorShadow = constants.ColorShadow
	ColorAccentYellow = constants.ColorAccentYellow
	ColorAccentCobalt = constants.ColorAccentCobalt
	ColorAccentCobaltDark = constants.ColorAccentCobaltDark
	ColorAccentCobaltDim = constants.ColorAccentCobaltDim
	ColorAccentCyan = constants.ColorAccentCyan
	ColorAccentCyanDim = constants.ColorAccentCyanDim
	ColorSuccess = constants.ColorSuccess
	ColorWarning = constants.ColorWarning
	ColorWarningTangerine = constants.ColorWarningTangerine
	ColorDanger = constants.ColorDanger
	ColorTechIndigo = constants.ColorTechIndigo
	ColorInfo = constants.ColorInfo
	ColorTextPrimary = constants.ColorTextPrimary
	ColorTextSecondary = constants.ColorTextSecondary
	ColorTextMuted = constants.ColorTextMuted
	ColorTextDisabled = constants.ColorTextDisabled

	ColorCardBgCyan = constants.ColorCardBgCyan
	ColorCardBgGreen = constants.ColorCardBgGreen
	ColorCardBgYellow = constants.ColorCardBgYellow
	ColorCardBgOrange = constants.ColorCardBgOrange
	ColorCardBgPurple = constants.ColorCardBgPurple
	ColorCardBgCoral = constants.ColorCardBgCoral
	ColorCardBgBlue = constants.ColorCardBgBlue
}

// === TYPOGRAPHY SCALE ===
const (
	FontSizeDisplay = constants.FontSizeDisplay
	FontSizeH1      = constants.FontSizeH1
	FontSizeH2      = constants.FontSizeH2
	FontSizeH3      = constants.FontSizeH3
	FontSizeBody    = constants.FontSizeBody
	FontSizeSmall   = constants.FontSizeSmall
	FontSizeLabel   = constants.FontSizeLabel
)

// === SPACING SYSTEM ===
const (
	SpaceXS  = constants.SpaceXS
	SpaceSM  = constants.SpaceSM
	SpaceMD  = constants.SpaceMD
	SpaceLG  = constants.SpaceLG
	SpaceXL  = constants.SpaceXL
	SpaceXXL = constants.SpaceXXL
)

// CustomCyanTheme provides a dynamic engineering workbench theme supporting both light and dark
type CustomCyanTheme struct {
	isDark bool
}

var _ fyne.Theme = (*CustomCyanTheme)(nil)

func NewCustomCyanTheme() fyne.Theme {
	return NewCustomTheme(constants.IsDarkTheme)
}

func NewCustomTheme(isDark bool) fyne.Theme {
	return &CustomCyanTheme{isDark: isDark}
}

func (m *CustomCyanTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if m.isDark {
		switch name {
		case theme.ColorNamePrimary:
			return ColorAccentYellow
		case theme.ColorNameHover:
			return ColorBgHover
		case theme.ColorNameFocus:
			return ColorAccentYellow
		case theme.ColorNameSelection:
			return ColorAccentCyanDim
		case theme.ColorNameHyperlink:
			return ColorAccentCyan
		case theme.ColorNameBackground:
			return ColorBgBase
		case theme.ColorNameMenuBackground:
			return ColorBgSidebar
		case theme.ColorNameOverlayBackground:
			return ColorBgCard
		case theme.ColorNameInputBackground:
			return ColorBgCardInner
		case theme.ColorNameButton:
			return ColorBgCard
		case theme.ColorNameForeground:
			return ColorTextPrimary
		case theme.ColorNamePlaceHolder:
			return ColorTextMuted
		case theme.ColorNameDisabled:
			return ColorTextDisabled
		case theme.ColorNameSeparator:
			return ColorBorderSubtle
		case theme.ColorNameSuccess:
			return ColorSuccess
		case theme.ColorNameWarning:
			return ColorWarning
		case theme.ColorNameError:
			return ColorDanger
		case theme.ColorNameShadow:
			return ColorShadow
		default:
			return theme.DefaultTheme().Color(name, theme.VariantDark)
		}
	}

	// Light mode - Neo-Brutalism
	switch name {
	case theme.ColorNamePrimary:
		return color.Black
	case theme.ColorNameHover:
		return ColorBgHover
	case theme.ColorNameFocus:
		return ColorAccentYellow
	case theme.ColorNameSelection:
		return ColorAccentCyanDim
	case theme.ColorNameHyperlink:
		return ColorAccentCobalt
	case theme.ColorNameBackground:
		return ColorBgBase
	case theme.ColorNameMenuBackground:
		return ColorBgSidebar
	case theme.ColorNameOverlayBackground:
		return ColorBgCard
	case theme.ColorNameInputBackground:
		return ColorBgCard
	case theme.ColorNameButton:
		return ColorBgCard
	case theme.ColorNameForeground:
		return ColorTextPrimary
	case theme.ColorNamePlaceHolder:
		return ColorTextMuted
	case theme.ColorNameDisabled:
		return ColorTextDisabled
	case theme.ColorNameSeparator:
		return ColorBorderSubtle
	case theme.ColorNameSuccess:
		return ColorSuccess
	case theme.ColorNameWarning:
		return ColorWarning
	case theme.ColorNameError:
		return ColorDanger
	case theme.ColorNameShadow:
		return ColorShadow
	default:
		return theme.DefaultTheme().Color(name, theme.VariantLight)
	}
}

func (m *CustomCyanTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		if len(fontNotoMonoBytes) > 0 {
			return ResourceFontMono
		}
	}
	if style.Bold {
		if len(fontGilroyBoldBytes) > 0 {
			return ResourceFontBold
		}
	}
	if len(fontGoogleSansBytes) > 0 {
		return ResourceFontRegular
	}
	return theme.DefaultTheme().Font(style)
}

func (m *CustomCyanTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m *CustomCyanTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return FontSizeBody
	case theme.SizeNameHeadingText:
		return FontSizeH1
	case theme.SizeNameSubHeadingText:
		return FontSizeH2
	case theme.SizeNameCaptionText:
		return FontSizeSmall
	case theme.SizeNameButtonRadius, theme.SizeNameInputRadius:
		return 4 // Semi-sharp blocky
	case theme.SizeNameCardRadius:
		return 4 // Semi-sharp blocky
	case theme.SizeNameSelectionRadius:
		return 2
	case theme.SizeNamePadding:
		return SpaceSM
	case theme.SizeNameInnerPadding:
		return SpaceSM
	case theme.SizeNameScrollBar:
		return 6
	case theme.SizeNameScrollBarSmall:
		return 4
	default:
		return theme.DefaultTheme().Size(name)
	}
}
