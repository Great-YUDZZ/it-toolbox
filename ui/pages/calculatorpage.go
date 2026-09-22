package pages

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/core/calculators"
	"github.com/yudz/it-toolbox/core/formatters"
	"github.com/yudz/it-toolbox/database"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

type CalculatorPage struct {
	window fyne.Window
}

func NewCalculatorPage(win fyne.Window) *CalculatorPage {
	return &CalculatorPage{window: win}
}

func (p *CalculatorPage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavToolbox,
		"Alat praktis untuk subnet, konversi bilangan, hashing, dan formatter teks.",
		components.BadgeCyan("NETWORK & CODE"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabSubnet, theme.ComputerIcon(), p.buildSubnetTab()),
		container.NewTabItemWithIcon(constants.TabConverter, theme.MenuIcon(), p.buildConverterTab()),
		container.NewTabItemWithIcon(constants.TabHashGen, theme.VisibilityIcon(), p.buildHashGenTab()),
		container.NewTabItemWithIcon(constants.TabFormatter, theme.DocumentCreateIcon(), p.buildFormatterTab()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
}

func (p *CalculatorPage) copyToClip(txt string) {
	p.window.Clipboard().SetContent(txt)
	dialog.ShowInformation("Clipboard", constants.StatusCopied, p.window)
}

// ----------------------------------------------------------------------------
// 1. Tab Subnet & IP Calculator (Termasuk Fitur Rekomendasi Host)
// ----------------------------------------------------------------------------
func (p *CalculatorPage) buildSubnetTab() fyne.CanvasObject {
	// Input form
	hostCountEntry := widget.NewEntry()
	hostCountEntry.SetPlaceHolder(constants.HostInputPlaceholder)
	hostCountEntry.SetText("50")

	baseIPEntry := widget.NewEntry()
	baseIPEntry.SetPlaceHolder(constants.BaseIPPlaceholder)
	baseIPEntry.SetText("192.168.1.0")

	// Stat KPI Cards with colorful Neo-Brutalist backgrounds and solid black fonts
	statPrefix := components.NewStatCard("REKOMENDASI PREFIX & NETMASK", "/26", constants.ColorInfo)
	statRange := components.NewStatCard("RENTANG HOST USABLE", "-", constants.ColorSuccess)
	statCapacity := components.NewStatCard("TOTAL HOST & EFISIENSI", "-", constants.ColorAccentYellow)
	statBroadcast := components.NewStatCard("BROADCAST & WILDCARD", "-", constants.ColorWarning)

	hintLabel := widget.NewLabel("-")
	hintLabel.Wrapping = fyne.TextWrapWord
	hintLabel.TextStyle = fyne.TextStyle{Bold: true}

	hintBg := canvas.NewRectangle(constants.ColorBgCardInner)
	hintBg.CornerRadius = constants.CurrentBadgeRadius
	hintBg.StrokeColor = constants.ColorBorderSubtle
	hintBg.StrokeWidth = constants.CurrentBorderWidth
	hintCallout := container.NewStack(hintBg, container.NewPadded(hintLabel))

	calcRecommendation := func() {
		hStr := strings.TrimSpace(hostCountEntry.Text)
		if hStr == "" {
			return
		}
		needed, err := strconv.Atoi(hStr)
		if err != nil || needed <= 0 {
			statPrefix.SetValue("Input tidak valid")
			statPrefix.SetColor(constants.ColorDanger)
			return
		}

		rec, err := calculators.FindSubnetForHosts(needed, baseIPEntry.Text)
		if err != nil {
			statPrefix.SetValue("Error: " + err.Error())
			statPrefix.SetColor(constants.ColorDanger)
			return
		}

		statPrefix.SetColor(constants.ColorInfo)
		statPrefix.SetValue(rec.CIDR)
		statPrefix.SetSubtext(fmt.Sprintf("Netmask: %s (Paling Hemat)", rec.Netmask))

		statRange.SetColor(constants.ColorSuccess)
		statRange.SetValue(fmt.Sprintf("%s -> %s", rec.FirstHost, rec.LastHost))
		statRange.SetSubtext(fmt.Sprintf("Tersedia %d IP Usable", rec.AllocatedHosts))

		statCapacity.SetColor(constants.ColorAccentYellow)
		statCapacity.SetValue(fmt.Sprintf("%d Host (%.1f%%)", rec.AllocatedHosts, rec.Efficiency))
		statCapacity.SetSubtext(fmt.Sprintf("%d Terpakai | %d Sisa", rec.NeededHosts, rec.WastedHosts))

		statBroadcast.SetColor(constants.ColorWarning)
		statBroadcast.SetValue(rec.Broadcast)
		statBroadcast.SetSubtext(fmt.Sprintf("Wildcard: %s", rec.WildcardMask))

		hintLabel.SetText(fmt.Sprintf("Saran Alokasi: %s", rec.ClassHint))
	}

	hostCountEntry.OnChanged = func(string) { calcRecommendation() }
	baseIPEntry.OnChanged = func(string) { calcRecommendation() }

	// Preset Quick Chips
	makeHostChip := func(h int) *widget.Button {
		btn := widget.NewButton(fmt.Sprintf("%d Host", h), func() {
			hostCountEntry.SetText(strconv.Itoa(h))
		})
		btn.Importance = widget.LowImportance
		return btn
	}

	hostChips := container.NewHBox(
		canvas.NewText("Preset:", constants.ColorTextSecondary),
		makeHostChip(10),
		makeHostChip(50),
		makeHostChip(100),
		makeHostChip(250),
	)

	copyRecBtn := widget.NewButtonWithIcon("Salin Ringkasan", theme.ContentCopyIcon(), func() {
		hStr := strings.TrimSpace(hostCountEntry.Text)
		needed, _ := strconv.Atoi(hStr)
		rec, err := calculators.FindSubnetForHosts(needed, baseIPEntry.Text)
		if err == nil {
			summary := fmt.Sprintf("Alokasi Subnet untuk %d Host:\nCIDR: %s\nNetmask: %s\nRentang Host: %s s/d %s\nTotal Usable: %d Host (Efisiensi: %.1f%%)\nBroadcast: %s",
				rec.NeededHosts, rec.CIDR, rec.Netmask, rec.FirstHost, rec.LastHost, rec.AllocatedHosts, rec.Efficiency, rec.Broadcast)
			p.copyToClip(summary)
		}
	})
	copyRecBtn.Importance = widget.LowImportance

	statGrid := container.NewGridWithColumns(2,
		statPrefix.Widget,
		statRange.Widget,
		statCapacity.Widget,
		statBroadcast.Widget,
	)

	hostLbl := canvas.NewText("Jumlah Host yang Dibutuhkan:", constants.ColorTextPrimary)
	hostLbl.TextSize = constants.FontSizeSmall
	hostLbl.TextStyle = fyne.TextStyle{Bold: true}

	ipLbl := canvas.NewText("IP Jaringan Awal (Opsional):", constants.ColorTextPrimary)
	ipLbl.TextSize = constants.FontSizeSmall
	ipLbl.TextStyle = fyne.TextStyle{Bold: true}

	calcHeaderBadge := container.NewHBox(components.BadgeMuted("HASIL KALKULASI OTOMATIS"))
	calcHeader := container.NewBorder(nil, nil, calcHeaderBadge, copyRecBtn)

	recommenderContent := container.NewVBox(
		container.NewGridWithColumns(2,
			container.NewVBox(hostLbl, hostCountEntry),
			container.NewVBox(ipLbl, baseIPEntry),
		),
		hostChips,
		widget.NewSeparator(),
		calcHeader,
		statGrid,
		container.NewPadded(hintCallout),
	)

	recommenderCard := components.NewPlainCard(recommenderContent)

	// ---------------------------------------------------------
	// FITUR 2: Kalkulator CIDR / Subnet IPv4 Standar
	// ---------------------------------------------------------
	cidrEntry := widget.NewEntry()
	cidrEntry.SetPlaceHolder("192.168.1.0/24")
	cidrEntry.SetText("192.168.1.0/24")

	statNet := components.NewStatCard("NETWORK ADDRESS", "-", constants.ColorInfo)
	statMask := components.NewStatCard("SUBNET MASK", "-", constants.ColorTechIndigo)
	statFirstHost := components.NewStatCard("HOST PERTAMA USABLE", "-", constants.ColorSuccess)
	statLastHost := components.NewStatCard("HOST TERAKHIR USABLE", "-", constants.ColorSuccess)
	statTotalHosts := components.NewStatCard("TOTAL USABLE HOST", "-", constants.ColorAccentYellow)
	statBcast := components.NewStatCard("BROADCAST ADDRESS", "-", constants.ColorWarning)

	calcCIDR := func() {
		info, err := calculators.ParseCIDR(cidrEntry.Text)
		if err != nil {
			statNet.SetValue("CIDR Tidak Valid")
			statNet.SetColor(constants.ColorDanger)
			return
		}
		statNet.SetColor(constants.ColorInfo)
		statNet.SetValue(info.NetworkAddress)
		statMask.SetColor(constants.ColorTechIndigo)
		statMask.SetValue(info.Netmask)
		statFirstHost.SetColor(constants.ColorSuccess)
		statFirstHost.SetValue(info.FirstHost)
		statLastHost.SetColor(constants.ColorSuccess)
		statLastHost.SetValue(info.LastHost)
		statTotalHosts.SetColor(constants.ColorAccentYellow)
		statTotalHosts.SetValue(fmt.Sprintf("%d Host", info.TotalHosts))
		statBcast.SetColor(constants.ColorWarning)
		statBcast.SetValue(info.BroadcastAddress)
	}

	cidrEntry.OnChanged = func(string) { calcCIDR() }

	makeCidrChip := func(c string) *widget.Button {
		btn := widget.NewButton(c, func() {
			cidrEntry.SetText(c)
		})
		btn.Importance = widget.LowImportance
		return btn
	}

	cidrChips := container.NewHBox(
		canvas.NewText("Preset CIDR:", constants.ColorTextSecondary),
		makeCidrChip("192.168.1.0/24"),
		makeCidrChip("192.168.0.0/22"),
		makeCidrChip("172.16.0.0/20"),
		makeCidrChip("10.0.0.0/16"),
		makeCidrChip("10.0.0.0/30"),
	)

	cidrGrid := container.NewGridWithColumns(2,
		statNet.Widget,
		statMask.Widget,
		statFirstHost.Widget,
		statLastHost.Widget,
		statTotalHosts.Widget,
		statBcast.Widget,
	)

	cidrInputLbl := canvas.NewText("Alamat CIDR Target:", constants.ColorTextPrimary)
	cidrInputLbl.TextSize = constants.FontSizeSmall
	cidrInputLbl.TextStyle = fyne.TextStyle{Bold: true}

	cidrContent := container.NewVBox(
		container.NewBorder(nil, nil, cidrInputLbl, nil, cidrEntry),
		cidrChips,
		widget.NewSeparator(),
		cidrGrid,
	)
	cidrCard := components.NewPlainCard(cidrContent)

	// ---------------------------------------------------------
	// FITUR 3: Perhitungan VLSM
	// ---------------------------------------------------------
	vlsmNetEntry := widget.NewEntry()
	vlsmNetEntry.SetPlaceHolder("192.168.1.0/24")
	vlsmNetEntry.SetText("192.168.1.0/24")

	vlsmHostsEntry := widget.NewEntry()
	vlsmHostsEntry.SetPlaceHolder("50, 20, 10")
	vlsmHostsEntry.SetText("50, 20, 10")

	vlsmResultArea := widget.NewMultiLineEntry()
	vlsmResultArea.SetMinRowsVisible(6)
	vlsmResultArea.TextStyle = fyne.TextStyle{Monospace: true}
	vlsmResultArea.Wrapping = fyne.TextWrapOff
	vlsmResultArea.Scroll = fyne.ScrollNone

	calcVLSM := func() {
		parts := strings.Split(vlsmHostsEntry.Text, ",")
		var hosts []int
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			h, err := strconv.Atoi(part)
			if err != nil || h <= 0 {
				dialog.ShowError(fmt.Errorf("kebutuhan host harus bilangan positif: %q", part), p.window)
				return
			}
			hosts = append(hosts, h)
		}

		res, err := calculators.CalculateVLSM(vlsmNetEntry.Text, hosts)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("=== Alokasi VLSM untuk %s (%d Subnets) ===\n\n", vlsmNetEntry.Text, len(res)))
		for i, sub := range res {
			sb.WriteString(fmt.Sprintf("Subnet #%d:\n", i+1))
			sb.WriteString(fmt.Sprintf("  CIDR: %s | Mask: %s\n", sub.CIDR, sub.Netmask))
			sb.WriteString(fmt.Sprintf("  Usable Hosts: %s s/d %s (%d hosts)\n", sub.FirstHost, sub.LastHost, sub.TotalHosts))
			sb.WriteString(fmt.Sprintf("  Broadcast: %s\n\n", sub.BroadcastAddress))
		}
		vlsmResultArea.SetText(sb.String())
	}

	vlsmBtn := widget.NewButtonWithIcon("Alokasikan VLSM", theme.ConfirmIcon(), calcVLSM)
	vlsmBtn.Importance = widget.HighImportance

	makeVlsmChip := func(txt string) *widget.Button {
		btn := widget.NewButton(txt, func() {
			vlsmHostsEntry.SetText(txt)
			calcVLSM()
		})
		btn.Importance = widget.LowImportance
		return btn
	}

	vlsmChips := container.NewHBox(
		canvas.NewText("Contoh Tugas Lab:", constants.ColorTextSecondary),
		makeVlsmChip("50, 20, 10"),
		makeVlsmChip("100, 50, 25, 12"),
		makeVlsmChip("30, 15, 6, 2"),
	)

	vlsmContent := container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabelWithStyle("Base Network:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), nil, vlsmNetEntry),
		container.NewBorder(nil, nil, widget.NewLabelWithStyle("Kebutuhan Host:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), vlsmBtn, vlsmHostsEntry),
		vlsmChips,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Hasil Pembagian Subnet:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		vlsmResultArea,
	)
	vlsmCard := components.NewPlainCard(vlsmContent)

	// Initial calculation
	calcRecommendation()
	calcCIDR()

	return container.NewVScroll(container.NewVBox(
		recommenderCard,
		cidrCard,
		vlsmCard,
	))
}

// ----------------------------------------------------------------------------
// 2. Tab Number & Data Converter (Bi-directional Live)
// ----------------------------------------------------------------------------
func (p *CalculatorPage) buildConverterTab() fyne.CanvasObject {
	decEntry := widget.NewEntry()
	decEntry.SetPlaceHolder("255")
	binEntry := widget.NewEntry()
	binEntry.SetPlaceHolder("11111111")
	hexEntry := widget.NewEntry()
	hexEntry.SetPlaceHolder("FF")
	octEntry := widget.NewEntry()
	octEntry.SetPlaceHolder("377")

	isUpdating := false

	setAll := func(val int, skipField string) {
		isUpdating = true
		defer func() { isUpdating = false }()

		if skipField != "dec" {
			decEntry.SetText(strconv.Itoa(val))
		}
		if skipField != "bin" {
			binEntry.SetText(calculators.DecToBin(val))
		}
		if skipField != "hex" {
			hexEntry.SetText(calculators.DecToHex(val))
		}
		if skipField != "oct" {
			octEntry.SetText(strconv.FormatInt(int64(val), 8))
		}
	}

	decEntry.OnChanged = func(s string) {
		if isUpdating {
			return
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		val, err := strconv.Atoi(s)
		if err == nil {
			setAll(val, "dec")
		}
	}

	binEntry.OnChanged = func(s string) {
		if isUpdating {
			return
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		val, err := calculators.BinToDec(s)
		if err == nil {
			setAll(val, "bin")
		}
	}

	hexEntry.OnChanged = func(s string) {
		if isUpdating {
			return
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		val, err := calculators.HexToDec(s)
		if err == nil {
			setAll(val, "hex")
		}
	}

	octEntry.OnChanged = func(s string) {
		if isUpdating {
			return
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}
		val, err := strconv.ParseInt(s, 8, 64)
		if err == nil {
			setAll(int(val), "oct")
		}
	}

	makeNumChip := func(v int) *widget.Button {
		btn := widget.NewButton(strconv.Itoa(v), func() {
			setAll(v, "")
		})
		btn.Importance = widget.LowImportance
		return btn
	}

	numChips := container.NewHBox(
		canvas.NewText("Angka Populer:", constants.ColorTextSecondary),
		makeNumChip(16),
		makeNumChip(255),
		makeNumChip(1024),
		makeNumChip(65535),
	)

	numContent := container.NewVBox(
		widget.NewLabel("Ketik pada salah satu kolom di bawah untuk mengonversi secara live ke seluruh basis:"),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Desimal (Basis 10):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), decEntry),
			container.NewVBox(widget.NewLabelWithStyle("Biner (Basis 2):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), binEntry),
			container.NewVBox(widget.NewLabelWithStyle("Heksadesimal (Basis 16):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), hexEntry),
			container.NewVBox(widget.NewLabelWithStyle("Oktal (Basis 8):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), octEntry),
		),
		numChips,
	)
	numCard := components.NewPlainCard(numContent)

	// Data Size Converter
	units := []string{"bit", "Byte", "KB", "MB", "GB", "TB"}
	dataValEntry := widget.NewEntry()
	dataValEntry.SetPlaceHolder("1024")
	dataValEntry.SetText("1024")

	fromSelect := widget.NewSelect(units, nil)
	fromSelect.SetSelected("MB")
	toSelect := widget.NewSelect(units, nil)
	toSelect.SetSelected("GB")

	resultStat := components.NewStatCard("HASIL KONVERSI UKURAN", "1.0000 GB", constants.ColorSuccess)

	calcData := func() {
		valStr := strings.TrimSpace(dataValEntry.Text)
		if valStr == "" {
			resultStat.SetValue("-")
			return
		}
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			resultStat.SetValue("Input tidak valid")
			resultStat.SetColor(constants.ColorDanger)
			return
		}
		resultStat.SetColor(constants.ColorSuccess)
		res := calculators.ConvertDataSize(val, fromSelect.Selected, toSelect.Selected)
		resultStat.SetValue(fmt.Sprintf("%.4f %s", res, toSelect.Selected))
		resultStat.SetSubtext(fmt.Sprintf("%s -> %s", fromSelect.Selected, toSelect.Selected))
	}

	dataValEntry.OnChanged = func(string) { calcData() }
	fromSelect.OnChanged = func(string) { calcData() }
	toSelect.OnChanged = func(string) { calcData() }

	makeDataChip := func(label, val, from, to string) *widget.Button {
		btn := widget.NewButton(label, func() {
			dataValEntry.SetText(val)
			fromSelect.SetSelected(from)
			toSelect.SetSelected(to)
		})
		btn.Importance = widget.LowImportance
		return btn
	}

	dataChips := container.NewHBox(
		canvas.NewText("Preset Ukuran:", constants.ColorTextSecondary),
		makeDataChip("1024 MB -> GB", "1024", "MB", "GB"),
		makeDataChip("4096 MB -> GB", "4096", "MB", "GB"),
		makeDataChip("1 TB -> GB", "1", "TB", "GB"),
		makeDataChip("8 bit -> Byte", "8", "bit", "Byte"),
	)

	dataContent := container.NewVBox(
		widget.NewLabel("Konversi kapasitas penyimpanan file dan memori secara presisi:"),
		container.NewGridWithColumns(3,
			container.NewVBox(widget.NewLabelWithStyle("Nilai:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), dataValEntry),
			container.NewVBox(widget.NewLabelWithStyle("Dari:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), fromSelect),
			container.NewVBox(widget.NewLabelWithStyle("Ke:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), toSelect),
		),
		dataChips,
		widget.NewSeparator(),
		resultStat.Widget,
	)
	dataCard := components.NewPlainCard(dataContent)

	// Set initial values
	setAll(255, "")

	return container.NewVScroll(container.NewVBox(
		numCard,
		dataCard,
	))
}

// ----------------------------------------------------------------------------
// 3. Tab Hash & Password Generator + Hashed Password Vault
// ----------------------------------------------------------------------------
func getAlgorithmBadge(algo string) fyne.CanvasObject {
	switch strings.ToUpper(strings.TrimSpace(algo)) {
	case "BCRYPT":
		return components.BadgeIndigo("BCRYPT")
	case "SHA-256", "SHA256":
		return components.BadgeCyan("SHA-256")
	case "SHA-512", "SHA512":
		return components.BadgeYellow("SHA-512")
	case "MD5":
		return components.BadgeMuted("MD5")
	case "SHA-1", "SHA1":
		return components.BadgeWarning("SHA-1")
	default:
		return components.BadgeCyan(algo)
	}
}

func (p *CalculatorPage) showSaveHashedPasswordDialog(defaultPlain, defaultAlgo, defaultHash string, onSaved func()) {
	titleEntry := components.NewScrollableEntry()
	titleEntry.SetPlaceHolder("Contoh: Admin Server, API Secret, Akun DB")

	if defaultAlgo == "" {
		defaultAlgo = "bcrypt"
	}

	plainEntry := components.NewScrollableEntry()
	plainEntry.SetPlaceHolder("Teks sandi plaintext (opsional jika mengisi hash langsung)")
	plainEntry.SetText(defaultPlain)

	savePlainCheck := widget.NewCheck("Simpan plaintext sandi juga (opsional, untuk catatan referensi)", nil)
	if defaultPlain != "" {
		savePlainCheck.SetChecked(true)
	}

	saltEntry := components.NewScrollableEntry()
	saltEntry.SetPlaceHolder("Salt opsional (khusus algoritma SHA-256, SHA-512, MD5, SHA-1)")

	hashEntry := components.NewScrollableEntry()
	hashEntry.SetPlaceHolder("Nilai hash akan dihitung otomatis atau masukkan hash manual...")
	hashEntry.TextStyle = fyne.TextStyle{Monospace: true}
	if defaultHash != "" {
		hashEntry.SetText(defaultHash)
	}

	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetPlaceHolder("Catatan keterangan penggunaan, host, atau lingkungan...")
	notesEntry.SetMinRowsVisible(2)

	var algoSelect *widget.Select
	recompute := func() {
		text := plainEntry.Text
		if text == "" {
			return
		}
		algo := "bcrypt"
		if algoSelect != nil && algoSelect.Selected != "" {
			algo = algoSelect.Selected
		}
		salt := saltEntry.Text

		switch algo {
		case "bcrypt":
			h, err := calculators.GenerateBcrypt(text, 10)
			if err == nil {
				hashEntry.SetText(h)
			}
		case "SHA-256":
			hashEntry.SetText(calculators.GenerateSHA256(text + salt))
		case "SHA-512":
			hashEntry.SetText(calculators.GenerateSHA512(text + salt))
		case "MD5":
			hashEntry.SetText(calculators.GenerateMD5(text + salt))
		case "SHA-1":
			hashEntry.SetText(calculators.GenerateSHA1(text + salt))
		}
	}

	algoSelect = widget.NewSelect([]string{"bcrypt", "SHA-256", "SHA-512", "MD5", "SHA-1"}, func(s string) {
		recompute()
	})
	algoSelect.SetSelected(defaultAlgo)

	plainEntry.OnChanged = func(s string) {
		recompute()
	}
	saltEntry.OnChanged = func(s string) {
		recompute()
	}

	if defaultPlain != "" && defaultHash == "" {
		recompute()
	}

	form := container.NewVBox(
		widget.NewLabelWithStyle("Nama / Judul Kredensial:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		titleEntry,
		container.NewGridWithColumns(2,
			container.NewVBox(
				widget.NewLabelWithStyle("Algoritma Hash:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				algoSelect,
			),
			container.NewVBox(
				widget.NewLabelWithStyle("Salt (Opsional):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				saltEntry,
			),
		),
		widget.NewLabelWithStyle("Sandi Plaintext (Otomatis Dihash):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		plainEntry,
		savePlainCheck,
		widget.NewLabelWithStyle("Nilai Hash yang Dihasilkan:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		hashEntry,
		widget.NewLabelWithStyle("Catatan Tambahan:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		notesEntry,
	)

	components.ShowBrutalistFormDialog(
		p.window,
		"VAULT",
		constants.ColorTechIndigo,
		"Simpan Sandi Ter-Hash ke Vault",
		"Simpan kredensial ber-hash secara persisten ke database lokal aman.",
		form,
		"Simpan ke Vault",
		func() {
			title := strings.TrimSpace(titleEntry.Text)
			if title == "" {
				dialog.ShowError(fmt.Errorf("nama / judul kredensial tidak boleh kosong"), p.window)
				return
			}
			hashVal := strings.TrimSpace(hashEntry.Text)
			if hashVal == "" {
				dialog.ShowError(fmt.Errorf("nilai hash tidak boleh kosong"), p.window)
				return
			}

			hp := &database.HashedPassword{
				Title:     title,
				Algorithm: algoSelect.Selected,
				HashValue: hashVal,
				Salt:      strings.TrimSpace(saltEntry.Text),
				Notes:     strings.TrimSpace(notesEntry.Text),
			}
			if savePlainCheck.Checked {
				hp.PlainPassword = plainEntry.Text
			}

			_, err := database.CreateHashedPassword(hp)
			if err != nil {
				dialog.ShowError(fmt.Errorf("gagal menyimpan ke vault: %w", err), p.window)
				return
			}

			dialog.ShowInformation("Vault", "Sandi ber-hash berhasil disimpan ke Vault!", p.window)
			if onSaved != nil {
				onSaved()
			}
		},
	)
}

func (p *CalculatorPage) showVerifyPasswordDialog(hp *database.HashedPassword) {
	testEntry := components.NewScrollableEntry()
	testEntry.SetPlaceHolder("Ketik teks sandi plaintext yang ingin diverifikasi...")

	resultBadgeContainer := container.NewHBox()
	resultText := canvas.NewText("Ketik sandi plaintext di atas untuk menguji validitas hash.", constants.ColorTextMuted)
	resultText.TextSize = constants.FontSizeBody

	verifyBox := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Target Kredensial:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(hp.Title),
			getAlgorithmBadge(hp.Algorithm),
		),
		widget.NewLabelWithStyle("Nilai Hash Tersimpan:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(hp.HashValue),
	)

	if hp.Salt != "" {
		verifyBox.Add(container.NewHBox(
			widget.NewLabelWithStyle("Salt Tersimpan:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel(hp.Salt),
		))
	}

	testEntry.OnChanged = func(s string) {
		resultBadgeContainer.Objects = nil
		if s == "" {
			resultText.Text = "Ketik sandi plaintext di atas untuk menguji validitas hash."
			resultText.Color = constants.ColorTextMuted
			resultText.Refresh()
			return
		}

		matched := calculators.VerifyHash(hp.Algorithm, hp.HashValue, s, hp.Salt)
		if matched {
			resultBadgeContainer.Add(components.BadgeSuccess("COCOK / VALID"))
			resultText.Text = "Sandi plaintext yang Anda masukkan COCOK dengan hash tersimpan!"
			resultText.Color = constants.ColorSuccess
		} else {
			resultBadgeContainer.Add(components.BadgeDanger("TIDAK COCOK"))
			resultText.Text = "Sandi plaintext TIDAK SESUAI dengan nilai hash ini."
			resultText.Color = constants.ColorDanger
		}
		resultBadgeContainer.Refresh()
		resultText.Refresh()
	}

	verifyContent := container.NewVBox(
		verifyBox,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Uji Sandi Plaintext:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		testEntry,
		container.NewVBox(
			resultBadgeContainer,
			resultText,
		),
	)

	components.ShowStyledInformationDialog(
		p.window,
		"VERIFIKASI HASH",
		constants.ColorAccentCobalt,
		"Verifikasi Sandi Plaintext",
		"Uji apakah suatu sandi plaintext menghasilkan hash yang identik.",
		verifyContent,
		"Selesai",
		nil,
	)
}

func (p *CalculatorPage) exportVaultJSON(list []database.HashedPassword) {
	if len(list) == 0 {
		dialog.ShowInformation("Ekspor Vault", "Tidak ada data sandi ber-hash untuk diekspor.", p.window)
		return
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		dialog.ShowError(fmt.Errorf("gagal memformat JSON: %w", err), p.window)
		return
	}

	home, _ := os.UserHomeDir()
	outDir := filepath.Join(home, ".it-toolbox")
	_ = os.MkdirAll(outDir, 0755)
	fileName := fmt.Sprintf("hashed_passwords_export_%d.json", time.Now().Unix())
	filePath := filepath.Join(outDir, fileName)

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		dialog.ShowError(fmt.Errorf("gagal menulis file JSON: %w", err), p.window)
		return
	}

	p.window.Clipboard().SetContent(string(data))
	msg := fmt.Sprintf("Berhasil mengekspor %d kredensial hash ke:\n%s\n\nData telah disalin ke clipboard!", len(list), filePath)
	dialog.ShowInformation("Ekspor JSON Berhasil", msg, p.window)
}

func (p *CalculatorPage) exportVaultCSV(list []database.HashedPassword) {
	if len(list) == 0 {
		dialog.ShowInformation("Ekspor Vault", "Tidak ada data sandi ber-hash untuk diekspor.", p.window)
		return
	}

	var sb strings.Builder
	w := csv.NewWriter(&sb)
	_ = w.Write([]string{"ID", "Title", "Algorithm", "HashValue", "PlainPassword", "Salt", "Notes", "CreatedAt"})

	for _, hp := range list {
		_ = w.Write([]string{
			strconv.FormatInt(hp.ID, 10),
			hp.Title,
			hp.Algorithm,
			hp.HashValue,
			hp.PlainPassword,
			hp.Salt,
			hp.Notes,
			hp.CreatedAt,
		})
	}
	w.Flush()

	csvData := sb.String()
	home, _ := os.UserHomeDir()
	outDir := filepath.Join(home, ".it-toolbox")
	_ = os.MkdirAll(outDir, 0755)
	fileName := fmt.Sprintf("hashed_passwords_export_%d.csv", time.Now().Unix())
	filePath := filepath.Join(outDir, fileName)

	err := os.WriteFile(filePath, []byte(csvData), 0644)
	if err != nil {
		dialog.ShowError(fmt.Errorf("gagal menulis file CSV: %w", err), p.window)
		return
	}

	p.window.Clipboard().SetContent(csvData)
	msg := fmt.Sprintf("Berhasil mengekspor %d kredensial hash ke:\n%s\n\nData telah disalin ke clipboard!", len(list), filePath)
	dialog.ShowInformation("Ekspor CSV Berhasil", msg, p.window)
}

func (p *CalculatorPage) buildHashGenTab() fyne.CanvasObject {
	hashInput := components.NewScrollableEntry()
	hashInput.SetPlaceHolder("Ketik teks untuk dihitung nilai hash-nya...")
	hashInput.SetText("Hello IT Toolbox 2026")

	md5Out := components.NewScrollableEntry()
	md5Out.TextStyle = fyne.TextStyle{Monospace: true}

	sha1Out := components.NewScrollableEntry()
	sha1Out.TextStyle = fyne.TextStyle{Monospace: true}

	sha256Out := components.NewScrollableEntry()
	sha256Out.TextStyle = fyne.TextStyle{Monospace: true}

	sha512Out := components.NewScrollableEntry()
	sha512Out.TextStyle = fyne.TextStyle{Monospace: true}

	updateHashes := func(s string) {
		if s == "" {
			md5Out.SetText("")
			sha1Out.SetText("")
			sha256Out.SetText("")
			sha512Out.SetText("")
			return
		}
		md5Out.SetText(calculators.GenerateMD5(s))
		sha1Out.SetText(calculators.GenerateSHA1(s))
		sha256Out.SetText(calculators.GenerateSHA256(s))
		sha512Out.SetText(calculators.GenerateSHA512(s))
	}

	hashInput.OnChanged = updateHashes

	var refreshVault func()

	makeHashRow := func(algoName string, entry *widget.Entry) fyne.CanvasObject {
		btnSave := widget.NewButtonWithIcon("Simpan", theme.DocumentSaveIcon(), func() {
			if entry.Text == "" {
				dialog.ShowError(fmt.Errorf("nilai hash belum tersedia untuk disimpan"), p.window)
				return
			}
			p.showSaveHashedPasswordDialog(hashInput.Text, algoName, entry.Text, refreshVault)
		})
		btnCopy := widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() {
			p.copyToClip(entry.Text)
		})
		return container.NewBorder(nil, nil, nil, container.NewHBox(btnSave, btnCopy), entry)
	}

	hashContent := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Teks Input:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			components.BadgeIndigo("MULTI-DIGEST"),
		),
		hashInput,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("MD5 (128-bit Digest):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		makeHashRow("MD5", md5Out),
		widget.NewLabelWithStyle("SHA-1 (160-bit Digest):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		makeHashRow("SHA-1", sha1Out),
		widget.NewLabelWithStyle("SHA-256 (256-bit Secure Digest):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		makeHashRow("SHA-256", sha256Out),
		widget.NewLabelWithStyle("SHA-512 (512-bit High Security Digest):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		makeHashRow("SHA-512", sha512Out),
	)
	hashCard := components.NewPlainCardWithAccent(hashContent, constants.ColorTechIndigo)

	// Password Generator
	pwdLenEntry := components.NewScrollableEntry()
	pwdLenEntry.SetText("16")
	pwdSymCheck := widget.NewCheck("Gunakan Simbol (!@#$%^&*)", nil)
	pwdSymCheck.SetChecked(true)
	pwdResult := components.NewScrollableEntry()
	pwdResult.TextStyle = fyne.TextStyle{Monospace: true}

	genPwd := func() {
		length, err := strconv.Atoi(pwdLenEntry.Text)
		if err != nil || length <= 0 {
			dialog.ShowError(fmt.Errorf("panjang password harus angka positif"), p.window)
			return
		}
		pwd := calculators.GeneratePassword(length, pwdSymCheck.Checked)
		pwdResult.SetText(pwd)
	}

	genPwdBtn := widget.NewButtonWithIcon("Generate Password Baru", theme.ViewRefreshIcon(), genPwd)
	genPwdBtn.Importance = widget.HighImportance

	makeLenChip := func(l int) *widget.Button {
		btn := widget.NewButton(fmt.Sprintf("%d", l), func() {
			pwdLenEntry.SetText(strconv.Itoa(l))
			genPwd()
		})
		btn.Importance = widget.LowImportance
		return btn
	}

	pwdChips := container.NewHBox(
		canvas.NewText("Panjang Standar:", constants.ColorTextSecondary),
		makeLenChip(8),
		makeLenChip(12),
		makeLenChip(16),
		makeLenChip(24),
		makeLenChip(32),
	)

	btnSavePwdToVault := widget.NewButtonWithIcon("Simpan ke Vault", theme.DocumentSaveIcon(), func() {
		if pwdResult.Text == "" {
			dialog.ShowError(fmt.Errorf("generate password terlebih dahulu"), p.window)
			return
		}
		p.showSaveHashedPasswordDialog(pwdResult.Text, "bcrypt", "", refreshVault)
	})
	btnSavePwdToVault.Importance = widget.MediumImportance

	pwdContent := container.NewVBox(
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Panjang Karakter:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), pwdLenEntry),
			container.NewVBox(widget.NewLabelWithStyle("Karakter Khusus:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), pwdSymCheck),
		),
		pwdChips,
		genPwdBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Hasil Password Terenkripsi (High Entropy):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, container.NewHBox(
			btnSavePwdToVault,
			widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() { p.copyToClip(pwdResult.Text) }),
		), pwdResult),
	)
	pwdCard := components.NewPlainCardWithAccent(pwdContent, constants.ColorTechIndigo)

	// Vault Card
	vaultItemsBox := container.NewVBox()
	countBadgeContainer := container.NewHBox(components.BadgeCyan("0 ITEM"))

	var currentQuery string
	var currentAlgo string = "Semua"
	var currentList []database.HashedPassword

	refreshVault = func() {
		list, err := database.SearchHashedPasswords(currentQuery, currentAlgo)
		if err != nil {
			dialog.ShowError(fmt.Errorf("gagal memuat data vault: %w", err), p.window)
			return
		}
		currentList = list

		countBadgeContainer.Objects = []fyne.CanvasObject{
			components.BadgeCyan(fmt.Sprintf("%d ITEM", len(list))),
		}
		countBadgeContainer.Refresh()

		vaultItemsBox.Objects = nil

		if len(list) == 0 {
			emptyCard := container.NewVBox(
				container.NewCenter(widget.NewIcon(theme.FolderOpenIcon())),
				container.NewCenter(widget.NewLabelWithStyle("Belum ada sandi ter-hash tersimpan.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
				container.NewCenter(canvas.NewText("Gunakan tombol 'Simpan' pada generator di atas atau 'Tambah Hash Manual' di bawah.", constants.ColorTextMuted)),
			)
			vaultItemsBox.Add(container.NewPadded(emptyCard))
			vaultItemsBox.Refresh()
			return
		}

		for _, item := range list {
			itemCopy := item // Capture for closure

			titleLabel := widget.NewLabelWithStyle(itemCopy.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			algoBadge := getAlgorithmBadge(itemCopy.Algorithm)
			dateText := canvas.NewText(itemCopy.CreatedAt, constants.ColorTextMuted)
			dateText.TextSize = constants.FontSizeSmall

			headerRow := container.NewBorder(
				nil, nil,
				container.NewHBox(titleLabel, algoBadge),
				dateText,
			)

			hashEntry := components.NewScrollableEntry()
			hashEntry.TextStyle = fyne.TextStyle{Monospace: true}
			hashEntry.SetText(itemCopy.HashValue)

			copyHashBtn := widget.NewButtonWithIcon("Salin Hash", theme.ContentCopyIcon(), func() {
				p.copyToClip(itemCopy.HashValue)
			})

			hashRow := container.NewBorder(
				nil, nil,
				canvas.NewText("Hash: ", constants.ColorTextSecondary),
				copyHashBtn,
				hashEntry,
			)

			cardRows := []fyne.CanvasObject{
				headerRow,
				hashRow,
			}

			if itemCopy.Salt != "" {
				saltRow := container.NewHBox(
					canvas.NewText("Salt: ", constants.ColorTextSecondary),
					canvas.NewText(itemCopy.Salt, constants.ColorTextPrimary),
				)
				cardRows = append(cardRows, saltRow)
			}

			if itemCopy.PlainPassword != "" {
				plainEntry := components.NewScrollableEntry()
				plainEntry.TextStyle = fyne.TextStyle{Monospace: true}
				plainEntry.SetText(strings.Repeat("*", len(itemCopy.PlainPassword)))
				isRevealed := false

				var btnToggle *widget.Button
				btnToggle = widget.NewButtonWithIcon("Lihat", theme.VisibilityIcon(), func() {
					if isRevealed {
						plainEntry.SetText(strings.Repeat("*", len(itemCopy.PlainPassword)))
						btnToggle.SetText("Lihat")
						isRevealed = false
					} else {
						plainEntry.SetText(itemCopy.PlainPassword)
						btnToggle.SetText("Sembunyi")
						isRevealed = true
					}
				})

				btnCopyPlain := widget.NewButtonWithIcon("Salin Sandi", theme.ContentCopyIcon(), func() {
					p.copyToClip(itemCopy.PlainPassword)
				})

				plainRow := container.NewBorder(
					nil, nil,
					canvas.NewText("Plain: ", constants.ColorTextSecondary),
					container.NewHBox(btnToggle, btnCopyPlain),
					plainEntry,
				)
				cardRows = append(cardRows, plainRow)
			}

			if itemCopy.Notes != "" {
				notesLabel := canvas.NewText("Catatan: "+itemCopy.Notes, constants.ColorTextMuted)
				notesLabel.TextSize = constants.FontSizeSmall
				cardRows = append(cardRows, notesLabel)
			}

			btnVerify := widget.NewButtonWithIcon("Verifikasi Sandi", theme.ConfirmIcon(), func() {
				p.showVerifyPasswordDialog(&itemCopy)
			})
			btnVerify.Importance = widget.MediumImportance

			btnDelete := widget.NewButtonWithIcon("Hapus", theme.DeleteIcon(), func() {
				components.ShowStyledConfirmDialog(
					p.window,
					"HAPUS",
					constants.ColorDanger,
					"Hapus Sandi Ter-Hash",
					"Apakah Anda yakin ingin menghapus data hash ini dari Vault?",
					widget.NewLabel(fmt.Sprintf("Judul: %s\nAlgoritma: %s", itemCopy.Title, itemCopy.Algorithm)),
					"Batal",
					"Hapus",
					func() {
						_ = database.DeleteHashedPassword(itemCopy.ID)
						refreshVault()
					},
				)
			})
			btnDelete.Importance = widget.LowImportance

			cardRows = append(cardRows, container.NewBorder(nil, nil, nil, container.NewHBox(btnVerify, btnDelete)))

			itemBox := container.NewVBox(cardRows...)
			itemCard := components.NewPlainCardWithAccent(itemBox, constants.ColorAccentCyan)
			vaultItemsBox.Add(itemCard)
		}
		vaultItemsBox.Refresh()
	}

	searchBar := components.NewSearchBar("Cari nama kredensial, catatan, atau sandi...", func(query string) {
		currentQuery = query
		refreshVault()
	})

	algoSelect := widget.NewSelect([]string{"Semua", "bcrypt", "SHA-256", "SHA-512", "MD5", "SHA-1"}, func(s string) {
		currentAlgo = s
		refreshVault()
	})
	algoSelect.SetSelected("Semua")

	btnAddManual := widget.NewButtonWithIcon("Tambah Hash", theme.ContentAddIcon(), func() {
		p.showSaveHashedPasswordDialog("", "bcrypt", "", refreshVault)
	})
	btnAddManual.Importance = widget.HighImportance

	btnExportJSON := widget.NewButtonWithIcon("Ekspor JSON", theme.DownloadIcon(), func() {
		p.exportVaultJSON(currentList)
	})

	btnExportCSV := widget.NewButtonWithIcon("Ekspor CSV", theme.DocumentSaveIcon(), func() {
		p.exportVaultCSV(currentList)
	})

	vaultHeader := container.NewVBox(
		container.NewBorder(
			nil, nil,
			container.NewHBox(
				widget.NewLabelWithStyle("Hashed Password Vault", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				countBadgeContainer,
			),
			btnAddManual,
		),
		canvas.NewText("Penyimpanan lokal kredensial ber-hash aman dengan verifikasi dan ekspor.", constants.ColorTextMuted),
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Pencarian Vault:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), searchBar.Container),
			container.NewVBox(widget.NewLabelWithStyle("Filter Algoritma:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), algoSelect),
		),
		container.NewHBox(
			canvas.NewText("Alat Ekspor Vault:", constants.ColorTextSecondary),
			btnExportJSON,
			btnExportCSV,
		),
		widget.NewSeparator(),
		vaultItemsBox,
	)
	vaultCard := components.NewPlainCardWithAccent(vaultHeader, constants.ColorTechIndigo)

	// Initialize
	updateHashes(hashInput.Text)
	genPwd()
	refreshVault()

	return container.NewVScroll(container.NewVBox(
		hashCard,
		pwdCard,
		vaultCard,
	))
}


// ----------------------------------------------------------------------------
// 4. Tab Text Formatter & Encoder
// ----------------------------------------------------------------------------
func (p *CalculatorPage) buildFormatterTab() fyne.CanvasObject {
	inputArea := widget.NewMultiLineEntry()
	inputArea.SetPlaceHolder("Tempelkan teks, JSON mentah, YAML, XML, atau string target regex...")
	inputArea.SetMinRowsVisible(6)
	inputArea.TextStyle = fyne.TextStyle{Monospace: true}
	inputArea.SetText("{\"app\":\"it-toolbox\",\"status\":\"active\",\"tags\":[\"go\",\"fyne\",\"sqlite\"]}")

	outputArea := widget.NewMultiLineEntry()
	outputArea.SetPlaceHolder("Hasil formatting akan ditampilkan di sini...")
	outputArea.SetMinRowsVisible(6)
	outputArea.TextStyle = fyne.TextStyle{Monospace: true}

	copyBtn := widget.NewButtonWithIcon(constants.BtnCopy, theme.ContentCopyIcon(), func() {
		p.copyToClip(outputArea.Text)
	})
	copyBtn.Importance = widget.HighImportance

	clearBtn := widget.NewButtonWithIcon(constants.BtnClear, theme.CancelIcon(), func() {
		inputArea.SetText("")
		outputArea.SetText("")
	})

	formatJSONBtn := widget.NewButtonWithIcon("Format JSON", theme.DocumentIcon(), func() {
		out, err := formatters.FormatJSON(inputArea.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		outputArea.SetText(out)
	})

	formatYAMLBtn := widget.NewButtonWithIcon("Format YAML", theme.DocumentIcon(), func() {
		out, err := formatters.FormatYAML(inputArea.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		outputArea.SetText(out)
	})

	formatXMLBtn := widget.NewButtonWithIcon("Format XML", theme.DocumentIcon(), func() {
		out, err := formatters.FormatXML(inputArea.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		outputArea.SetText(out)
	})

	encodeB64Btn := widget.NewButton("Encode Base64", func() {
		outputArea.SetText(formatters.EncodeBase64(inputArea.Text))
	})

	decodeB64Btn := widget.NewButton("Decode Base64", func() {
		out, err := formatters.DecodeBase64(inputArea.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		outputArea.SetText(out)
	})

	encodeURLBtn := widget.NewButton("Encode URL", func() {
		outputArea.SetText(formatters.EncodeURL(inputArea.Text))
	})

	decodeURLBtn := widget.NewButton("Decode URL", func() {
		out, err := formatters.DecodeURL(inputArea.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		outputArea.SetText(out)
	})

	// Regex Tester
	regexEntry := widget.NewEntry()
	regexEntry.SetPlaceHolder(`cth: [a-z]+@[a-z0-9.]+\.[a-z]{2,} atau \d+`)
	regexEntry.SetText(`[a-zA-Z0-9_-]+`)
	regexEntry.TextStyle = fyne.TextStyle{Monospace: true}

	testRegexBtn := widget.NewButtonWithIcon("Uji Regex", theme.SearchIcon(), func() {
		matches, err := formatters.TestRegex(regexEntry.Text, inputArea.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		if len(matches) == 0 {
			outputArea.SetText("Tidak ada pola yang cocok.")
			return
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Ditemukan %d match:\n\n", len(matches)))
		for i, m := range matches {
			sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, m))
		}
		outputArea.SetText(sb.String())
	})
	testRegexBtn.Importance = widget.HighImportance

	makeSampleChip := func(label, sample string) *widget.Button {
		btn := widget.NewButton(label, func() {
			inputArea.SetText(sample)
		})
		btn.Importance = widget.LowImportance
		return btn
	}

	presetChips := container.NewHBox(
		canvas.NewText("Contoh Data:", constants.ColorTextSecondary),
		makeSampleChip("JSON Sample", `{"title":"Subnetting","hosts":50,"vlan":10,"active":true}`),
		makeSampleChip("YAML Sample", "version: '3.8'\nservices:\n  web:\n    image: nginx:alpine\n    ports:\n      - 80:80"),
		makeSampleChip("XML Sample", "<network><device type=\"router\"><name>Core-R1</name><ip>192.168.1.1</ip></device></network>"),
	)

	btnGrid := container.NewGridWithColumns(4,
		formatJSONBtn, formatYAMLBtn, formatXMLBtn,
		encodeB64Btn, decodeB64Btn, encodeURLBtn, decodeURLBtn,
	)

	regexBox := container.NewBorder(nil, nil, widget.NewLabelWithStyle("Pola Regex:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), testRegexBtn, regexEntry)

	content := container.NewVBox(
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("Input Teks / Payload:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			clearBtn,
		),
		inputArea,
		presetChips,
		btnGrid,
		regexBox,
		widget.NewSeparator(),
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("Hasil Formatting / Encoding:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			copyBtn,
		),
		outputArea,
	)

	return container.NewVScroll(components.NewPlainCardWithAccent(content, constants.ColorTechIndigo))
}
