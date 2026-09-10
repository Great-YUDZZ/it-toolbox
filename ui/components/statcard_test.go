package components_test

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/yudz/it-toolbox/ui"
	"github.com/yudz/it-toolbox/ui/components"
)

func TestStatCardBoundedLayout(t *testing.T) {
	app := test.NewApp()
	app.Settings().SetTheme(ui.NewCustomCyanTheme())

	card := components.NewStatCard("RENTANG HOST USABLE", "192.168.1.1 ➔ 192.168.1.62", color.White)
	card.SetSubtext("Tersedia 62 IP Usable")

	minSize := card.Widget.MinSize()
	t.Logf("StatCard MinSize with long IP range: %+v", minSize)

	// Verify that MinSize.Width is bounded (does not explode to 400+ px)
	if minSize.Width > 200 {
		t.Errorf("Expected StatCard MinSize.Width <= 200, got %f", minSize.Width)
	}

	// Update with even longer text
	card.SetValue("172.16.0.1 ➔ 172.16.15.254 (Subnet Panjang)")
	minSizeAfter := card.Widget.MinSize()
	t.Logf("StatCard MinSize after long text update: %+v", minSizeAfter)

	if minSizeAfter.Width > 200 {
		t.Errorf("Expected StatCard MinSize.Width to remain <= 200 after update, got %f", minSizeAfter.Width)
	}
}
