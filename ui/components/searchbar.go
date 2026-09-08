package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SearchBar is a reusable search input with a clear button and live change callback
type SearchBar struct {
	Entry     *widget.Entry
	Container *fyne.Container
	OnSearch  func(query string)
}

func NewSearchBar(placeholder string, onSearch func(query string)) *SearchBar {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)

	sb := &SearchBar{
		Entry:    entry,
		OnSearch: onSearch,
	}

	clearBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		entry.SetText("")
		if sb.OnSearch != nil {
			sb.OnSearch("")
		}
	})

	entry.OnChanged = func(s string) {
		if sb.OnSearch != nil {
			sb.OnSearch(s)
		}
	}

	sb.Container = container.NewBorder(nil, nil, widget.NewIcon(theme.SearchIcon()), clearBtn, entry)
	return sb
}

func (s *SearchBar) GetText() string {
	return s.Entry.Text
}

func (s *SearchBar) SetText(text string) {
	s.Entry.SetText(text)
}
