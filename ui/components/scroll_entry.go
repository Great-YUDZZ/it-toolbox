package components

import (
	"reflect"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ScrollableEntry wraps widget.Entry so mouse wheel scrolling flows freely to parent scrollers
// while maintaining a compact, bounded MinSize and full internal horizontal scrolling.
type ScrollableEntry struct {
	widget.Entry
}

func NewScrollableEntry() *widget.Entry {
	e := &ScrollableEntry{}
	e.ExtendBaseWidget(e)
	return &e.Entry
}

func (e *ScrollableEntry) CreateRenderer() fyne.WidgetRenderer {
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

// scrollShield wraps a CanvasObject to hide the internal Scrollable implementation from Fyne's
// driver hit-testing, allowing ScrollableMultiLineEntry and ScrollableEntry to receive or pass the Scrolled event.
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
type ScrollableMultiLineEntry struct {
	widget.Entry
	parentScroller *container.Scroll
}

func NewScrollableMultiLineEntry(parent *container.Scroll) *ScrollableMultiLineEntry {
	e := &ScrollableMultiLineEntry{parentScroller: parent}
	e.MultiLine = true
	e.Wrapping = fyne.TextWrapWord
	e.ExtendBaseWidget(e)
	return e
}

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

	if contentH <= viewH {
		if e.parentScroller != nil {
			e.parentScroller.Scrolled(ev)
		}
		return
	}

	oldY := internal.Offset.Y
	internal.Scrolled(ev)

	if internal.Offset.Y == oldY && e.parentScroller != nil {
		e.parentScroller.Scrolled(ev)
	}
}
