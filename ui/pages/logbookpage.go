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

	"github.com/yudz/it-toolbox/database"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

type LogbookPage struct {
	window fyne.Window
}

func NewLogbookPage(win fyne.Window) *LogbookPage {
	return &LogbookPage{window: win}
}

func (p *LogbookPage) copyToClip(txt string) {
	if p.window != nil {
		p.window.Clipboard().SetContent(txt)
		dialog.ShowInformation("Clipboard", constants.StatusCopied, p.window)
	}
}

func (p *LogbookPage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavLogbook,
		"Pencatatan Error Logbook insiden debugging, repositori Code Snippet siap pakai, dan Checklist verifikasi.",
		components.BadgeCyan("LOGBOOK & KNOWLEDGE"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabErrorLog, theme.WarningIcon(), p.buildErrorLogTab()),
		container.NewTabItemWithIcon(constants.TabSnippets, theme.DocumentCreateIcon(), p.buildSnippetsTab()),
		container.NewTabItemWithIcon(constants.TabChecklists, theme.CheckButtonCheckedIcon(), p.buildChecklistsTab()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
}

// 1. Tab Error Logbook
func (p *LogbookPage) buildErrorLogTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

	// Form entries for inline Create Error Log
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("cth: NullPointerException saat query database / Connection refused port 8080")

	tagsEntry := widget.NewEntry()
	tagsEntry.SetPlaceHolder("cth: go, sql, cisco, docker, spring")

	errEntry := widget.NewMultiLineEntry()
	errEntry.SetPlaceHolder("Salin pesan error, panic output, atau stack trace di sini...")
	errEntry.SetMinRowsVisible(3)
	errEntry.TextStyle = fyne.TextStyle{Monospace: true}

	solEntry := widget.NewMultiLineEntry()
	solEntry.SetPlaceHolder("Tulis langkah perbaikan atau solusi teknis yang berhasil mengatasi error...")
	solEntry.SetMinRowsVisible(3)

	var formCard fyne.CanvasObject
	var reloadLogs func(query string)

	reloadLogs = func(query string) {
		listContainer.Objects = nil
		var logs []database.ErrorLog
		var err error

		if query == "" {
			logs, err = database.GetAllLogs()
		} else {
			logs, err = database.SearchLogs(query)
		}

		if err != nil {
			listContainer.Add(widget.NewLabel(fmt.Sprintf("Error memuat log: %v", err)))
			listContainer.Refresh()
			return
		}

		if len(logs) == 0 {
			emptyTitle := canvas.NewText("Belum Ada Catatan Error Logbook", constants.ColorTextPrimary)
			emptyTitle.TextSize = constants.FontSizeH2
			emptyTitle.TextStyle = fyne.TextStyle{Bold: true}

			emptySub := canvas.NewText("Dokumentasikan pesan error dan solusi yang Anda temukan agar mudah dicari kembali di masa depan.", constants.ColorTextMuted)
			emptySub.TextSize = constants.FontSizeBody

			emptyBadge := components.BadgeYellow("KNOWLEDGE BASE")
			emptyHeader := container.NewHBox(emptyTitle, emptyBadge)

			emptyBox := container.NewVBox(
				emptyHeader,
				emptySub,
			)
			listContainer.Add(components.NewPlainCardWithAccent(emptyBox, constants.ColorDanger))
			listContainer.Refresh()
			return
		}

		for _, l := range logs {
			item := l

			delBtn := widget.NewButtonWithIcon("Hapus", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Log", fmt.Sprintf("Yakin ingin menghapus catatan error %q?", item.Title), func(ok bool) {
					if ok {
						_ = database.DeleteLog(item.ID)
						reloadLogs("")
					}
				}, p.window)
			})
			delBtn.Importance = widget.LowImportance

			editBtn := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
				eTitle := widget.NewEntry()
				eTitle.SetText(item.Title)
				eTags := widget.NewEntry()
				eTags.SetText(item.Tags)
				eErr := widget.NewMultiLineEntry()
				eErr.SetText(item.ErrorMessage)
				eErr.SetMinRowsVisible(3)
				eErr.TextStyle = fyne.TextStyle{Monospace: true}
				eSol := widget.NewMultiLineEntry()
				eSol.SetText(item.Solution)
				eSol.SetMinRowsVisible(3)

				lblT := canvas.NewText("JUDUL MASALAH", constants.ColorTextPrimary)
				lblT.TextSize = constants.FontSizeLabel
				lblT.TextStyle = fyne.TextStyle{Bold: true}

				lblTags := canvas.NewText("TAGS (DIPISAH KOMA)", constants.ColorTextPrimary)
				lblTags.TextSize = constants.FontSizeLabel
				lblTags.TextStyle = fyne.TextStyle{Bold: true}

				lblE := canvas.NewText("PESAN ERROR / STACK TRACE", constants.ColorTextPrimary)
				lblE.TextSize = constants.FontSizeLabel
				lblE.TextStyle = fyne.TextStyle{Bold: true}

				lblS := canvas.NewText("SOLUSI TERVERIFIKASI", constants.ColorTextPrimary)
				lblS.TextSize = constants.FontSizeLabel
				lblS.TextStyle = fyne.TextStyle{Bold: true}

				fContent := container.NewVBox(
					lblT, eTitle,
					lblTags, eTags,
					lblE, eErr,
					lblS, eSol,
				)

				components.ShowBrutalistFormDialog(
					p.window,
					"EDIT ERROR LOG",
					constants.ColorDanger,
					"Edit Catatan Insiden Error",
					"Perbarui pesan error, solusi teknis, atau tag pencarian",
					fContent,
					"Simpan Perubahan",
					func() {
						if eTitle.Text != "" {
							_ = database.UpdateLog(database.ErrorLog{
								ID:           item.ID,
								Title:        eTitle.Text,
								ErrorMessage: eErr.Text,
								Solution:     eSol.Text,
								Tags:         eTags.Text,
							})
							reloadLogs("")
						}
					},
				)
			})
			editBtn.Importance = widget.LowImportance

			copySolBtn := widget.NewButtonWithIcon("Salin Solusi", theme.ContentCopyIcon(), func() {
				p.copyToClip(item.Solution)
			})
			copySolBtn.Importance = widget.LowImportance

			titleTxt := canvas.NewText(item.Title, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeH2
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			errBadge := components.BadgeDanger("ERROR LOG")
			tagBadge := components.BadgeYellow(item.Tags)
			dateBadge := components.BadgeMuted(item.CreatedAt)

			headerLeft := container.NewHBox(titleTxt, errBadge)
			if item.Tags != "" {
				headerLeft.Add(tagBadge)
			}
			headerLeft.Add(dateBadge)

			headerRight := container.NewHBox(copySolBtn, editBtn, delBtn)
			header := container.NewBorder(nil, nil, headerLeft, headerRight)

			// Pesan Error Box (Monospace Dark Panel)
			errBg := canvas.NewRectangle(constants.ColorBgCardInner)
			errBg.StrokeColor = constants.ColorDanger
			errBg.StrokeWidth = 1.5
			errBg.CornerRadius = constants.CornerRadiusBrutal

			errLabel := widget.NewLabel(item.ErrorMessage)
			errLabel.Wrapping = fyne.TextWrapWord
			errLabel.TextStyle = fyne.TextStyle{Monospace: true}

			copyErrBtn := widget.NewButtonWithIcon("Salin Error", theme.ContentCopyIcon(), func() {
				p.copyToClip(item.ErrorMessage)
			})
			copyErrBtn.Importance = widget.LowImportance

			errTop := container.NewBorder(nil, nil,
				components.BadgeDanger("🚨 PESAN ERROR / STACK TRACE"),
				copyErrBtn,
			)
			errBox := container.NewStack(errBg, container.NewPadded(container.NewVBox(errTop, errLabel)))

			// Solusi Box (Green/Mint Panel)
			solBg := canvas.NewRectangle(constants.ColorBgCardInner)
			solBg.StrokeColor = constants.ColorSuccess
			solBg.StrokeWidth = 1.5
			solBg.CornerRadius = constants.CornerRadiusBrutal

			solLabel := widget.NewLabel(item.Solution)
			solLabel.Wrapping = fyne.TextWrapWord

			solTop := container.NewHBox(components.BadgeSuccess("💡 SOLUSI TERVERIFIKASI"))
			solBox := container.NewStack(solBg, container.NewPadded(container.NewVBox(solTop, solLabel)))

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				errBox,
				solBox,
			)
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorDanger))
		}
		listContainer.Refresh()
	}

	// Inline Neo-Brutalist Form Card for Adding New Error Log
	formTitle := canvas.NewText("🚨 CATAT INSIDEN ERROR BARU", constants.ColorTextPrimary)
	formTitle.TextSize = constants.FontSizeH2
	formTitle.TextStyle = fyne.TextStyle{Bold: true}

	formBadge := components.BadgeDanger("DEBUGGING LOG")
	formHeader := container.NewHBox(formTitle, formBadge)

	formSub := canvas.NewText("Dokumentasikan kendala error dan solusinya agar menjadi aset pengetahuan mandiri Anda.", constants.ColorTextMuted)
	formSub.TextSize = constants.FontSizeSmall

	lblT := canvas.NewText("JUDUL MASALAH / ERROR", constants.ColorTextPrimary)
	lblT.TextSize = constants.FontSizeLabel
	lblT.TextStyle = fyne.TextStyle{Bold: true}

	lblTags := canvas.NewText("TAGS / KATA KUNCI (DIPISAH KOMA)", constants.ColorTextPrimary)
	lblTags.TextSize = constants.FontSizeLabel
	lblTags.TextStyle = fyne.TextStyle{Bold: true}

	lblE := canvas.NewText("PESAN ERROR (OUTPUT TERMINAL / STACK TRACE)", constants.ColorTextPrimary)
	lblE.TextSize = constants.FontSizeLabel
	lblE.TextStyle = fyne.TextStyle{Bold: true}

	lblS := canvas.NewText("SOLUSI TERVERIFIKASI / CARA MEMPERBAIKI", constants.ColorTextPrimary)
	lblS.TextSize = constants.FontSizeLabel
	lblS.TextStyle = fyne.TextStyle{Bold: true}

	btnCancel := widget.NewButtonWithIcon("Batal", theme.CancelIcon(), func() {
		if formCard != nil {
			formCard.Hide()
		}
	})
	btnCancel.Importance = widget.LowImportance

	btnSave := widget.NewButtonWithIcon("Simpan Catatan Error", theme.DocumentSaveIcon(), func() {
		if titleEntry.Text != "" {
			_ = database.CreateLog(database.ErrorLog{
				Title:        titleEntry.Text,
				ErrorMessage: errEntry.Text,
				Solution:     solEntry.Text,
				Tags:         tagsEntry.Text,
			})
			titleEntry.SetText("")
			tagsEntry.SetText("")
			errEntry.SetText("")
			solEntry.SetText("")
			if formCard != nil {
				formCard.Hide()
			}
			reloadLogs("")
		}
	})
	btnSave.Importance = widget.HighImportance

	formActions := container.NewBorder(nil, nil, nil, container.NewHBox(btnCancel, btnSave))

	formInner := container.NewVBox(
		formHeader,
		formSub,
		widget.NewSeparator(),
		lblT, titleEntry,
		lblTags, tagsEntry,
		lblE, errEntry,
		lblS, solEntry,
		widget.NewSeparator(),
		formActions,
	)
	formCard = components.NewPlainCardWithAccent(formInner, constants.ColorDanger)
	formCard.Hide()

	searchBar := components.NewSearchBar(constants.SearchPlaceholder, reloadLogs)

	toggleAddBtn := widget.NewButtonWithIcon("+ Catat Error Baru", theme.ContentAddIcon(), func() {
		if formCard.Visible() {
			formCard.Hide()
		} else {
			formCard.Show()
		}
	})
	toggleAddBtn.Importance = widget.HighImportance

	topBar := container.NewBorder(nil, nil, nil, toggleAddBtn, searchBar.Container)
	reloadLogs("")

	mainContent := container.NewVBox(
		container.NewPadded(topBar),
		formCard,
		listContainer,
	)

	return container.NewVScroll(container.NewPadded(mainContent))
}

// 2. Tab Code Snippets
func (p *LogbookPage) buildSnippetsTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

	// Form inputs for inline Create Snippet Card
	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("cth: Database Connection Pool / Router OSPF Baseline Script")

	langEntry := widget.NewEntry()
	langEntry.SetPlaceHolder("cth: go, bash, python, sql, cisco cli")

	descEntry := widget.NewEntry()
	descEntry.SetPlaceHolder("Deskripsi singkat kegunaan snippet ini...")

	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetPlaceHolder("Tulis atau tempel kode script siap pakai di sini...")
	contentEntry.SetMinRowsVisible(6)
	contentEntry.TextStyle = fyne.TextStyle{Monospace: true}

	var formCard fyne.CanvasObject
	var reloadSnippets func()

	reloadSnippets = func() {
		listContainer.Objects = nil
		snippets, err := database.GetAllSnippets()
		if err != nil {
			listContainer.Add(widget.NewLabel(fmt.Sprintf("Error memuat snippets: %v", err)))
			listContainer.Refresh()
			return
		}

		if len(snippets) == 0 {
			emptyTitle := canvas.NewText("Belum Ada Koleksi Snippet Kode", constants.ColorTextPrimary)
			emptyTitle.TextSize = constants.FontSizeH2
			emptyTitle.TextStyle = fyne.TextStyle{Bold: true}

			emptySub := canvas.NewText("Simpan fungsi reusable, template CLI, atau script automation Anda agar siap digunakan kapan saja.", constants.ColorTextMuted)
			emptySub.TextSize = constants.FontSizeBody

			emptyBadge := components.BadgeCyan("KODE SNIPPET")
			emptyHeader := container.NewHBox(emptyTitle, emptyBadge)

			emptyBox := container.NewVBox(
				emptyHeader,
				emptySub,
			)
			listContainer.Add(components.NewPlainCardWithAccent(emptyBox, constants.ColorInfo))
			listContainer.Refresh()
			return
		}

		for _, s := range snippets {
			snip := s

			copyBtn := widget.NewButtonWithIcon("Salin Kode", theme.ContentCopyIcon(), func() {
				p.copyToClip(snip.Content)
			})
			copyBtn.Importance = widget.LowImportance

			delBtn := widget.NewButtonWithIcon("Hapus", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Snippet", fmt.Sprintf("Hapus snippet %q?", snip.Title), func(ok bool) {
					if ok {
						_ = database.DeleteSnippet(snip.ID)
						reloadSnippets()
					}
				}, p.window)
			})
			delBtn.Importance = widget.LowImportance

			editBtn := widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
				eTitle := widget.NewEntry()
				eTitle.SetText(snip.Title)
				eLang := widget.NewEntry()
				eLang.SetText(snip.Language)
				eDesc := widget.NewEntry()
				eDesc.SetText(snip.Description)
				eCode := widget.NewMultiLineEntry()
				eCode.SetText(snip.Content)
				eCode.SetMinRowsVisible(6)
				eCode.TextStyle = fyne.TextStyle{Monospace: true}

				lblT := canvas.NewText("JUDUL SNIPPET", constants.ColorTextPrimary)
				lblT.TextSize = constants.FontSizeLabel
				lblT.TextStyle = fyne.TextStyle{Bold: true}

				lblL := canvas.NewText("BAHASA / TOOLS", constants.ColorTextPrimary)
				lblL.TextSize = constants.FontSizeLabel
				lblL.TextStyle = fyne.TextStyle{Bold: true}

				lblD := canvas.NewText("DESKRIPSI", constants.ColorTextPrimary)
				lblD.TextSize = constants.FontSizeLabel
				lblD.TextStyle = fyne.TextStyle{Bold: true}

				lblC := canvas.NewText("KODE / SCRIPT", constants.ColorTextPrimary)
				lblC.TextSize = constants.FontSizeLabel
				lblC.TextStyle = fyne.TextStyle{Bold: true}

				fContent := container.NewVBox(
					lblT, eTitle,
					lblL, eLang,
					lblD, eDesc,
					lblC, eCode,
				)

				components.ShowBrutalistFormDialog(
					p.window,
					"EDIT SNIPPET",
					constants.ColorInfo,
					"Edit Code Snippet",
					"Perbarui judul, bahasa, atau isi kode script",
					fContent,
					"Simpan Snippet",
					func() {
						if eTitle.Text != "" && eCode.Text != "" {
							_ = database.UpdateSnippet(database.Snippet{
								ID:          snip.ID,
								Title:       eTitle.Text,
								Content:     eCode.Text,
								Language:    eLang.Text,
								Description: eDesc.Text,
							})
							reloadSnippets()
						}
					},
				)
			})
			editBtn.Importance = widget.LowImportance

			titleTxt := canvas.NewText(snip.Title, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeH2
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			langBadge := components.BadgeCyan(strings.ToUpper(snip.Language))

			headerLeft := container.NewHBox(titleTxt, langBadge)
			headerRight := container.NewHBox(copyBtn, editBtn, delBtn)
			header := container.NewBorder(nil, nil, headerLeft, headerRight)

			// Terminal Monospace Code Box
			codeBg := canvas.NewRectangle(constants.ColorBgCardInner)
			codeBg.StrokeColor = constants.ColorBorderSubtle
			codeBg.StrokeWidth = 1.5
			codeBg.CornerRadius = constants.CornerRadiusBrutal

			codeLabel := widget.NewLabel(snip.Content)
			codeLabel.Wrapping = fyne.TextWrapWord
			codeLabel.TextStyle = fyne.TextStyle{Monospace: true}

			codeHeader := container.NewBorder(nil, nil,
				components.BadgeMuted("SCRIPT CONTENT:"),
				nil,
			)
			codePanel := container.NewStack(codeBg, container.NewPadded(container.NewVBox(codeHeader, codeLabel)))

			var cardItems []fyne.CanvasObject
			cardItems = append(cardItems, header)

			if snip.Description != "" {
				descLabel := widget.NewLabel(snip.Description)
				descLabel.Wrapping = fyne.TextWrapWord
				cardItems = append(cardItems, descLabel)
			}
			cardItems = append(cardItems, widget.NewSeparator(), codePanel)

			cardContent := container.NewVBox(cardItems...)
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorInfo))
		}
		listContainer.Refresh()
	}

	// Inline Neo-Brutalist Form Card for Adding New Snippet
	formTitle := canvas.NewText("💻 SIMPAN SNIPPET KODE BARU", constants.ColorTextPrimary)
	formTitle.TextSize = constants.FontSizeH2
	formTitle.TextStyle = fyne.TextStyle{Bold: true}

	formBadge := components.BadgeCyan("CODE REPO")
	formHeader := container.NewHBox(formTitle, formBadge)

	formSub := canvas.NewText("Koleksi script reusable untuk mempercepat pengerjaan praktikum dan debugging Anda.", constants.ColorTextMuted)
	formSub.TextSize = constants.FontSizeSmall

	lblT := canvas.NewText("JUDUL SNIPPET", constants.ColorTextPrimary)
	lblT.TextSize = constants.FontSizeLabel
	lblT.TextStyle = fyne.TextStyle{Bold: true}

	lblL := canvas.NewText("BAHASA / TOOLS (cth: go, bash, python, sql)", constants.ColorTextPrimary)
	lblL.TextSize = constants.FontSizeLabel
	lblL.TextStyle = fyne.TextStyle{Bold: true}

	lblD := canvas.NewText("DESKRIPSI (OPSIONAL)", constants.ColorTextPrimary)
	lblD.TextSize = constants.FontSizeLabel
	lblD.TextStyle = fyne.TextStyle{Bold: true}

	lblC := canvas.NewText("KODE / SCRIPT", constants.ColorTextPrimary)
	lblC.TextSize = constants.FontSizeLabel
	lblC.TextStyle = fyne.TextStyle{Bold: true}

	btnCancel := widget.NewButtonWithIcon("Batal", theme.CancelIcon(), func() {
		if formCard != nil {
			formCard.Hide()
		}
	})
	btnCancel.Importance = widget.LowImportance

	btnSave := widget.NewButtonWithIcon("Simpan Snippet", theme.DocumentSaveIcon(), func() {
		if titleEntry.Text != "" && contentEntry.Text != "" {
			_ = database.CreateSnippet(database.Snippet{
				Title:       titleEntry.Text,
				Content:     contentEntry.Text,
				Language:    langEntry.Text,
				Description: descEntry.Text,
			})
			titleEntry.SetText("")
			langEntry.SetText("")
			descEntry.SetText("")
			contentEntry.SetText("")
			if formCard != nil {
				formCard.Hide()
			}
			reloadSnippets()
		}
	})
	btnSave.Importance = widget.HighImportance

	formActions := container.NewBorder(nil, nil, nil, container.NewHBox(btnCancel, btnSave))

	formInner := container.NewVBox(
		formHeader,
		formSub,
		widget.NewSeparator(),
		lblT, titleEntry,
		lblL, langEntry,
		lblD, descEntry,
		lblC, contentEntry,
		widget.NewSeparator(),
		formActions,
	)
	formCard = components.NewPlainCardWithAccent(formInner, constants.ColorInfo)
	formCard.Hide()

	titleHdr := canvas.NewText("Koleksi Snippet Kode Siap Pakai", constants.ColorTextPrimary)
	titleHdr.TextSize = constants.FontSizeH2
	titleHdr.TextStyle = fyne.TextStyle{Bold: true}

	snippetBadge := components.BadgeIndigo("SNIPPETS")
	titleBox := container.NewHBox(titleHdr, snippetBadge)

	toggleAddBtn := widget.NewButtonWithIcon("+ Tambah Snippet Baru", theme.ContentAddIcon(), func() {
		if formCard.Visible() {
			formCard.Hide()
		} else {
			formCard.Show()
		}
	})
	toggleAddBtn.Importance = widget.HighImportance

	topBar := container.NewBorder(nil, nil,
		titleBox,
		toggleAddBtn,
	)
	reloadSnippets()

	mainContent := container.NewVBox(
		container.NewPadded(topBar),
		formCard,
		listContainer,
	)

	return container.NewVScroll(container.NewPadded(mainContent))
}

// 3. Tab Checklists
func (p *LogbookPage) buildChecklistsTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

	newChecklistEntry := widget.NewEntry()
	newChecklistEntry.SetPlaceHolder("Nama Checklist Baru (cth: Pre-Deployment Server / Uji Topologi)...")

	var formCard fyne.CanvasObject
	var reloadChecklists func()

	reloadChecklists = func() {
		listContainer.Objects = nil
		checklists, err := database.GetChecklistsWithItems()
		if err != nil {
			listContainer.Add(widget.NewLabel(fmt.Sprintf("Error memuat checklist: %v", err)))
			listContainer.Refresh()
			return
		}

		if len(checklists) == 0 {
			emptyTitle := canvas.NewText("Belum Ada Checklist", constants.ColorTextPrimary)
			emptyTitle.TextSize = constants.FontSizeH2
			emptyTitle.TextStyle = fyne.TextStyle{Bold: true}

			emptySub := canvas.NewText("Buat daftar periksa checklist tugas, konfigurasi, atau tahapan praktikum.", constants.ColorTextMuted)
			emptySub.TextSize = constants.FontSizeBody

			emptyBadge := components.BadgeSuccess("CHECKLIST")
			emptyHeader := container.NewHBox(emptyTitle, emptyBadge)

			emptyBox := container.NewVBox(
				emptyHeader,
				emptySub,
			)
			listContainer.Add(components.NewPlainCardWithAccent(emptyBox, constants.ColorSuccess))
			listContainer.Refresh()
			return
		}

		for _, cl := range checklists {
			currentCL := cl
			delCLBtn := widget.NewButtonWithIcon("Hapus", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Checklist", fmt.Sprintf("Hapus checklist %q?", currentCL.Name), func(ok bool) {
					if ok {
						_ = database.DeleteChecklist(currentCL.ID)
						reloadChecklists()
					}
				}, p.window)
			})
			delCLBtn.Importance = widget.LowImportance

			titleTxt := canvas.NewText(currentCL.Name, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeH2
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			// Calculation
			completed := 0
			for _, it := range currentCL.Items {
				if it.Checked {
					completed++
				}
			}

			var itemBadge fyne.CanvasObject
			if len(currentCL.Items) > 0 && completed == len(currentCL.Items) {
				itemBadge = components.BadgeSuccess(fmt.Sprintf("%d/%d Selesai ✓", completed, len(currentCL.Items)))
			} else {
				itemBadge = components.BadgeYellow(fmt.Sprintf("%d/%d Selesai", completed, len(currentCL.Items)))
			}

			header := container.NewBorder(nil, nil,
				container.NewHBox(titleTxt, itemBadge),
				delCLBtn,
			)

			itemsBox := container.NewVBox()
			for _, it := range currentCL.Items {
				itemObj := it
				chk := widget.NewCheck(itemObj.Item, func(bool) {
					_ = database.ToggleItem(itemObj.ID)
					reloadChecklists()
				})
				chk.SetChecked(itemObj.Checked)

				delItemBtn := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
					_ = database.DeleteChecklistItem(itemObj.ID)
					reloadChecklists()
				})
				delItemBtn.Importance = widget.LowImportance

				itemRow := container.NewBorder(nil, nil, chk, delItemBtn)
				itemsBox.Add(itemRow)
			}

			// Add item input row with brutalist badge
			newItemEntry := widget.NewEntry()
			newItemEntry.SetPlaceHolder("Tulis item checklist baru...")
			addItemBtn := widget.NewButtonWithIcon("Tambah", theme.ContentAddIcon(), func() {
				if newItemEntry.Text != "" {
					_ = database.AddItem(currentCL.ID, newItemEntry.Text)
					newItemEntry.SetText("")
					reloadChecklists()
				}
			})
			addItemBtn.Importance = widget.MediumImportance

			addBadge := components.BadgeYellow("+ ITEM")
			addRow := container.NewBorder(nil, nil, addBadge, addItemBtn, newItemEntry)

			cardContent := container.NewVBox(header, widget.NewSeparator(), itemsBox, addRow)
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorSuccess))
		}
		listContainer.Refresh()
	}

	// Inline Neo-Brutalist Form Card for Adding New Checklist
	formTitle := canvas.NewText("📋 BUAT CHECKLIST BARU", constants.ColorTextPrimary)
	formTitle.TextSize = constants.FontSizeH2
	formTitle.TextStyle = fyne.TextStyle{Bold: true}

	formBadge := components.BadgeSuccess("VERIFIKASI")
	formHeader := container.NewHBox(formTitle, formBadge)

	lblCL := canvas.NewText("NAMA CHECKLIST BARU", constants.ColorTextPrimary)
	lblCL.TextSize = constants.FontSizeLabel
	lblCL.TextStyle = fyne.TextStyle{Bold: true}

	btnCancelCL := widget.NewButtonWithIcon("Batal", theme.CancelIcon(), func() {
		if formCard != nil {
			formCard.Hide()
		}
	})
	btnCancelCL.Importance = widget.LowImportance

	btnSaveCL := widget.NewButtonWithIcon("Simpan Checklist", theme.DocumentSaveIcon(), func() {
		if newChecklistEntry.Text != "" {
			_ = database.CreateChecklist(newChecklistEntry.Text)
			newChecklistEntry.SetText("")
			if formCard != nil {
				formCard.Hide()
			}
			reloadChecklists()
		}
	})
	btnSaveCL.Importance = widget.HighImportance

	formActions := container.NewBorder(nil, nil, nil, container.NewHBox(btnCancelCL, btnSaveCL))

	formInner := container.NewVBox(
		formHeader,
		widget.NewSeparator(),
		lblCL, newChecklistEntry,
		widget.NewSeparator(),
		formActions,
	)
	formCard = components.NewPlainCardWithAccent(formInner, constants.ColorSuccess)
	formCard.Hide()

	titleHdr := canvas.NewText("Daftar Checklist & Prosedur Praktikum", constants.ColorTextPrimary)
	titleHdr.TextSize = constants.FontSizeH2
	titleHdr.TextStyle = fyne.TextStyle{Bold: true}

	clBadge := components.BadgeSuccess("CHECKLIST")
	titleBox := container.NewHBox(titleHdr, clBadge)

	toggleCLBtn := widget.NewButtonWithIcon("+ Buat Checklist Baru", theme.ContentAddIcon(), func() {
		if formCard.Visible() {
			formCard.Hide()
		} else {
			formCard.Show()
		}
	})
	toggleCLBtn.Importance = widget.HighImportance

	topBar := container.NewBorder(nil, nil, titleBox, toggleCLBtn)
	reloadChecklists()

	mainContent := container.NewVBox(
		container.NewPadded(topBar),
		formCard,
		listContainer,
	)

	return container.NewVScroll(container.NewPadded(mainContent))
}
