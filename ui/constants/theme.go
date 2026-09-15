package constants

import "image/color"

// ============================================================================
// DESIGN TOKEN SYSTEM — Neo-Brutalism & Neumorphism / Glassmorphism Multi-Theme
// ============================================================================

const (
	ThemeNeoBrutalism     = "neobrutalism"
	ThemeNeumorphismLight = "neumorphism_light"
	ThemeNeumorphismDark  = "neumorphism_dark"
)

var (
	ActiveTheme         = ThemeNeoBrutalism
	IsNeumorphism       = false
	IsDarkTheme         = false
	CurrentCornerRadius = CornerRadiusBrutal
	CurrentBorderWidth  = BorderWidthHeavy
	CurrentBadgeRadius  = CornerRadiusBrutal

	ColorNeumorphLightShadow = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // Top-left highlight
	ColorNeumorphDarkShadow  = color.RGBA{R: 0xBD, G: 0xCE, B: 0xE2, A: 0x90} // Bottom-right shadow
)

// === GEOMETRY CONSTANTS ===
const (
	BorderWidthHeavy     = float32(2.5)  // Bold prominent strokes for Neo-Brutalism
	BorderWidthMedium    = float32(2.0)  // Buttons, badges, inputs
	BorderWidthThin      = float32(1.5)  // Separators, sub-elements
	BorderWidthNeumorph  = float32(1.0)  // Fine crystalline glass border for Neumorph/Glass
	ShadowOffsetHeavy    = float32(4.0)  // Hard offset drop shadow for cards (+4px, +4px)
	ShadowOffsetMedium   = float32(3.0)  // Hard offset drop shadow for badges & KPIs
	CornerRadiusBrutal   = float32(4.0)  // Semi-sharp blocky corner radius
	CornerRadiusNeumorph = float32(14.0) // Smooth soft UI rounded corners
	CornerRadiusPill     = float32(12.0) // Rounded pill for tags & badges
)

// === ACTIVE TOKENS (Dynamically Swapped via SetTheme) ===
var (
	ColorBgBase      = color.RGBA{R: 0xFF, G: 0xFD, B: 0xF8, A: 0xFF} // Warm retro paper canvas
	ColorBgSidebar   = color.RGBA{R: 0xF4, G: 0xEF, B: 0xE6, A: 0xFF} // Tinted brutalist sidebar
	ColorBgCard      = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // Crisp white card surface
	ColorBgCardInner = color.RGBA{R: 0xFF, G: 0xFD, B: 0xF0, A: 0xFF} // Warm cream inner container
	ColorBgHover     = color.RGBA{R: 0xEE, G: 0xE8, B: 0xDD, A: 0xFF} // Hover tint

	ColorBorderSubtle = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Solid Jet Black 2.5px border
	ColorBorderActive = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Solid Jet Black active border
	ColorShadow       = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Solid 100% black hard offset shadow

	ColorAccentYellow     = color.RGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF} // Iconic Electric Acid / Gumroad Yellow
	ColorAccentCobalt     = color.RGBA{R: 0x25, G: 0x63, B: 0xEB, A: 0xFF} // Electric Cobalt Blue
	ColorAccentCobaltDark = color.RGBA{R: 0x1D, G: 0x4E, B: 0xD8, A: 0xFF}
	ColorAccentCobaltDim  = color.RGBA{R: 0xDB, G: 0xEA, B: 0xFE, A: 0xFF}

	ColorAccentCyan    = color.RGBA{R: 0x00, G: 0xE5, B: 0xFF, A: 0xFF} // Vivid Cyber Cyan
	ColorAccentCyanDim = color.RGBA{R: 0xBA, G: 0xE6, B: 0xFD, A: 0xFF}

	ColorSuccess          = color.RGBA{R: 0x00, G: 0xF0, B: 0x76, A: 0xFF} // Vivid Electric Neo Green
	ColorWarning          = color.RGBA{R: 0xFF, G: 0x90, B: 0x00, A: 0xFF} // Vivid Electric Tangerine
	ColorWarningTangerine = color.RGBA{R: 0xFF, G: 0x77, B: 0x00, A: 0xFF} // Punchy Orange
	ColorDanger           = color.RGBA{R: 0xFF, G: 0x38, B: 0x5C, A: 0xFF} // Vivid Electric Coral / Red
	ColorTechIndigo       = color.RGBA{R: 0xA8, G: 0x55, B: 0xF7, A: 0xFF} // Electric Violet / Purple
	ColorInfo             = color.RGBA{R: 0x00, G: 0xE0, B: 0xFF, A: 0xFF} // Electric Sky Cyan

	ColorTextPrimary   = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Pitch Black
	ColorTextSecondary = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Pitch Black
	ColorTextMuted     = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Pitch Black
	ColorTextDisabled  = color.RGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xFF} // Dark Slate

	// === VIBRANT CARD BACKGROUNDS ===
	ColorCardBgCyan   = color.RGBA{R: 0x00, G: 0xE0, B: 0xFF, A: 0xFF} // Vivid Electric Cyan (#00E0FF)
	ColorCardBgGreen  = color.RGBA{R: 0x22, G: 0xE5, B: 0x65, A: 0xFF} // Vivid Electric Kelly Green (#22E565)
	ColorCardBgYellow = color.RGBA{R: 0xFF, G: 0xE6, B: 0xFF, A: 0xFF} // Vivid Electric Acid Yellow (#FFE600)
	ColorCardBgOrange = color.RGBA{R: 0xFF, G: 0x90, B: 0x00, A: 0xFF} // Vivid Electric Tangerine (#FF9000)
	ColorCardBgPurple = color.RGBA{R: 0xB8, G: 0x7C, B: 0xF8, A: 0xFF} // Vivid Electric Purple (#B87CF8)
	ColorCardBgCoral  = color.RGBA{R: 0xFF, G: 0x38, B: 0x5C, A: 0xFF} // Vivid Electric Coral Red (#FF385C)
	ColorCardBgBlue   = color.RGBA{R: 0x38, G: 0x82, B: 0xF6, A: 0xFF} // Vivid Electric Royal Blue (#3882F6)
)

// SetTheme switches all design tokens between Neo-Brutalism and Neumorphism/Glassmorphism (Light / Dark)
func SetTheme(themeName string) {
	ActiveTheme = themeName
	switch themeName {
	case ThemeNeumorphismLight:
		IsNeumorphism = true
		IsDarkTheme = false
		CurrentCornerRadius = CornerRadiusNeumorph
		CurrentBorderWidth = BorderWidthNeumorph
		CurrentBadgeRadius = CornerRadiusPill

		// Glassmorphic Neumorphism (Soft Glass UI) — Luminous Frosted Glass over Cool Ambient Ice
		ColorBgBase = color.RGBA{R: 0xEF, G: 0xF3, B: 0xF8, A: 0xFF}      // Ethereal ice/pearl canvas
		ColorBgSidebar = color.RGBA{R: 0xE6, G: 0xEC, B: 0xF4, A: 0xFF}   // Soft frosted sidebar
		ColorBgCard = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}      // Luminous frosted white glass card
		ColorBgCardInner = color.RGBA{R: 0xF4, G: 0xF7, B: 0xFB, A: 0xFF} // Inset frosted container
		ColorBgHover = color.RGBA{R: 0xE2, G: 0xEB, B: 0xF5, A: 0xFF}     // Frosted hover tint

		ColorBorderSubtle = color.RGBA{R: 0xE2, G: 0xE8, B: 0xF0, A: 0xFF} // Fine crystalline glass hairline
		ColorBorderActive = color.RGBA{R: 0x3B, G: 0x82, B: 0xF6, A: 0xFF} // Luminous Blue focus
		ColorShadow = color.RGBA{R: 0xBD, G: 0xCE, B: 0xE2, A: 0x88}       // Soft airy ambient shadow

		ColorNeumorphLightShadow = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // Top-left brilliant white highlight
		ColorNeumorphDarkShadow = color.RGBA{R: 0xBD, G: 0xCE, B: 0xE2, A: 0x90}  // Bottom-right soft cool depth

		ColorAccentYellow = color.RGBA{R: 0xD9, G: 0x77, B: 0x06, A: 0xFF}     // Warm Amber
		ColorAccentCobalt = color.RGBA{R: 0x25, G: 0x63, B: 0xEB, A: 0xFF}     // Royal Cobalt
		ColorAccentCobaltDark = color.RGBA{R: 0x1D, G: 0x4E, B: 0xD8, A: 0xFF}
		ColorAccentCobaltDim = color.RGBA{R: 0xDB, G: 0xEA, B: 0xFE, A: 0xFF}

		ColorAccentCyan = color.RGBA{R: 0x02, G: 0x84, B: 0xC7, A: 0xFF}    // Deep Sky Cyan
		ColorAccentCyanDim = color.RGBA{R: 0xBA, G: 0xE6, B: 0xFD, A: 0x80}

		ColorSuccess = color.RGBA{R: 0x05, G: 0x96, B: 0x69, A: 0xFF}          // Emerald Green
		ColorWarning = color.RGBA{R: 0xD9, G: 0x77, B: 0x06, A: 0xFF}          // Amber
		ColorWarningTangerine = color.RGBA{R: 0xEA, G: 0x58, B: 0x0C, A: 0xFF} // Burnt Orange
		ColorDanger = color.RGBA{R: 0xDC, G: 0x26, B: 0x26, A: 0xFF}           // Coral Red
		ColorTechIndigo = color.RGBA{R: 0x7C, G: 0x3A, B: 0xED, A: 0xFF}       // Soft Violet
		ColorInfo = color.RGBA{R: 0x02, G: 0x84, B: 0xC7, A: 0xFF}             // Sky Blue

		// High-contrast sharp Slate typography for crystal clear legibility on frosted glass
		ColorTextPrimary = color.RGBA{R: 0x0F, G: 0x17, B: 0x2A, A: 0xFF}   // Slate 900
		ColorTextSecondary = color.RGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xFF} // Slate 700
		ColorTextMuted = color.RGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF}     // Slate 500
		ColorTextDisabled = color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF}  // Slate 400

		ColorCardBgCyan = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		ColorCardBgGreen = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		ColorCardBgYellow = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		ColorCardBgOrange = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		ColorCardBgPurple = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		ColorCardBgCoral = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
		ColorCardBgBlue = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	case ThemeNeumorphismDark:
		IsNeumorphism = true
		IsDarkTheme = true
		CurrentCornerRadius = CornerRadiusNeumorph
		CurrentBorderWidth = BorderWidthNeumorph
		CurrentBadgeRadius = CornerRadiusPill

		// Neumorphism Dark Monochromatic palette (#21242B)
		ColorBgBase = color.RGBA{R: 0x1A, G: 0x1C, B: 0x22, A: 0xFF}      // Deep dark base
		ColorBgSidebar = color.RGBA{R: 0x16, G: 0x18, B: 0x1D, A: 0xFF}   // Inset sidebar
		ColorBgCard = color.RGBA{R: 0x21, G: 0x24, B: 0x2B, A: 0xFF}      // Monochromatic dark card
		ColorBgCardInner = color.RGBA{R: 0x18, G: 0x1A, B: 0x20, A: 0xFF} // Inner container
		ColorBgHover = color.RGBA{R: 0x28, G: 0x2C, B: 0x35, A: 0xFF}     // Hover tint

		ColorBorderSubtle = color.RGBA{R: 0x2D, G: 0x32, B: 0x3C, A: 0xFF} // 1.0px subtle border
		ColorBorderActive = color.RGBA{R: 0x81, G: 0x8C, B: 0xF8, A: 0xFF} // Soft Indigo focus
		ColorShadow = color.RGBA{R: 0x10, G: 0x12, B: 0x16, A: 0xFF}

		ColorNeumorphLightShadow = color.RGBA{R: 0x2E, G: 0x33, B: 0x3E, A: 0xDD} // Top-left glow
		ColorNeumorphDarkShadow = color.RGBA{R: 0x10, G: 0x12, B: 0x16, A: 0xF0}  // Bottom-right deep shadow

		ColorAccentYellow = color.RGBA{R: 0xFC, G: 0xD3, B: 0x4D, A: 0xFF}     // Pastel Amber
		ColorAccentCobalt = color.RGBA{R: 0x60, G: 0xA5, B: 0xFA, A: 0xFF}     // Pastel Blue
		ColorAccentCobaltDark = color.RGBA{R: 0x3B, G: 0x82, B: 0xF6, A: 0xFF}
		ColorAccentCobaltDim = color.RGBA{R: 0x3B, G: 0x82, B: 0xF6, A: 0x33}

		ColorAccentCyan = color.RGBA{R: 0x38, G: 0xBD, B: 0xF8, A: 0xFF}    // Vivid Cyan
		ColorAccentCyanDim = color.RGBA{R: 0x02, G: 0x84, B: 0xC7, A: 0x33}

		ColorSuccess = color.RGBA{R: 0x34, G: 0xD3, B: 0x99, A: 0xFF}          // Mint Emerald
		ColorWarning = color.RGBA{R: 0xFB, G: 0xBF, B: 0x24, A: 0xFF}          // Soft Amber
		ColorWarningTangerine = color.RGBA{R: 0xFB, G: 0x92, B: 0x3C, A: 0xFF} // Orange
		ColorDanger = color.RGBA{R: 0xF8, G: 0x71, B: 0xF8, A: 0xFF}           // Soft Coral
		ColorTechIndigo = color.RGBA{R: 0xA7, G: 0x8B, B: 0xFA, A: 0xFF}       // Soft Violet
		ColorInfo = color.RGBA{R: 0x38, G: 0xBD, B: 0xF8, A: 0xFF}             // Sky Blue

		ColorTextPrimary = color.RGBA{R: 0xF1, G: 0xF5, B: 0xF9, A: 0xFF}   // Crisp Soft White
		ColorTextSecondary = color.RGBA{R: 0xCB, G: 0xD5, B: 0xE1, A: 0xFF} // Zinc Light
		ColorTextMuted = color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF}     // Slate Muted
		ColorTextDisabled = color.RGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF}

		ColorCardBgCyan = color.RGBA{R: 0x21, G: 0x24, B: 0x2B, A: 0xFF}
		ColorCardBgGreen = color.RGBA{R: 0x21, G: 0x24, B: 0x2B, A: 0xFF}
		ColorCardBgYellow = color.RGBA{R: 0x21, G: 0x24, B: 0x2B, A: 0xFF}
		ColorCardBgOrange = color.RGBA{R: 0x21, G: 0x24, B: 0x2B, A: 0xFF}
		ColorCardBgPurple = color.RGBA{R: 0x21, G: 0x24, B: 0x2B, A: 0xFF}
		ColorCardBgCoral = color.RGBA{R: 0x21, G: 0x24, B: 0x2B, A: 0xFF}
		ColorCardBgBlue = color.RGBA{R: 0x21, G: 0x24, B: 0x2B, A: 0xFF}

	default: // ThemeNeoBrutalism (Signature Single Mode)
		ActiveTheme = ThemeNeoBrutalism
		IsNeumorphism = false
		IsDarkTheme = false
		CurrentCornerRadius = CornerRadiusBrutal
		CurrentBorderWidth = BorderWidthHeavy
		CurrentBadgeRadius = CornerRadiusBrutal

		// Classic Neo-Brutalism (Warm Retro Paper, Solid 2.5px Jet Black Borders & Hard Shadows)
		ColorBgBase = color.RGBA{R: 0xFF, G: 0xFD, B: 0xF8, A: 0xFF}      // Warm Retro Paper Canvas
		ColorBgSidebar = color.RGBA{R: 0xF4, G: 0xEF, B: 0xE6, A: 0xFF}   // Tinted Brutalist Sidebar
		ColorBgCard = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}      // Solid White Card Surface
		ColorBgCardInner = color.RGBA{R: 0xFF, G: 0xFD, B: 0xF0, A: 0xFF} // Warm Cream Inner Container
		ColorBgHover = color.RGBA{R: 0xEE, G: 0xE8, B: 0xDD, A: 0xFF}     // Hover tint

		ColorBorderSubtle = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Solid Jet Black 2.5px border
		ColorBorderActive = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Solid Jet Black
		ColorShadow = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}       // Solid 100% Black Hard Shadow

		ColorAccentYellow = color.RGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF}     // Gumroad Electric Yellow
		ColorAccentCobalt = color.RGBA{R: 0x25, G: 0x63, B: 0xEB, A: 0xFF}     // Electric Cobalt
		ColorAccentCobaltDark = color.RGBA{R: 0x1D, G: 0x4E, B: 0xD8, A: 0xFF} // Cobalt Dark
		ColorAccentCobaltDim = color.RGBA{R: 0xDB, G: 0xEA, B: 0xFE, A: 0xFF}  // Blue Tint

		ColorAccentCyan = color.RGBA{R: 0x00, G: 0xE5, B: 0xFF, A: 0xFF}    // Vivid Cyan
		ColorAccentCyanDim = color.RGBA{R: 0xBA, G: 0xE6, B: 0xFD, A: 0xFF} // Cyan Tint

		ColorSuccess = color.RGBA{R: 0x00, G: 0xF0, B: 0x76, A: 0xFF}          // Vivid Electric Neo Green
		ColorWarning = color.RGBA{R: 0xFF, G: 0x90, B: 0x00, A: 0xFF}          // Vivid Electric Tangerine
		ColorWarningTangerine = color.RGBA{R: 0xFF, G: 0x77, B: 0x00, A: 0xFF} // Punchy Orange
		ColorDanger = color.RGBA{R: 0xFF, G: 0x38, B: 0x5C, A: 0xFF}           // Vivid Electric Coral / Red
		ColorTechIndigo = color.RGBA{R: 0xA8, G: 0x55, B: 0xF7, A: 0xFF}       // Electric Violet
		ColorInfo = color.RGBA{R: 0x00, G: 0xE0, B: 0xFF, A: 0xFF}             // Electric Sky Cyan

		// All fonts in Neo-Brutalism: 100% PITCH BLACK
		ColorTextPrimary = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
		ColorTextSecondary = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
		ColorTextMuted = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
		ColorTextDisabled = color.RGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xFF}

		ColorCardBgCyan = color.RGBA{R: 0x00, G: 0xE0, B: 0xFF, A: 0xFF}
		ColorCardBgGreen = color.RGBA{R: 0x22, G: 0xE5, B: 0x65, A: 0xFF}
		ColorCardBgYellow = color.RGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF}
		ColorCardBgOrange = color.RGBA{R: 0xFF, G: 0x90, B: 0x00, A: 0xFF}
		ColorCardBgPurple = color.RGBA{R: 0xB8, G: 0x7C, B: 0xF8, A: 0xFF}
		ColorCardBgCoral = color.RGBA{R: 0xFF, G: 0x38, B: 0x5C, A: 0xFF}
		ColorCardBgBlue = color.RGBA{R: 0x38, G: 0x82, B: 0xF6, A: 0xFF}
	}
}

// SetThemeBool provides backward compatibility for boolean toggling
func SetThemeBool(isDark bool) {
	if isDark {
		SetTheme(ThemeNeumorphismDark)
	} else {
		SetTheme(ThemeNeoBrutalism)
	}
}

// === TYPOGRAPHY SCALE ===
const (
	FontSizeDisplay = float32(22)
	FontSizeH1      = float32(18)
	FontSizeH2      = float32(15)
	FontSizeH3      = float32(13)
	FontSizeBody    = float32(12)
	FontSizeSmall   = float32(11)
	FontSizeLabel   = float32(9.5)
)

// === SPACING SYSTEM ===
const (
	SpaceXS  = float32(4)
	SpaceSM  = float32(8)
	SpaceMD  = float32(12)
	SpaceLG  = float32(16)
	SpaceXL  = float32(24)
	SpaceXXL = float32(32)
)
