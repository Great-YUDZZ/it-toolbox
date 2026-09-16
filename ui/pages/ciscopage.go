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
		container.NewTabItemWithIcon(constants.TabCiscoLibrary, theme.FolderOpenIcon(), p.buildLibraryView()),
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
	codeEntry.Wrapping = fyne.TextWrapOff
	codeEntry.Scroll = fyne.ScrollNone

	lines := strings.Split(initialCommands, "\n")
	lineCount := len(lines)
	if lineCount < 4 {
		lineCount = 4
	}
	minHeight := float32(lineCount*21 + 28)
	if minHeight < 120 {
		minHeight = 120
	}

	codeSpacer := canvas.NewRectangle(color.Transparent)
	codeSpacer.SetMinSize(fyne.NewSize(0, minHeight))
	codeStack := container.NewStack(codeSpacer, codeEntry)

	codeLabel := canvas.NewText("PERINTAH CLI CISCO IOS (SIAP SALIN):", constants.ColorTextPrimary)
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

		paramBadge := components.BadgeYellow("PARAMETER KUSTOMISASI TOPOLOGI")
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
		toggleBtn = widget.NewButtonWithIcon("Kustomisasi Parameter (Hostname, IP, VLAN...)", theme.SettingsIcon(), func() {
			if paramPanel.Visible() {
				paramPanel.Hide()
				toggleBtn.SetText("Kustomisasi Parameter (Hostname, IP, VLAN...)")
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
	var mainScroller *container.Scroll
	listContainer := container.NewVBox()
	editingStepIDs := make(map[int]bool)

	// Form inputs for inline Create Topology Card
	tTitleEntry := widget.NewEntry()
	tTitleEntry.SetPlaceHolder("cth: Topologi 2 Router 2 Switch dengan DHCP & OSPF")
	tDescEntry := components.NewScrollableMultiLineEntry(nil)
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
				progBadge = components.BadgeSuccess(fmt.Sprintf("%d/%d Langkah Selesai", completedCount, len(currTopo.Steps)))
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
				eDescEntry := components.NewScrollableMultiLineEntry(nil)
				eDescEntry.SetText(currTopo.Description)
				eDescEntry.SetMinRowsVisible(5)

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

				descTitle := canvas.NewText("SKENARIO / TARGET LAB:", constants.ColorTextPrimary)
				descTitle.TextSize = constants.FontSizeLabel
				descTitle.TextStyle = fyne.TextStyle{Bold: true}

				descBox := container.NewStack(descBg, container.NewPadded(container.NewVBox(descTitle, descLabel)))
				cardBody = append(cardBody, descBox)
			}
			cardBody = append(cardBody, widget.NewSeparator())

			// Steps section header
			stepsHdr := canvas.NewText("TAHAPAN & LANGKAH KONFIGURASI:", constants.ColorTextPrimary)
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
					stepBadge = components.BadgeSuccess(fmt.Sprintf("Langkah %d", currStep.StepNumber))
				} else {
					stepBadge = components.BadgeIndigo(fmt.Sprintf("Langkah %d", currStep.StepNumber))
				}

				isEditing := editingStepIDs[currStep.ID]

				if isEditing {
					// --- INLINE EDIT MODE LANGSUNG DI KARTU LANGKAH ---
					sTitleEntry := widget.NewEntry()
					sTitleEntry.SetText(currStep.Title)
					sTitleEntry.SetPlaceHolder("Judul langkah...")

					sDetailEntry := components.NewScrollableMultiLineEntry(mainScroller)
					sDetailEntry.SetText(currStep.Detail)
					sDetailEntry.TextStyle = fyne.TextStyle{Monospace: true}
					sDetailEntry.SetMinRowsVisible(6)
					sDetailEntry.SetPlaceHolder("Perintah CLI / catatan konfigurasi langkah ini...")

					saveBtn := widget.NewButtonWithIcon("Simpan", theme.ConfirmIcon(), func() {
						if sTitleEntry.Text != "" {
							_ = database.UpdateCiscoTopologyStep(currStep.ID, sTitleEntry.Text, sDetailEntry.Text)
							delete(editingStepIDs, currStep.ID)
							reloadTopologies()
						}
					})
					saveBtn.Importance = widget.HighImportance

					cancelBtn := widget.NewButtonWithIcon("Batal", theme.CancelIcon(), func() {
						delete(editingStepIDs, currStep.ID)
						reloadTopologies()
					})
					cancelBtn.Importance = widget.LowImportance

					editHdrLeft := container.NewHBox(stepBadge, components.BadgeYellow("SEDANG MENGEDIT LANGSUNG"))
					editHdrRight := container.NewHBox(cancelBtn, saveBtn)
					editTop := container.NewBorder(nil, nil, editHdrLeft, editHdrRight)

					lblTitle := canvas.NewText("JUDUL LANGKAH", constants.ColorTextPrimary)
					lblTitle.TextSize = constants.FontSizeLabel
					lblTitle.TextStyle = fyne.TextStyle{Bold: true}

					lblDetail := canvas.NewText("CATATAN / PERINTAH CLI (OPSIONAL)", constants.ColorTextPrimary)
					lblDetail.TextSize = constants.FontSizeLabel
					lblDetail.TextStyle = fyne.TextStyle{Bold: true}

					editBox := container.NewVBox(
						editTop,
						lblTitle,
						sTitleEntry,
						lblDetail,
						sDetailEntry,
					)

					editBg := canvas.NewRectangle(constants.ColorBgCardInner)
					editBg.StrokeColor = constants.ColorAccentYellow
					editBg.StrokeWidth = 2
					editBg.CornerRadius = constants.CornerRadiusBrutal

					stepFull := container.NewStack(editBg, container.NewPadded(editBox))
					stepsContainer.Add(stepFull)
				} else {
					// --- NORMAL VIEW MODE ---
					stepTitle := widget.NewLabel(currStep.Title)
					stepTitle.TextStyle = fyne.TextStyle{Bold: true}
					stepTitle.Wrapping = fyne.TextWrapWord

					editStepBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
						editingStepIDs[currStep.ID] = true
						reloadTopologies()
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
	formTitle := canvas.NewText("TAMBAH CATATAN TOPOLOGI BARU", constants.ColorTextPrimary)
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

	mainScroller = container.NewVScroll(container.NewPadded(mainContent))
	tDescEntry.SetParentScroller(mainScroller)
	return mainScroller
}

// ============================================================================
// PERPUSTAKAAN KODE & KONFIGURASI CISCO (CISCO CONFIG LIBRARY)
// ============================================================================

func (p *CiscoPage) buildLibraryView() fyne.CanvasObject {
	cardsContainer := container.NewVBox()

	libQuery := ""
	libDevice := cisco.DeviceAll
	libCategory := cisco.CategoryAll
	libSource := "Semua Sumber"

	statsLabel := widget.NewLabel("Memuat perpustakaan...")
	statsLabel.TextStyle = fyne.TextStyle{Bold: true}

	var currentFilteredList []cisco.CiscoCommand

	var renderLibrary func()
	renderLibrary = func() {
		cardsContainer.Objects = nil

		builtin := cisco.GetAllCommands()
		customList, _ := database.GetAllCiscoCustomSnippets()

		var allItems []cisco.CiscoCommand
		if libSource != "Kustom Saya" {
			allItems = append(allItems, builtin...)
		}
		if libSource != "Koleksi Bawaan" {
			for _, cs := range customList {
				allItems = append(allItems, cisco.CiscoCommand{
					ID:          fmt.Sprintf("custom-%d", cs.ID),
					Title:       cs.Title,
					Device:      cisco.DeviceType(cs.Device),
					Category:    cisco.Category(cs.Category),
					Mode:        cisco.CLIMode(cs.Mode),
					Description: cs.Description,
					Commands:    cs.Commands,
					Verification: cs.Verification,
					IsCustom:    true,
					CustomID:    cs.ID,
				})
			}
		}

		// Update Stats
		tot := len(builtin) + len(customList)
		statsLabel.SetText(fmt.Sprintf("%d Total Resep Konfigurasi (%d Koleksi Bawaan, %d Resep Kustom)", tot, len(builtin), len(customList)))

		// Filter
		q := strings.ToLower(strings.TrimSpace(libQuery))
		var filtered []cisco.CiscoCommand

		for _, item := range allItems {
			// Device Filter
			if libDevice != cisco.DeviceAll && item.Device != libDevice && item.Device != cisco.DeviceAll {
				continue
			}

			// Category Filter
			if libCategory != cisco.CategoryAll && item.Category != libCategory {
				continue
			}

			// Text Search
			if q != "" {
				mTitle := strings.Contains(strings.ToLower(item.Title), q)
				mDesc := strings.Contains(strings.ToLower(item.Description), q)
				mCmds := strings.Contains(strings.ToLower(item.Commands), q)
				mVerif := strings.Contains(strings.ToLower(item.Verification), q)
				mTips := strings.Contains(strings.ToLower(item.TroubleshootingTips), q)
				mCat := strings.Contains(strings.ToLower(string(item.Category)), q)
				mDev := strings.Contains(strings.ToLower(string(item.Device)), q)

				mTags := false
				for _, tg := range item.Tags {
					if strings.Contains(strings.ToLower(tg), q) {
						mTags = true
						break
					}
				}

				mExpl := false
				for _, ex := range item.Explanation {
					if strings.Contains(strings.ToLower(ex.Command), q) || strings.Contains(strings.ToLower(ex.Explanation), q) {
						mExpl = true
						break
					}
				}

				if !mTitle && !mDesc && !mCmds && !mVerif && !mTips && !mCat && !mDev && !mTags && !mExpl {
					continue
				}
			}

			filtered = append(filtered, item)
		}

		currentFilteredList = filtered

		if len(filtered) == 0 {
			emptyTitle := canvas.NewText("Tidak Ada Kode Konfigurasi yang Ditemukan", constants.ColorTextPrimary)
			emptyTitle.TextSize = constants.FontSizeH2
			emptyTitle.TextStyle = fyne.TextStyle{Bold: true}

			emptySub := canvas.NewText("Coba ubah kata kunci pencarian atau sesuaikan filter perangkat dan kategori.", constants.ColorTextMuted)
			emptySub.TextSize = constants.FontSizeBody

			resetBtn := widget.NewButtonWithIcon("Reset Semua Filter", theme.ViewRefreshIcon(), func() {
				libQuery = ""
				libDevice = cisco.DeviceAll
				libCategory = cisco.CategoryAll
				libSource = "Semua Sumber"
				renderLibrary()
			})
			resetBtn.Importance = widget.MediumImportance

			emptyBox := container.NewVBox(
				container.NewHBox(emptyTitle, components.BadgeWarning("HASIL KOSONG")),
				emptySub,
				widget.NewSeparator(),
				container.NewHBox(resetBtn),
			)
			cardsContainer.Add(components.NewPlainCardWithAccent(emptyBox, constants.ColorWarning))
			cardsContainer.Refresh()
			return
		}

		for _, item := range filtered {
			cmd := item
			card := p.buildLibraryCard(cmd, renderLibrary)
			cardsContainer.Add(card)
		}
		cardsContainer.Refresh()
	}

	// 1. Search Bar
	searchBar := components.NewSearchBar(constants.SearchCiscoLibraryPlaceholder, func(q string) {
		libQuery = q
		renderLibrary()
	})

	// 2. Device Selector
	deviceOptions := []string{
		string(cisco.DeviceAll),
		string(cisco.DeviceRouter),
		string(cisco.DeviceSwitchL2),
		string(cisco.DeviceSwitchL3),
		string(cisco.DevicePC),
	}
	devSelect := widget.NewSelect(deviceOptions, func(val string) {
		libDevice = cisco.DeviceType(val)
		renderLibrary()
	})
	devSelect.SetSelected(string(cisco.DeviceAll))

	devSpacer := canvas.NewRectangle(color.Transparent)
	devSpacer.SetMinSize(fyne.NewSize(140, 36))
	devSelectBox := container.NewStack(devSpacer, devSelect)

	// 3. Category Selector
	categoryOptions := []string{
		string(cisco.CategoryAll),
		string(cisco.CategoryBasic),
		string(cisco.CategoryInterface),
		string(cisco.CategoryVLAN),
		string(cisco.CategorySpanningTree),
		string(cisco.CategoryRouting),
		string(cisco.CategoryServices),
		string(cisco.CategorySecurity),
		string(cisco.CategoryNAT),
		string(cisco.CategoryRedundancy),
		string(cisco.CategoryWAN),
		string(cisco.CategoryHardening),
		string(cisco.CategoryRecovery),
		string(cisco.CategoryShowDiag),
	}
	catSelect := widget.NewSelect(categoryOptions, func(val string) {
		libCategory = cisco.Category(val)
		renderLibrary()
	})
	catSelect.SetSelected(string(cisco.CategoryAll))

	catSpacer := canvas.NewRectangle(color.Transparent)
	catSpacer.SetMinSize(fyne.NewSize(200, 36))
	catSelectBox := container.NewStack(catSpacer, catSelect)

	// 4. Source Selector
	sourceOptions := []string{
		"Semua Sumber",
		"Koleksi Bawaan",
		"Kustom Saya",
	}
	sourceSelect := widget.NewSelect(sourceOptions, func(val string) {
		libSource = val
		renderLibrary()
	})
	sourceSelect.SetSelected("Semua Sumber")

	sourceSpacer := canvas.NewRectangle(color.Transparent)
	sourceSpacer.SetMinSize(fyne.NewSize(130, 36))
	sourceSelectBox := container.NewStack(sourceSpacer, sourceSelect)

	// Action Buttons
	addCustomBtn := widget.NewButtonWithIcon("+ Tambah Resep Kustom", theme.ContentAddIcon(), func() {
		p.showAddCustomSnippetDialog(renderLibrary)
	})
	addCustomBtn.Importance = widget.HighImportance

	exportBtn := widget.NewButtonWithIcon("Ekspor Koleksi (.txt)", theme.DocumentSaveIcon(), func() {
		p.exportLibraryCheatSheet(currentFilteredList)
	})
	exportBtn.Importance = widget.LowImportance

	filterRow1 := container.NewHBox(
		container.NewCenter(components.BadgeCyan("PERANGKAT:")),
		devSelectBox,
		widget.NewSeparator(),
		container.NewCenter(components.BadgeYellow("KATEGORI:")),
		catSelectBox,
	)

	filterRow2 := container.NewBorder(nil, nil,
		container.NewHBox(
			container.NewCenter(components.BadgeIndigo("SUMBER:")),
			sourceSelectBox,
		),
		container.NewHBox(exportBtn, addCustomBtn),
	)

	// Hero Header Card for Perpustakaan
	libHeroTitle := canvas.NewText("PERPUSTAKAAN KODE KONFIGURASI CISCO", constants.ColorTextPrimary)
	libHeroTitle.TextSize = constants.FontSizeH2
	libHeroTitle.TextStyle = fyne.TextStyle{Bold: true}

	libHeroBadge := components.BadgeCyan("CCNA & PACKET TRACER")
	libHeroDesc := canvas.NewText("Ensiklopedia lengkap seluruh kode perintah konfigurasi Cisco IOS Router, Switch L2/L3, dan PC. Dilengkapi penjelasan baris demi baris, verifikasi, dan penyimpanan resep kustom.", constants.ColorTextMuted)
	libHeroDesc.TextSize = constants.FontSizeBody

	libHeaderBox := container.NewVBox(
		container.NewHBox(libHeroTitle, libHeroBadge),
		libHeroDesc,
		widget.NewSeparator(),
		statsLabel,
	)
	libHeaderCard := components.NewPlainCardWithAccent(libHeaderBox, constants.ColorTechIndigo)

	headerBox := container.NewVBox(
		libHeaderCard,
		searchBar.Container,
		filterRow1,
		filterRow2,
		widget.NewSeparator(),
	)

	renderLibrary()

	scrollList := container.NewVScroll(container.NewPadded(cardsContainer))
	return container.NewBorder(headerBox, nil, nil, nil, scrollList)
}

func (p *CiscoPage) buildLibraryCard(cmd cisco.CiscoCommand, refreshFn func()) fyne.CanvasObject {
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

	var srcBadge fyne.CanvasObject
	if cmd.IsCustom {
		srcBadge = components.BadgeWarning("KUSTOM SAYA")
	} else {
		srcBadge = components.BadgeSuccess("BAWAAN")
	}

	// Title
	titleLabel := widget.NewLabel(cmd.Title)
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Wrapping = fyne.TextWrapWord

	// Copy Script Button
	copyBtn := widget.NewButtonWithIcon("Salin Kode CLI", theme.ContentCopyIcon(), func() {
		p.copyToClip(cmd.Commands)
	})
	copyBtn.Importance = widget.HighImportance

	headerBadges := container.NewHBox(devBadge, modeBadge, catBadge, srcBadge)
	headerLeft := container.NewVBox(headerBadges, titleLabel)

	var topRightButtons []fyne.CanvasObject
	if cmd.IsCustom {
		editBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			p.showEditCustomSnippetDialog(cmd, refreshFn)
		})
		editBtn.Importance = widget.LowImportance

		delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			dialog.ShowConfirm("Hapus Resep Kustom", fmt.Sprintf("Hapus resep konfigurasi %q dari perpustakaan Anda?", cmd.Title), func(ok bool) {
				if ok {
					_ = database.DeleteCiscoCustomSnippet(cmd.CustomID)
					refreshFn()
				}
			}, p.window)
		})
		delBtn.Importance = widget.LowImportance
		topRightButtons = append(topRightButtons, editBtn, delBtn)
	}
	topRightButtons = append(topRightButtons, copyBtn)
	headerRight := container.NewHBox(topRightButtons...)

	topHeader := container.NewBorder(nil, nil, headerLeft, headerRight)

	// Description
	descLabel := widget.NewLabel(cmd.Description)
	descLabel.Wrapping = fyne.TextWrapWord

	// Code Entry Box
	codeEntry := widget.NewMultiLineEntry()
	codeEntry.SetText(cmd.Commands)
	codeEntry.TextStyle = fyne.TextStyle{Monospace: true}
	codeEntry.Wrapping = fyne.TextWrapOff
	codeEntry.Scroll = fyne.ScrollNone

	lines := strings.Split(cmd.Commands, "\n")
	lineCount := len(lines)
	if lineCount < 4 {
		lineCount = 4
	}
	minHeight := float32(lineCount*21 + 28)
	if minHeight < 110 {
		minHeight = 110
	}

	codeSpacer := canvas.NewRectangle(color.Transparent)
	codeSpacer.SetMinSize(fyne.NewSize(0, minHeight))
	codeStack := container.NewStack(codeSpacer, codeEntry)

	codeLabel := canvas.NewText("SKRIP KONFIGURASI LENGKAP:", constants.ColorTextPrimary)
	codeLabel.TextSize = constants.FontSizeLabel
	codeLabel.TextStyle = fyne.TextStyle{Bold: true}

	lineInfo := canvas.NewText(fmt.Sprintf("%d baris kode", len(lines)), constants.ColorTextMuted)
	lineInfo.TextSize = constants.FontSizeLabel
	lineInfo.TextStyle = fyne.TextStyle{Monospace: true}

	codeHeader := container.NewBorder(nil, nil, codeLabel, lineInfo)

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

	// Card Actions Row
	explainBtn := widget.NewButtonWithIcon("Detail & Penjelasan Tiap Baris", theme.InfoIcon(), func() {
		p.showLineExplanationsDialog(cmd)
	})
	explainBtn.Importance = widget.MediumImportance

	// Verification Box (Expandable)
	var verifPanel *fyne.Container
	if cmd.Verification != "" || cmd.TroubleshootingTips != "" {
		var verifItems []fyne.CanvasObject
		if cmd.Verification != "" {
			vBadge := components.BadgeSuccess("CARA VERIFIKASI / TES")
			vLbl := widget.NewLabel(cmd.Verification)
			vLbl.Wrapping = fyne.TextWrapWord
			verifItems = append(verifItems, container.NewHBox(vBadge), vLbl)
		}
		if cmd.TroubleshootingTips != "" {
			tBadge := components.BadgeDanger("TIPS & JEBAKAN UMUM")
			tLbl := widget.NewLabel(cmd.TroubleshootingTips)
			tLbl.Wrapping = fyne.TextWrapWord
			verifItems = append(verifItems, widget.NewSeparator(), container.NewHBox(tBadge), tLbl)
		}

		vBg := canvas.NewRectangle(constants.ColorBgCardInner)
		vBg.StrokeColor = constants.ColorSuccess
		vBg.StrokeWidth = 1.5
		vBg.CornerRadius = constants.CornerRadiusBrutal

		verifPanel = container.NewStack(vBg, container.NewPadded(container.NewVBox(verifItems...)))
		verifPanel.Hide()
	}

	var toggleVerifBtn *widget.Button
	if verifPanel != nil {
		toggleVerifBtn = widget.NewButtonWithIcon("Cara Verifikasi & Tips", theme.ConfirmIcon(), func() {
			if verifPanel.Visible() {
				verifPanel.Hide()
				toggleVerifBtn.SetText("Cara Verifikasi & Tips")
				toggleVerifBtn.SetIcon(theme.ConfirmIcon())
			} else {
				verifPanel.Show()
				toggleVerifBtn.SetText("▲ Sembunyikan Verifikasi")
				toggleVerifBtn.SetIcon(theme.MenuDropUpIcon())
			}
		})
		toggleVerifBtn.Importance = widget.LowImportance
	}

	var actionRow fyne.CanvasObject
	if toggleVerifBtn != nil {
		actionRow = container.NewHBox(explainBtn, toggleVerifBtn)
	} else {
		actionRow = container.NewHBox(explainBtn)
	}

	cardItems := []fyne.CanvasObject{
		topHeader,
		widget.NewSeparator(),
		descLabel,
		terminalPanel,
		actionRow,
	}
	if verifPanel != nil {
		cardItems = append(cardItems, verifPanel)
	}

	fullCardContent := container.NewVBox(cardItems...)
	return components.NewPlainCardWithAccent(fullCardContent, accentColor)
}

func (p *CiscoPage) showLineExplanationsDialog(cmd cisco.CiscoCommand) {
	var listItems []fyne.CanvasObject

	headerTitle := canvas.NewText("Penjelasan Teknis Perintah CLI", constants.ColorTextPrimary)
	headerTitle.TextSize = constants.FontSizeH2
	headerTitle.TextStyle = fyne.TextStyle{Bold: true}

	subTitle := canvas.NewText(cmd.Title, constants.ColorTextMuted)
	subTitle.TextSize = constants.FontSizeSmall
	subTitle.TextStyle = fyne.TextStyle{Bold: true}

	listItems = append(listItems, headerTitle, subTitle, widget.NewSeparator())

	if len(cmd.Explanation) > 0 {
		for i, ex := range cmd.Explanation {
			lineNumBadge := components.BadgeIndigo(fmt.Sprintf("Baris %d", i+1))

			cmdText := canvas.NewText(ex.Command, constants.ColorTextPrimary)
			cmdText.TextSize = constants.FontSizeBody
			cmdText.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

			topRow := container.NewHBox(lineNumBadge, cmdText)

			explLabel := widget.NewLabel(ex.Explanation)
			explLabel.Wrapping = fyne.TextWrapWord

			explBg := canvas.NewRectangle(constants.ColorBgCardInner)
			explBg.StrokeColor = constants.ColorBorderSubtle
			explBg.StrokeWidth = 1
			explBg.CornerRadius = constants.CornerRadiusBrutal

			boxContent := container.NewVBox(topRow, explLabel)
			box := container.NewStack(explBg, container.NewPadded(boxContent))
			listItems = append(listItems, box)
		}
	} else {
		// Fallback: breakdown lines of cmd.Commands
		rawLines := strings.Split(cmd.Commands, "\n")
		validIndex := 1
		for _, rawLine := range rawLines {
			line := strings.TrimSpace(rawLine)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			numBadge := components.BadgeIndigo(fmt.Sprintf("Baris %d", validIndex))
			validIndex++

			cmdTxt := canvas.NewText(line, constants.ColorTextPrimary)
			cmdTxt.TextSize = constants.FontSizeBody
			cmdTxt.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

			desc := "Perintah konfigurasi aktif untuk dieksekusi pada mode prompt Cisco IOS terkait."
			if strings.HasPrefix(line, "enable") {
				desc = "Beralih dari User EXEC mode (>) ke Privileged EXEC mode (#) untuk akses konfigurasi tingkat lanjut."
			} else if strings.HasPrefix(line, "configure terminal") || line == "conf t" {
				desc = "Masuk ke mode Global Configuration (config)# untuk mengubah setelan sistem perangkat."
			} else if strings.HasPrefix(line, "exit") {
				desc = "Keluar dari sub-mode saat ini ke mode satu tingkat di atasnya."
			} else if strings.HasPrefix(line, "no shutdown") {
				desc = "Menghidupkan port antarmuka fisik (mengubah status dari administratively down menjadi UP)."
			} else if strings.HasPrefix(line, "write memory") || strings.HasPrefix(line, "copy running-config") {
				desc = "Menyimpan konfigurasi aktif di RAM ke NVRAM agar tidak hilang saat perangkat mati atau reboot."
			}

			lbl := widget.NewLabel(desc)
			lbl.Wrapping = fyne.TextWrapWord

			bg := canvas.NewRectangle(constants.ColorBgCardInner)
			bg.StrokeColor = constants.ColorBorderSubtle
			bg.StrokeWidth = 1
			bg.CornerRadius = constants.CornerRadiusBrutal

			itemBox := container.NewStack(bg, container.NewPadded(container.NewVBox(container.NewHBox(numBadge, cmdTxt), lbl)))
			listItems = append(listItems, itemBox)
		}
	}

	scrollContent := container.NewVScroll(container.NewVBox(listItems...))
	scrollContent.SetMinSize(fyne.NewSize(580, 400))

	var d dialog.Dialog
	closeBtn := widget.NewButtonWithIcon("Tutup", theme.CancelIcon(), func() {
		if d != nil {
			d.Hide()
		}
	})
	closeBtn.Importance = widget.HighImportance

	dialogContent := container.NewBorder(nil, container.NewCenter(closeBtn), nil, nil, scrollContent)
	d = dialog.NewCustom("BEDAH KODE CISCO", "Tutup", dialogContent, p.window)
	d.Resize(fyne.NewSize(620, 480))
	d.Show()
}

func (p *CiscoPage) showAddCustomSnippetDialog(refreshFn func()) {
	tTitleEntry := widget.NewEntry()
	tTitleEntry.SetPlaceHolder("cth: Setup EtherChannel LACP Switch Core")

	deviceOptions := []string{
		string(cisco.DeviceRouter),
		string(cisco.DeviceSwitchL2),
		string(cisco.DeviceSwitchL3),
		string(cisco.DevicePC),
	}
	devSelect := widget.NewSelect(deviceOptions, nil)
	devSelect.SetSelected(string(cisco.DeviceRouter))

	categoryOptions := []string{
		string(cisco.CategoryBasic),
		string(cisco.CategoryInterface),
		string(cisco.CategoryVLAN),
		string(cisco.CategorySpanningTree),
		string(cisco.CategoryRouting),
		string(cisco.CategoryServices),
		string(cisco.CategorySecurity),
		string(cisco.CategoryNAT),
		string(cisco.CategoryRedundancy),
		string(cisco.CategoryWAN),
		string(cisco.CategoryHardening),
		string(cisco.CategoryRecovery),
		string(cisco.CategoryShowDiag),
	}
	catSelect := widget.NewSelect(categoryOptions, nil)
	catSelect.SetSelected(string(cisco.CategoryBasic))

	modeOptions := []string{
		string(cisco.ModeGlobalConfig),
		string(cisco.ModePrivExec),
		string(cisco.ModeInterface),
		string(cisco.ModeVLAN),
		string(cisco.ModeLine),
		string(cisco.ModeRouterConfig),
		string(cisco.ModeDHCPConfig),
		string(cisco.ModeUserExec),
		string(cisco.ModePCTerminal),
	}
	modeSelect := widget.NewSelect(modeOptions, nil)
	modeSelect.SetSelected(string(cisco.ModeGlobalConfig))

	tDescEntry := components.NewScrollableMultiLineEntry(nil)
	tDescEntry.SetPlaceHolder("cth: Skrip konfigurasi untuk menggabungkan port trunk antar switch...")
	tDescEntry.SetMinRowsVisible(3)

	tCmdsEntry := components.NewScrollableMultiLineEntry(nil)
	tCmdsEntry.SetPlaceHolder("Ketik kode konfigurasi Cisco IOS di sini (tiap baris perintah)...")
	tCmdsEntry.TextStyle = fyne.TextStyle{Monospace: true}
	tCmdsEntry.SetMinRowsVisible(8)

	tVerifEntry := components.NewScrollableMultiLineEntry(nil)
	tVerifEntry.SetPlaceHolder("cth: show etherchannel summary / show running-config")
	tVerifEntry.SetMinRowsVisible(3)

	makeLbl := func(txt string) *canvas.Text {
		t := canvas.NewText(txt, constants.ColorTextPrimary)
		t.TextSize = constants.FontSizeLabel
		t.TextStyle = fyne.TextStyle{Bold: true}
		return t
	}

	formContent := container.NewVBox(
		makeLbl("NAMA / JUDUL RESEP KONFIGURASI"),
		tTitleEntry,
		container.NewGridWithColumns(3,
			container.NewVBox(makeLbl("PERANGKAT"), devSelect),
			container.NewVBox(makeLbl("KATEGORI"), catSelect),
			container.NewVBox(makeLbl("MODE PROMPT CLI"), modeSelect),
		),
		makeLbl("DESKRIPSI / TUJUAN"),
		tDescEntry,
		makeLbl("SKRIP PERINTAH CLI (MONOSPACE)"),
		tCmdsEntry,
		makeLbl("CARA VERIFIKASI / CATATAN"),
		tVerifEntry,
	)

	components.ShowBrutalistFormDialog(
		p.window,
		"TAMBAH KODE",
		constants.ColorAccentCyan,
		"Tambah Resep Cisco Baru",
		"Simpan kode konfigurasi kustom Anda ke perpustakaan lokal",
		formContent,
		"Simpan ke Perpustakaan",
		func() {
			if strings.TrimSpace(tTitleEntry.Text) == "" || strings.TrimSpace(tCmdsEntry.Text) == "" {
				dialog.ShowError(fmt.Errorf("Judul dan Skrip Perintah tidak boleh kosong"), p.window)
				return
			}
			_, err := database.CreateCiscoCustomSnippet(database.CiscoCustomSnippet{
				Title:        tTitleEntry.Text,
				Device:       devSelect.Selected,
				Category:     catSelect.Selected,
				Mode:         modeSelect.Selected,
				Description:  tDescEntry.Text,
				Commands:     tCmdsEntry.Text,
				Verification: tVerifEntry.Text,
			})
			if err != nil {
				dialog.ShowError(fmt.Errorf("Gagal menyimpan resep: %w", err), p.window)
				return
			}
			dialog.ShowInformation("Tersimpan", "Resep konfigurasi berhasil ditambahkan ke perpustakaan!", p.window)
			refreshFn()
		},
	)
}

func (p *CiscoPage) showEditCustomSnippetDialog(cmd cisco.CiscoCommand, refreshFn func()) {
	tTitleEntry := widget.NewEntry()
	tTitleEntry.SetText(cmd.Title)

	deviceOptions := []string{
		string(cisco.DeviceRouter),
		string(cisco.DeviceSwitchL2),
		string(cisco.DeviceSwitchL3),
		string(cisco.DevicePC),
	}
	devSelect := widget.NewSelect(deviceOptions, nil)
	devSelect.SetSelected(string(cmd.Device))

	categoryOptions := []string{
		string(cisco.CategoryBasic),
		string(cisco.CategoryInterface),
		string(cisco.CategoryVLAN),
		string(cisco.CategorySpanningTree),
		string(cisco.CategoryRouting),
		string(cisco.CategoryServices),
		string(cisco.CategorySecurity),
		string(cisco.CategoryNAT),
		string(cisco.CategoryRedundancy),
		string(cisco.CategoryWAN),
		string(cisco.CategoryHardening),
		string(cisco.CategoryRecovery),
		string(cisco.CategoryShowDiag),
	}
	catSelect := widget.NewSelect(categoryOptions, nil)
	catSelect.SetSelected(string(cmd.Category))

	modeOptions := []string{
		string(cisco.ModeGlobalConfig),
		string(cisco.ModePrivExec),
		string(cisco.ModeInterface),
		string(cisco.ModeVLAN),
		string(cisco.ModeLine),
		string(cisco.ModeRouterConfig),
		string(cisco.ModeDHCPConfig),
		string(cisco.ModeUserExec),
		string(cisco.ModePCTerminal),
	}
	modeSelect := widget.NewSelect(modeOptions, nil)
	modeSelect.SetSelected(string(cmd.Mode))

	tDescEntry := components.NewScrollableMultiLineEntry(nil)
	tDescEntry.SetText(cmd.Description)
	tDescEntry.SetMinRowsVisible(3)

	tCmdsEntry := components.NewScrollableMultiLineEntry(nil)
	tCmdsEntry.SetText(cmd.Commands)
	tCmdsEntry.TextStyle = fyne.TextStyle{Monospace: true}
	tCmdsEntry.SetMinRowsVisible(8)

	tVerifEntry := components.NewScrollableMultiLineEntry(nil)
	tVerifEntry.SetText(cmd.Verification)
	tVerifEntry.SetMinRowsVisible(3)

	makeLbl := func(txt string) *canvas.Text {
		t := canvas.NewText(txt, constants.ColorTextPrimary)
		t.TextSize = constants.FontSizeLabel
		t.TextStyle = fyne.TextStyle{Bold: true}
		return t
	}

	formContent := container.NewVBox(
		makeLbl("NAMA / JUDUL RESEP KONFIGURASI"),
		tTitleEntry,
		container.NewGridWithColumns(3,
			container.NewVBox(makeLbl("PERANGKAT"), devSelect),
			container.NewVBox(makeLbl("KATEGORI"), catSelect),
			container.NewVBox(makeLbl("MODE PROMPT CLI"), modeSelect),
		),
		makeLbl("DESKRIPSI / TUJUAN"),
		tDescEntry,
		makeLbl("SKRIP PERINTAH CLI (MONOSPACE)"),
		tCmdsEntry,
		makeLbl("CARA VERIFIKASI / CATATAN"),
		tVerifEntry,
	)

	components.ShowBrutalistFormDialog(
		p.window,
		"EDIT RESEP",
		constants.ColorAccentYellow,
		"Edit Resep Kustom",
		"Perbarui rincian resep konfigurasi perpustakaan Anda",
		formContent,
		"Simpan Perubahan",
		func() {
			if strings.TrimSpace(tTitleEntry.Text) == "" || strings.TrimSpace(tCmdsEntry.Text) == "" {
				dialog.ShowError(fmt.Errorf("Judul dan Skrip Perintah tidak boleh kosong"), p.window)
				return
			}
			err := database.UpdateCiscoCustomSnippet(database.CiscoCustomSnippet{
				ID:           cmd.CustomID,
				Title:        tTitleEntry.Text,
				Device:       devSelect.Selected,
				Category:     catSelect.Selected,
				Mode:         modeSelect.Selected,
				Description:  tDescEntry.Text,
				Commands:     tCmdsEntry.Text,
				Verification: tVerifEntry.Text,
			})
			if err != nil {
				dialog.ShowError(fmt.Errorf("Gagal memperbarui resep: %w", err), p.window)
				return
			}
			dialog.ShowInformation("Tersimpan", "Perubahan berhasil disimpan!", p.window)
			refreshFn()
		},
	)
}

func (p *CiscoPage) exportLibraryCheatSheet(items []cisco.CiscoCommand) {
	if len(items) == 0 {
		dialog.ShowInformation("Ekspor", "Tidak ada resep yang sesuai untuk diekspor.", p.window)
		return
	}

	var sb strings.Builder
	sb.WriteString("================================================================================\n")
	sb.WriteString("IT TOOLBOX - PERPUSTAKAAN KODE KONFIGURASI CISCO PACKET TRACER\n")
	sb.WriteString(fmt.Sprintf("Jumlah Resep: %d Konfigurasi\n", len(items)))
	sb.WriteString("================================================================================\n\n")

	for i, it := range items {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, it.Title))
		sb.WriteString(fmt.Sprintf("Perangkat : %s\n", it.Device))
		sb.WriteString(fmt.Sprintf("Kategori  : %s\n", it.Category))
		sb.WriteString(fmt.Sprintf("Mode CLI  : %s\n", it.Mode))
		if it.Description != "" {
			sb.WriteString(fmt.Sprintf("Deskripsi : %s\n", it.Description))
		}
		sb.WriteString("\n--- SKRIP KONFIGURASI CLI ---\n")
		sb.WriteString(it.Commands)
		sb.WriteString("\n")
		if it.Verification != "" {
			sb.WriteString("\n--- CARA VERIFIKASI ---\n")
			sb.WriteString(it.Verification)
			sb.WriteString("\n")
		}
		sb.WriteString("\n--------------------------------------------------------------------------------\n\n")
	}

	p.copyToClip(sb.String())
	dialog.ShowInformation("Ekspor Berhasil", fmt.Sprintf("%d resep konfigurasi Cisco berhasil diekspor dan disalin ke clipboard!\nAnda dapat langsung mem-paste ke berkas teks atau Notepad.", len(items)), p.window)
}

