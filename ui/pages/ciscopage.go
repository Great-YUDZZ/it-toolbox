package pages

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/core/cisco"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

type CiscoPage struct {
	window fyne.Window

	searchQuery      string
	selectedDevice   cisco.DeviceType
	selectedMode     cisco.CLIMode
	selectedCategory cisco.Category
}

func NewCiscoPage(win fyne.Window) *CiscoPage {
	return &CiscoPage{
		window:           win,
		selectedDevice:   cisco.DeviceAll,
		selectedMode:     cisco.ModeAll,
		selectedCategory: cisco.CategoryAll,
	}
}

func (p *CiscoPage) copyToClip(txt string) {
	if p.window != nil {
		p.window.Clipboard().SetContent(txt)
		dialog.ShowInformation("Clipboard", constants.StatusCopied, p.window)
	}
}

func (p *CiscoPage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavCisco,
		"Katalog lengkap perintah CLI Cisco IOS untuk Router, Switch, dan PC. Dikelompokkan berdasarkan perangkat, mode akses, dan tujuan konfigurasi serta panduan verifikasi topologi.",
		components.BadgeCyan("CISCO PACKET TRACER"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabCiscoAll, theme.ListIcon(), p.buildCategoryView(cisco.CategoryAll)),
		container.NewTabItemWithIcon(constants.TabCiscoBasic, theme.SettingsIcon(), p.buildCategoryView(cisco.CategoryBasic)),
		container.NewTabItemWithIcon(constants.TabCiscoIP, theme.RadioButtonIcon(), p.buildCategoryView(cisco.CategoryInterface)),
		container.NewTabItemWithIcon(constants.TabCiscoVLAN, theme.FolderIcon(), p.buildCategoryView(cisco.CategoryVLAN)),
		container.NewTabItemWithIcon(constants.TabCiscoRouting, theme.NavigateNextIcon(), p.buildCategoryView(cisco.CategoryRouting)),
		container.NewTabItemWithIcon(constants.TabCiscoServices, theme.StorageIcon(), p.buildCategoryView(cisco.CategoryServices)),
		container.NewTabItemWithIcon(constants.TabCiscoSecurity, theme.VisibilityIcon(), p.buildCategoryView(cisco.CategorySecurity)),
		container.NewTabItemWithIcon(constants.TabCiscoNAT, theme.ComputerIcon(), p.buildCategoryView(cisco.CategoryNAT)),
		container.NewTabItemWithIcon(constants.TabCiscoDiag, theme.SearchIcon(), p.buildCategoryView(cisco.CategoryShowDiag)),
		container.NewTabItemWithIcon(constants.TabCiscoVerify, theme.ConfirmIcon(), p.buildVerificationGuideView()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
}

func (p *CiscoPage) buildCategoryView(cat cisco.Category) fyne.CanvasObject {
	listContainer := container.NewVBox()

	query := p.searchQuery
	devFilter := p.selectedDevice
	modeFilter := p.selectedMode

	renderList := func() {
		listContainer.Objects = nil

		items := cisco.SearchCommands(cisco.FilterCriteria{
			Query:    query,
			Device:   devFilter,
			Mode:     modeFilter,
			Category: cat,
		})

		if len(items) == 0 {
			emptyCard := container.NewVBox(
				canvas.NewText("Tidak ada perintah Cisco yang cocok dengan kriteria filter.", constants.ColorTextMuted),
				canvas.NewText("Coba ubah kata kunci pencarian atau reset filter perangkat/mode ke 'Semua'.", constants.ColorTextSecondary),
			)
			listContainer.Add(components.NewPlainCard(emptyCard))
			listContainer.Refresh()
			return
		}

		for _, item := range items {
			cmd := item
			card := p.buildCommandCard(cmd)
			listContainer.Add(card)
		}
		listContainer.Refresh()
	}

	// Device selector
	deviceOptions := []string{
		string(cisco.DeviceAll),
		string(cisco.DeviceRouter),
		string(cisco.DeviceSwitchL2),
		string(cisco.DeviceSwitchL3),
		string(cisco.DevicePC),
	}
	devSelect := widget.NewSelect(deviceOptions, func(val string) {
		devFilter = cisco.DeviceType(val)
		renderList()
	})
	devSelect.SetSelected(string(cisco.DeviceAll))

	// Mode selector
	modeOptions := []string{
		string(cisco.ModeAll),
		string(cisco.ModeUserExec),
		string(cisco.ModePrivExec),
		string(cisco.ModeGlobalConfig),
		string(cisco.ModeInterface),
		string(cisco.ModeSubInterface),
		string(cisco.ModeVLAN),
		string(cisco.ModeLine),
		string(cisco.ModeRouterConfig),
		string(cisco.ModeDHCPConfig),
		string(cisco.ModePCTerminal),
	}
	modeSelect := widget.NewSelect(modeOptions, func(val string) {
		modeFilter = cisco.CLIMode(val)
		renderList()
	})
	modeSelect.SetSelected(string(cisco.ModeAll))

	// Search bar
	searchBar := components.NewSearchBar(constants.SearchCiscoPlaceholder, func(q string) {
		query = q
		renderList()
	})

	// Filter Bar with labels
	filterRow := container.NewHBox(
		components.BadgeCyan("PERANGKAT:"),
		devSelect,
		components.BadgeIndigo("MODE CLI:"),
		modeSelect,
	)

	headerBox := container.NewVBox(
		searchBar.Container,
		filterRow,
		widget.NewSeparator(),
	)

	renderList()

	scrollList := container.NewVScroll(container.NewPadded(listContainer))
	return container.NewBorder(headerBox, nil, nil, nil, scrollList)
}

func (p *CiscoPage) buildCommandCard(cmd cisco.CiscoCommand) fyne.CanvasObject {
	// 1. Accent and Badges
	var accentColor color.Color
	var devBadge fyne.CanvasObject

	switch cmd.Device {
	case cisco.DeviceRouter:
		accentColor = constants.ColorAccentCyan
		devBadge = components.BadgeCyan("ROUTER")
	case cisco.DeviceSwitchL2:
		accentColor = constants.ColorAccentYellow
		devBadge = components.BadgeYellow("SWITCH L2")
	case cisco.DeviceSwitchL3:
		accentColor = constants.ColorWarning
		devBadge = components.BadgeWarning("SWITCH L3")
	case cisco.DevicePC:
		accentColor = constants.ColorSuccess
		devBadge = components.BadgeSuccess("PC / END DEVICE")
	default:
		accentColor = constants.ColorTechIndigo
		devBadge = components.BadgeIndigo("ROUTER & SWITCH")
	}

	modeBadge := components.BadgeIndigo(string(cmd.Mode))
	catBadge := components.BadgeMuted(string(cmd.Category))

	// Title
	titleText := canvas.NewText(cmd.Title, constants.ColorTextPrimary)
	titleText.TextSize = constants.FontSizeH2
	titleText.TextStyle = fyne.TextStyle{Bold: true}

	// Copy Script Button
	copyBtn := widget.NewButtonWithIcon("Salin Script CLI", theme.ContentCopyIcon(), func() {
		p.copyToClip(cmd.Commands)
	})
	copyBtn.Importance = widget.HighImportance

	headerLeft := container.NewVBox(
		container.NewHBox(devBadge, modeBadge, catBadge),
		titleText,
	)
	topHeader := container.NewBorder(nil, nil, headerLeft, copyBtn)

	// Description
	descLabel := widget.NewLabel(cmd.Description)
	descLabel.Wrapping = fyne.TextWrapWord

	cardItems := []fyne.CanvasObject{
		topHeader,
		widget.NewSeparator(),
		descLabel,
	}

	// IP Example (if available)
	if cmd.IPExample != "" {
		ipBadge := components.BadgeWarning("SKEMA IP / TOPOLOGI")
		ipText := canvas.NewText(cmd.IPExample, constants.ColorTextPrimary)
		ipText.TextSize = constants.FontSizeSmall
		ipText.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

		ipBox := container.NewHBox(ipBadge, ipText)
		cardItems = append(cardItems, ipBox)
	}

	// CLI Script Panel
	codeEntry := widget.NewMultiLineEntry()
	codeEntry.SetText(cmd.Commands)
	codeEntry.TextStyle = fyne.TextStyle{Monospace: true}
	codeEntry.Wrapping = fyne.TextWrapWord
	codeEntry.Disable() // Read-only

	codeLabel := canvas.NewText("⌨ PERINTAH CLI CISCO IOS (SIAP SALIN):", constants.ColorTextPrimary)
	codeLabel.TextSize = constants.FontSizeLabel
	codeLabel.TextStyle = fyne.TextStyle{Bold: true}

	codeBox := container.NewVBox(
		codeLabel,
		codeEntry,
	)
	cardItems = append(cardItems, codeBox)

	// Verification Guide Box
	verifBadge := components.BadgeSuccess("CARA VERIFIKASI / TES BEKERJA")
	verifLabel := widget.NewLabel(cmd.Verification)
	verifLabel.Wrapping = fyne.TextWrapWord

	verifContent := container.NewVBox(
		container.NewHBox(verifBadge),
		verifLabel,
	)
	cardItems = append(cardItems, verifContent)

	// Troubleshooting Tips Box
	if cmd.TroubleshootingTips != "" {
		tipsBadge := components.BadgeDanger("TIPS & JEBAKAN UMUM")
		tipsLabel := widget.NewLabel(cmd.TroubleshootingTips)
		tipsLabel.Wrapping = fyne.TextWrapWord

		tipsContent := container.NewVBox(
			container.NewHBox(tipsBadge),
			tipsLabel,
		)
		cardItems = append(cardItems, tipsContent)
	}

	fullCardContent := container.NewVBox(cardItems...)
	return components.NewPlainCardWithAccent(fullCardContent, accentColor)
}

func (p *CiscoPage) buildVerificationGuideView() fyne.CanvasObject {
	contentList := container.NewVBox()

	// Section 1: Diagnosa Lampu Indikator Fisik
	sec1Title := canvas.NewText("1. DIAGNOSA LAMPU INDIKATOR LINK (LAYER 1 FISIK)", constants.ColorTextPrimary)
	sec1Title.TextSize = constants.FontSizeH2
	sec1Title.TextStyle = fyne.TextStyle{Bold: true}

	lampuHijau := container.NewHBox(
		components.BadgeSuccess("HIJAU (SOLID)"),
		widget.NewLabel("Link Layer 1 & 2 Normal (Port UP dan Protokol UP). Siap mengirim data."),
	)
	lampuOranye := container.NewHBox(
		components.BadgeWarning("ORANYE (BLINK/SOLID)"),
		widget.NewLabel("Spanning Tree Protocol (STP) sedang Listening/Learning. Tunggu 30-50 detik atau klik tombol 'Fast Forward Time' (Alt + D) 2x di Packet Tracer."),
	)
	lampuMerah := container.NewHBox(
		components.BadgeDanger("MERAH (SOLID)"),
		widget.NewLabel("Port Mati atau Kabel Salah! Cek: (a) Ketik 'no shutdown' di interface router, (b) Cek jenis kabel (Gunakan Straight-Through antar PC ke Switch, Crossover antar PC ke PC / Router ke PC)."),
	)

	card1 := components.NewPlainCardWithAccent(
		container.NewVBox(sec1Title, widget.NewSeparator(), lampuHijau, lampuOranye, lampuMerah),
		constants.ColorSuccess,
	)
	contentList.Add(card1)

	// Section 2: Prosedur Pengujian Bertingkat (Step-by-Step Testing)
	sec2Title := canvas.NewText("2. PROSEDUR 6 LANGKAH PENGUJIAN KONEKTIVITAS (PING WORKFLOW)", constants.ColorTextPrimary)
	sec2Title.TextSize = constants.FontSizeH2
	sec2Title.TextStyle = fyne.TextStyle{Bold: true}

	step1 := widget.NewLabel("Langkah 1: Cek IP PC di Command Prompt (Ketik: ipconfig /all). Pastikan IP, Subnet Mask, dan Default Gateway terisi benar.")
	step1.Wrapping = fyne.TextWrapWord

	step2 := widget.NewLabel("Langkah 2: Tes TCP/IP Stack Lokal (Ketik: ping 127.0.0.1). Jika reply, berarti kartu jaringan/NIC komputer berfungsi normal.")
	step2.Wrapping = fyne.TextWrapWord

	step3 := widget.NewLabel("Langkah 3: Ping Default Gateway (Ketik: ping <ip_gateway_router>, misal 192.168.1.1). Memastikan PC terhubung lancar ke router lokal.")
	step3.Wrapping = fyne.TextWrapWord

	step4 := widget.NewLabel("Langkah 4: Ping Interface Router Lawan / Next-Hop (Ketik: ping <ip_serial_lawan>, misal 10.10.10.2). Memastikan link WAN antar-router terhubung.")
	step4.Wrapping = fyne.TextWrapWord

	step5 := widget.NewLabel("Langkah 5: Ping End-to-End ke PC / Server Tujuan (Ketik: ping <ip_tujuan>, misal 192.168.20.10). Memastikan tabel routing sudah konvergen.")
	step5.Wrapping = fyne.TextWrapWord

	step6 := widget.NewLabel("Langkah 6: Traceroute Jalur (Ketik: tracert <ip_tujuan>). Melacak hop mana yang meneruskan paket dan di mana paket terhenti jika gagal.")
	step6.Wrapping = fyne.TextWrapWord

	card2 := components.NewPlainCardWithAccent(
		container.NewVBox(sec2Title, widget.NewSeparator(), step1, step2, step3, step4, step5, step6),
		constants.ColorAccentCyan,
	)
	contentList.Add(card2)

	// Section 3: Cara Menggunakan Mode Simulasi (Add Simple PDU)
	sec3Title := canvas.NewText("3. MELACAK PAKET MENGGUNAKAN SIMULASI PDU (AMPLOP SURAT)", constants.ColorTextPrimary)
	sec3Title.TextSize = constants.FontSizeH2
	sec3Title.TextStyle = fyne.TextStyle{Bold: true}

	pduDesc := widget.NewLabel("Jika ping gagal di mode Realtime, beralihlah ke Mode Simulasi (tekan Shift + S atau klik tab Simulation di kanan bawah).\n\n" +
		"1. Klik ikon Amplop Surat Tertutup (Add Simple PDU) atau tekan tombol keyboard 'P'.\n" +
		"2. Klik perangkat sumber (PC asal), lalu klik perangkat tujuan (PC/Server tujuan).\n" +
		"3. Klik tombol 'Capture / Forward' berkali-kali untuk melihat pergerakan amplop loncat dari PC -> Switch -> Router -> Tujuan.\n" +
		"4. Jika amplop memunculkan tanda Centang Hijau, maka koneksi berhasil.\n" +
		"5. Jika amplop memunculkan tanda Silang Merah di suatu perangkat, KLIK amplop tersebut! Buka tab 'In Layers' dan 'Out Layers' untuk membaca alasan mengapa paket didrop (misal: 'Device drops the packet because routing table has no route').")
	pduDesc.Wrapping = fyne.TextWrapWord

	card3 := components.NewPlainCardWithAccent(
		container.NewVBox(sec3Title, widget.NewSeparator(), pduDesc),
		constants.ColorWarning,
	)
	contentList.Add(card3)

	// Section 4: Tabel Arti Pesan Error Ping
	sec4Title := canvas.NewText("4. ARTI PESAN HASIL PING & SOLUSINYA", constants.ColorTextPrimary)
	sec4Title.TextSize = constants.FontSizeH2
	sec4Title.TextStyle = fyne.TextStyle{Bold: true}

	rto := container.NewHBox(
		components.BadgeDanger("Request timed out (RTO)"),
		widget.NewLabel("Paket terkirim tapi tidak ada balasan. Cek: (1) Default Gateway di PC tujuan, (2) Routing balik dari router lawan, (3) Jalur terblokir ACL."),
	)
	dhu := container.NewHBox(
		components.BadgeWarning("Destination Host Unreachable"),
		widget.NewLabel("Router tidak menemukan rute ke subnet tujuan. Cek tabel routing dengan 'show ip route' dan pastikan rute terdaftar."),
	)
	firstDrop := container.NewHBox(
		components.BadgeYellow("Ping 1 Gagal, lalu Reply ( . ! ! ! ! )"),
		widget.NewLabel("Ini NORMAL di Cisco Packet Tracer! Paket pertama dipakai untuk proses broadcast ARP (Address Resolution Protocol) mencari MAC address."),
	)

	card4 := components.NewPlainCardWithAccent(
		container.NewVBox(sec4Title, widget.NewSeparator(), rto, dhu, firstDrop),
		constants.ColorTechIndigo,
	)
	contentList.Add(card4)

	return container.NewVScroll(container.NewPadded(contentList))
}
