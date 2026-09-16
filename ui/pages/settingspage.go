package pages

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/core/shortcut"
	"github.com/yudz/it-toolbox/core/updater"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

// SettingsPage provides application configuration, GitHub release checking, and system information
type SettingsPage struct {
	window        fyne.Window
	onThemeChange func(string)

	isChecking       bool
	lastResult       *updater.UpdateCheckResult
	lastError        error
	renderResultFunc func()
}

// NewSettingsPage constructs a new SettingsPage instance
func NewSettingsPage(win fyne.Window, onThemeChange func(string)) *SettingsPage {
	return &SettingsPage{
		window:        win,
		onThemeChange: onThemeChange,
	}
}

// SetResult pre-populates the update check result (e.g. from background check)
func (p *SettingsPage) SetResult(res *updater.UpdateCheckResult) {
	p.lastResult = res
	p.lastError = nil
	if p.renderResultFunc != nil {
		p.renderResultFunc()
	}
}

// SetError updates the UI with an error state (e.g. offline)
func (p *SettingsPage) SetError(err error) {
	p.lastError = err
	p.lastResult = nil
	if p.renderResultFunc != nil {
		p.renderResultFunc()
	}
}

// BuildLayout renders the entire Settings view
func (p *SettingsPage) BuildLayout() fyne.CanvasObject {
	// 1. Hero Header
	title := canvas.NewText("Pengaturan & Pembaruan Sistem", constants.ColorTextPrimary)
	title.TextSize = constants.FontSizeH1
	title.TextStyle = fyne.TextStyle{Bold: true}

	desc := canvas.NewText("Kelola preferensi aplikasi, cek update versi terbaru dari GitHub, dan informasi sistem.", constants.ColorTextSecondary)
	desc.TextSize = constants.FontSizeBody

	heroLeft := container.NewVBox(title, desc)
	heroBadge := components.BadgeCyan("PENGATURAN")
	hero := container.NewBorder(nil, nil, nil, container.NewCenter(heroBadge), heroLeft)

	// 2. Main Content Cards
	updateCard := p.buildUpdateCard()
	themeCard := p.buildThemeCard()
	systemCard := p.buildSystemCard()

	contentList := container.NewVBox(
		container.NewPadded(hero),
		widget.NewSeparator(),
		updateCard,
		themeCard,
		systemCard,
	)

	return container.NewVScroll(container.NewPadded(contentList))
}

// ----------------------------------------------------------------------------
// 1. Card: Pembaruan Aplikasi (GitHub Release Checker)
// ----------------------------------------------------------------------------
func (p *SettingsPage) buildUpdateCard() fyne.CanvasObject {
	cardTitle := canvas.NewText("Pembaruan Aplikasi (GitHub Release)", constants.ColorTextPrimary)
	cardTitle.TextSize = constants.FontSizeH2
	cardTitle.TextStyle = fyne.TextStyle{Bold: true}

	currBadge := components.BadgeIndigo("Versi Terpasang: v" + constants.AppVersion)
	header := container.NewBorder(nil, nil, cardTitle, currBadge)

	// Dynamic result container
	resultContainer := container.NewVBox()

	var checkBtn *widget.Button

	renderResult := func() {
		resultContainer.Objects = nil

		if p.isChecking {
			spinner := widget.NewProgressBarInfinite()
			spinnerText := canvas.NewText("Sedang menghubungkan ke GitHub dan memeriksa rilis terbaru...", constants.ColorTextMuted)
			spinnerText.TextSize = constants.FontSizeBody
			resultContainer.Add(container.NewVBox(spinnerText, spinner))
			resultContainer.Refresh()
			return
		}

		if p.lastError != nil {
			errBadge := components.BadgeDanger("GAGAL MEMERIKSA")
			errMsg := canvas.NewText(fmt.Sprintf("Gagal memeriksa: %v", p.lastError), constants.ColorDanger)
			errMsg.TextSize = constants.FontSizeBody
			errMsg.TextStyle = fyne.TextStyle{Bold: true}

			hint := canvas.NewText("Pastikan komputer Anda terhubung ke jaringan internet lalu coba lagi.", constants.ColorTextMuted)
			hint.TextSize = constants.FontSizeSmall

			errBox := container.NewVBox(
				container.NewHBox(errBadge),
				errMsg,
				hint,
			)
			resultContainer.Add(container.NewPadded(errBox))
			resultContainer.Refresh()
			return
		}

		if p.lastResult != nil {
			res := p.lastResult
			if res.HasUpdate {
				// UPDATE AVAILABLE!
				upBadge := components.BadgeSuccess("UPDATE BARU TERSEDIA")
				upTitle := canvas.NewText(fmt.Sprintf("%s (%s)", res.Release.Name, res.Release.TagName), constants.ColorSuccess)
				upTitle.TextSize = constants.FontSizeH2
				upTitle.TextStyle = fyne.TextStyle{Bold: true}

				notesHdr := canvas.NewText("Catatan Pembaruan (Changelog):", constants.ColorTextPrimary)
				notesHdr.TextSize = constants.FontSizeLabel
				notesHdr.TextStyle = fyne.TextStyle{Bold: true}

				notesArea := widget.NewMultiLineEntry()
				notesArea.SetText(res.Release.Body)
				notesArea.TextStyle = fyne.TextStyle{Monospace: true}
				notesArea.Wrapping = fyne.TextWrapOff
				notesArea.Scroll = fyne.ScrollNone

				openReleaseBtn := widget.NewButtonWithIcon("Buka Halaman Rilis GitHub", theme.NavigateNextIcon(), func() {
					if u, err := url.Parse(res.Release.HTMLURL); err == nil {
						_ = fyne.CurrentApp().OpenURL(u)
					}
				})
				openReleaseBtn.Importance = widget.HighImportance

				// Find best asset for current OS
				var directDownloadBtn *widget.Button
				for _, asset := range res.Release.Assets {
					a := asset
					isWin := runtime.GOOS == "windows" && (strings.HasSuffix(a.Name, ".exe") || strings.HasSuffix(a.Name, ".zip"))
					isLinux := runtime.GOOS == "linux" && strings.HasSuffix(a.Name, ".deb")
					if isWin || isLinux {
						directDownloadBtn = widget.NewButtonWithIcon(fmt.Sprintf("Unduh Langsung: %s", a.Name), theme.DownloadIcon(), func() {
							if u, err := url.Parse(a.BrowserDownloadURL); err == nil {
								_ = fyne.CurrentApp().OpenURL(u)
							}
						})
						directDownloadBtn.Importance = widget.MediumImportance
						break
					}
				}

				var actionRow fyne.CanvasObject
				if directDownloadBtn != nil {
					actionRow = container.NewHBox(openReleaseBtn, directDownloadBtn)
				} else {
					actionRow = container.NewHBox(openReleaseBtn)
				}

				upBox := container.NewVBox(
					container.NewHBox(upBadge),
					upTitle,
					widget.NewSeparator(),
					notesHdr,
					container.NewPadded(notesArea),
					widget.NewSeparator(),
					actionRow,
				)
				resultContainer.Add(components.NewPlainCardWithAccent(upBox, constants.ColorSuccess))
			} else {
				// ALREADY LATEST VERSION!
				okBadge := components.BadgeSuccess("VERSI TERBARU")
				okText := canvas.NewText(fmt.Sprintf("Aplikasi Anda sudah menggunakan versi paling mutakhir (%s). Belum ada update baru.", res.LatestVersion), constants.ColorTextPrimary)
				okText.TextSize = constants.FontSizeBody

				timeInfo := canvas.NewText(fmt.Sprintf("Terakhir diperiksa: %s", res.CheckedAt.Format("15:04:05 WIB (02 Jan 2006)")), constants.ColorTextMuted)
				timeInfo.TextSize = constants.FontSizeSmall

				okBox := container.NewVBox(
					container.NewHBox(okBadge),
					okText,
					timeInfo,
				)
				resultContainer.Add(container.NewPadded(okBox))
			}
			resultContainer.Refresh()
			return
		}

		// Initial state before check
		initialText := canvas.NewText("Klik tombol 'Periksa Pembaruan Sekarang' untuk mengecek ketersediaan versi terbaru di GitHub.", constants.ColorTextMuted)
		initialText.TextSize = constants.FontSizeBody
		resultContainer.Add(initialText)
		resultContainer.Refresh()
	}

	checkBtn = widget.NewButtonWithIcon("Periksa Pembaruan Sekarang", theme.ViewRefreshIcon(), func() {
		p.isChecking = true
		p.lastError = nil
		renderResult()

		go func() {
			res, err := updater.CheckLatestRelease(constants.AppVersion)
			fyne.Do(func() {
				p.isChecking = false
				p.lastResult = res
				p.lastError = err
				renderResult()
			})
		}()
	})
	checkBtn.Importance = widget.HighImportance

	// Auto-check preference checkbox
	appPref := fyne.CurrentApp().Preferences()
	autoCheckKey := "auto_check_update"
	autoCheckVal := appPref.BoolWithFallback(autoCheckKey, true)

	autoCheckChk := widget.NewCheck("Periksa pembaruan secara otomatis saat aplikasi dibuka (jika terhubung internet)", func(checked bool) {
		appPref.SetBool(autoCheckKey, checked)
	})
	autoCheckChk.SetChecked(autoCheckVal)

	p.renderResultFunc = renderResult
	renderResult()

	body := container.NewVBox(
		header,
		widget.NewSeparator(),
		container.NewHBox(checkBtn),
		autoCheckChk,
		widget.NewSeparator(),
		resultContainer,
	)

	return components.NewPlainCardWithAccent(body, constants.ColorTechIndigo)
}

// ----------------------------------------------------------------------------
// 2. Card: Pengaturan Tampilan & Tema
// ----------------------------------------------------------------------------
func (p *SettingsPage) buildThemeCard() fyne.CanvasObject {
	cardTitle := canvas.NewText("Tampilan & Gaya Tema (Theme Switcher)", constants.ColorTextPrimary)
	cardTitle.TextSize = constants.FontSizeH2
	cardTitle.TextStyle = fyne.TextStyle{Bold: true}

	subTitle := canvas.NewText("Pilih gaya antarmuka yang paling nyaman untuk lingkungan kerja Anda:", constants.ColorTextSecondary)
	subTitle.TextSize = constants.FontSizeBody

	btnNeo := widget.NewButtonWithIcon("Neo-Brutalism (Clean Contrast)", theme.VisibilityIcon(), func() {
		if p.onThemeChange != nil {
			p.onThemeChange(constants.ThemeNeoBrutalism)
		}
	})
	btnNeo.Importance = widget.MediumImportance

	btnLight := widget.NewButtonWithIcon("Neumorphism Light (Soft Glass)", theme.ColorPaletteIcon(), func() {
		if p.onThemeChange != nil {
			p.onThemeChange(constants.ThemeNeumorphismLight)
		}
	})
	btnLight.Importance = widget.MediumImportance

	btnDark := widget.NewButtonWithIcon("Neumorphism Dark (Deep Midnight)", theme.StorageIcon(), func() {
		if p.onThemeChange != nil {
			p.onThemeChange(constants.ThemeNeumorphismDark)
		}
	})
	btnDark.Importance = widget.MediumImportance

	themeGrid := container.NewGridWithColumns(3, btnNeo, btnLight, btnDark)

	body := container.NewVBox(
		cardTitle,
		subTitle,
		widget.NewSeparator(),
		themeGrid,
	)

	return components.NewPlainCardWithAccent(body, constants.ColorAccentYellow)
}

// ----------------------------------------------------------------------------
// 3. Card: Pintasan Desktop & Informasi Sistem
// ----------------------------------------------------------------------------
func (p *SettingsPage) buildSystemCard() fyne.CanvasObject {
	cardTitle := canvas.NewText("Pintasan Sistem & Informasi Aplikasi", constants.ColorTextPrimary)
	cardTitle.TextSize = constants.FontSizeH2
	cardTitle.TextStyle = fyne.TextStyle{Bold: true}

	shortcutBtn := widget.NewButtonWithIcon("Buat / Perbarui Pintasan (Desktop & Start Menu)", theme.ContentAddIcon(), func() {
		t1 := canvas.NewText("Pintasan aplikasi 'IT Toolbox' akan dipasang di dua lokasi sistem:", constants.ColorTextPrimary)
		t1.TextSize = constants.FontSizeBody
		t2 := canvas.NewText("  • Layar Desktop (Desktop Shortcut)", constants.ColorTextSecondary)
		t2.TextSize = constants.FontSizeBody
		t3 := canvas.NewText("  • Start Menu / Peluncur Aplikasi (Start Menu Launcher)", constants.ColorTextSecondary)
		t3.TextSize = constants.FontSizeBody
		t4 := canvas.NewText("Apakah Anda ingin membuat / memperbarui pintasan sekarang?", constants.ColorTextMuted)
		t4.TextSize = constants.FontSizeSmall

		msgContent := container.NewVBox(
			t1,
			container.NewVBox(t2, t3),
			widget.NewSeparator(),
			t4,
		)

		components.ShowStyledConfirmDialog(
			p.window,
			"DESKTOP & START",
			constants.ColorAccentCobalt,
			"Pasang Pintasan Sistem",
			"Integrasi peluncur sistem operasi",
			msgContent,
			"Batal",
			"Pasang Sekarang",
			func() {
				shortcut.EnsureDesktopShortcut()
				s1 := canvas.NewText("Pintasan 'IT Toolbox' telah berhasil dibuat!", constants.ColorSuccess)
				s1.TextSize = constants.FontSizeBody
				s1.TextStyle = fyne.TextStyle{Bold: true}
				s2 := canvas.NewText("Aplikasi kini dapat diluncurkan langsung dari Desktop maupun Start Menu sistem Anda.", constants.ColorTextSecondary)
				s2.TextSize = constants.FontSizeBody
				successContent := container.NewVBox(s1, s2)
				components.ShowStyledInformationDialog(
					p.window,
					"SUKSES",
					constants.ColorSuccess,
					"Pintasan Berhasil Dipasang",
					"Telah ditambahkan ke Desktop dan Start Menu",
				successContent,
					"Tutup",
					nil,
				)
			},
		)
	})
	shortcutBtn.Importance = widget.MediumImportance

	// System Information Grid
	infoItems := container.NewVBox(
		container.NewHBox(
			canvas.NewText("Nama Aplikasi:", constants.ColorTextMuted),
			canvas.NewText(constants.AppTitle, constants.ColorTextPrimary),
		),
		container.NewHBox(
			canvas.NewText("Versi Terpasang:", constants.ColorTextMuted),
			canvas.NewText("v"+constants.AppVersion, constants.ColorTextPrimary),
		),
		container.NewHBox(
			canvas.NewText("Sistem Operasi:", constants.ColorTextMuted),
			canvas.NewText(fmt.Sprintf("%s (%s)", runtime.GOOS, runtime.GOARCH), constants.ColorTextPrimary),
		),
		container.NewHBox(
			canvas.NewText("Compiler / Runtime:", constants.ColorTextMuted),
			canvas.NewText(runtime.Version(), constants.ColorTextPrimary),
		),
		container.NewHBox(
			canvas.NewText("Database Penyimpanan:", constants.ColorTextMuted),
			canvas.NewText("SQLite Embedded Database (Lokal & Offline)", constants.ColorTextPrimary),
		),
		container.NewHBox(
			canvas.NewText("Repositori GitHub:", constants.ColorTextMuted),
			canvas.NewText("github.com/Great-YUDZZ/it-toolbox", constants.ColorTextPrimary),
		),
	)

	gitRepoBtn := widget.NewButtonWithIcon("Buka Repositori GitHub", theme.NavigateNextIcon(), func() {
		if u, err := url.Parse("https://github.com/Great-YUDZZ/it-toolbox"); err == nil {
			_ = fyne.CurrentApp().OpenURL(u)
		}
	})
	gitRepoBtn.Importance = widget.LowImportance

	body := container.NewVBox(
		cardTitle,
		widget.NewSeparator(),
		container.NewHBox(shortcutBtn),
		widget.NewSeparator(),
		infoItems,
		widget.NewSeparator(),
		container.NewHBox(gitRepoBtn),
	)

	return components.NewPlainCardWithAccent(body, constants.ColorAccentCyan)
}
