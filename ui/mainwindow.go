package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
			if constants.IsDarkTheme {
				// Glassmorphic Neumorphism Dark: Frosted dark indigo-blue glass pill with specular rim
				n.bg.FillColor = color.RGBA{R: 0x1A, G: 0x24, B: 0x38, A: 0xEE}
				n.bg.StrokeColor = color.RGBA{R: 0x3E, G: 0x54, B: 0x7A, A: 0xDD}
				n.bg.StrokeWidth = constants.CurrentBorderWidth
				n.leftBar.FillColor = constants.ColorAccentCobalt
				n.labelTxt.Color = constants.ColorTextPrimary
				n.labelTxt.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				// Glassmorphic Neumorphism: Luminous Ice-Blue active glass pill
				n.bg.FillColor = color.RGBA{R: 0xDB, G: 0xEA, B: 0xFE, A: 0xFF}
				n.bg.StrokeColor = color.RGBA{R: 0x93, G: 0xC5, B: 0xFD, A: 0xFF}
				n.bg.StrokeWidth = constants.CurrentBorderWidth
				n.leftBar.FillColor = constants.ColorAccentCobalt
				n.labelTxt.Color = color.RGBA{R: 0x1D, G: 0x4E, B: 0xD8, A: 0xFF}
				n.labelTxt.TextStyle = fyne.TextStyle{Bold: true}
			}
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
		if constants.ActiveTheme == constants.ThemeNeumorphismLight {
			n.bg.FillColor = color.RGBA{R: 0xEE, G: 0xF4, B: 0xFC, A: 0xFF}
			n.bg.StrokeColor = color.RGBA{R: 0xBF, G: 0xDB, B: 0xFE, A: 0xFF}
		} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
			n.bg.FillColor = color.RGBA{R: 0x16, G: 0x20, B: 0x32, A: 0xBB}
			n.bg.StrokeColor = color.RGBA{R: 0x2D, G: 0x3C, B: 0x58, A: 0xAA}
		} else {
			n.bg.FillColor = constants.ColorBgHover
			n.bg.StrokeColor = constants.ColorBorderSubtle
		}
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

// ShowThemeMenu opens a non-modal popup menu anchored directly at the trigger button
func (m *MainWindow) ShowThemeMenu(anchor fyne.CanvasObject) {
	headerTxt := canvas.NewText("🎨 GANTI TEMA TAMPILAN", constants.ColorTextPrimary)
	headerTxt.TextSize = constants.FontSizeSmall
	headerTxt.TextStyle = fyne.TextStyle{Bold: true}

	var pop *widget.PopUp

	makeThemeRow := func(k, name, desc string, icon fyne.Resource) fyne.CanvasObject {
		titleTxt := canvas.NewText(name, constants.ColorTextPrimary)
		titleTxt.TextSize = constants.FontSizeBody
		titleTxt.TextStyle = fyne.TextStyle{Bold: true}

		descTxt := canvas.NewText(desc, constants.ColorTextMuted)
		descTxt.TextSize = constants.FontSizeLabel

		textCol := container.NewVBox(titleTxt, descTxt)

		var rightBadge fyne.CanvasObject
		if m.CurrentTheme == k {
			rightBadge = components.BadgeSuccess("AKTIF ✓")
		} else {
			rightBadge = canvas.NewText("", constants.ColorTextMuted)
		}

		iconWidget := widget.NewIcon(icon)
		rowContent := container.NewBorder(nil, nil, iconWidget, rightBadge, container.NewPadded(textCol))

		btn := widget.NewButton("", func() {
			if pop != nil {
				pop.Hide()
			}
			m.SwitchTheme(k)
		})
		btn.Importance = widget.LowImportance

		return container.NewStack(btn, rowContent)
	}

	rowBrutal := makeThemeRow(constants.ThemeNeoBrutalism, "Neo-Brutalism", "Gaya retro paper & hard shadow", theme.SettingsIcon())
	rowLight := makeThemeRow(constants.ThemeNeumorphismLight, "Neumorphism Glass", "Soft glass & ambient aurora", theme.ColorPaletteIcon())
	rowDark := makeThemeRow(constants.ThemeNeumorphismDark, "Neumorphism Glass (Gelap)", "Frosted obsidian glass & ambient midnight", theme.HomeIcon())

	popInner := container.NewVBox(
		container.NewPadded(headerTxt),
		widget.NewSeparator(),
		rowBrutal,
		rowLight,
		rowDark,
	)

	card := components.NewPlainCard(popInner)
	pop = widget.NewPopUp(card, m.Window.Canvas())

	cardMin := card.MinSize()
	var popX, popY float32

	if anchor != nil {
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(anchor)
		aSize := anchor.Size()

		popX = pos.X

		// If anchor is in lower half of screen, pop UP above the button
		if pos.Y > cardMin.Height+20 {
			popY = pos.Y - cardMin.Height - 8
		} else {
			// Pop DOWN below the button
			popY = pos.Y + aSize.Height + 8
		}

		// Ensure popup stays within canvas boundaries
		canvasW := m.Window.Canvas().Size().Width
		if popX+cardMin.Width > canvasW-8 {
			popX = canvasW - cardMin.Width - 8
		}
		if popX < 8 {
			popX = 8
		}
		if popY < 8 {
			popY = 8
		}
	} else {
		popX = 16
		popY = m.Window.Canvas().Size().Height - cardMin.Height - 60
	}

	pop.ShowAtPosition(fyne.NewPos(popX, popY))
}

// ShowThemeDialog provides backward compatibility, forwarding to ShowThemeMenu
func (m *MainWindow) ShowThemeDialog() {
	m.ShowThemeMenu(nil)
}

func (m *MainWindow) buildLayout() fyne.CanvasObject {
	// ------------------------------------------------------------------------
	// 1. Sidebar Header (Branding & Quick Theme Toggle)
	// ------------------------------------------------------------------------
	var brandFill color.Color = constants.ColorAccentYellow
	var brandTextColor color.Color = color.Black
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		brandFill = constants.ColorAccentCobalt
		brandTextColor = color.White
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		brandFill = color.RGBA{R: 0x16, G: 0x1E, B: 0x2E, A: 0xFF}
		brandTextColor = color.RGBA{R: 0xF8, G: 0xFA, B: 0xFC, A: 0xFF}
	}

	brandBg := canvas.NewRectangle(brandFill)
	brandBg.StrokeColor = constants.ColorBorderSubtle
	brandBg.StrokeWidth = constants.CurrentBorderWidth
	brandBg.CornerRadius = constants.CurrentCornerRadius

	brandTitle := canvas.NewText("⚡ IT TOOLBOX", brandTextColor)
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

	var quickThemeBtn *widget.Button
	quickThemeBtn = widget.NewButtonWithIcon("", theme.ColorPaletteIcon(), func() {
		m.ShowThemeMenu(quickThemeBtn)
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
		themeBtnText = "🫧 Neumorph Glass"
	case constants.ThemeNeumorphismDark:
		themeBtnText = "🌙 Dark Glass Neumorph"
	default:
		themeBtnText = "⚡ Neo-Brutalism"
	}
	var themeBtn *widget.Button
	themeBtn = widget.NewButtonWithIcon(themeBtnText, theme.ColorPaletteIcon(), func() {
		m.ShowThemeMenu(themeBtn)
	})
	themeBtn.Importance = widget.LowImportance

	var pillFill color.Color = constants.ColorSuccess
	var pillBorder color.Color = constants.ColorBorderSubtle
	var pillTextColor color.Color = color.Black
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		pillFill = color.RGBA{R: 0xD1, G: 0xFA, B: 0xE5, A: 0xFF}
		pillBorder = color.RGBA{R: 0x6E, G: 0xE7, B: 0xB7, A: 0xFF}
		pillTextColor = color.RGBA{R: 0x06, G: 0x5F, B: 0x46, A: 0xFF}
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		pillFill = color.RGBA{R: 0x06, G: 0x4E, B: 0x3B, A: 0xEE}
		pillBorder = color.RGBA{R: 0x05, G: 0x96, B: 0x69, A: 0xEE}
		pillTextColor = color.RGBA{R: 0x6E, G: 0xEE, B: 0xB7, A: 0xFF}
	}

	statusPillBg := canvas.NewRectangle(pillFill)
	statusPillBg.StrokeColor = pillBorder
	statusPillBg.StrokeWidth = constants.CurrentBorderWidth
	statusPillBg.CornerRadius = constants.CurrentBadgeRadius
	m.StatusLabel.Color = pillTextColor
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
	var bgSidebar fyne.CanvasObject
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		grad := canvas.NewVerticalGradient(
			color.RGBA{R: 0xED, G: 0xF3, B: 0xFB, A: 0xFF},
			color.RGBA{R: 0xE2, G: 0xEA, B: 0xF5, A: 0xFF},
		)
		grad.SetMinSize(fyne.NewSize(float32(constants.SidebarWidth), 0))
		bgSidebar = grad
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		grad := canvas.NewVerticalGradient(
			color.RGBA{R: 0x11, G: 0x16, B: 0x24, A: 0xFF}, // Dark frosted glass top
			color.RGBA{R: 0x0A, G: 0x0E, B: 0x17, A: 0xFF}, // Deep space midnight bottom
		)
		grad.SetMinSize(fyne.NewSize(float32(constants.SidebarWidth), 0))
		bgSidebar = grad
	} else {
		rect := canvas.NewRectangle(constants.ColorBgSidebar)
		rect.SetMinSize(fyne.NewSize(float32(constants.SidebarWidth), 0))
		bgSidebar = rect
	}
	sidebarWrapper := container.NewMax(bgSidebar, sidebarContent)
	sidebarWithSep := container.NewBorder(nil, nil, nil, widget.NewSeparator(), sidebarWrapper)

	// Content Area with ambient luminous backdrop for Neumorphism Light (Glassmorphic)
	var contentAreaWrapper fyne.CanvasObject
	if constants.ActiveTheme == constants.ThemeNeumorphismLight {
		ambientBg := canvas.NewHorizontalGradient(
			color.RGBA{R: 0xD6, G: 0xE6, B: 0xFD, A: 0xFF}, // Soft Sky Cyan (#D6E6FD)
			color.RGBA{R: 0xEE, G: 0xE2, B: 0xFD, A: 0xFF}, // Soft Dreamy Lavender (#EEE2FD)
		)
		contentAreaWrapper = container.NewStack(ambientBg, container.NewPadded(m.ContentArea))
	} else if constants.ActiveTheme == constants.ThemeNeumorphismDark {
		ambientBg := canvas.NewHorizontalGradient(
			color.RGBA{R: 0x0A, G: 0x0E, B: 0x18, A: 0xFF}, // Deep space midnight navy (#0A0E18)
			color.RGBA{R: 0x14, G: 0x1A, B: 0x2D, A: 0xFF}, // Ambient Midnight Indigo Glass (#141A2D)
		)
		contentAreaWrapper = container.NewStack(ambientBg, container.NewPadded(m.ContentArea))
	} else {
		bgContent := canvas.NewRectangle(constants.ColorBgBase)
		contentAreaWrapper = container.NewStack(bgContent, container.NewPadded(m.ContentArea))
	}

	// Main Layout: Sidebar on Left, Content Area in Center
	mainLayout := container.NewBorder(nil, nil, sidebarWithSep, nil, contentAreaWrapper)
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

// ShowPage switches the active page
func (m *MainWindow) ShowPage(name string) {
	m.showPage(name)
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
