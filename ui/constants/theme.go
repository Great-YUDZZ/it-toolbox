package constants

import "image/color"

// ============================================================================
// DESIGN TOKEN SYSTEM — Neo-Brutalism Design System (Light & Dark)
// ============================================================================

var IsDarkTheme = false

// === NEO-BRUTALISM GEOMETRY CONSTANTS ===
const (
	BorderWidthHeavy   = float32(2.5) // Bold prominent strokes (cards, panels)
	BorderWidthMedium  = float32(2.0) // Buttons, badges, inputs
	BorderWidthThin    = float32(1.5) // Separators, sub-elements
	ShadowOffsetHeavy  = float32(4.0) // Hard offset drop shadow for cards (+4px, +4px)
	ShadowOffsetMedium = float32(3.0) // Hard offset drop shadow for badges & KPIs
	CornerRadiusBrutal = float32(4.0) // Semi-sharp blocky corner radius
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

	ColorAccentYellow = color.RGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF} // Iconic Electric Acid / Gumroad Yellow
	ColorAccentCobalt = color.RGBA{R: 0x25, G: 0x63, B: 0xEB, A: 0xFF} // Electric Cobalt Blue
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

	// === NEO-BRUTALISM VIBRANT CARD BACKGROUNDS ===
	ColorCardBgCyan   = color.RGBA{R: 0x00, G: 0xE0, B: 0xFF, A: 0xFF} // Vivid Electric Cyan (#00E0FF)
	ColorCardBgGreen  = color.RGBA{R: 0x22, G: 0xE5, B: 0x65, A: 0xFF} // Vivid Electric Kelly Green (#22E565)
	ColorCardBgYellow = color.RGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF} // Vivid Electric Acid Yellow (#FFE600)
	ColorCardBgOrange = color.RGBA{R: 0xFF, G: 0x90, B: 0x00, A: 0xFF} // Vivid Electric Tangerine (#FF9000)
	ColorCardBgPurple = color.RGBA{R: 0xB8, G: 0x7C, B: 0xF8, A: 0xFF} // Vivid Electric Purple (#B87CF8)
	ColorCardBgCoral  = color.RGBA{R: 0xFF, G: 0x38, B: 0x5C, A: 0xFF} // Vivid Electric Coral Red (#FF385C)
	ColorCardBgBlue   = color.RGBA{R: 0x38, G: 0x82, B: 0xF6, A: 0xFF} // Vivid Electric Royal Blue (#3882F6)
)

// SetTheme switches all token values between Light and Dark mode
func SetTheme(isDark bool) {
	IsDarkTheme = isDark
	if isDark {
		// Cyber-Brutalism Dark Mode (Deep Obsidian Canvas, Stark White Outlines, Saturated Neon)
		ColorBgBase = color.RGBA{R: 0x12, G: 0x12, B: 0x14, A: 0xFF}      // Deep Obsidian Canvas
		ColorBgSidebar = color.RGBA{R: 0x18, G: 0x18, B: 0x1B, A: 0xFF}   // Zinc 900 Sidebar
		ColorBgCard = color.RGBA{R: 0x1E, G: 0x1E, B: 0x24, A: 0xFF}      // Dark brutalist card face
		ColorBgCardInner = color.RGBA{R: 0x27, G: 0x27, B: 0x30, A: 0xFF} // Inner panel
		ColorBgHover = color.RGBA{R: 0x2E, G: 0x2E, B: 0x3A, A: 0xFF}     // Hover state

		ColorBorderSubtle = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // Stark White 2.5px border
		ColorBorderActive = color.RGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF} // Electric Yellow Active Border
		ColorShadow = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}       // Solid Pitch Black Hard Shadow

		ColorAccentYellow = color.RGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF}     // Electric Acid Yellow
		ColorAccentCobalt = color.RGBA{R: 0x60, G: 0xA5, B: 0xFA, A: 0xFF}     // Bright Blue
		ColorAccentCobaltDark = color.RGBA{R: 0x3B, G: 0x82, B: 0xF6, A: 0xFF} // Blue
		ColorAccentCobaltDim = color.RGBA{R: 0x3B, G: 0x82, B: 0xF6, A: 0x40}  // Blue tint

		ColorAccentCyan = color.RGBA{R: 0x00, G: 0xE5, B: 0xFF, A: 0xFF}    // Vivid Cyan
		ColorAccentCyanDim = color.RGBA{R: 0x00, G: 0xE5, B: 0xFF, A: 0x30} // Cyan tint

		ColorSuccess = color.RGBA{R: 0x00, G: 0xF0, B: 0xA0, A: 0xFF}          // Vivid Neo Mint
		ColorWarning = color.RGBA{R: 0xFF, G: 0xB8, B: 0x00, A: 0xFF}          // Vivid Amber
		ColorWarningTangerine = color.RGBA{R: 0xFF, G: 0x8C, B: 0x32, A: 0xFF} // Tangerine
		ColorDanger = color.RGBA{R: 0xFF, G: 0x53, B: 0x53, A: 0xFF}           // Neon Red
		ColorTechIndigo = color.RGBA{R: 0xC0, G: 0x84, B: 0xFC, A: 0xFF}       // Neon Violet
		ColorInfo = color.RGBA{R: 0x60, G: 0xA5, B: 0xFA, A: 0xFF}             // Sky Blue

		ColorTextPrimary = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}   // Pure White
		ColorTextSecondary = color.RGBA{R: 0xE4, G: 0xE4, B: 0xE7, A: 0xFF} // Zinc 200
		ColorTextMuted = color.RGBA{R: 0xA1, G: 0xA1, B: 0xAA, A: 0xFF}     // Zinc 400
		ColorTextDisabled = color.RGBA{R: 0x71, G: 0x71, B: 0x7A, A: 0xFF}  // Zinc 500

		ColorCardBgCyan = color.RGBA{R: 0x0E, G: 0x2A, B: 0x38, A: 0xFF}
		ColorCardBgGreen = color.RGBA{R: 0x0C, G: 0x2E, B: 0x22, A: 0xFF}
		ColorCardBgYellow = color.RGBA{R: 0x2E, G: 0x26, B: 0x0C, A: 0xFF}
		ColorCardBgOrange = color.RGBA{R: 0x33, G: 0x1D, B: 0x0C, A: 0xFF}
		ColorCardBgPurple = color.RGBA{R: 0x24, G: 0x18, B: 0x36, A: 0xFF}
		ColorCardBgCoral = color.RGBA{R: 0x33, G: 0x12, B: 0x16, A: 0xFF}
		ColorCardBgBlue = color.RGBA{R: 0x1A, G: 0x2B, B: 0x4C, A: 0xFF}
	} else {
		// Classic Neo-Brutalism Light Mode (Warm Retro Paper, Solid 2.5px Jet Black Borders & Hard Shadows)
		ColorBgBase = color.RGBA{R: 0xFF, G: 0xFD, B: 0xF8, A: 0xFF}      // Warm Retro Paper Canvas
		ColorBgSidebar = color.RGBA{R: 0xF4, G: 0xEF, B: 0xE6, A: 0xFF}   // Tinted Brutalist Sidebar
		ColorBgCard = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}      // Solid White Card Surface
		ColorBgCardInner = color.RGBA{R: 0xFF, G: 0xFD, B: 0xF0, A: 0xFF} // Warm Cream Inner Container
		ColorBgHover = color.RGBA{R: 0xEE, G: 0xE8, B: 0xDD, A: 0xFF}     // Hover tint

		ColorBorderSubtle = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Solid Jet Black
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

		// All fonts in light mode: 100% PITCH BLACK
		ColorTextPrimary = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}   // Pitch Black
		ColorTextSecondary = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF} // Pitch Black
		ColorTextMuted = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}     // Pitch Black
		ColorTextDisabled = color.RGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xFF}  // Dark Slate

		ColorCardBgCyan = color.RGBA{R: 0x00, G: 0xE0, B: 0xFF, A: 0xFF}
		ColorCardBgGreen = color.RGBA{R: 0x22, G: 0xE5, B: 0x65, A: 0xFF}
		ColorCardBgYellow = color.RGBA{R: 0xFF, G: 0xE6, B: 0x00, A: 0xFF}
		ColorCardBgOrange = color.RGBA{R: 0xFF, G: 0x90, B: 0x00, A: 0xFF}
		ColorCardBgPurple = color.RGBA{R: 0xB8, G: 0x7C, B: 0xF8, A: 0xFF}
		ColorCardBgCoral = color.RGBA{R: 0xFF, G: 0x38, B: 0x5C, A: 0xFF}
		ColorCardBgBlue = color.RGBA{R: 0x38, G: 0x82, B: 0xF6, A: 0xFF}
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
