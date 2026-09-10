package pages

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/core/calculators"
	"github.com/yudz/it-toolbox/core/formatters"
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

	hintLabel := canvas.NewText("-", color.Black)
	hintLabel.TextSize = constants.FontSizeSmall
	hintLabel.TextStyle = fyne.TextStyle{Bold: true}

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
		statRange.SetValue(fmt.Sprintf("%s ➔ %s", rec.FirstHost, rec.LastHost))
		statRange.SetSubtext(fmt.Sprintf("Tersedia %d IP Usable", rec.AllocatedHosts))

		statCapacity.SetColor(constants.ColorAccentYellow)
		statCapacity.SetValue(fmt.Sprintf("%d Host (%.1f%%)", rec.AllocatedHosts, rec.Efficiency))
		statCapacity.SetSubtext(fmt.Sprintf("%d Terpakai | %d Sisa", rec.NeededHosts, rec.WastedHosts))

		statBroadcast.SetColor(constants.ColorWarning)
		statBroadcast.SetValue(rec.Broadcast)
		statBroadcast.SetSubtext(fmt.Sprintf("Wildcard: %s", rec.WildcardMask))

		hintLabel.Text = fmt.Sprintf("💡 Saran Alokasi: %s", rec.ClassHint)
		hintLabel.Color = color.Black
		hintLabel.Refresh()
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
		canvas.NewText("Preset Cepat:", constants.ColorTextSecondary),
		makeHostChip(10),
		makeHostChip(30),
		makeHostChip(50),
		makeHostChip(100),
		makeHostChip(250),
		makeHostChip(500),
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
		container.NewPadded(hintLabel),
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
		resultStat.SetSubtext(fmt.Sprintf("%s ➔ %s", fromSelect.Selected, toSelect.Selected))
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
		makeDataChip("1024 MB ➔ GB", "1024", "MB", "GB"),
		makeDataChip("4096 MB ➔ GB", "4096", "MB", "GB"),
		makeDataChip("1 TB ➔ GB", "1", "TB", "GB"),
		makeDataChip("8 bit ➔ Byte", "8", "bit", "Byte"),
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
// 3. Tab Hash & Password Generator
// ----------------------------------------------------------------------------
func (p *CalculatorPage) buildHashGenTab() fyne.CanvasObject {
	hashInput := widget.NewEntry()
	hashInput.SetPlaceHolder("Ketik teks untuk dihitung nilai hash-nya...")
	hashInput.SetText("Hello IT Toolbox 2026")

	md5Out := widget.NewEntry()
	md5Out.TextStyle = fyne.TextStyle{Monospace: true}
	sha256Out := widget.NewEntry()
	sha256Out.TextStyle = fyne.TextStyle{Monospace: true}

	updateHashes := func(s string) {
		if s == "" {
			md5Out.SetText("")
			sha256Out.SetText("")
			return
		}
		md5Out.SetText(calculators.GenerateMD5(s))
		sha256Out.SetText(calculators.GenerateSHA256(s))
	}

	hashInput.OnChanged = updateHashes

	hashContent := container.NewVBox(
		widget.NewLabelWithStyle("Teks Input:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		hashInput,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("MD5 (128-bit Digest):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() { p.copyToClip(md5Out.Text) }), md5Out),
		widget.NewLabelWithStyle("SHA-256 (256-bit Secure Digest):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() { p.copyToClip(sha256Out.Text) }), sha256Out),
	)
	hashCard := components.NewPlainCardWithAccent(hashContent, constants.ColorTechIndigo)

	// Password Generator
	pwdLenEntry := widget.NewEntry()
	pwdLenEntry.SetText("16")
	pwdSymCheck := widget.NewCheck("Gunakan Simbol (!@#$%^&*)", nil)
	pwdSymCheck.SetChecked(true)
	pwdResult := widget.NewEntry()
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
		btn := widget.NewButton(fmt.Sprintf("%d Karakter", l), func() {
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

	pwdContent := container.NewVBox(
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Panjang Karakter:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), pwdLenEntry),
			container.NewVBox(widget.NewLabelWithStyle("Karakter Khusus:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), pwdSymCheck),
		),
		pwdChips,
		genPwdBtn,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Hasil Password Terenkripsi (High Entropy):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() { p.copyToClip(pwdResult.Text) }), pwdResult),
	)
	pwdCard := components.NewPlainCardWithAccent(pwdContent, constants.ColorTechIndigo)

	// Initialize
	updateHashes(hashInput.Text)
	genPwd()

	return container.NewVScroll(container.NewVBox(
		hashCard,
		pwdCard,
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
