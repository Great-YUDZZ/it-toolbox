package pages

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/core/cisco"
	"github.com/yudz/it-toolbox/database"
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
		container.NewTabItemWithIcon(constants.TabCiscoSwitch, theme.FolderIcon(), p.buildCategoryView(cisco.CategoryVLAN)),
		container.NewTabItemWithIcon(constants.TabCiscoRouting, theme.NavigateNextIcon(), p.buildCategoryView(cisco.CategoryRouting)),
		container.NewTabItemWithIcon(constants.TabCiscoServices, theme.StorageIcon(), p.buildCategoryView(cisco.CategoryServices)),
		container.NewTabItemWithIcon(constants.TabCiscoVerify, theme.ConfirmIcon(), p.buildVerificationGuideView()),
		container.NewTabItemWithIcon(constants.TabCiscoTopologyNotes, theme.DocumentCreateIcon(), p.buildTopologyNotesView()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
}

func (p *CiscoPage) buildCategoryView(cat cisco.Category) fyne.CanvasObject {
	listContainer := container.NewVBox()

	query := p.searchQuery
	devFilter := p.selectedDevice
	modeFilter := p.selectedMode
	activeCat := cat

	renderList := func() {
		listContainer.Objects = nil

		items := cisco.SearchCommands(cisco.FilterCriteria{
			Query:    query,
			Device:   devFilter,
			Mode:     modeFilter,
			Category: activeCat,
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

	// Device selector with comfortable width
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

	devSpacer := canvas.NewRectangle(color.Transparent)
	devSpacer.SetMinSize(fyne.NewSize(140, 36))
	devSelectBox := container.NewStack(devSpacer, devSelect)

	// Mode selector with comfortable width
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

	modeSpacer := canvas.NewRectangle(color.Transparent)
	modeSpacer.SetMinSize(fyne.NewSize(160, 36))
	modeSelectBox := container.NewStack(modeSpacer, modeSelect)

	// Search bar
	searchBar := components.NewSearchBar(constants.SearchCiscoPlaceholder, func(q string) {
		query = q
		renderList()
	})

	// Filter Bar items with centered badges and proper spacing
	filterRow1 := container.NewHBox(
		container.NewCenter(components.BadgeCyan("PERANGKAT:")),
		devSelectBox,
		widget.NewSeparator(),
		container.NewCenter(components.BadgeIndigo("MODE CLI:")),
		modeSelectBox,
	)

	var filterBox fyne.CanvasObject = filterRow1

	// Add Category selector if in All Commands tab
	if cat == cisco.CategoryAll {
		catOptions := []string{
			string(cisco.CategoryAll),
			string(cisco.CategoryBasic),
			string(cisco.CategoryInterface),
			string(cisco.CategoryVLAN),
			string(cisco.CategoryRouting),
			string(cisco.CategoryServices),
			string(cisco.CategorySecurity),
			string(cisco.CategoryNAT),
			string(cisco.CategoryShowDiag),
		}
		catSelect := widget.NewSelect(catOptions, func(val string) {
			activeCat = cisco.Category(val)
			renderList()
		})
		catSelect.SetSelected(string(cisco.CategoryAll))

		catSpacer := canvas.NewRectangle(color.Transparent)
		catSpacer.SetMinSize(fyne.NewSize(200, 36))
		catSelectBox := container.NewStack(catSpacer, catSelect)

		filterRow2 := container.NewHBox(
			container.NewCenter(components.BadgeYellow("KATEGORI:")),
			catSelectBox,
		)
		filterBox = container.NewVBox(filterRow1, filterRow2)
	}

	headerBox := container.NewVBox(
		searchBar.Container,
		filterBox,
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
	titleLabel := widget.NewLabel(cmd.Title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Wrapping = fyne.TextWrapWord

	// State for dynamic parameters
	userParams := make(map[string]string)
	for _, param := range cmd.Parameters {
		userParams[param.Key] = param.DefaultValue
	}

	initialCommands := cmd.RenderCommands(userParams)
	initialIPExample := cmd.RenderIPExample(userParams)

	// CLI Script Entry (Monospace)
	codeEntry := widget.NewMultiLineEntry()
	codeEntry.SetText(initialCommands)
	codeEntry.TextStyle = fyne.TextStyle{Monospace: true}
	codeEntry.Wrapping = fyne.TextWrapWord

	lines := strings.Split(initialCommands, "\n")
	lineCount := len(lines)
	if lineCount < 4 {
		lineCount = 4
	}
	minHeight := float32(lineCount*21 + 28)
	if minHeight < 120 {
		minHeight = 120
	}
	if minHeight > 360 {
		minHeight = 360
	}

	codeSpacer := canvas.NewRectangle(color.Transparent)
	codeSpacer.SetMinSize(fyne.NewSize(0, minHeight))
	codeStack := container.NewStack(codeSpacer, codeEntry)

	codeLabel := canvas.NewText("⌨ PERINTAH CLI CISCO IOS (SIAP SALIN):", constants.ColorTextPrimary)
	codeLabel.TextSize = constants.FontSizeLabel
	codeLabel.TextStyle = fyne.TextStyle{Bold: true}

	lineInfo := canvas.NewText(fmt.Sprintf("%d baris perintah", len(lines)), constants.ColorTextMuted)
	lineInfo.TextSize = constants.FontSizeLabel
	lineInfo.TextStyle = fyne.TextStyle{Monospace: true}

	codeHeader := container.NewBorder(nil, nil, codeLabel, lineInfo)

	// Copy Script Button (Always copies current codeEntry.Text)
	copyBtn := widget.NewButtonWithIcon("Salin Script CLI", theme.ContentCopyIcon(), func() {
		p.copyToClip(codeEntry.Text)
	})
	copyBtn.Importance = widget.HighImportance

	headerLeft := container.NewVBox(
		container.NewHBox(devBadge, modeBadge, catBadge),
		titleLabel,
	)
	topHeader := container.NewBorder(nil, nil, nil, container.NewCenter(copyBtn), headerLeft)

	// Description
	descLabel := widget.NewLabel(cmd.Description)
	descLabel.Wrapping = fyne.TextWrapWord

	cardItems := []fyne.CanvasObject{
		topHeader,
		widget.NewSeparator(),
		descLabel,
	}

	// IP Example (if available) with clean aligned metadata strip
	var ipLabel *widget.Label
	if cmd.IPExample != "" {
		ipBadge := components.BadgeWarning("SKEMA IP / TOPOLOGI")
		ipLabel = widget.NewLabel(initialIPExample)
		ipLabel.Wrapping = fyne.TextWrapWord
		ipLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

		ipBox := container.NewBorder(nil, nil, container.NewCenter(ipBadge), nil, ipLabel)

		ipBg := canvas.NewRectangle(constants.ColorBgCardInner)
		ipBg.StrokeColor = constants.ColorBorderSubtle
		ipBg.StrokeWidth = 1
		ipBg.CornerRadius = constants.CornerRadiusBrutal

		ipCard := container.NewStack(ipBg, container.NewPadded(ipBox))
		cardItems = append(cardItems, ipCard)
	}

	// Update output function
	updateOutput := func() {
		rendered := cmd.RenderCommands(userParams)
		codeEntry.SetText(rendered)
		if ipLabel != nil {
			ipLabel.SetText(cmd.RenderIPExample(userParams))
		}
		currLines := strings.Split(rendered, "\n")
		lineInfo.Text = fmt.Sprintf("%d baris perintah", len(currLines))
		lineInfo.Refresh()
	}

	var fullCardContent *fyne.Container

	// Dynamic Parameter Customization Panel (if parameters exist)
	if len(cmd.Parameters) > 0 {
		var paramInputs []*widget.Entry
		var paramGridItems []fyne.CanvasObject

		for _, param := range cmd.Parameters {
			pKey := param.Key
			pDef := param.DefaultValue

			pLbl := canvas.NewText(param.Label, constants.ColorTextPrimary)
			pLbl.TextSize = constants.FontSizeLabel
			pLbl.TextStyle = fyne.TextStyle{Bold: true}

			pEnt := widget.NewEntry()
			pEnt.SetText(pDef)
			pEnt.SetPlaceHolder(param.Placeholder)
			paramInputs = append(paramInputs, pEnt)

			pEnt.OnChanged = func(val string) {
				userParams[pKey] = val
				updateOutput()
			}

			box := container.NewVBox(pLbl, pEnt)
			paramGridItems = append(paramGridItems, box)
		}

		cols := 2
		if len(paramGridItems) == 1 {
			cols = 1
		}
		paramGrid := container.NewGridWithColumns(cols, paramGridItems...)

		resetBtn := widget.NewButtonWithIcon("Reset Default", theme.ViewRefreshIcon(), func() {
			for i, pr := range cmd.Parameters {
				userParams[pr.Key] = pr.DefaultValue
				if i < len(paramInputs) {
					paramInputs[i].SetText(pr.DefaultValue)
				}
			}
			updateOutput()
		})
		resetBtn.Importance = widget.LowImportance

		paramBadge := components.BadgeYellow("⚙️ PARAMETER KUSTOMISASI TOPOLOGI")
		paramTop := container.NewBorder(nil, nil, paramBadge, resetBtn)

		paramBg := canvas.NewRectangle(constants.ColorBgCardInner)
		paramBg.StrokeColor = constants.ColorAccentYellow
		paramBg.StrokeWidth = 1.5
		paramBg.CornerRadius = constants.CornerRadiusBrutal

		paramContent := container.NewVBox(
			paramTop,
			widget.NewSeparator(),
			paramGrid,
		)
		paramPanel := container.NewStack(paramBg, container.NewPadded(paramContent))
		paramPanel.Hide()

		var toggleBtn *widget.Button
		toggleBtn = widget.NewButtonWithIcon("⚙️ Kustomisasi Parameter (Hostname, IP, VLAN...)", theme.SettingsIcon(), func() {
			if paramPanel.Visible() {
				paramPanel.Hide()
				toggleBtn.SetText("⚙️ Kustomisasi Parameter (Hostname, IP, VLAN...)")
				toggleBtn.SetIcon(theme.SettingsIcon())
			} else {
				paramPanel.Show()
				toggleBtn.SetText("▲ Sembunyikan Form Parameter")
				toggleBtn.SetIcon(theme.MenuDropUpIcon())
			}
			if fullCardContent != nil {
				fullCardContent.Refresh()
			}
		})
		toggleBtn.Importance = widget.MediumImportance

		cardItems = append(cardItems, toggleBtn, paramPanel)
	}

	terminalBg := canvas.NewRectangle(constants.ColorBgCardInner)
	terminalBg.StrokeColor = constants.ColorBorderSubtle
	terminalBg.StrokeWidth = constants.BorderWidthMedium
	terminalBg.CornerRadius = constants.CornerRadiusBrutal

	terminalContent := container.NewVBox(
		container.NewPadded(codeHeader),
		widget.NewSeparator(),
		container.NewPadded(codeStack),
	)
	terminalPanel := container.NewStack(terminalBg, terminalContent)
	cardItems = append(cardItems, terminalPanel)

	// Verification Guide Box with clean framed card
	verifBadge := components.BadgeSuccess("CARA VERIFIKASI / TES BEKERJA")
	verifLabel := widget.NewLabel(cmd.Verification)
	verifLabel.Wrapping = fyne.TextWrapWord

	verifBg := canvas.NewRectangle(constants.ColorBgCardInner)
	verifBg.StrokeColor = constants.ColorSuccess
	verifBg.StrokeWidth = constants.BorderWidthMedium
	verifBg.CornerRadius = constants.CornerRadiusBrutal

	verifContent := container.NewVBox(
		container.NewHBox(verifBadge),
		verifLabel,
	)
	verifCard := container.NewStack(verifBg, container.NewPadded(verifContent))
	cardItems = append(cardItems, verifCard)

	// Troubleshooting Tips Box with clean framed card
	if cmd.TroubleshootingTips != "" {
		tipsBadge := components.BadgeDanger("TIPS & JEBAKAN UMUM")
		tipsLabel := widget.NewLabel(cmd.TroubleshootingTips)
		tipsLabel.Wrapping = fyne.TextWrapWord

		tipsBg := canvas.NewRectangle(constants.ColorBgCardInner)
		tipsBg.StrokeColor = constants.ColorDanger
		tipsBg.StrokeWidth = constants.BorderWidthMedium
		tipsBg.CornerRadius = constants.CornerRadiusBrutal

		tipsContent := container.NewVBox(
			container.NewHBox(tipsBadge),
			tipsLabel,
		)
		tipsCard := container.NewStack(tipsBg, container.NewPadded(tipsContent))
		cardItems = append(cardItems, tipsCard)
	}

	fullCardContent = container.NewVBox(cardItems...)
	return components.NewPlainCardWithAccent(fullCardContent, accentColor)
}

func (p *CiscoPage) buildVerificationGuideView() fyne.CanvasObject {
	contentList := container.NewVBox()

	// Section 1: Diagnosa Lampu Indikator Fisik
	sec1Title := canvas.NewText("1. DIAGNOSA LAMPU INDIKATOR LINK (LAYER 1 FISIK)", constants.ColorTextPrimary)
	sec1Title.TextSize = constants.FontSizeH2
	sec1Title.TextStyle = fyne.TextStyle{Bold: true}

	makeInfoRow := func(badge fyne.CanvasObject, text string) fyne.CanvasObject {
		lbl := widget.NewLabel(text)
		lbl.Wrapping = fyne.TextWrapWord
		return container.NewBorder(nil, nil, container.NewCenter(badge), nil, lbl)
	}

	lampuHijau := makeInfoRow(components.BadgeSuccess("HIJAU (SOLID)"), "Link Layer 1 & 2 Normal (Port UP dan Protokol UP). Siap mengirim data.")
	lampuOranye := makeInfoRow(components.BadgeWarning("ORANYE (BLINK/SOLID)"), "Spanning Tree Protocol (STP) sedang Listening/Learning. Tunggu 30-50 detik atau klik tombol 'Fast Forward Time' (Alt + D) 2x di Packet Tracer.")
	lampuMerah := makeInfoRow(components.BadgeDanger("MERAH (SOLID)"), "Port Mati atau Kabel Salah! Cek: (a) Ketik 'no shutdown' di interface router, (b) Cek jenis kabel (Gunakan Straight-Through antar PC ke Switch, Crossover antar PC ke PC / Router ke PC).")

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

	rto := makeInfoRow(components.BadgeDanger("Request timed out (RTO)"), "Paket terkirim tapi tidak ada balasan. Cek: (1) Default Gateway di PC tujuan, (2) Routing balik dari router lawan, (3) Jalur terblokir ACL.")
	dhu := makeInfoRow(components.BadgeWarning("Destination Host Unreachable"), "Router tidak menemukan rute ke subnet tujuan. Cek tabel routing dengan 'show ip route' dan pastikan rute terdaftar.")
	firstDrop := makeInfoRow(components.BadgeYellow("Ping 1 Gagal, lalu Reply ( . ! ! ! ! )"), "Ini NORMAL di Cisco Packet Tracer! Paket pertama dipakai untuk proses broadcast ARP (Address Resolution Protocol) mencari MAC address.")

	card4 := components.NewPlainCardWithAccent(
		container.NewVBox(sec4Title, widget.NewSeparator(), rto, dhu, firstDrop),
		constants.ColorTechIndigo,
	)
	contentList.Add(card4)

	return container.NewVScroll(container.NewPadded(contentList))
}


func (p *CiscoPage) buildTopologyNotesView() fyne.CanvasObject {
	listContainer := container.NewVBox()

	// Form inputs for inline Create Topology Card
	tTitleEntry := widget.NewEntry()
	tTitleEntry.SetPlaceHolder("cth: Topologi 2 Router 2 Switch dengan DHCP & OSPF")
	tDescEntry := widget.NewMultiLineEntry()
	tDescEntry.SetPlaceHolder("cth: Hubungkan LAN Teknik (VLAN 10) dan LAN Keuangan (VLAN 20) lintas router WAN...")
	tDescEntry.SetMinRowsVisible(3)

	var formCard fyne.CanvasObject
	var reloadTopologies func()

	reloadTopologies = func() {
		listContainer.Objects = nil

		topologies, err := database.GetAllCiscoTopologies()
		if err != nil {
			listContainer.Add(widget.NewLabel(fmt.Sprintf("Gagal memuat catatan topologi: %v", err)))
			listContainer.Refresh()
			return
		}

		if len(topologies) == 0 {
			emptyTitle := canvas.NewText("Belum Ada Catatan Langkah Topologi", constants.ColorTextPrimary)
			emptyTitle.TextSize = constants.FontSizeH2
			emptyTitle.TextStyle = fyne.TextStyle{Bold: true}

			emptySub := canvas.NewText("Mulai dengan membuat rencana topologi baru atau muat template standar untuk referensi praktikum.", constants.ColorTextMuted)
			emptySub.TextSize = constants.FontSizeBody

			loadTemplateBtn := widget.NewButtonWithIcon("Muat Template Standar (ROAS & Routing)", theme.FolderIcon(), func() {
				_ = database.SeedDefaultCiscoTopologyTemplates()
				reloadTopologies()
			})
			loadTemplateBtn.Importance = widget.HighImportance

			emptyBadge := components.BadgeYellow("LAB STANDAR")
			emptyHeader := container.NewHBox(emptyTitle, emptyBadge)

			emptyBox := container.NewVBox(
				emptyHeader,
				emptySub,
				widget.NewSeparator(),
				container.NewHBox(loadTemplateBtn),
			)
			listContainer.Add(components.NewPlainCardWithAccent(emptyBox, constants.ColorTechIndigo))
			listContainer.Refresh()
			return
		}

		for _, t := range topologies {
			currTopo := t

			// Progress badge calculation
			completedCount := 0
			for _, st := range currTopo.Steps {
				if st.IsCompleted {
					completedCount++
				}
			}

			var progBadge fyne.CanvasObject
			var accentColor color.Color = constants.ColorTechIndigo

			if len(currTopo.Steps) == 0 {
				progBadge = components.BadgeMuted("Belum ada langkah")
			} else if completedCount == len(currTopo.Steps) {
				progBadge = components.BadgeSuccess(fmt.Sprintf("%d/%d Langkah Selesai ✓", completedCount, len(currTopo.Steps)))
				accentColor = constants.ColorSuccess
			} else {
				progBadge = components.BadgeYellow(fmt.Sprintf("%d/%d Selesai", completedCount, len(currTopo.Steps)))
				accentColor = constants.ColorAccentYellow
			}

			// Copy checklist to clipboard
			copyBtn := widget.NewButtonWithIcon("Salin Langkah", theme.ContentCopyIcon(), func() {
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("# %s\n", currTopo.Title))
				if currTopo.Description != "" {
					sb.WriteString(fmt.Sprintf("%s\n\n", currTopo.Description))
				} else {
					sb.WriteString("\n")
				}
				sb.WriteString("Langkah-Langkah Pembuatan Topologi:\n")
				for _, st := range currTopo.Steps {
					chkMark := "[ ]"
					if st.IsCompleted {
						chkMark = "[x]"
					}
					if st.Detail != "" {
						sb.WriteString(fmt.Sprintf("%d. %s %s - %s\n", st.StepNumber, chkMark, st.Title, st.Detail))
					} else {
						sb.WriteString(fmt.Sprintf("%d. %s %s\n", st.StepNumber, chkMark, st.Title))
					}
				}
				p.copyToClip(sb.String())
			})
			copyBtn.Importance = widget.LowImportance

			// Edit Topology Button using Neo-Brutalist modal
			editTopoBtn := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
				eTitleEntry := widget.NewEntry()
				eTitleEntry.SetText(currTopo.Title)
				eDescEntry := widget.NewMultiLineEntry()
				eDescEntry.SetText(currTopo.Description)
				eDescEntry.SetMinRowsVisible(3)

				lbl1 := canvas.NewText("NAMA / JUDUL TOPOLOGI", constants.ColorTextPrimary)
				lbl1.TextSize = constants.FontSizeLabel
				lbl1.TextStyle = fyne.TextStyle{Bold: true}

				lbl2 := canvas.NewText("DESKRIPSI / SKENARIO LAB", constants.ColorTextPrimary)
				lbl2.TextSize = constants.FontSizeLabel
				lbl2.TextStyle = fyne.TextStyle{Bold: true}

				formContent := container.NewVBox(
					lbl1, eTitleEntry,
					lbl2, eDescEntry,
				)

				components.ShowBrutalistFormDialog(
					p.window,
					"EDIT TOPOLOGI",
					constants.ColorAccentCyan,
					"Edit Catatan Topologi",
					"Perbarui nama atau target skenario topologi Cisco Packet Tracer",
					formContent,
					"Simpan Perubahan",
					func() {
						if eTitleEntry.Text != "" {
							_ = database.UpdateCiscoTopology(currTopo.ID, eTitleEntry.Text, eDescEntry.Text)
							reloadTopologies()
						}
					},
				)
			})
			editTopoBtn.Importance = widget.LowImportance

			// Delete Topology Button
			delTopoBtn := widget.NewButtonWithIcon("Hapus", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Topologi", fmt.Sprintf("Hapus topologi %q beserta semua langkahnya?", currTopo.Title), func(ok bool) {
					if ok {
						_ = database.DeleteCiscoTopology(currTopo.ID)
						reloadTopologies()
					}
				}, p.window)
			})
			delTopoBtn.Importance = widget.LowImportance

			// Header row
			titleTxt := canvas.NewText(currTopo.Title, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeH2
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			topoBadge := components.BadgeCyan("CISCO TOPOLOGY")
			headerLeft := container.NewHBox(titleTxt, topoBadge, progBadge)
			headerRight := container.NewHBox(copyBtn, editTopoBtn, delTopoBtn)
			headerRow := container.NewBorder(nil, nil, headerLeft, headerRight)

			var cardBody []fyne.CanvasObject
			cardBody = append(cardBody, headerRow)

			if currTopo.Description != "" {
				descBg := canvas.NewRectangle(constants.ColorBgCardInner)
				descBg.StrokeColor = constants.ColorBorderSubtle
				descBg.StrokeWidth = 1
				descBg.CornerRadius = constants.CornerRadiusBrutal

				descLabel := widget.NewLabel(currTopo.Description)
				descLabel.Wrapping = fyne.TextWrapWord

				descTitle := canvas.NewText("🎯 SKENARIO / TARGET LAB:", constants.ColorTextPrimary)
				descTitle.TextSize = constants.FontSizeLabel
				descTitle.TextStyle = fyne.TextStyle{Bold: true}

				descBox := container.NewStack(descBg, container.NewPadded(container.NewVBox(descTitle, descLabel)))
				cardBody = append(cardBody, descBox)
			}
			cardBody = append(cardBody, widget.NewSeparator())

			// Steps section header
			stepsHdr := canvas.NewText("📋 TAHAPAN & LANGKAH KONFIGURASI:", constants.ColorTextPrimary)
			stepsHdr.TextSize = constants.FontSizeLabel
			stepsHdr.TextStyle = fyne.TextStyle{Bold: true}
			cardBody = append(cardBody, stepsHdr)

			// Steps container
			stepsContainer := container.NewVBox()
			for _, st := range currTopo.Steps {
				currStep := st

				chk := widget.NewCheck("", func(bool) {
					_ = database.ToggleCiscoTopologyStep(currStep.ID)
					reloadTopologies()
				})
				chk.SetChecked(currStep.IsCompleted)

				var stepBadge fyne.CanvasObject
				if currStep.IsCompleted {
					stepBadge = components.BadgeSuccess(fmt.Sprintf("✓ Langkah %d", currStep.StepNumber))
				} else {
					stepBadge = components.BadgeIndigo(fmt.Sprintf("Langkah %d", currStep.StepNumber))
				}

				stepTitle := widget.NewLabel(currStep.Title)
				stepTitle.TextStyle = fyne.TextStyle{Bold: true}
				stepTitle.Wrapping = fyne.TextWrapWord

				editStepBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
					sTitleEntry := widget.NewEntry()
					sTitleEntry.SetText(currStep.Title)
					sDetailEntry := widget.NewMultiLineEntry()
					sDetailEntry.SetText(currStep.Detail)
					sDetailEntry.SetMinRowsVisible(3)

					lbl1 := canvas.NewText("JUDUL LANGKAH", constants.ColorTextPrimary)
					lbl1.TextSize = constants.FontSizeLabel
					lbl1.TextStyle = fyne.TextStyle{Bold: true}

					lbl2 := canvas.NewText("CATATAN / PERINTAH CLI (OPSIONAL)", constants.ColorTextPrimary)
					lbl2.TextSize = constants.FontSizeLabel
					lbl2.TextStyle = fyne.TextStyle{Bold: true}

					fContent := container.NewVBox(
						lbl1, sTitleEntry,
						lbl2, sDetailEntry,
					)

					components.ShowBrutalistFormDialog(
						p.window,
						fmt.Sprintf("LANGKAH %d", currStep.StepNumber),
						constants.ColorAccentYellow,
						fmt.Sprintf("Edit Langkah %d", currStep.StepNumber),
						"Sesuaikan judul atau perintah CLI untuk langkah ini",
						fContent,
						"Simpan Langkah",
						func() {
							if sTitleEntry.Text != "" {
								_ = database.UpdateCiscoTopologyStep(currStep.ID, sTitleEntry.Text, sDetailEntry.Text)
								reloadTopologies()
							}
						},
					)
				})
				editStepBtn.Importance = widget.LowImportance

				delStepBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
					_ = database.DeleteCiscoTopologyStep(currStep.ID)
					reloadTopologies()
				})
				delStepBtn.Importance = widget.LowImportance

				stepTopRight := container.NewHBox(editStepBtn, delStepBtn)
				stepTopLeft := container.NewHBox(chk, stepBadge)
				stepRow := container.NewBorder(nil, nil, stepTopLeft, stepTopRight, stepTitle)

				if currStep.Detail != "" {
					detailBg := canvas.NewRectangle(constants.ColorBgCardInner)
					detailBg.StrokeColor = constants.ColorBorderSubtle
					detailBg.StrokeWidth = 1
					detailBg.CornerRadius = constants.CornerRadiusBrutal

					detailLbl := widget.NewLabel(currStep.Detail)
					detailLbl.Wrapping = fyne.TextWrapWord
					detailLbl.TextStyle = fyne.TextStyle{Monospace: true}

					copyDetailBtn := widget.NewButtonWithIcon("Salin CLI", theme.ContentCopyIcon(), func() {
						p.copyToClip(currStep.Detail)
					})
					copyDetailBtn.Importance = widget.LowImportance

					detailHdr := container.NewBorder(nil, nil,
						components.BadgeMuted("DETAIL / PERINTAH CLI:"),
						copyDetailBtn,
					)

					detailBox := container.NewStack(detailBg, container.NewPadded(container.NewVBox(detailHdr, detailLbl)))
					stepFull := container.NewVBox(stepRow, detailBox)
					stepsContainer.Add(stepFull)
				} else {
					stepsContainer.Add(stepRow)
				}
			}

			cardBody = append(cardBody, stepsContainer)

			// Quick Add Step Form Strip
			newStepTitle := widget.NewEntry()
			newStepTitle.SetPlaceHolder("Judul langkah baru (misal: Pasang kabel / Konfigurasi VLAN)...")

			newStepDetail := widget.NewEntry()
			newStepDetail.SetPlaceHolder("Catatan / perintah CLI singkat (opsional)...")

			addStepBtn := widget.NewButtonWithIcon("Tambah Langkah", theme.ContentAddIcon(), func() {
				if newStepTitle.Text != "" {
					_ = database.AddCiscoTopologyStep(currTopo.ID, newStepTitle.Text, newStepDetail.Text)
					newStepTitle.SetText("")
					newStepDetail.SetText("")
					reloadTopologies()
				}
			})
			addStepBtn.Importance = widget.MediumImportance

			addBadge := components.BadgeYellow("+ LANGKAH")
			addStepInputs := container.NewGridWithColumns(2, newStepTitle, newStepDetail)
			addStepBox := container.NewBorder(nil, nil, addBadge, addStepBtn, addStepInputs)

			cardBody = append(cardBody, widget.NewSeparator(), addStepBox)

			fullCardContent := container.NewVBox(cardBody...)
			listContainer.Add(components.NewPlainCardWithAccent(fullCardContent, accentColor))
		}

		listContainer.Refresh()
	}

	// Inline Neo-Brutalist Form Card for Creating New Topology
	formTitle := canvas.NewText("📝 TAMBAH CATATAN TOPOLOGI BARU", constants.ColorTextPrimary)
	formTitle.TextSize = constants.FontSizeH2
	formTitle.TextStyle = fyne.TextStyle{Bold: true}

	formBadge := components.BadgeCyan("FORM INPUT")
	formHeader := container.NewHBox(formTitle, formBadge)

	formSub := canvas.NewText("Catat rancangan dan kebutuhan skenario topologi Cisco Packet Tracer baru Anda.", constants.ColorTextMuted)
	formSub.TextSize = constants.FontSizeSmall

	lblT := canvas.NewText("NAMA / JUDUL TOPOLOGI", constants.ColorTextPrimary)
	lblT.TextSize = constants.FontSizeLabel
	lblT.TextStyle = fyne.TextStyle{Bold: true}

	lblD := canvas.NewText("DESKRIPSI / TARGET SKENARIO LAB", constants.ColorTextPrimary)
	lblD.TextSize = constants.FontSizeLabel
	lblD.TextStyle = fyne.TextStyle{Bold: true}

	btnCancelCreate := widget.NewButtonWithIcon("Batal", theme.CancelIcon(), func() {
		if formCard != nil {
			formCard.Hide()
		}
	})
	btnCancelCreate.Importance = widget.LowImportance

	btnSaveCreate := widget.NewButtonWithIcon("Simpan Topologi", theme.DocumentSaveIcon(), func() {
		if tTitleEntry.Text != "" {
			_, _ = database.CreateCiscoTopology(tTitleEntry.Text, tDescEntry.Text)
			tTitleEntry.SetText("")
			tDescEntry.SetText("")
			if formCard != nil {
				formCard.Hide()
			}
			reloadTopologies()
		}
	})
	btnSaveCreate.Importance = widget.HighImportance

	formActions := container.NewBorder(nil, nil, nil, container.NewHBox(btnCancelCreate, btnSaveCreate))

	formInner := container.NewVBox(
		formHeader,
		formSub,
		widget.NewSeparator(),
		lblT,
		tTitleEntry,
		lblD,
		tDescEntry,
		widget.NewSeparator(),
		formActions,
	)
	formCard = components.NewPlainCardWithAccent(formInner, constants.ColorAccentCyan)
	formCard.Hide()

	// Top Action Bar
	topBarTitle := canvas.NewText("Catatan Topologi Baru", constants.ColorTextPrimary)
	topBarTitle.TextSize = constants.FontSizeH2
	topBarTitle.TextStyle = fyne.TextStyle{Bold: true}

	topBadge := components.BadgeCyan("CISCO LAB")
	topTitleBox := container.NewHBox(topBarTitle, topBadge)

	toggleCreateBtn := widget.NewButtonWithIcon("+ Topologi Baru", theme.ContentAddIcon(), func() {
		if formCard.Visible() {
			formCard.Hide()
		} else {
			formCard.Show()
		}
	})
	toggleCreateBtn.Importance = widget.HighImportance

	templateBtn := widget.NewButtonWithIcon("Template", theme.FolderIcon(), func() {
		_ = database.SeedDefaultCiscoTopologyTemplates()
		reloadTopologies()
	})
	templateBtn.Importance = widget.LowImportance

	topBarRight := container.NewHBox(templateBtn, toggleCreateBtn)
	topBar := container.NewBorder(nil, nil, topTitleBox, topBarRight)

	reloadTopologies()

	mainContent := container.NewVBox(
		container.NewPadded(topBar),
		formCard,
		listContainer,
	)

	return container.NewVScroll(container.NewPadded(mainContent))
}
