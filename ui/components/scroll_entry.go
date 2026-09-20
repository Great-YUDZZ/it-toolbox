package components

import (
	"reflect"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewScrollableEntry creates a single-line entry that transparently forwards mouse wheel scrolling
// to its parent container without getting stuck or blocking page scrolling.
func NewScrollableEntry() *widget.Entry {
	e := widget.NewEntry()
	e.Wrapping = fyne.TextWrapOff
	e.Scroll = fyne.ScrollNone
	return e
}

// scrollShield wraps a CanvasObject to hide the internal Scrollable implementation from Fyne's
// driver hit-testing, allowing ScrollableMultiLineEntry to receive the Scrolled event directly.
type scrollShield struct {
	fyne.CanvasObject
}

type scrollableEntryRenderer struct {
	base    fyne.WidgetRenderer
	objects []fyne.CanvasObject
}

func (r *scrollableEntryRenderer) Destroy()                     { r.base.Destroy() }
func (r *scrollableEntryRenderer) Layout(s fyne.Size)           { r.base.Layout(s) }
func (r *scrollableEntryRenderer) MinSize() fyne.Size           { return r.base.MinSize() }
func (r *scrollableEntryRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *scrollableEntryRenderer) Refresh()                     { r.base.Refresh() }

// ScrollableMultiLineEntry wraps widget.Entry to implement fyne.Scrollable.
// In standard Fyne, mouse-scrolling while hovering over an Entry is consumed and not passed to
// parent scroll containers. ScrollableMultiLineEntry allows smooth scrolling:
// 1. If text fits within visible height, scrolling immediately scrolls the parent container.
// 2. If text exceeds visible rows, it scrolls internally until reaching the boundary,
//    then seamlessly propagates the scroll event to the parent container.
type ScrollableMultiLineEntry struct {
	widget.Entry
	parentScroller *container.Scroll
}

// NewScrollableMultiLineEntry creates a multiline entry that forwards scroll events to parentScroller.
func NewScrollableMultiLineEntry(parent *container.Scroll) *ScrollableMultiLineEntry {
	e := &ScrollableMultiLineEntry{parentScroller: parent}
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	e.ExtendBaseWidget(e)
	return e
}

// CreateRenderer intercepts renderer objects and shields the inner scroll from swallowing wheel events.
func (e *ScrollableMultiLineEntry) CreateRenderer() fyne.WidgetRenderer {
	base := e.Entry.CreateRenderer()
	rawObjs := base.Objects()
	wrapped := make([]fyne.CanvasObject, len(rawObjs))
	for i, o := range rawObjs {
		if _, ok := o.(fyne.Scrollable); ok {
			wrapped[i] = &scrollShield{CanvasObject: o}
		} else {
			wrapped[i] = o
		}
	}
	return &scrollableEntryRenderer{base: base, objects: wrapped}
}

// SetParentScroller sets or updates the parent scroll container for scroll propagation.
func (e *ScrollableMultiLineEntry) SetParentScroller(parent *container.Scroll) {
	e.parentScroller = parent
}

func (e *ScrollableMultiLineEntry) getInternalScroll() *container.Scroll {
	val := reflect.ValueOf(&e.Entry).Elem()
	field := val.FieldByName("scroll")
	if !field.IsValid() {
		return nil
	}
	ptr := unsafe.Pointer(field.UnsafeAddr())
	return *(**container.Scroll)(ptr)
}

// Scrolled implements fyne.Scrollable so mouse wheel events over this entry don't get stuck.
func (e *ScrollableMultiLineEntry) Scrolled(ev *fyne.ScrollEvent) {
	internal := e.getInternalScroll()
	if internal == nil || internal.Content == nil {
		if e.parentScroller != nil {
			e.parentScroller.Scrolled(ev)
		}
		return
	}

	contentH := internal.Content.MinSize().Height
	viewH := internal.Size().Height

	// If entry content fits completely inside the entry height without needing internal scroll,
	// pass scroll directly to parent container!
	if contentH <= viewH {
		if e.parentScroller != nil {
			e.parentScroller.Scrolled(ev)
		}
		return
	}

	oldY := internal.Offset.Y
	// Scroll internally
	internal.Scrolled(ev)

	// If internal scroll reached boundary (top when scrolling up, or bottom when scrolling down),
	// propagate the scroll to the parent container so user is never trapped!
	if internal.Offset.Y == oldY && e.parentScroller != nil {
		e.parentScroller.Scrolled(ev)
	}
}
