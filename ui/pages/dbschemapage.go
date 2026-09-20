package pages

import (
	"fmt"
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
	hero := components.NewHeroHeader(
		constants.NavDatabase,
		"Perancang dan generator skema database lengkap untuk MySQL, MariaDB, PostgreSQL, SQLite, SQL Server, dan Oracle. Dilengkapi perpustakaan skema siap pakai, katalog perintah DDL, dan kamus tipe data.",
		components.BadgeCyan("DATABASE ARCHITECT & DDL GENERATOR"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabDBSchemas, theme.FolderOpenIcon(), p.buildSchemasTab()),
		container.NewTabItemWithIcon(constants.TabDBCatalog, theme.ListIcon(), p.buildCatalogTab()),
		container.NewTabItemWithIcon(constants.TabDBGenerator, theme.DocumentCreateIcon(), p.buildGeneratorTab()),
		container.NewTabItemWithIcon(constants.TabDBDataTypes, theme.InfoIcon(), p.buildDataTypesTab()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
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

	dbNameEntry := widget.NewEntry()
	dbNameEntry.SetText("db_aplikasi_utama")

	prefixEntry := widget.NewEntry()
	prefixEntry.SetText("")
	prefixEntry.SetPlaceHolder("cth: tbl_")

	sqlPreview := widget.NewMultiLineEntry()
	sqlPreview.TextStyle = fyne.TextStyle{Monospace: true}
	sqlPreview.Wrapping = fyne.TextWrapOff

	tableBadgesBox := container.NewHBox()
	descLabel := widget.NewLabel("")
	descLabel.Wrapping = fyne.TextWrapWord

	notesLabel := canvas.NewText("", constants.ColorTextSecondary)
	notesLabel.TextSize = constants.FontSizeSmall

	refreshPreview := func() {
		tmpl := templates[selectedIdx]
		descLabel.SetText(tmpl.Description)
		notesLabel.Text = "Catatan: " + tmpl.Notes
		notesLabel.Refresh()

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

	controlGrid := container.NewGridWithColumns(4,
		container.NewVBox(canvas.NewText("Pilih Topik Skema:", constants.ColorTextPrimary), topicSelect),
		container.NewVBox(canvas.NewText("Target Engine:", constants.ColorTextPrimary), engineSelect),
		container.NewVBox(canvas.NewText("Nama Basis Data:", constants.ColorTextPrimary), dbNameEntry),
		container.NewVBox(canvas.NewText("Prefix Tabel (Opsional):", constants.ColorTextPrimary), prefixEntry),
	)

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

	return container.NewBorder(headerPanel, nil, nil, nil, container.NewPadded(sqlPreview))
}

// ----------------------------------------------------------------------------
// TAB 2: KATALOG PERINTAH DDL & ADMINISTRASI
// ----------------------------------------------------------------------------
func (p *DBSchemaPage) buildCatalogTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

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
			titleTxt := canvas.NewText(cmd.Title, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeBody
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			engineBadge := components.BadgeCyan(string(cmd.Engine))
			catBadge := components.BadgeYellow(string(cmd.Category))

			btnCopy := widget.NewButtonWithIcon("Salin", theme.ContentCopyIcon(), func() {
				p.copyToClip(cmd.RenderCommand(nil))
			})
			btnCopy.Importance = widget.LowImportance

			btnCustomize := widget.NewButtonWithIcon("Sesuaikan Parameter", theme.SettingsIcon(), func() {
				p.showParameterModal(cmd)
			})
			btnCustomize.Importance = widget.MediumImportance

			actionBox := container.NewHBox(btnCustomize, btnCopy)
			header := container.NewBorder(nil, nil,
				container.NewHBox(engineBadge, catBadge, titleTxt),
				actionBox,
			)

			descLabel := widget.NewLabel(cmd.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			codeBox := widget.NewMultiLineEntry()
			codeBox.SetText(cmd.RenderCommand(nil))
			codeBox.TextStyle = fyne.TextStyle{Monospace: true}
			codeBox.Wrapping = fyne.TextWrapOff

			notesTxt := canvas.NewText("Petunjuk: "+cmd.Notes, constants.ColorTextSecondary)
			notesTxt.TextSize = constants.FontSizeSmall

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				descLabel,
				container.NewPadded(codeBox),
				notesTxt,
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

	filterBar := container.NewGridWithColumns(3,
		container.NewVBox(canvas.NewText("Pencarian Perintah:", constants.ColorTextPrimary), searchBar.Container),
		container.NewVBox(canvas.NewText("Filter Engine Database:", constants.ColorTextPrimary), engineSelect),
		container.NewVBox(canvas.NewText("Filter Kategori:", constants.ColorTextPrimary), catSelect),
	)

	renderList()
	scroll := container.NewVScroll(listContainer)
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
// TAB 3: PEMBUAT TABEL KUSTOM (INTERACTIVE DDL BUILDER)
// ----------------------------------------------------------------------------
func (p *DBSchemaPage) buildGeneratorTab() fyne.CanvasObject {
	dbNameEntry := widget.NewEntry()
	dbNameEntry.SetText(p.builderTableDef.DatabaseName)

	tableNameEntry := widget.NewEntry()
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

	sqlOutput := widget.NewMultiLineEntry()
	sqlOutput.TextStyle = fyne.TextStyle{Monospace: true}
	sqlOutput.Wrapping = fyne.TextWrapOff

	columnsBox := container.NewVBox()

	updateSQL := func() {
		p.builderTableDef.DatabaseName = strings.TrimSpace(dbNameEntry.Text)
		p.builderTableDef.TableName = strings.TrimSpace(tableNameEntry.Text)
		sql := dbschema.GenerateTableDDL(p.builderTableDef)
		sqlOutput.SetText(sql)
	}

	var renderColumns func()
	renderColumns = func() {
		columnsBox.Objects = nil
		for i, col := range p.builderTableDef.Columns {
			idx := i
			c := col

			nameEntry := widget.NewEntry()
			nameEntry.SetText(c.Name)
			nameEntry.SetPlaceHolder("nama_kolom")
			nameEntry.OnChanged = func(val string) {
				p.builderTableDef.Columns[idx].Name = val
				updateSQL()
			}

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

			lenEntry := widget.NewEntry()
			lenEntry.SetText(c.Length)
			lenEntry.SetPlaceHolder("Panjang (255)")
			lenEntry.OnChanged = func(val string) {
				p.builderTableDef.Columns[idx].Length = val
				updateSQL()
			}

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

			colRow := container.NewGridWithColumns(6,
				nameEntry,
				typeSelect,
				lenEntry,
				notNullChk,
				uniqueChk,
				btnDeleteCol,
			)
			columnsBox.Add(colRow)
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

	btnPresetStandard := widget.NewButton("Preset Standar", func() {
		p.builderTableDef.Columns = dbschema.PresetColumnsStandard()
		renderColumns()
	})
	btnPresetUser := widget.NewButton("Preset Akun Pengguna", func() {
		p.builderTableDef.Columns = dbschema.PresetColumnsUser()
		renderColumns()
	})
	btnPresetProduct := widget.NewButton("Preset Produk E-Commerce", func() {
		p.builderTableDef.Columns = dbschema.PresetColumnsProduct()
		renderColumns()
	})

	presetBar := container.NewHBox(
		canvas.NewText("Preset Kolom Cepat:", constants.ColorTextMuted),
		btnPresetStandard,
		btnPresetUser,
		btnPresetProduct,
		btnAddCol,
	)

	btnCopy := widget.NewButtonWithIcon("Salin DDL SQL", theme.ContentCopyIcon(), func() {
		p.copyToClip(sqlOutput.Text)
	})
	btnCopy.Importance = widget.HighImportance

	btnExport := widget.NewButtonWithIcon("Ekspor .sql", theme.DocumentSaveIcon(), func() {
		fname := strings.TrimSpace(tableNameEntry.Text)
		if fname == "" {
			fname = "tabel_kustom"
		}
		p.exportSQLFile(fname+".sql", sqlOutput.Text)
	})
	btnExport.Importance = widget.MediumImportance

	topSettings := container.NewGridWithColumns(4,
		container.NewVBox(canvas.NewText("Nama Basis Data:", constants.ColorTextPrimary), dbNameEntry),
		container.NewVBox(canvas.NewText("Nama Tabel:", constants.ColorTextPrimary), tableNameEntry),
		container.NewVBox(canvas.NewText("Target Engine:", constants.ColorTextPrimary), engineSelect),
		container.NewVBox(canvas.NewText("Opsi Tambahan:", constants.ColorTextPrimary), dropCheck),
	)

	leftPanel := container.NewBorder(
		container.NewVBox(
			topSettings,
			widget.NewSeparator(),
			presetBar,
			widget.NewSeparator(),
			container.NewGridWithColumns(6,
				canvas.NewText("Nama Kolom", constants.ColorTextPrimary),
				canvas.NewText("Tipe Data", constants.ColorTextPrimary),
				canvas.NewText("Panjang", constants.ColorTextPrimary),
				canvas.NewText("Nullability", constants.ColorTextPrimary),
				canvas.NewText("Unique", constants.ColorTextPrimary),
				canvas.NewText("Aksi", constants.ColorTextPrimary),
			),
		),
		nil,
		nil,
		nil,
		container.NewVScroll(columnsBox),
	)

	rightPanel := container.NewBorder(
		container.NewBorder(nil, nil, canvas.NewText("Preview DDL SQL:", constants.ColorTextPrimary), container.NewHBox(btnExport, btnCopy)),
		nil,
		nil,
		nil,
		container.NewPadded(sqlOutput),
	)

	renderColumns()
	updateSQL()

	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.55)
	return container.NewPadded(split)
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

			grid := container.NewGridWithColumns(3,
				container.NewVBox(canvas.NewText("MySQL / MariaDB:", constants.ColorTextMuted), widget.NewLabel(it.MySQL)),
				container.NewVBox(canvas.NewText("PostgreSQL:", constants.ColorTextMuted), widget.NewLabel(it.PostgreSQL)),
				container.NewVBox(canvas.NewText("SQLite:", constants.ColorTextMuted), widget.NewLabel(it.SQLite)),
				container.NewVBox(canvas.NewText("SQL Server (T-SQL):", constants.ColorTextMuted), widget.NewLabel(it.SQLServer)),
				container.NewVBox(canvas.NewText("Oracle Database:", constants.ColorTextMuted), widget.NewLabel(it.Oracle)),
				container.NewVBox(canvas.NewText("Rekomendasi Penggunaan:", constants.ColorTextMuted), widget.NewLabel(it.BestUse)),
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
