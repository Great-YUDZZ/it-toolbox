package pages

import (
	"fmt"
	"net/url"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
		initialLbl := widget.NewLabel("Klik tombol 'Periksa Pembaruan Sekarang' untuk mengecek rilis terbaru di GitHub.")
		initialLbl.Wrapping = fyne.TextWrapWord
		resultContainer.Add(initialLbl)
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

	autoCheckChk := widget.NewCheck("Periksa pembaruan otomatis saat aplikasi dibuka", func(checked bool) {
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

	btnNeo := widget.NewButtonWithIcon("Neo-Brutalism", theme.VisibilityIcon(), func() {
		if p.onThemeChange != nil {
			p.onThemeChange(constants.ThemeNeoBrutalism)
		}
	})
	btnNeo.Importance = widget.MediumImportance

	btnLight := widget.NewButtonWithIcon("Neumorphism Light", theme.ColorPaletteIcon(), func() {
		if p.onThemeChange != nil {
			p.onThemeChange(constants.ThemeNeumorphismLight)
		}
	})
	btnLight.Importance = widget.MediumImportance

	btnDark := widget.NewButtonWithIcon("Neumorphism Dark", theme.StorageIcon(), func() {
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

	shortcutBtn := widget.NewButtonWithIcon("Buat / Perbarui Pintasan Sistem", theme.ContentAddIcon(), func() {
		var d dialog.Dialog

		titleTxt := canvas.NewText("Pemasangan Pintasan Sistem", constants.ColorTextPrimary)
		titleTxt.TextSize = constants.FontSizeH2
		titleTxt.TextStyle = fyne.TextStyle{Bold: true}

		badge := components.BadgeCyan("PINTASAN SISTEM")
		headerBox := container.NewVBox(
			container.NewHBox(titleTxt, badge),
			canvas.NewText("Pilih lokasi pembuatan pintasan aplikasi IT Toolbox:", constants.ColorTextMuted),
			widget.NewSeparator(),
		)

		chkDesktop := widget.NewCheck("Pintasan Layar Desktop (Desktop Shortcut)", nil)
		chkDesktop.SetChecked(true)

		chkStart := widget.NewCheck("Pintasan Start Menu / Menu Aplikasi (Start Menu Launcher)", nil)
		chkStart.SetChecked(true)

		bodyContent := container.NewVBox(
			chkDesktop,
			chkStart,
		)

		btnCancel := widget.NewButtonWithIcon("Batal", theme.CancelIcon(), func() {
			if d != nil {
				d.Hide()
			}
		})
		btnCancel.Importance = widget.LowImportance

		btnInstall := widget.NewButtonWithIcon("Pasang Sekarang", theme.ConfirmIcon(), func() {
			if d != nil {
				d.Hide()
			}
			shortcut.CreateShortcuts(chkDesktop.Checked, chkStart.Checked)

			var locs []string
			if chkDesktop.Checked {
				locs = append(locs, "Desktop")
			}
			if chkStart.Checked {
				locs = append(locs, "Start Menu")
			}
			locText := strings.Join(locs, " dan ")
			if len(locs) == 0 {
				locText = "Tidak ada lokasi yang dipilih"
			}

			s1 := canvas.NewText(fmt.Sprintf("Pintasan IT Toolbox berhasil dipasang pada: %s.", locText), constants.ColorSuccess)
			s1.TextSize = constants.FontSizeBody
			s1.TextStyle = fyne.TextStyle{Bold: true}
			s2 := canvas.NewText("Aplikasi kini siap diakses dengan cepat dari sistem Anda.", constants.ColorTextSecondary)
			s2.TextSize = constants.FontSizeBody
			successContent := container.NewVBox(s1, s2)
			components.ShowStyledInformationDialog(
				p.window,
				"SUKSES",
				constants.ColorSuccess,
				"Pintasan Berhasil Dipasang",
				"Integrasi sistem operasi selesai",
				successContent,
				"Tutup",
				nil,
			)
		})
		btnInstall.Importance = widget.HighImportance

		actionBar := container.NewBorder(nil, nil, nil, container.NewHBox(btnCancel, btnInstall))

		dialogBody := container.NewBorder(
			headerBox,
			container.NewVBox(widget.NewSeparator(), actionBar),
			nil,
			nil,
			container.NewPadded(bodyContent),
		)

		card := components.NewPlainCardWithAccent(dialogBody, constants.ColorAccentCobalt)
		d = dialog.NewCustomWithoutButtons("", card, p.window)
		d.Resize(fyne.NewSize(580, 280))
		d.Show()
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
