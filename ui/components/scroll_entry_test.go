package components_test

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/yudz/it-toolbox/ui/components"
)

func TestScrollableEntryTransparentScroll(t *testing.T) {
	w := test.NewWindow(nil)

	content := canvas.NewRectangle(theme.PrimaryColor())
	content.SetMinSize(fyne.NewSize(300, 1000))

	entry := components.NewScrollableEntry()
	entry.SetText("Search input test")

	scroll := container.NewVScroll(container.NewVBox(entry, content))
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(300, 300))

	initOffset := scroll.Offset.Y
	test.Scroll(w.Canvas(), fyne.NewPos(50, 15), 0, -40)

	if scroll.Offset.Y <= initOffset {
		t.Fatalf("Expected scroll offset to increase from %f, remained %f", initOffset, scroll.Offset.Y)
	}
}

func TestSearchBarTransparentScroll(t *testing.T) {
	w := test.NewWindow(nil)

	content := canvas.NewRectangle(theme.PrimaryColor())
	content.SetMinSize(fyne.NewSize(300, 1000))

	searchBar := components.NewSearchBar("Cari...", nil)

	scroll := container.NewVScroll(container.NewVBox(searchBar.Container, content))
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(300, 300))

	initOffset := scroll.Offset.Y
	test.Scroll(w.Canvas(), fyne.NewPos(50, 15), 0, -40)

	if scroll.Offset.Y <= initOffset {
		t.Fatalf("Expected search bar scroll offset to increase from %f, remained %f", initOffset, scroll.Offset.Y)
	}
}

func TestScrollableMultiLineEntryTransparentScroll(t *testing.T) {
	w := test.NewWindow(nil)

	content := canvas.NewRectangle(theme.PrimaryColor())
	content.SetMinSize(fyne.NewSize(300, 1000))

	scroll := container.NewVScroll(nil)
	entry := components.NewScrollableMultiLineEntry(scroll)
	entry.SetText("Line 1\nLine 2\nLine 3")

	scroll.Content = container.NewVBox(entry, content)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(300, 300))

	initOffset := scroll.Offset.Y
	test.Scroll(w.Canvas(), fyne.NewPos(50, 15), 0, -40)

	if scroll.Offset.Y <= initOffset {
		t.Fatalf("Expected multiline entry scroll offset to increase from %f, remained %f", initOffset, scroll.Offset.Y)
	}
}
