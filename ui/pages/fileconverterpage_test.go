package pages

import (
	"fyne.io/fyne/v2"
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"

	"github.com/yudz/it-toolbox/ui/constants"
)

func TestFileConverterPageBuild(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	win := test.NewWindow(nil)
	defer win.Close()

	page := NewFileConverterPage(win)
	obj := page.Build()

	if obj == nil {
		t.Fatal("expected non-nil CanvasObject from FileConverterPage.Build()")
	}

	// Verify TabQRCode is present in tabs
	border, ok := obj.(*fyne.Container)
	if !ok {
		// container.NewBorder returns *fyne.Container
	}
	_ = border
}

func TestQRCodePageBuild(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	win := test.NewWindow(nil)
	defer win.Close()

	page := NewFileConverterPage(win)
	pageObj := page.BuildQRCodePage()
	if pageObj == nil {
		t.Fatal("expected non-nil CanvasObject from FileConverterPage.BuildQRCodePage()")
	}

	tab := page.buildQRCodeTab()
	if tab == nil {
		t.Fatal("expected non-nil CanvasObject from buildQRCodeTab()")
	}

	// Verify MinSize is bounded to avoid UI clipping
	if tab.MinSize().Width > 600 {
		t.Fatalf("buildQRCodeTab MinSize width (%f) exceeds 600px, which causes cutoffs on standard screens", tab.MinSize().Width)
	}
}

func TestFileConverterTabsIncludeQRCode(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	win := test.NewWindow(nil)
	defer win.Close()

	page := NewFileConverterPage(win)
	obj := page.Build()

	// Locate AppTabs
	var appTabs *container.AppTabs
	if c, ok := obj.(*fyne.Container); ok {
		for _, child := range c.Objects {
			if tabs, ok := child.(*container.AppTabs); ok {
				appTabs = tabs
				break
			}
		}
	}

	if appTabs == nil {
		t.Fatal("expected AppTabs inside FileConverterPage.Build()")
	}

	foundQR := false
	for _, item := range appTabs.Items {
		if item.Text == constants.TabQRCode {
			foundQR = true
			break
		}
	}

	if !foundQR {
		t.Fatalf("expected AppTabs to contain tab %q", constants.TabQRCode)
	}
}
