package components_test

import (
	"testing"

		"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/yudz/it-toolbox/ui"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

func TestThemeSwitching(t *testing.T) {
	app := test.NewApp()

	// 1. Test Neo-Brutalism (Signature)
	ui.ApplyTheme(constants.ThemeNeoBrutalism)
	app.Settings().SetTheme(ui.NewCustomCyanTheme())

	if constants.IsNeumorphism {
		t.Errorf("Expected IsNeumorphism=false in NeoBrutalism")
	}
	if constants.IsDarkTheme {
		t.Errorf("Expected IsDarkTheme=false in NeoBrutalism")
	}
	if constants.CurrentCornerRadius != constants.CornerRadiusBrutal {
		t.Errorf("Expected CornerRadius=4 in NeoBrutalism, got %f", constants.CurrentCornerRadius)
	}
	if constants.CurrentBorderWidth != constants.BorderWidthHeavy {
		t.Errorf("Expected BorderWidth=2.5 in NeoBrutalism, got %f", constants.CurrentBorderWidth)
	}

	cardBrutal := components.NewPlainCard(widget.NewLabel("Test Card"))
	minSizeBrutal := cardBrutal.MinSize()
	t.Logf("Neo-Brutalism PlainCard MinSize: %+v", minSizeBrutal)

	// 2. Test Neumorphism Mode Terang (Light Soft UI)
	ui.ApplyTheme(constants.ThemeNeumorphismLight)
	app.Settings().SetTheme(ui.NewCustomCyanTheme())

	if !constants.IsNeumorphism {
		t.Errorf("Expected IsNeumorphism=true in NeumorphismLight")
	}
	if constants.IsDarkTheme {
		t.Errorf("Expected IsDarkTheme=false in NeumorphismLight")
	}
	if constants.CurrentCornerRadius != constants.CornerRadiusNeumorph {
		t.Errorf("Expected CornerRadius=14 in NeumorphismLight, got %f", constants.CurrentCornerRadius)
	}
	if constants.CurrentBorderWidth != constants.BorderWidthNeumorph {
		t.Errorf("Expected BorderWidth=0.8 in NeumorphismLight, got %f", constants.CurrentBorderWidth)
	}

	cardNeumorphLight := components.NewPlainCard(widget.NewLabel("Test Card"))
	minSizeLight := cardNeumorphLight.MinSize()
	t.Logf("Neumorphism Light PlainCard MinSize: %+v", minSizeLight)

	// 3. Test Neumorphism Mode Gelap (Dark Soft UI)
	ui.ApplyTheme(constants.ThemeNeumorphismDark)
	app.Settings().SetTheme(ui.NewCustomCyanTheme())

	if !constants.IsNeumorphism {
		t.Errorf("Expected IsNeumorphism=true in NeumorphismDark")
	}
	if !constants.IsDarkTheme {
		t.Errorf("Expected IsDarkTheme=true in NeumorphismDark")
	}
	if constants.CurrentCornerRadius != constants.CornerRadiusNeumorph {
		t.Errorf("Expected CornerRadius=14 in NeumorphismDark, got %f", constants.CurrentCornerRadius)
	}

	cardNeumorphDark := components.NewPlainCard(widget.NewLabel("Test Card"))
	minSizeDark := cardNeumorphDark.MinSize()
	t.Logf("Neumorphism Dark PlainCard MinSize: %+v", minSizeDark)

	// Verify MinSizes are sane and do not balloon
	if minSizeLight.Width > 500 || minSizeDark.Width > 500 {
		t.Errorf("Card min sizes should remain compact, got light: %+v, dark: %+v", minSizeLight, minSizeDark)
	}

	// 4. Test Badges under all themes
	bYellow := components.BadgeYellow("TAG")
	if bYellow.MinSize().Width == 0 {
		t.Errorf("Badge minSize width should be > 0")
	}

	// Restore Neo-Brutalism as default
	ui.ApplyTheme(constants.ThemeNeoBrutalism)
	app.Settings().SetTheme(ui.NewCustomCyanTheme())
}
