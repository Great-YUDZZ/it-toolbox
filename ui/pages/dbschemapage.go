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

	"github.com/yudz/it-toolbox/core/dbschema"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

type DBSchemaPage struct {
	window fyne.Window

	cachedView fyne.CanvasObject

	// Builder State
	builderTableDef dbschema.TableDef

	// Catalog filter state
	selectedEngine   dbschema.DatabaseEngine
	selectedCategory dbschema.CommandCategory
	searchQuery      string
}

func NewDBSchemaPage(win fyne.Window) *DBSchemaPage {
	return &DBSchemaPage{
		window:           win,
		selectedEngine:   dbschema.EngineAll,
		selectedCategory: dbschema.CatAll,
		builderTableDef: dbschema.TableDef{
			DatabaseName: "db_aplikasi",
			TableName:    "pengguna",
			Engine:       dbschema.EngineMySQL,
			IncludeDrop:  true,
			Columns:      dbschema.PresetColumnsUser(),
		},
	}
}

func (p *DBSchemaPage) copyToClip(txt string) {
	if p.window != nil {
		p.window.Clipboard().SetContent(txt)
		dialog.ShowInformation("Clipboard", constants.StatusCopied, p.window)
	}
}

func (p *DBSchemaPage) exportSQLFile(defaultName, sqlContent string) {
	if p.window == nil {
		return
	}
	saveDlg := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if err != nil || uc == nil {
			return
		}
		defer uc.Close()
		_, _ = uc.Write([]byte(sqlContent))
		dialog.ShowInformation("Ekspor Berhasil", fmt.Sprintf("Berkas SQL berhasil disimpan:\n%s", uc.URI().Path()), p.window)
	}, p.window)

	if !strings.HasSuffix(strings.ToLower(defaultName), ".sql") {
		defaultName += ".sql"
	}
	saveDlg.SetFileName(defaultName)
	saveDlg.Show()
}

func (p *DBSchemaPage) Build() fyne.CanvasObject {
	if p.cachedView != nil {
		return p.cachedView
	}

	hero := components.NewHeroHeader(
		constants.NavDatabase,
		"Perancang dan generator skema database lengkap untuk MySQL, MariaDB, PostgreSQL, SQLite, SQL Server, dan Oracle. Dilengkapi perpustakaan skema siap pakai, katalog perintah DDL, dan kamus tipe data.",
		components.BadgeCyan("DATABASE ARCHITECT & DDL GENERATOR"),
	)

	tab1 := container.NewTabItemWithIcon(constants.TabDBSchemas, theme.FolderOpenIcon(), p.buildSchemasTab())

	tabCatalogBox := container.NewStack()
	tabGeneratorBox := container.NewStack()
	tabDataTypesBox := container.NewStack()

	tab2 := container.NewTabItemWithIcon(constants.TabDBCatalog, theme.ListIcon(), tabCatalogBox)
	tab3 := container.NewTabItemWithIcon(constants.TabDBGenerator, theme.DocumentCreateIcon(), tabGeneratorBox)
	tab4 := container.NewTabItemWithIcon(constants.TabDBDataTypes, theme.InfoIcon(), tabDataTypesBox)

	tabs := container.NewAppTabs(tab1, tab2, tab3, tab4)
	tabs.OnSelected = func(t *container.TabItem) {
		switch t {
		case tab2:
			if len(tabCatalogBox.Objects) == 0 {
				tabCatalogBox.Objects = []fyne.CanvasObject{p.buildCatalogTab()}
				tabCatalogBox.Refresh()
			}
		case tab3:
			if len(tabGeneratorBox.Objects) == 0 {
				tabGeneratorBox.Objects = []fyne.CanvasObject{p.buildGeneratorTab()}
				tabGeneratorBox.Refresh()
			}
		case tab4:
			if len(tabDataTypesBox.Objects) == 0 {
				tabDataTypesBox.Objects = []fyne.CanvasObject{p.buildDataTypesTab()}
				tabDataTypesBox.Refresh()
			}
		}
	}

	p.cachedView = container.NewBorder(hero, nil, nil, nil, tabs)
	return p.cachedView
}

// ----------------------------------------------------------------------------
// TAB 1: PERPUSTAKAAN SKEMA SIAP PAKAI
// ----------------------------------------------------------------------------
func (p *DBSchemaPage) buildSchemasTab() fyne.CanvasObject {
	templates := dbschema.GetSchemaTemplates()
	if len(templates) == 0 {
		return widget.NewLabel("Belum ada template skema.")
	}

	selectedIdx := 0
	selectedEngine := dbschema.EngineMySQL

	var titles []string
	for _, t := range templates {
		titles = append(titles, t.Title)
	}

	dbNameEntry := components.NewScrollableEntry()
	dbNameEntry.SetText("db_aplikasi_utama")

	prefixEntry := components.NewScrollableEntry()
	prefixEntry.SetText("")
	prefixEntry.SetPlaceHolder("cth: tbl_")

	tabScroll := container.NewVScroll(nil)
	sqlPreview := components.NewScrollableMultiLineEntry(tabScroll)
	sqlPreview.TextStyle = fyne.TextStyle{Monospace: true}
	sqlPreview.Wrapping = fyne.TextWrapOff

	tableBadgesBox := container.NewHBox()
	descLabel := widget.NewLabel("")
	descLabel.Wrapping = fyne.TextWrapWord

	notesLabel := widget.NewLabel("")
	notesLabel.Wrapping = fyne.TextWrapWord

	refreshPreview := func() {
		tmpl := templates[selectedIdx]
		descLabel.SetText(tmpl.Description)
		notesLabel.SetText("Catatan: " + tmpl.Notes)

		tableBadgesBox.Objects = nil
		for _, tbl := range tmpl.Tables {
			tableBadgesBox.Add(components.BadgeIndigo(tbl))
		}
		tableBadgesBox.Refresh()

		dbName := strings.TrimSpace(dbNameEntry.Text)
		prefix := strings.TrimSpace(prefixEntry.Text)
		rendered := tmpl.RenderScript(selectedEngine, dbName, prefix)
		sqlPreview.SetText(rendered)
	}

	topicSelect := widget.NewSelect(titles, func(selected string) {
		for i, t := range templates {
			if t.Title == selected {
				selectedIdx = i
				break
			}
		}
		refreshPreview()
	})
	topicSelect.SetSelectedIndex(0)

	engineSelect := widget.NewSelect([]string{"MySQL", "MariaDB", "PostgreSQL", "SQLite"}, func(val string) {
		switch val {
		case "PostgreSQL":
			selectedEngine = dbschema.EnginePostgres
		case "SQLite":
			selectedEngine = dbschema.EngineSQLite
		case "MariaDB":
			selectedEngine = dbschema.EngineMariaDB
		default:
			selectedEngine = dbschema.EngineMySQL
		}
		refreshPreview()
	})
	engineSelect.SetSelected("MySQL")

	dbNameEntry.OnChanged = func(_ string) { refreshPreview() }
	prefixEntry.OnChanged = func(_ string) { refreshPreview() }

	btnCopy := widget.NewButtonWithIcon("Salin SQL Penuh", theme.ContentCopyIcon(), func() {
		p.copyToClip(sqlPreview.Text)
	})
	btnCopy.Importance = widget.HighImportance

	btnExport := widget.NewButtonWithIcon("Ekspor File (.sql)", theme.DocumentSaveIcon(), func() {
		tmpl := templates[selectedIdx]
		filename := fmt.Sprintf("%s_%s.sql", tmpl.ID, strings.ToLower(string(selectedEngine)))
		p.exportSQLFile(filename, sqlPreview.Text)
	})
	btnExport.Importance = widget.MediumImportance

	controlRow1 := container.NewGridWithColumns(2,
		container.NewVBox(canvas.NewText("Pilih Topik Skema:", constants.ColorTextPrimary), topicSelect),
		container.NewVBox(canvas.NewText("Target Engine Database:", constants.ColorTextPrimary), engineSelect),
	)
	controlRow2 := container.NewGridWithColumns(2,
		container.NewVBox(canvas.NewText("Nama Basis Data:", constants.ColorTextPrimary), dbNameEntry),
		container.NewVBox(canvas.NewText("Prefix Tabel (Opsional):", constants.ColorTextPrimary), prefixEntry),
	)
	controlGrid := container.NewVBox(controlRow1, controlRow2)

	actionBar := container.NewBorder(nil, nil, nil, container.NewHBox(btnExport, btnCopy))

	headerPanel := container.NewVBox(
		components.NewPlainCardWithAccent(controlGrid, constants.ColorAccentCobalt),
		container.NewPadded(container.NewVBox(
			descLabel,
			container.NewHBox(canvas.NewText("Tabel Terintegrasi:", constants.ColorTextMuted), tableBadgesBox),
			notesLabel,
		)),
		widget.NewSeparator(),
		actionBar,
	)

	refreshPreview()

	previewSpacer := canvas.NewRectangle(color.Transparent)
	previewSpacer.SetMinSize(fyne.NewSize(0, 320))
	previewBox := container.NewStack(previewSpacer, container.NewPadded(sqlPreview))

	tabScroll.Content = container.NewVBox(headerPanel, previewBox)
	return tabScroll
}

// ----------------------------------------------------------------------------
// TAB 2: KATALOG PERINTAH DDL & ADMINISTRASI
// ----------------------------------------------------------------------------
func (p *DBSchemaPage) buildCatalogTab() fyne.CanvasObject {
	listContainer := container.NewVBox()
	scroll := container.NewVScroll(listContainer)

	renderList := func() {
		listContainer.Objects = nil
		items := dbschema.SearchCommands(p.selectedEngine, p.selectedCategory, p.searchQuery)
		if len(items) == 0 {
			listContainer.Add(components.NewPlainCard(widget.NewLabel(constants.NoDataMessage)))
			listContainer.Refresh()
			return
		}

		for _, item := range items {
			cmd := item

			// Header items
			titleLabel := widget.NewLabel(cmd.Title)
			titleLabel.Wrapping = fyne.TextWrapWord
			titleLabel.TextStyle = fyne.TextStyle{Bold: true}

			engineBadge := components.BadgeCyan(string(cmd.Engine))
			catBadge := components.BadgeYellow(string(cmd.Category))

			btnCopy := widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() {
				p.copyToClip(cmd.RenderCommand(nil))
			})
			btnCopy.Importance = widget.LowImportance

			btnCustomize := widget.NewButtonWithIcon("Sesuaikan", theme.SettingsIcon(), func() {
				p.showParameterModal(cmd)
			})
			btnCustomize.Importance = widget.MediumImportance

			actionBox := container.NewHBox(btnCustomize, btnCopy)
			topInfo := container.NewVBox(container.NewHBox(engineBadge, catBadge), titleLabel)
			header := container.NewBorder(nil, nil, topInfo, actionBox)

			descLabel := widget.NewLabel(cmd.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			codeBox := components.NewScrollableMultiLineEntry(scroll)
			codeBox.SetText(cmd.RenderCommand(nil))
			codeBox.TextStyle = fyne.TextStyle{Monospace: true}
			codeBox.Wrapping = fyne.TextWrapOff

			notesLabel := widget.NewLabel("Petunjuk: " + cmd.Notes)
			notesLabel.Wrapping = fyne.TextWrapWord

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				descLabel,
				container.NewPadded(codeBox),
				notesLabel,
			)

			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorAccentCobalt))
		}
		listContainer.Refresh()
	}

	searchBar := components.NewSearchBar(constants.SearchDBCatalogPlaceholder, func(q string) {
		p.searchQuery = q
		renderList()
	})

	engineSelect := widget.NewSelect([]string{
		string(dbschema.EngineAll),
		string(dbschema.EngineMySQL),
		string(dbschema.EngineMariaDB),
		string(dbschema.EnginePostgres),
		string(dbschema.EngineSQLite),
		string(dbschema.EngineSQLServer),
		string(dbschema.EngineOracle),
	}, func(val string) {
		p.selectedEngine = dbschema.DatabaseEngine(val)
		renderList()
	})
	engineSelect.SetSelected(string(dbschema.EngineAll))

	catSelect := widget.NewSelect([]string{
		string(dbschema.CatAll),
		string(dbschema.CatDatabaseSetup),
		string(dbschema.CatUserPrivileges),
		string(dbschema.CatTableDDL),
		string(dbschema.CatBackupRestore),
		string(dbschema.CatMaintenance),
	}, func(val string) {
		p.selectedCategory = dbschema.CommandCategory(val)
		renderList()
	})
	catSelect.SetSelected(string(dbschema.CatAll))

	filterBar := container.NewVBox(
		container.NewVBox(canvas.NewText("Pencarian Perintah DDL & Administrasi:", constants.ColorTextPrimary), searchBar.Container),
		container.NewGridWithColumns(2,
			container.NewVBox(canvas.NewText("Filter Engine Database:", constants.ColorTextPrimary), engineSelect),
			container.NewVBox(canvas.NewText("Filter Kategori:", constants.ColorTextPrimary), catSelect),
		),
	)

	renderList()
	return container.NewBorder(container.NewPadded(components.NewPlainCard(filterBar)), nil, nil, nil, scroll)
}

func (p *DBSchemaPage) showParameterModal(cmd dbschema.DBCommandItem) {
	var d dialog.Dialog
	values := make(map[string]string)
	var formRows []fyne.CanvasObject

	titleTxt := canvas.NewText("Sesuaikan Parameter: "+cmd.Title, constants.ColorTextPrimary)
	titleTxt.TextSize = constants.FontSizeH2
	titleTxt.TextStyle = fyne.TextStyle{Bold: true}

	header := container.NewVBox(
		titleTxt,
		canvas.NewText("Ubah nilai parameter sesuai kebutuhan sistem Anda lalu salin perintah SQL yang dihasilkan:", constants.ColorTextMuted),
		widget.NewSeparator(),
	)

	previewEntry := widget.NewMultiLineEntry()
	previewEntry.TextStyle = fyne.TextStyle{Monospace: true}
	previewEntry.Wrapping = fyne.TextWrapOff

	updatePreview := func() {
		rendered := cmd.RenderCommand(values)
		previewEntry.SetText(rendered)
	}

	for _, param := range cmd.Parameters {
		pm := param
		values[pm.Key] = pm.DefaultValue

		entry := widget.NewEntry()
		entry.SetText(pm.DefaultValue)
		entry.SetPlaceHolder(pm.Placeholder)
		entry.OnChanged = func(val string) {
			values[pm.Key] = val
			updatePreview()
		}

		lbl := canvas.NewText(pm.Label+" ("+pm.Key+"):", constants.ColorTextPrimary)
		lbl.TextStyle = fyne.TextStyle{Bold: true}
		formRows = append(formRows, container.NewVBox(lbl, entry))
	}

	updatePreview()

	btnCancel := widget.NewButtonWithIcon("Tutup", theme.CancelIcon(), func() {
		if d != nil {
			d.Hide()
		}
	})
	btnCancel.Importance = widget.LowImportance

	btnCopy := widget.NewButtonWithIcon("Salin Perintah", theme.ContentCopyIcon(), func() {
		p.copyToClip(previewEntry.Text)
		if d != nil {
			d.Hide()
		}
	})
	btnCopy.Importance = widget.HighImportance

	actionBar := container.NewBorder(nil, nil, nil, container.NewHBox(btnCancel, btnCopy))

	body := container.NewGridWithColumns(2,
		container.NewVScroll(container.NewVBox(formRows...)),
		container.NewVBox(canvas.NewText("Hasil Perintah yang Disesuaikan:", constants.ColorTextPrimary), container.NewPadded(previewEntry)),
	)

	dialogContent := container.NewBorder(
		header,
		container.NewVBox(widget.NewSeparator(), actionBar),
		nil,
		nil,
		body,
	)

	card := components.NewPlainCardWithAccent(dialogContent, constants.ColorAccentCobalt)
	d = dialog.NewCustomWithoutButtons("", card, p.window)
	d.Resize(fyne.NewSize(740, 480))
	d.Show()
}

// ----------------------------------------------------------------------------
// MODAL: PREVIEW DDL SQL
// ----------------------------------------------------------------------------
func (p *DBSchemaPage) showSQLPreviewModal(tableName, sqlText, engine string, colCount int) {
	var d dialog.Dialog

	title := canvas.NewText(fmt.Sprintf("Preview DDL SQL — %s", tableName), constants.ColorTextPrimary)
	title.TextSize = constants.FontSizeH2
	title.TextStyle = fyne.TextStyle{Bold: true}

	engineBadge := components.BadgeCyan(engine)
	colBadge := components.BadgeIndigo(fmt.Sprintf("%d Kolom Terdefinisi", colCount))

	headerBox := container.NewVBox(
		container.NewBorder(nil, nil, title, container.NewHBox(engineBadge, colBadge)),
		canvas.NewText("Skrip DDL SQL lengkap siap pakai untuk membuat struktur tabel pada database target.", constants.ColorTextMuted),
		widget.NewSeparator(),
	)

	sqlEntry := widget.NewMultiLineEntry()
	sqlEntry.SetText(sqlText)
	sqlEntry.TextStyle = fyne.TextStyle{Monospace: true}
	sqlEntry.Wrapping = fyne.TextWrapOff

	btnCancel := widget.NewButtonWithIcon("Tutup", theme.CancelIcon(), func() {
		if d != nil {
			d.Hide()
		}
	})
	btnCancel.Importance = widget.LowImportance

	btnExport := widget.NewButtonWithIcon("Ekspor File (.sql)", theme.DocumentSaveIcon(), func() {
		fname := strings.TrimSpace(tableName)
		if fname == "" {
			fname = "tabel_kustom"
		}
		p.exportSQLFile(fname+".sql", sqlText)
	})
	btnExport.Importance = widget.MediumImportance

	btnCopy := widget.NewButtonWithIcon("Salin DDL SQL", theme.ContentCopyIcon(), func() {
		p.copyToClip(sqlText)
	})
	btnCopy.Importance = widget.HighImportance

	lineCount := len(strings.Split(sqlText, "\n"))
	lineInfo := canvas.NewText(fmt.Sprintf("Total %d baris kode SQL DDL", lineCount), constants.ColorTextSecondary)
	lineInfo.TextSize = constants.FontSizeSmall

	actionBar := container.NewBorder(nil, nil, lineInfo, container.NewHBox(btnCancel, btnExport, btnCopy))

	dialogContent := container.NewBorder(
		headerBox,
		container.NewVBox(widget.NewSeparator(), actionBar),
		nil,
		nil,
		container.NewPadded(sqlEntry),
	)

	card := components.NewPlainCardWithAccent(dialogContent, constants.ColorAccentCobalt)
	d = dialog.NewCustomWithoutButtons("", card, p.window)
	d.Resize(fyne.NewSize(760, 520))
	d.Show()
}

// ----------------------------------------------------------------------------
// TAB 3: PEMBUAT TABEL KUSTOM (INTERACTIVE DDL BUILDER)
// ----------------------------------------------------------------------------
func (p *DBSchemaPage) buildGeneratorTab() fyne.CanvasObject {
	dbNameEntry := components.NewScrollableEntry()
	dbNameEntry.SetText(p.builderTableDef.DatabaseName)

	tableNameEntry := components.NewScrollableEntry()
	tableNameEntry.SetText(p.builderTableDef.TableName)

	engineSelect := widget.NewSelect([]string{
		string(dbschema.EngineMySQL),
		string(dbschema.EngineMariaDB),
		string(dbschema.EnginePostgres),
		string(dbschema.EngineSQLite),
		string(dbschema.EngineSQLServer),
		string(dbschema.EngineOracle),
	}, func(val string) {
		p.builderTableDef.Engine = dbschema.DatabaseEngine(val)
	})
	engineSelect.SetSelected(string(p.builderTableDef.Engine))

	dropCheck := widget.NewCheck("Sertakan DROP TABLE IF EXISTS", func(chk bool) {
		p.builderTableDef.IncludeDrop = chk
	})
	dropCheck.SetChecked(p.builderTableDef.IncludeDrop)

	columnsBox := container.NewVBox()

	currentSQL := ""
	updateSQL := func() {
		p.builderTableDef.DatabaseName = strings.TrimSpace(dbNameEntry.Text)
		p.builderTableDef.TableName = strings.TrimSpace(tableNameEntry.Text)
		currentSQL = dbschema.GenerateTableDDL(p.builderTableDef)
	}

	var renderColumns func()
	renderColumns = func() {
		columnsBox.Objects = nil
		for i, col := range p.builderTableDef.Columns {
			idx := i
			c := col

			nameEntry := components.NewScrollableEntry()
			nameEntry.SetText(c.Name)
			nameEntry.SetPlaceHolder("nama_kolom")
			nameEntry.OnChanged = func(val string) {
				p.builderTableDef.Columns[idx].Name = val
				updateSQL()
			}
			nameSpacer := canvas.NewRectangle(color.Transparent)
			nameSpacer.SetMinSize(fyne.NewSize(180, 36))
			nameBox := container.NewStack(nameSpacer, nameEntry)

			typeSelect := widget.NewSelect([]string{
				string(dbschema.TypeAutoID),
				string(dbschema.TypeUUID),
				string(dbschema.TypeVarchar),
				string(dbschema.TypeText),
				string(dbschema.TypeInteger),
				string(dbschema.TypeBigInt),
				string(dbschema.TypeDecimal),
				string(dbschema.TypeBoolean),
				string(dbschema.TypeDateTime),
				string(dbschema.TypeDate),
				string(dbschema.TypeJSON),
			}, func(val string) {
				p.builderTableDef.Columns[idx].Type = dbschema.ColumnType(val)
				updateSQL()
			})
			typeSelect.SetSelected(string(c.Type))

			lenEntry := components.NewScrollableEntry()
			lenEntry.SetText(c.Length)
			lenEntry.SetPlaceHolder("Panjang")
			lenEntry.OnChanged = func(val string) {
				p.builderTableDef.Columns[idx].Length = val
				updateSQL()
			}
			lenSpacer := canvas.NewRectangle(color.Transparent)
			lenSpacer.SetMinSize(fyne.NewSize(75, 36))
			lenBox := container.NewStack(lenSpacer, lenEntry)

			notNullChk := widget.NewCheck("Not Null", func(chk bool) {
				p.builderTableDef.Columns[idx].IsNotNull = chk
				updateSQL()
			})
			notNullChk.SetChecked(c.IsNotNull)

			uniqueChk := widget.NewCheck("Unique", func(chk bool) {
				p.builderTableDef.Columns[idx].IsUnique = chk
				updateSQL()
			})
			uniqueChk.SetChecked(c.IsUnique)

			btnDeleteCol := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				if len(p.builderTableDef.Columns) > 1 {
					p.builderTableDef.Columns = append(p.builderTableDef.Columns[:idx], p.builderTableDef.Columns[idx+1:]...)
					renderColumns()
					updateSQL()
				}
			})
			btnDeleteCol.Importance = widget.DangerImportance

			rightControls := container.NewHBox(lenBox, notNullChk, uniqueChk, btnDeleteCol)
			colRow := container.NewBorder(nil, nil, nameBox, rightControls, typeSelect)
			columnsBox.Add(container.NewVBox(colRow, widget.NewSeparator()))
		}
		columnsBox.Refresh()
		updateSQL()
	}

	dbNameEntry.OnChanged = func(_ string) { updateSQL() }
	tableNameEntry.OnChanged = func(_ string) { updateSQL() }
	engineSelect.OnChanged = func(_ string) { updateSQL() }
	dropCheck.OnChanged = func(_ bool) { updateSQL() }

	btnAddCol := widget.NewButtonWithIcon("Tambah Kolom", theme.ContentAddIcon(), func() {
		p.builderTableDef.Columns = append(p.builderTableDef.Columns, dbschema.ColumnDef{
			Name:      fmt.Sprintf("kolom_%d", len(p.builderTableDef.Columns)+1),
			Type:      dbschema.TypeVarchar,
			Length:    "255",
			IsNotNull: false,
		})
		renderColumns()
	})
	btnAddCol.Importance = widget.MediumImportance

	presetSelect := widget.NewSelect([]string{"Pilih Preset...", "Preset Standar", "Preset Pengguna", "Preset Produk"}, func(val string) {
		switch val {
		case "Preset Standar":
			p.builderTableDef.Columns = dbschema.PresetColumnsStandard()
			renderColumns()
		case "Preset Pengguna":
			p.builderTableDef.Columns = dbschema.PresetColumnsUser()
			renderColumns()
		case "Preset Produk":
			p.builderTableDef.Columns = dbschema.PresetColumnsProduct()
			renderColumns()
		}
	})
	presetSelect.SetSelected("Pilih Preset...")

	presetBar := container.NewHBox(
		presetSelect,
		btnAddCol,
	)

	btnPreviewSQL := widget.NewButtonWithIcon("Preview SQL", theme.VisibilityIcon(), func() {
		updateSQL()
		p.showSQLPreviewModal(tableNameEntry.Text, currentSQL, string(p.builderTableDef.Engine), len(p.builderTableDef.Columns))
	})
	btnPreviewSQL.Importance = widget.HighImportance

	btnCopy := widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() {
		updateSQL()
		p.copyToClip(currentSQL)
	})
	btnCopy.Importance = widget.MediumImportance

	btnExport := widget.NewButtonWithIcon("Ekspor", theme.DocumentSaveIcon(), func() {
		updateSQL()
		fname := strings.TrimSpace(tableNameEntry.Text)
		if fname == "" {
			fname = "tabel_kustom"
		}
		p.exportSQLFile(fname+".sql", currentSQL)
	})
	btnExport.Importance = widget.LowImportance

	topRow1 := container.NewGridWithColumns(2,
		container.NewVBox(canvas.NewText("Nama Basis Data:", constants.ColorTextPrimary), dbNameEntry),
		container.NewVBox(canvas.NewText("Nama Tabel:", constants.ColorTextPrimary), tableNameEntry),
	)
	topRow2 := container.NewGridWithColumns(2,
		container.NewVBox(canvas.NewText("Target Engine Database:", constants.ColorTextPrimary), engineSelect),
		container.NewVBox(canvas.NewText("Opsi Tambahan:", constants.ColorTextPrimary), dropCheck),
	)
	topSettings := container.NewVBox(topRow1, topRow2)

	actionsRow := container.NewBorder(nil, nil,
		presetBar,
		container.NewHBox(btnPreviewSQL, btnCopy, btnExport),
	)

	tableHeader := container.NewBorder(nil, nil,
		canvas.NewText("Nama Kolom (Field)", constants.ColorTextPrimary),
		canvas.NewText("Panjang, Nullability, Unique & Hapus", constants.ColorTextPrimary),
		container.NewCenter(canvas.NewText("Tipe Data", constants.ColorTextPrimary)),
	)

	topHeader := container.NewVBox(
		components.NewPlainCardWithAccent(topSettings, constants.ColorAccentCobalt),
		actionsRow,
		widget.NewSeparator(),
		tableHeader,
		widget.NewSeparator(),
	)

	renderColumns()
	updateSQL()

	return container.NewPadded(
		container.NewBorder(
			topHeader,
			nil,
			nil,
			nil,
			container.NewVScroll(columnsBox),
		),
	)
}
// ----------------------------------------------------------------------------
// TAB 4: KAMUS & KOMPARASI TIPE DATA
// ----------------------------------------------------------------------------
func (p *DBSchemaPage) buildDataTypesTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

	renderList := func(q string) {
		listContainer.Objects = nil
		items := dbschema.SearchDataTypes(q)
		if len(items) == 0 {
			listContainer.Add(components.NewPlainCard(widget.NewLabel(constants.NoDataMessage)))
			listContainer.Refresh()
			return
		}

		for _, item := range items {
			it := item

			titleTxt := canvas.NewText(it.LogicalType, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeBody
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			badge := components.BadgeIndigo("KOMPARASI TIPE")
			header := container.NewBorder(nil, nil, container.NewHBox(titleTxt, badge), nil)

			descLabel := widget.NewLabel(it.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			makeVal := func(txt string) *widget.Label {
				lbl := widget.NewLabel(txt)
				lbl.Wrapping = fyne.TextWrapWord
				return lbl
			}
			grid := container.NewGridWithColumns(2,
				container.NewVBox(canvas.NewText("MySQL / MariaDB:", constants.ColorTextMuted), makeVal(it.MySQL)),
				container.NewVBox(canvas.NewText("PostgreSQL:", constants.ColorTextMuted), makeVal(it.PostgreSQL)),
				container.NewVBox(canvas.NewText("SQLite:", constants.ColorTextMuted), makeVal(it.SQLite)),
				container.NewVBox(canvas.NewText("SQL Server (T-SQL):", constants.ColorTextMuted), makeVal(it.SQLServer)),
				container.NewVBox(canvas.NewText("Oracle Database:", constants.ColorTextMuted), makeVal(it.Oracle)),
				container.NewVBox(canvas.NewText("Rekomendasi Penggunaan:", constants.ColorTextMuted), makeVal(it.BestUse)),
			)

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				descLabel,
				grid,
			)

			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorSuccess))
		}
		listContainer.Refresh()
	}

	searchBar := components.NewSearchBar(constants.SearchDBDataTypesPlaceholder, renderList)
	renderList("")

	scroll := container.NewVScroll(listContainer)
	return container.NewBorder(container.NewPadded(components.NewPlainCard(searchBar.Container)), nil, nil, nil, scroll)
}
