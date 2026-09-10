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

func (p *LogbookPage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavLogbook,
		"Pencatatan Error Logbook, repositori Code Snippet, dan Checklist debugging.",
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
			listContainer.Add(widget.NewLabel(constants.NoDataMessage))
			listContainer.Refresh()
			return
		}

		for _, l := range logs {
			item := l
			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Log", fmt.Sprintf("Yakin ingin menghapus %q?", item.Title), func(ok bool) {
					if ok {
						_ = database.DeleteLog(item.ID)
						reloadLogs("")
					}
				}, p.window)
			})
			delBtn.Importance = widget.LowImportance

			titleTxt := canvas.NewText(item.Title, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeBody
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			tagBadge := components.BadgeYellow(item.Tags)
			dateTxt := canvas.NewText(item.CreatedAt, constants.ColorTextMuted)
			dateTxt.TextSize = constants.FontSizeLabel

			header := container.NewBorder(nil, nil,
				container.NewHBox(titleTxt, tagBadge, dateTxt),
				delBtn,
			)

			errTitle := container.NewHBox(components.BadgeDanger("PESAN ERROR"))
			errBox := widget.NewLabel(item.ErrorMessage)
			errBox.Wrapping = fyne.TextWrapWord

			solTitle := container.NewHBox(components.BadgeSuccess("SOLUSI TERVERIFIKASI"))
			solBox := widget.NewLabel(item.Solution)
			solBox.Wrapping = fyne.TextWrapWord

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				errTitle,
				errBox,
				widget.NewSeparator(),
				solTitle,
				solBox,
			)
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorDanger))
		}
		listContainer.Refresh()
	}

	searchBar := components.NewSearchBar(constants.SearchPlaceholder, reloadLogs)

	addBtn := widget.NewButtonWithIcon("Tambah Log Baru", theme.ContentAddIcon(), func() {
		titleEntry := widget.NewEntry()
		errEntry := widget.NewMultiLineEntry()
		errEntry.SetMinRowsVisible(3)
		solEntry := widget.NewMultiLineEntry()
		solEntry.SetMinRowsVisible(3)
		tagsEntry := widget.NewEntry()
		tagsEntry.SetPlaceHolder("cth: go, sql, docker")

		formItems := []*widget.FormItem{
			{Text: "Judul Masalah", Widget: titleEntry},
			{Text: "Pesan Error", Widget: errEntry},
			{Text: "Solusi / Solved By", Widget: solEntry},
			{Text: "Tags (dipisah koma)", Widget: tagsEntry},
		}

		dialog.ShowForm("Tambah Catatan Error", "Simpan", "Batal", formItems, func(ok bool) {
			if !ok || titleEntry.Text == "" {
				return
			}
			err := database.CreateLog(database.ErrorLog{
				Title:        titleEntry.Text,
				ErrorMessage: errEntry.Text,
				Solution:     solEntry.Text,
				Tags:         tagsEntry.Text,
			})
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			reloadLogs("")
		}, p.window)
	})
	addBtn.Importance = widget.HighImportance

	topBar := container.NewBorder(nil, nil, nil, addBtn, searchBar.Container)
	reloadLogs("")

	return container.NewBorder(container.NewPadded(topBar), nil, nil, nil, container.NewVScroll(listContainer))
}

// 2. Tab Code Snippets
func (p *LogbookPage) buildSnippetsTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

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
			listContainer.Add(widget.NewLabel(constants.NoDataMessage))
			listContainer.Refresh()
			return
		}

		for _, s := range snippets {
			snip := s
			copyBtn := widget.NewButtonWithIcon("Salin Kode", theme.ContentCopyIcon(), func() {
				p.window.Clipboard().SetContent(snip.Content)
				dialog.ShowInformation("Clipboard", constants.StatusCopied, p.window)
			})
			copyBtn.Importance = widget.LowImportance

			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Snippet", fmt.Sprintf("Hapus snippet %q?", snip.Title), func(ok bool) {
					if ok {
						_ = database.DeleteSnippet(snip.ID)
						reloadSnippets()
					}
				}, p.window)
			})
			delBtn.Importance = widget.LowImportance

			titleTxt := canvas.NewText(snip.Title, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeBody
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			langBadge := components.BadgeCyan(strings.ToUpper(snip.Language))

			header := container.NewBorder(nil, nil,
				container.NewHBox(titleTxt, langBadge),
				container.NewHBox(copyBtn, delBtn),
			)

			descLabel := widget.NewLabel(snip.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			codeArea := widget.NewMultiLineEntry()
			codeArea.SetText(snip.Content)
			codeArea.Disable()
			codeArea.TextStyle = fyne.TextStyle{Monospace: true}
			codeArea.SetMinRowsVisible(4)

			cardContent := container.NewVBox(
				header,
				descLabel,
				widget.NewSeparator(),
				codeArea,
			)
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorInfo))
		}
		listContainer.Refresh()
	}

	addBtn := widget.NewButtonWithIcon("Tambah Snippet Baru", theme.ContentAddIcon(), func() {
		titleEntry := widget.NewEntry()
		langEntry := widget.NewEntry()
		langEntry.SetPlaceHolder("cth: go, bash, python, sql")
		descEntry := widget.NewEntry()
		contentEntry := widget.NewMultiLineEntry()
		contentEntry.SetMinRowsVisible(6)
		contentEntry.TextStyle = fyne.TextStyle{Monospace: true}

		formItems := []*widget.FormItem{
			{Text: "Judul Snippet", Widget: titleEntry},
			{Text: "Bahasa / Tools", Widget: langEntry},
			{Text: "Deskripsi", Widget: descEntry},
			{Text: "Kode / Script", Widget: contentEntry},
		}

		dialog.ShowForm("Tambah Snippet Baru", "Simpan", "Batal", formItems, func(ok bool) {
			if !ok || titleEntry.Text == "" || contentEntry.Text == "" {
				return
			}
			err := database.CreateSnippet(database.Snippet{
				Title:       titleEntry.Text,
				Content:     contentEntry.Text,
				Language:    langEntry.Text,
				Description: descEntry.Text,
			})
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			reloadSnippets()
		}, p.window)
	})
	addBtn.Importance = widget.HighImportance

	titleHdr := canvas.NewText("Koleksi Snippet Kode Siap Pakai", constants.ColorTextPrimary)
	titleHdr.TextSize = constants.FontSizeBody
	titleHdr.TextStyle = fyne.TextStyle{Bold: true}

	topBar := container.NewBorder(nil, nil,
		titleHdr,
		addBtn,
	)
	reloadSnippets()

	return container.NewBorder(container.NewPadded(topBar), nil, nil, nil, container.NewVScroll(listContainer))
}

// 3. Tab Checklists
func (p *LogbookPage) buildChecklistsTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

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
			listContainer.Add(widget.NewLabel(constants.NoDataMessage))
			listContainer.Refresh()
			return
		}

		for _, cl := range checklists {
			currentCL := cl
			delCLBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Checklist", fmt.Sprintf("Hapus checklist %q?", currentCL.Name), func(ok bool) {
					if ok {
						_ = database.DeleteChecklist(currentCL.ID)
						reloadChecklists()
					}
				}, p.window)
			})
			delCLBtn.Importance = widget.LowImportance

			titleTxt := canvas.NewText(currentCL.Name, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeBody
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			itemBadge := components.BadgeYellow(fmt.Sprintf("%d Items", len(currentCL.Items)))

			header := container.NewBorder(nil, nil,
				container.NewHBox(titleTxt, itemBadge),
				delCLBtn,
			)

			itemsBox := container.NewVBox()
			for _, it := range currentCL.Items {
				itemObj := it
				chk := widget.NewCheck(itemObj.Item, func(bool) {
					_ = database.ToggleItem(itemObj.ID)
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

			// Add item input row
			newItemEntry := widget.NewEntry()
			newItemEntry.SetPlaceHolder("Tambah item checklist baru...")
			addItemBtn := widget.NewButtonWithIcon("Tambah", theme.ContentAddIcon(), func() {
				if newItemEntry.Text != "" {
					_ = database.AddItem(currentCL.ID, newItemEntry.Text)
					newItemEntry.SetText("")
					reloadChecklists()
				}
			})
			addItemBtn.Importance = widget.LowImportance
			addRow := container.NewBorder(nil, nil, nil, addItemBtn, newItemEntry)

			cardContent := container.NewVBox(header, widget.NewSeparator(), itemsBox, addRow)
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorSuccess))
		}
		listContainer.Refresh()
	}

	newChecklistEntry := widget.NewEntry()
	newChecklistEntry.SetPlaceHolder("Nama Checklist Baru (cth: Pre-Deployment Checklist)...")
	createCLBtn := widget.NewButtonWithIcon("Buat Checklist", theme.ContentAddIcon(), func() {
		if newChecklistEntry.Text != "" {
			_ = database.CreateChecklist(newChecklistEntry.Text)
			newChecklistEntry.SetText("")
			reloadChecklists()
		}
	})
	createCLBtn.Importance = widget.HighImportance

	topBar := container.NewBorder(nil, nil, nil, createCLBtn, newChecklistEntry)
	reloadChecklists()

	return container.NewBorder(container.NewPadded(topBar), nil, nil, nil, container.NewVScroll(listContainer))
}
