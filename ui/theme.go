package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/yudz/it-toolbox/ui/constants"
)

// ============================================================================
// DESIGN TOKEN SYSTEM — Dynamic Theme Support (Neo-Brutalism & Neumorphism)
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
func ApplyTheme(themeName string) {
	constants.SetTheme(themeName)

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

// ApplyThemeBool provides backward compatibility for boolean toggling
func ApplyThemeBool(isDark bool) {
	if isDark {
		ApplyTheme(constants.ThemeNeumorphismDark)
	} else {
		ApplyTheme(constants.ThemeNeoBrutalism)
	}
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

// CustomCyanTheme provides a dynamic engineering workbench theme
type CustomCyanTheme struct {
	themeName string
}

var _ fyne.Theme = (*CustomCyanTheme)(nil)

func NewCustomCyanTheme() fyne.Theme {
	return NewCustomTheme(constants.ActiveTheme)
}

func NewCustomTheme(themeName string) fyne.Theme {
	return &CustomCyanTheme{themeName: themeName}
}

func NewCustomThemeBool(isDark bool) fyne.Theme {
	if isDark {
		return NewCustomTheme(constants.ThemeNeumorphismDark)
	}
	return NewCustomTheme(constants.ThemeNeoBrutalism)
}

func (m *CustomCyanTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch m.themeName {
	case constants.ThemeNeumorphismLight:
		switch name {
		case theme.ColorNameShadow:
			// Translucent dark-slate scrim for modal backdrop blur/dim effect
			return color.NRGBA{R: 0x0F, G: 0x17, B: 0x2A, A: 0x88}
		case theme.ColorNamePrimary:
			return ColorAccentCobalt
		case theme.ColorNameHover:
			return ColorBgHover
		case theme.ColorNameFocus:
			return ColorBorderActive
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
			return color.RGBA{R: 0xE5, G: 0xEB, B: 0xF2, A: 0xFF}
		case theme.ColorNameSuccess:
			return ColorSuccess
		case theme.ColorNameWarning:
			return ColorWarning
		case theme.ColorNameError:
			return ColorDanger
		default:
			return theme.DefaultTheme().Color(name, theme.VariantLight)
		}

	case constants.ThemeNeumorphismDark:
		switch name {
		case theme.ColorNameShadow:
			// Translucent deep obsidian scrim for dark-mode focused glass backdrop
			return color.NRGBA{R: 0x02, G: 0x06, B: 0x12, A: 0xB8}
		case theme.ColorNamePrimary:
			return ColorAccentCobalt
		case theme.ColorNameHover:
			return ColorBgHover
		case theme.ColorNameFocus:
			return ColorBorderActive
		case theme.ColorNameSelection:
			return ColorAccentCobaltDim
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
		default:
			return theme.DefaultTheme().Color(name, theme.VariantDark)
		}

	default: // Neo-Brutalism (Signature Single Mode)
		switch name {
		case theme.ColorNameShadow:
			// Translucent dark scrim for brutalist modal backdrop
			return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x77}
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
		default:
			return theme.DefaultTheme().Color(name, theme.VariantLight)
		}
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
		if constants.IsNeumorphism {
			return 10
		}
		return 4 // Semi-sharp blocky
	case theme.SizeNameCardRadius:
		if constants.IsNeumorphism {
			return 14
		}
		return 4 // Semi-sharp blocky
	case theme.SizeNameSelectionRadius:
		if constants.IsNeumorphism {
			return 6
		}
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
