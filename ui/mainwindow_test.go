package ui_test

import (
	"image/color"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/yudz/it-toolbox/ui"
	"github.com/yudz/it-toolbox/ui/constants"
)

type testSlidingSidebarLayout struct {
	width float32
}

func (l *testSlidingSidebarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}
	totalWidth := float32(constants.SidebarWidth)
	clip := objects[0].(*container.Clip)
	content := clip.Content

	clip.Resize(fyne.NewSize(l.width, size.Height))
	clip.Move(fyne.NewPos(0, 0))

	offsetX := l.width - totalWidth
	content.Move(fyne.NewPos(offsetX, 0))
	content.Resize(fyne.NewSize(totalWidth, size.Height))
}

func (l *testSlidingSidebarLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(l.width, 0)
}

func TestSlidingSidebarLayout(t *testing.T) {
	w := test.NewWindow(nil)
	content := canvas.NewRectangle(color.Black)
	totalW := float32(constants.SidebarWidth)
	content.SetMinSize(fyne.NewSize(totalW, 400))

	clip := container.NewClip(content)
	layout := &testSlidingSidebarLayout{width: totalW}
	slidingContainer := container.New(layout, clip)

	w.SetContent(slidingContainer)
	w.Resize(fyne.NewSize(800, 600))

	// When width is totalW
	slidingContainer.Resize(fyne.NewSize(totalW, 600))
	layout.Layout(slidingContainer.Objects, fyne.NewSize(totalW, 600))
	if content.Position().X != 0 {
		t.Errorf("Expected content X to be 0 when expanded, got %f", content.Position().X)
	}

	// When width is 100
	layout.width = 100
	slidingContainer.Resize(fyne.NewSize(100, 600))
	layout.Layout(slidingContainer.Objects, fyne.NewSize(100, 600))
	if content.Position().X != (100 - totalW) {
		t.Errorf("Expected content X to be %f when width=100, got %f", 100-totalW, content.Position().X)
	}

	// When width is 0
	layout.width = 0
	slidingContainer.Resize(fyne.NewSize(0, 600))
	layout.Layout(slidingContainer.Objects, fyne.NewSize(0, 600))
	if content.Position().X != -totalW {
		t.Errorf("Expected content X to be %f when width=0, got %f", -totalW, content.Position().X)
	}
}

func TestMainWindowToggleSidebarAnimated(t *testing.T) {
	app := test.NewApp()
	mw := ui.NewMainWindow(app)

	// Initially expanded
	if mw.RootContainer == nil {
		t.Fatal("RootContainer should not be nil")
	}

	// Toggle collapse
	mw.ToggleSidebar()
	time.Sleep(220 * time.Millisecond)

	// Toggle expand
	mw.ToggleSidebar()
	time.Sleep(220 * time.Millisecond)
}

func TestMainWindowRapidToggle(t *testing.T) {
	app := test.NewApp()
	mw := ui.NewMainWindow(app)

	// Rapidly toggle multiple times in flight
	for i := 0; i < 5; i++ {
		mw.ToggleSidebar()
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(250 * time.Millisecond)
}

func TestMainWindowShowPageTransitions(t *testing.T) {
	app := test.NewApp()
	mw := ui.NewMainWindow(app)

	pagesToTest := []string{
		constants.NavToolbox,
		constants.NavFileConverter,
		constants.NavYouTube,
		constants.NavCisco,
		constants.NavDatabase,
		constants.NavReference,
		constants.NavLogbook,
		constants.NavTracker,
		constants.NavSettings,
	}

	for _, p := range pagesToTest {
		mw.ShowPage(p)
		if mw.ActiveMenu != p {
			t.Errorf("Expected ActiveMenu to be %s, got %s", p, mw.ActiveMenu)
		}
		// Allow animation tick
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)
}
