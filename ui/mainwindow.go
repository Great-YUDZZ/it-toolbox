package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
	"github.com/yudz/it-toolbox/ui/pages"
)

type MainWindow struct {
	App         fyne.App
	Window      fyne.Window
	ContentArea *fyne.Container
	StatusLabel *canvas.Text
	ActiveMenu  string

	btnToolbox  *widget.Button
	btnFileConv *widget.Button
	btnRef      *widget.Button
	btnLogbook  *widget.Button
	btnTracker  *widget.Button

	calcPage     *pages.CalculatorPage
	fileConvPage *pages.FileConverterPage
	refPage      *pages.ReferencePage
	logbookPage  *pages.LogbookPage
	trackerPage  *pages.TrackerPage
}

func NewMainWindow(app fyne.App) *MainWindow {
	win := app.NewWindow(constants.AppTitle + " — " + constants.AppSubtitle)
	win.Resize(fyne.NewSize(constants.DefaultWinW, constants.DefaultWinH))
	win.CenterOnScreen()

	statusTxt := canvas.NewText("🟢 Sistem Siap • SQLite Terhubung", color.RGBA{R: 0x10, G: 0xB9, B: 0x81, A: 0xFF})
	statusTxt.TextSize = 11
	statusTxt.TextStyle = fyne.TextStyle{Bold: true}

	mw := &MainWindow{
		App:          app,
		Window:       win,
		ContentArea:  container.NewStack(),
		StatusLabel:  statusTxt,
		calcPage:     pages.NewCalculatorPage(win),
		fileConvPage: pages.NewFileConverterPage(win),
		refPage:      pages.NewReferencePage(win),
		logbookPage:  pages.NewLogbookPage(win),
		trackerPage:  pages.NewTrackerPage(win),
	}

	win.SetContent(mw.buildLayout())
	mw.showPage(constants.NavToolbox)
	return mw
}

func (m *MainWindow) buildLayout() fyne.CanvasObject {
	// ------------------------------------------------------------------------
	// 1. Sidebar Header (Branding & Identity)
	// ------------------------------------------------------------------------
	brandTitle := canvas.NewText("⚡ IT TOOLBOX", color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0xFF})
	brandTitle.TextSize = 17
	brandTitle.TextStyle = fyne.TextStyle{Bold: true}

	subTitle := canvas.NewText("ENGINEERING WORKBENCH", color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF})
	subTitle.TextSize = 9.5
	subTitle.TextStyle = fyne.TextStyle{Bold: true}

	verBadge := components.BadgeCyan("v" + constants.AppVersion)
	brandBox := container.NewVBox(
		container.NewHBox(brandTitle, verBadge),
		subTitle,
	)

	headerCard := container.NewVBox(
		container.NewPadded(brandBox),
		widget.NewSeparator(),
	)

	// ------------------------------------------------------------------------
	// 2. Navigation Section with Category Dividers
	// ------------------------------------------------------------------------
	secCore := canvas.NewText("ALAT & KALKULATOR", color.RGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF})
	secCore.TextSize = 9.5
	secCore.TextStyle = fyne.TextStyle{Bold: true}

	m.btnToolbox = widget.NewButtonWithIcon(constants.NavToolbox, theme.SettingsIcon(), func() {
		m.showPage(constants.NavToolbox)
	})
	m.btnFileConv = widget.NewButtonWithIcon(constants.NavFileConverter, theme.FolderOpenIcon(), func() {
		m.showPage(constants.NavFileConverter)
	})

	secDocs := canvas.NewText("PENGETAHUAN & LOG", color.RGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF})
	secDocs.TextSize = 9.5
	secDocs.TextStyle = fyne.TextStyle{Bold: true}

	m.btnRef = widget.NewButtonWithIcon(constants.NavReference, theme.HelpIcon(), func() {
		m.showPage(constants.NavReference)
	})
	m.btnLogbook = widget.NewButtonWithIcon(constants.NavLogbook, theme.DocumentIcon(), func() {
		m.showPage(constants.NavLogbook)
	})
	m.btnTracker = widget.NewButtonWithIcon(constants.NavTracker, theme.ConfirmIcon(), func() {
		m.showPage(constants.NavTracker)
	})

	navContainer := container.NewVBox(
		secCore,
		m.btnToolbox,
		m.btnFileConv,
		widget.NewSeparator(),
		secDocs,
		m.btnRef,
		m.btnLogbook,
		m.btnTracker,
	)

	// ------------------------------------------------------------------------
	// 3. Sidebar Footer (Status Pill & Environment Info)
	// ------------------------------------------------------------------------
	footerBox := container.NewVBox(
		widget.NewSeparator(),
		container.NewPadded(container.NewVBox(
			m.StatusLabel,
			canvas.NewText("Go 1.25 • Fyne v2.8 • Offline", color.RGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF}),
		)),
	)

	sidebarContent := container.NewBorder(headerCard, footerBox, nil, nil, container.NewPadded(navContainer))

	// Sidebar background with rich dark midnight slate
	bgSidebar := canvas.NewRectangle(color.RGBA{R: 0x0E, G: 0x13, B: 0x1D, A: 0xFF})
	sidebarWrapper := container.NewMax(bgSidebar, sidebarContent)
	sidebarWithSep := container.NewBorder(nil, nil, nil, widget.NewSeparator(), sidebarWrapper)

	// Main Layout: Sidebar on Left, Content Area in Center
	mainLayout := container.NewBorder(nil, nil, sidebarWithSep, nil, container.NewPadded(m.ContentArea))
	return mainLayout
}

func (m *MainWindow) updateNavHighlights(active string) {
	m.ActiveMenu = active

	resetBtn := func(btn *widget.Button, isActive bool) {
		if btn == nil {
			return
		}
		if isActive {
			btn.Importance = widget.HighImportance
		} else {
			btn.Importance = widget.MediumImportance
		}
		btn.Refresh()
	}

	resetBtn(m.btnToolbox, active == constants.NavToolbox)
	resetBtn(m.btnFileConv, active == constants.NavFileConverter)
	resetBtn(m.btnRef, active == constants.NavReference)
	resetBtn(m.btnLogbook, active == constants.NavLogbook)
	resetBtn(m.btnTracker, active == constants.NavTracker)

	m.StatusLabel.Text = fmt.Sprintf("🟢 %s", active)
	m.StatusLabel.Refresh()
}

func (m *MainWindow) showPage(name string) {
	m.ContentArea.Objects = nil

	var content fyne.CanvasObject
	switch name {
	case constants.NavToolbox:
		content = m.calcPage.Build()
	case constants.NavFileConverter:
		content = m.fileConvPage.Build()
	case constants.NavReference:
		content = m.refPage.Build()
	case constants.NavLogbook:
		content = m.logbookPage.Build()
	case constants.NavTracker:
		content = m.trackerPage.Build()
	default:
		content = m.calcPage.Build()
	}

	m.ContentArea.Objects = []fyne.CanvasObject{content}
	m.ContentArea.Refresh()
	m.updateNavHighlights(name)
}

func (m *MainWindow) ShowAndRun() {
	m.Window.ShowAndRun()
}
