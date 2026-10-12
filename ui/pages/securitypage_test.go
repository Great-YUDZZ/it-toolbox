package pages

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"

	"github.com/yudz/it-toolbox/ui/constants"
)

func TestSecurityPageBuild(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	win := test.NewWindow(nil)
	defer win.Close()

	page := NewSecurityPage(win)
	obj := page.Build()

	if obj == nil {
		t.Fatal("expected non-nil CanvasObject from SecurityPage.Build()")
	}

	// Locate AppTabs inside the returned container
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
		t.Fatal("expected AppTabs inside SecurityPage.Build()")
	}

	expectedTabs := []string{
		constants.TabSecurityScanner,
		constants.TabSecurityWebSSL,
		constants.TabSecurityForensics,
		constants.TabSecurityEncoders,
		constants.TabSecurityPlaybook,
	}

	for _, exp := range expectedTabs {
		found := false
		for _, item := range appTabs.Items {
			if item.Text == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected AppTabs to contain tab %q", exp)
		}
	}
}

func TestSecurityPageTabsMinSizes(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	win := test.NewWindow(nil)
	defer win.Close()

	page := NewSecurityPage(win)

	tabs := []struct {
		name string
		obj  fyne.CanvasObject
	}{
		{"PortScannerTab", page.buildPortScannerTab()},
		{"WebSSLTab", page.buildWebSSLTab()},
		{"ForensicsTab", page.buildForensicsTab()},
		{"EncodersASCIITab", page.buildEncodersASCIITab()},
		{"PlaybookTab", page.buildPlaybookTab()},
	}

	for _, tt := range tabs {
		if tt.obj == nil {
			t.Fatalf("tab %s produced nil CanvasObject", tt.name)
		}
		w := tt.obj.MinSize().Width
		if w > 600 {
			t.Errorf("tab %s MinSize width (%f) exceeds 600px threshold, risk of UI overflow", tt.name, w)
		}
	}
}
