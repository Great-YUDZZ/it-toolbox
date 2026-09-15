package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
	"github.com/yudz/it-toolbox/ui/pages"
)

// NavItem represents a custom sidebar navigation item with active indicator and hover state
type NavItem struct {
	widget.BaseWidget
	title      string
	iconRes    fyne.Resource
	active     bool
	onTap      func()

	bg         *canvas.Rectangle
	leftBar    *canvas.Rectangle
	iconWidget *widget.Icon
	labelTxt   *canvas.Text
	container  *fyne.Container
}

func NewNavItem(title string, iconRes fyne.Resource, onTap func()) *NavItem {
	n := &NavItem{
		title:   title,
		iconRes: iconRes,
		onTap:   onTap,
	}
	n.ExtendBaseWidget(n)

	n.bg = canvas.NewRectangle(color.Transparent)
	n.bg.CornerRadius = constants.CurrentCornerRadius

	n.leftBar = canvas.NewRectangle(color.Transparent)
	n.leftBar.SetMinSize(fyne.NewSize(3, 30))

	n.iconWidget = widget.NewIcon(iconRes)

	n.labelTxt = canvas.NewText(title, constants.ColorTextSecondary)
	n.labelTxt.TextSize = constants.FontSizeBody
	n.labelTxt.Alignment = fyne.TextAlignLeading

	innerRow := container.NewHBox(n.iconWidget, n.labelTxt)
	content := container.NewBorder(nil, nil, n.leftBar, nil, container.NewPadded(innerRow))
	n.container = container.NewStack(n.bg, content)

	return n
}

func (n *NavItem) SetActive(active bool) {
	n.active = active
	n.bg.CornerRadius = constants.CurrentCornerRadius
	if active {
		if constants.IsNeumorphism {
			n.bg.FillColor = constants.ColorBgCardInner
			n.bg.StrokeColor = constants.ColorBorderActive
			n.bg.StrokeWidth = constants.CurrentBorderWidth
			n.leftBar.FillColor = constants.ColorAccentCobalt
			n.labelTxt.Color = constants.ColorTextPrimary
			n.labelTxt.TextStyle = fyne.TextStyle{Bold: true}
		} else {
			// Neo-Brutalism (Signature)
			n.bg.FillColor = constants.ColorAccentYellow
			n.bg.StrokeColor = constants.ColorBorderSubtle
			n.bg.StrokeWidth = constants.BorderWidthMedium
			n.leftBar.FillColor = color.Black
			n.labelTxt.Color = color.Black
			n.labelTxt.TextStyle = fyne.TextStyle{Bold: true}
		}
	} else {
		n.bg.FillColor = color.Transparent
		n.bg.StrokeColor = color.Transparent
		n.bg.StrokeWidth = 0
		n.leftBar.FillColor = color.Transparent
		n.labelTxt.Color = constants.ColorTextSecondary
		n.labelTxt.TextStyle = fyne.TextStyle{Bold: false}
	}
	n.bg.Refresh()
	n.leftBar.Refresh()
	n.labelTxt.Refresh()
}

func (n *NavItem) Tapped(_ *fyne.PointEvent) {
	if n.onTap != nil {
		n.onTap()
	}
}

func (n *NavItem) MouseIn(_ *desktop.MouseEvent) {
	if !n.active {
		n.bg.FillColor = constants.ColorBgHover
		n.bg.StrokeColor = constants.ColorBorderSubtle
		n.bg.StrokeWidth = constants.CurrentBorderWidth
		n.bg.CornerRadius = constants.CurrentCornerRadius
		n.bg.Refresh()
	}
}

func (n *NavItem) MouseMoved(_ *desktop.MouseEvent) {}

func (n *NavItem) MouseOut() {
	if !n.active {
		n.bg.FillColor = color.Transparent
		n.bg.StrokeColor = color.Transparent
		n.bg.StrokeWidth = 0
		n.bg.Refresh()
	}
}

func (n *NavItem) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

func (n *NavItem) CreateRenderer() fyne.WidgetRenderer {
	return &navItemRenderer{item: n}
}

type navItemRenderer struct {
	item *NavItem
}

func (r *navItemRenderer) Layout(s fyne.Size) {
	r.item.container.Resize(s)
	r.item.container.Move(fyne.NewPos(0, 0))
}

func (r *navItemRenderer) MinSize() fyne.Size {
	return r.item.container.MinSize()
}

func (r *navItemRenderer) Refresh() {
	canvas.Refresh(r.item)
}

func (r *navItemRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.item.container}
}

func (r *navItemRenderer) Destroy() {}

// MainWindow manages the application shell layout, theme state, and routing
type MainWindow struct {
	App           fyne.App
	Window        fyne.Window
	RootContainer *fyne.Container
	ContentArea   *fyne.Container
	StatusLabel   *canvas.Text
	ActiveMenu    string
	CurrentTheme  string
	IsDark        bool

	navToolbox  *NavItem
	navFileConv *NavItem
	navYouTube  *NavItem
	navCisco    *NavItem
	navRef      *NavItem
	navLogbook  *NavItem
	navTracker  *NavItem

	calcPage     *pages.CalculatorPage
	fileConvPage *pages.FileConverterPage
	ciscoPage    *pages.CiscoPage
	refPage      *pages.ReferencePage
	logbookPage  *pages.LogbookPage
	trackerPage  *pages.TrackerPage
}

func NewMainWindow(app fyne.App) *MainWindow {
	win := app.NewWindow(constants.AppTitle + " — " + constants.AppSubtitle)
	win.Resize(fyne.NewSize(constants.DefaultWinW, constants.DefaultWinH))
	win.CenterOnScreen()

	themeName := app.Preferences().StringWithFallback("active_theme", constants.ThemeNeoBrutalism)
	if themeName == "" {
		themeName = constants.ThemeNeoBrutalism
	}
	ApplyTheme(themeName)
	app.Settings().SetTheme(NewCustomTheme(themeName))

	statusTxt := canvas.NewText("● Sistem Siap", constants.ColorSuccess)
	statusTxt.TextSize = constants.FontSizeSmall
	statusTxt.TextStyle = fyne.TextStyle{Bold: true}

	mw := &MainWindow{
		App:          app,
		Window:       win,
		ContentArea:  container.NewStack(),
		StatusLabel:  statusTxt,
		CurrentTheme: themeName,
		IsDark:       constants.IsDarkTheme,
		calcPage:     pages.NewCalculatorPage(win),
		fileConvPage: pages.NewFileConverterPage(win),
		ciscoPage:    pages.NewCiscoPage(win),
		refPage:      pages.NewReferencePage(win),
		logbookPage:  pages.NewLogbookPage(win),
		trackerPage:  pages.NewTrackerPage(win),
	}

	mw.RootContainer = container.NewStack(mw.buildLayout())
	win.SetContent(mw.RootContainer)
	mw.showPage(constants.NavToolbox)
	return mw
}

// SwitchTheme changes the active application theme and reloads all pages
func (m *MainWindow) SwitchTheme(themeName string) {
	m.CurrentTheme = themeName
	m.IsDark = constants.IsDarkTheme
	m.App.Preferences().SetString("active_theme", themeName)
	m.App.Preferences().SetBool("theme_dark", m.IsDark)
	ApplyTheme(themeName)
	m.App.Settings().SetTheme(NewCustomTheme(themeName))

	// Rebuild pages with updated theme objects
	m.calcPage = pages.NewCalculatorPage(m.Window)
	m.fileConvPage = pages.NewFileConverterPage(m.Window)
	m.ciscoPage = pages.NewCiscoPage(m.Window)
	m.refPage = pages.NewReferencePage(m.Window)
	m.logbookPage = pages.NewLogbookPage(m.Window)
	m.trackerPage = pages.NewTrackerPage(m.Window)

	// Rebuild window layout in place without resetting window geometry
	m.RootContainer.Objects = []fyne.CanvasObject{m.buildLayout()}
	m.RootContainer.Refresh()
	m.showPage(m.ActiveMenu)
}

// CycleTheme cycles through the 3 available theme options
func (m *MainWindow) CycleTheme() {
	var nextTheme string
	switch m.CurrentTheme {
	case constants.ThemeNeoBrutalism:
		nextTheme = constants.ThemeNeumorphismLight
	case constants.ThemeNeumorphismLight:
		nextTheme = constants.ThemeNeumorphismDark
	case constants.ThemeNeumorphismDark:
		nextTheme = constants.ThemeNeoBrutalism
	default:
		nextTheme = constants.ThemeNeoBrutalism
	}
	m.SwitchTheme(nextTheme)
}

// ToggleTheme provides backward compatibility
func (m *MainWindow) ToggleTheme() {
	m.CycleTheme()
}

// ShowThemeDialog opens a modal for choosing between Neo-Brutalism and Neumorphism Light/Dark
func (m *MainWindow) ShowThemeDialog() {
	var d dialog.Dialog

	title := canvas.NewText("🎨 PILIH TEMA TAMPILAN", constants.ColorTextPrimary)
	title.TextSize = constants.FontSizeH2
	title.TextStyle = fyne.TextStyle{Bold: true}

	sub := canvas.NewText("Pilih gaya visual antarmuka IT-Toolbox yang Anda inginkan:", constants.ColorTextMuted)
	sub.TextSize = constants.FontSizeSmall

	makeOptionCard := func(themeKey, name, desc, badgeLabel string, badgeFn func(string) fyne.CanvasObject) fyne.CanvasObject {
		optTitle := canvas.NewText(name, constants.ColorTextPrimary)
		optTitle.TextSize = constants.FontSizeH3
		optTitle.TextStyle = fyne.TextStyle{Bold: true}

		badge := badgeFn(badgeLabel)
		headerRow := container.NewHBox(optTitle, badge)
		if m.CurrentTheme == themeKey {
			activeBadge := components.BadgeSuccess("AKTIF ✓")
			headerRow.Add(activeBadge)
		}

		descText := widget.NewLabel(desc)
		descText.Wrapping = fyne.TextWrapWord

		selectBtn := widget.NewButton("Pilih Tema Ini", func() {
			if d != nil {
				d.Hide()
			}
			m.SwitchTheme(themeKey)
		})
		if m.CurrentTheme == themeKey {
			selectBtn.Importance = widget.HighImportance
			selectBtn.SetText("Tema Sedang Aktif ✓")
			selectBtn.Disable()
		} else {
			selectBtn.Importance = widget.MediumImportance
		}

		cardInner := container.NewVBox(
			headerRow,
			descText,
			selectBtn,
		)
		return components.NewPlainCard(cardInner)
	}

	optBrutal := makeOptionCard(
		constants.ThemeNeoBrutalism,
		"⚡ Neo-Brutalism (Signature)",
		"Gaya retro retro paper, border tegas 2.5px solid hitam, font pitch-black, & bayangan tajam (zero blur).",
		"SIGNATURE",
		components.BadgeYellow,
	)

	optNeumorphLight := makeOptionCard(
		constants.ThemeNeumorphismLight,
		"🫧 Neumorphism — Mode Terang",
		"Soft UI monokromatik abu-abu lembut dengan bayangan ganda (dual-tone shadow: highlight putih + bayangan halus) dan sudut melengkung 14px.",
		"SOFT UI",
		components.BadgeCyan,
	)

	optNeumorphDark := makeOptionCard(
		constants.ThemeNeumorphismDark,
		"🌙 Neumorphism — Mode Gelap",
		"Dark Soft UI monokromatik slate gelap, sudut lembut, bayangan ganda emboss, & teks lembut yang nyaman di mata.",
		"DARK SOFT",
		components.BadgeIndigo,
	)

	closeBtn := widget.NewButton("Tutup", func() {
		if d != nil {
			d.Hide()
		}
	})
	closeBtn.Importance = widget.LowImportance

	content := container.NewVBox(
		title,
		sub,
		widget.NewSeparator(),
		optBrutal,
		optNeumorphLight,
		optNeumorphDark,
		widget.NewSeparator(),
		container.NewCenter(closeBtn),
	)

	scrollContent := container.NewVScroll(container.NewPadded(content))
	scrollContent.SetMinSize(fyne.NewSize(480, 420))

	d = dialog.NewCustomWithoutButtons("Ganti Tema", scrollContent, m.Window)
	d.Show()
}

func (m *MainWindow) buildLayout() fyne.CanvasObject {
	// ------------------------------------------------------------------------
	// 1. Sidebar Header (Branding & Quick Theme Toggle)
	// ------------------------------------------------------------------------
	brandBg := canvas.NewRectangle(constants.ColorAccentYellow)
	brandBg.StrokeColor = constants.ColorBorderSubtle
	brandBg.StrokeWidth = constants.CurrentBorderWidth
	brandBg.CornerRadius = constants.CurrentCornerRadius

	brandTitle := canvas.NewText("⚡ IT TOOLBOX", color.Black)
	brandTitle.TextSize = constants.FontSizeH2
	brandTitle.TextStyle = fyne.TextStyle{Bold: true}

	brandBadge := container.NewStack(brandBg, container.NewPadded(brandTitle))
	verBadge := components.BadgeCyan("v" + constants.AppVersion)

	var fullScreenBtn *widget.Button
	fullScreenBtn = widget.NewButtonWithIcon("", theme.ViewFullScreenIcon(), func() {
		isFull := !m.Window.FullScreen()
		m.Window.SetFullScreen(isFull)
		if isFull {
			fullScreenBtn.SetIcon(theme.ViewRestoreIcon())
		} else {
			fullScreenBtn.SetIcon(theme.ViewFullScreenIcon())
		}
	})
	fullScreenBtn.Importance = widget.LowImportance

	quickThemeBtn := widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		m.ShowThemeDialog()
	})
	quickThemeBtn.Importance = widget.LowImportance

	headerActions := container.NewHBox(fullScreenBtn, quickThemeBtn)

	brandHeader := container.NewBorder(nil, nil,
		container.NewHBox(brandBadge, verBadge),
		headerActions,
	)

	subTitle := canvas.NewText("ENGINEERING WORKBENCH", constants.ColorTextMuted)
	subTitle.TextSize = constants.FontSizeLabel
	subTitle.TextStyle = fyne.TextStyle{Bold: true}

	brandBox := container.NewVBox(
		brandHeader,
		subTitle,
	)

	headerCard := container.NewVBox(
		container.NewPadded(brandBox),
		widget.NewSeparator(),
	)

	// ------------------------------------------------------------------------
	// 2. Navigation Section with Refined Custom Nav Items
	// ------------------------------------------------------------------------
	secCore := canvas.NewText("ALAT & KALKULATOR", constants.ColorTextPrimary)
	secCore.TextSize = constants.FontSizeLabel
	secCore.TextStyle = fyne.TextStyle{Bold: true}

	m.navToolbox = NewNavItem(constants.NavToolbox, theme.SettingsIcon(), func() {
		m.showPage(constants.NavToolbox)
	})
	m.navFileConv = NewNavItem(constants.NavFileConverter, theme.FolderOpenIcon(), func() {
		m.showPage(constants.NavFileConverter)
	})
	m.navYouTube = NewNavItem(constants.NavYouTube, theme.DownloadIcon(), func() {
		m.showPage(constants.NavYouTube)
	})

	secDocs := canvas.NewText("PENGETAHUAN & LOG", constants.ColorTextPrimary)
	secDocs.TextSize = constants.FontSizeLabel
	secDocs.TextStyle = fyne.TextStyle{Bold: true}

	m.navCisco = NewNavItem(constants.NavCisco, theme.ComputerIcon(), func() {
		m.showPage(constants.NavCisco)
	})
	m.navRef = NewNavItem(constants.NavReference, theme.InfoIcon(), func() {
		m.showPage(constants.NavReference)
	})
	m.navLogbook = NewNavItem(constants.NavLogbook, theme.DocumentCreateIcon(), func() {
		m.showPage(constants.NavLogbook)
	})
	m.navTracker = NewNavItem(constants.NavTracker, theme.ListIcon(), func() {
		m.showPage(constants.NavTracker)
	})

	navContainer := container.NewVBox(
		secCore,
		m.navToolbox,
		m.navFileConv,
		m.navYouTube,
		widget.NewSeparator(),
		secDocs,
		m.navCisco,
		m.navRef,
		m.navLogbook,
		m.navTracker,
	)

	// ------------------------------------------------------------------------
	// 3. Sidebar Footer (Theme Switcher, Status Pill & Environment Info)
	// ------------------------------------------------------------------------
	var themeBtnText string
	switch m.CurrentTheme {
	case constants.ThemeNeumorphismLight:
		themeBtnText = "🫧 Neumorph (Terang)"
	case constants.ThemeNeumorphismDark:
		themeBtnText = "🌙 Neumorph (Gelap)"
	default:
		themeBtnText = "⚡ Neo-Brutalism"
	}
	themeBtn := widget.NewButtonWithIcon(themeBtnText, theme.ColorPaletteIcon(), func() {
		m.ShowThemeDialog()
	})
	themeBtn.Importance = widget.LowImportance

	statusPillBg := canvas.NewRectangle(constants.ColorSuccess)
	statusPillBg.StrokeColor = constants.ColorBorderSubtle
	statusPillBg.StrokeWidth = constants.CurrentBorderWidth
	statusPillBg.CornerRadius = constants.CurrentBadgeRadius

	if constants.IsNeumorphism {
		if constants.IsDarkTheme {
			m.StatusLabel.Color = color.RGBA{R: 0x06, G: 0x4E, B: 0x3B, A: 0xFF}
		} else {
			m.StatusLabel.Color = color.RGBA{R: 0x06, G: 0x5F, B: 0x46, A: 0xFF}
		}
	} else {
		m.StatusLabel.Color = color.Black
	}
	m.StatusLabel.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	statusPill := container.NewStack(statusPillBg, container.NewPadded(m.StatusLabel))

	envTxt := canvas.NewText("Go 1.25 • Fyne v2.8 • Offline", constants.ColorTextMuted)
	envTxt.TextSize = constants.FontSizeLabel

	footerBox := container.NewVBox(
		widget.NewSeparator(),
		container.NewPadded(container.NewVBox(
			themeBtn,
			statusPill,
			envTxt,
		)),
	)

	sidebarContent := container.NewBorder(headerCard, footerBox, nil, nil, container.NewPadded(navContainer))

	// Sidebar background
	bgSidebar := canvas.NewRectangle(constants.ColorBgSidebar)
	bgSidebar.SetMinSize(fyne.NewSize(float32(constants.SidebarWidth), 0))
	sidebarWrapper := container.NewMax(bgSidebar, sidebarContent)
	sidebarWithSep := container.NewBorder(nil, nil, nil, widget.NewSeparator(), sidebarWrapper)

	// Main Layout: Sidebar on Left, Content Area in Center
	mainLayout := container.NewBorder(nil, nil, sidebarWithSep, nil, container.NewPadded(m.ContentArea))
	return mainLayout
}

func (m *MainWindow) updateNavHighlights(active string) {
	m.ActiveMenu = active

	m.navToolbox.SetActive(active == constants.NavToolbox)
	m.navFileConv.SetActive(active == constants.NavFileConverter)
	m.navYouTube.SetActive(active == constants.NavYouTube)
	m.navCisco.SetActive(active == constants.NavCisco)
	m.navRef.SetActive(active == constants.NavReference)
	m.navLogbook.SetActive(active == constants.NavLogbook)
	m.navTracker.SetActive(active == constants.NavTracker)

	m.StatusLabel.Text = fmt.Sprintf("● %s", active)
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
	case constants.NavYouTube:
		content = m.fileConvPage.BuildYouTubePage()
	case constants.NavCisco:
		content = m.ciscoPage.Build()
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
