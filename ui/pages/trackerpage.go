package pages

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"

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

type TrackerPage struct {
	window fyne.Window
}

func NewTrackerPage(win fyne.Window) *TrackerPage {
	return &TrackerPage{window: win}
}

func (p *TrackerPage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavTracker,
		"Papan Kanban manajemen tugas mahasiswa IT dan Showcase Portofolio Proyek.",
		components.BadgeCyan("PROJECT & TASKS"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabKanban, theme.ListIcon(), p.buildKanbanTab()),
		container.NewTabItemWithIcon(constants.TabProjects, theme.FolderOpenIcon(), p.buildProjectsTab()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
}

// 1. Tab Kanban Board (Belum Mulai, Dikerjakan, Selesai)
func (p *TrackerPage) buildKanbanTab() fyne.CanvasObject {
	todoBox := container.NewVBox()
	inprogBox := container.NewVBox()
	doneBox := container.NewVBox()

	var reloadKanban func()
	reloadKanban = func() {
		todoBox.Objects = nil
		inprogBox.Objects = nil
		doneBox.Objects = nil

		tasks, err := database.GetAllTasks()
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}

		for _, t := range tasks {
			task := t
			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Tugas", fmt.Sprintf("Hapus tugas %q?", task.Title), func(ok bool) {
					if ok {
						_ = database.DeleteTask(task.ID)
						reloadKanban()
					}
				}, p.window)
			})
			delBtn.Importance = widget.LowImportance

			actions := container.NewHBox()
			switch task.Status {
			case "todo":
				btnStart := widget.NewButtonWithIcon("Mulai", theme.MediaPlayIcon(), func() {
					_ = database.UpdateTaskStatus(task.ID, "inprogress")
					reloadKanban()
				})
				btnStart.Importance = widget.HighImportance
				actions.Add(btnStart)
			case "inprogress":
				btnCancel := widget.NewButtonWithIcon("◀ Batal", theme.NavigateBackIcon(), func() {
					_ = database.UpdateTaskStatus(task.ID, "todo")
					reloadKanban()
				})
				btnCancel.Importance = widget.LowImportance
				actions.Add(btnCancel)

				btnDone := widget.NewButtonWithIcon("Selesai ✔", theme.ConfirmIcon(), func() {
					_ = database.UpdateTaskStatus(task.ID, "done")
					reloadKanban()
				})
				btnDone.Importance = widget.HighImportance
				actions.Add(btnDone)
			case "done":
				btnReopen := widget.NewButtonWithIcon("◀ Re-open", theme.NavigateBackIcon(), func() {
					_ = database.UpdateTaskStatus(task.ID, "inprogress")
					reloadKanban()
				})
				btnReopen.Importance = widget.LowImportance
				actions.Add(btnReopen)
			}
			actions.Add(delBtn)

			titleTxt := canvas.NewText(task.Title, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
			titleTxt.TextSize = 13.5
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			cardHeader := container.NewBorder(nil, nil,
				titleTxt,
				actions,
			)

			catBadge := components.BadgeIndigo(task.Category)
			var dlBadge fyne.CanvasObject
			if task.Deadline != "" {
				dlBadge = components.BadgeWarning("DL: " + task.Deadline)
			}

			metaRow := container.NewHBox(catBadge)
			if dlBadge != nil {
				metaRow.Add(dlBadge)
			}

			descLabel := widget.NewLabel(task.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			cardContent := container.NewVBox(
				cardHeader,
				metaRow,
				descLabel,
			)
			card := components.NewPlainCard(cardContent)

			switch task.Status {
			case "inprogress":
				inprogBox.Add(card)
			case "done":
				doneBox.Add(card)
			default:
				todoBox.Add(card)
			}
		}

		todoBox.Refresh()
		inprogBox.Refresh()
		doneBox.Refresh()
	}

	addTaskBtn := widget.NewButtonWithIcon("Tambah Tugas Baru", theme.ContentAddIcon(), func() {
		titleEntry := widget.NewEntry()
		descEntry := widget.NewMultiLineEntry()
		descEntry.SetMinRowsVisible(3)
		catSelect := widget.NewSelect([]string{"Praktikum", "Tugas Kuliah", "Project", "Lainnya"}, nil)
		catSelect.SetSelected("Praktikum")
		deadlineEntry := widget.NewEntry()
		deadlineEntry.SetPlaceHolder("YYYY-MM-DD (cth: 2026-09-30)")

		formItems := []*widget.FormItem{
			{Text: "Judul Tugas", Widget: titleEntry},
			{Text: "Kategori", Widget: catSelect},
			{Text: "Deadline", Widget: deadlineEntry},
			{Text: "Deskripsi", Widget: descEntry},
		}

		dialog.ShowForm("Tambah Tugas Kanban", "Simpan", "Batal", formItems, func(ok bool) {
			if !ok || titleEntry.Text == "" {
				return
			}
			err := database.CreateTask(database.Task{
				Title:       titleEntry.Text,
				Description: descEntry.Text,
				Category:    catSelect.Selected,
				Deadline:    deadlineEntry.Text,
				Status:      "todo",
			})
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			reloadKanban()
		}, p.window)
	})
	addTaskBtn.Importance = widget.HighImportance

	topBar := container.NewBorder(nil, nil,
		canvas.NewText("Papan Pelacak Tugas Kanban", color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}),
		addTaskBtn,
	)

	col1Header := container.NewPadded(components.BadgeIndigo("📌 BELUM MULAI"))
	col2Header := container.NewPadded(components.BadgeWarning("⚡ SEDANG DIKERJAKAN"))
	col3Header := container.NewPadded(components.BadgeSuccess("✅ SELESAI"))

	col1 := container.NewBorder(col1Header, nil, nil, nil, container.NewVScroll(todoBox))
	col2 := container.NewBorder(col2Header, nil, nil, nil, container.NewVScroll(inprogBox))
	col3 := container.NewBorder(col3Header, nil, nil, nil, container.NewVScroll(doneBox))

	kanbanGrid := container.NewGridWithColumns(3, col1, col2, col3)

	reloadKanban()
	return container.NewBorder(container.NewPadded(topBar), nil, nil, nil, kanbanGrid)
}

// 2. Tab Projects Showcase & Portfolio Markdown Exporter
func (p *TrackerPage) buildProjectsTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

	var reloadProjects func()
	reloadProjects = func() {
		listContainer.Objects = nil
		projects, err := database.GetAllProjects()
		if err != nil {
			listContainer.Add(widget.NewLabel(fmt.Sprintf("Error memuat project: %v", err)))
			listContainer.Refresh()
			return
		}

		if len(projects) == 0 {
			listContainer.Add(widget.NewLabel(constants.NoDataMessage))
			listContainer.Refresh()
			return
		}

		for _, pr := range projects {
			proj := pr
			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				dialog.ShowConfirm("Hapus Proyek", fmt.Sprintf("Hapus proyek %q?", proj.Title), func(ok bool) {
					if ok {
						_ = database.DeleteProject(proj.ID)
						reloadProjects()
					}
				}, p.window)
			})
			delBtn.Importance = widget.LowImportance

			titleTxt := canvas.NewText(proj.Title, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
			titleTxt.TextSize = 14
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			header := container.NewBorder(nil, nil,
				titleTxt,
				delBtn,
			)

			techBadge := components.BadgeCyan(proj.Technologies)
			descLabel := widget.NewLabel(proj.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			var links []fyne.CanvasObject
			if proj.RepoURL != "" {
				links = append(links, canvas.NewText("🔗 Repo: "+proj.RepoURL, color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0xFF}))
			}
			if proj.LiveURL != "" {
				links = append(links, canvas.NewText("🌐 Demo: "+proj.LiveURL, color.RGBA{R: 0x34, G: 0xD3, B: 0x99, A: 0xFF}))
			}
			linksBox := container.NewVBox(links...)

			cardContent := container.NewVBox(
				header,
				container.NewHBox(techBadge),
				descLabel,
				widget.NewSeparator(),
				linksBox,
			)
			listContainer.Add(components.NewPlainCard(cardContent))
		}
		listContainer.Refresh()
	}

	addProjBtn := widget.NewButtonWithIcon("Tambah Proyek Baru", theme.ContentAddIcon(), func() {
		titleEntry := widget.NewEntry()
		descEntry := widget.NewMultiLineEntry()
		descEntry.SetMinRowsVisible(3)
		techEntry := widget.NewEntry()
		techEntry.SetPlaceHolder("cth: Go, Fyne, SQLite, Docker")
		repoEntry := widget.NewEntry()
		repoEntry.SetPlaceHolder("cth: https://github.com/username/repo")
		liveEntry := widget.NewEntry()
		liveEntry.SetPlaceHolder("cth: https://myproject.com")

		formItems := []*widget.FormItem{
			{Text: "Nama Proyek", Widget: titleEntry},
			{Text: "Teknologi (dipisah koma)", Widget: techEntry},
			{Text: "Repository URL", Widget: repoEntry},
			{Text: "Demo / Live URL", Widget: liveEntry},
			{Text: "Deskripsi", Widget: descEntry},
		}

		dialog.ShowForm("Tambah Proyek ke Portofolio", "Simpan", "Batal", formItems, func(ok bool) {
			if !ok || titleEntry.Text == "" {
				return
			}
			err := database.CreateProject(database.Project{
				Title:        titleEntry.Text,
				Description:  descEntry.Text,
				Technologies: techEntry.Text,
				RepoURL:      repoEntry.Text,
				LiveURL:      liveEntry.Text,
			})
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			reloadProjects()
		}, p.window)
	})
	addProjBtn.Importance = widget.HighImportance

	exportBtn := widget.NewButtonWithIcon(constants.BtnExportMD, theme.DocumentSaveIcon(), func() {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "."
		}
		exportFile := filepath.Join(cwd, "portfolio.md")
		err = database.ExportPortfolioMarkdown(exportFile)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		dialog.ShowInformation("Ekspor Berhasil", fmt.Sprintf("Dokumen portofolio berhasil disimpan ke:\n%s", exportFile), p.window)
	})
	exportBtn.Importance = widget.HighImportance

	topBar := container.NewBorder(nil, nil,
		canvas.NewText("Showcase & Profil Portofolio Proyek", color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}),
		container.NewHBox(addProjBtn, exportBtn),
	)

	reloadProjects()
	return container.NewBorder(container.NewPadded(topBar), nil, nil, nil, container.NewVScroll(listContainer))
}
