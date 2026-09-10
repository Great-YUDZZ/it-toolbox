package pages

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/core/converter"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

type FileConverterPage struct {
	window fyne.Window
}

func NewFileConverterPage(win fyne.Window) *FileConverterPage {
	return &FileConverterPage{window: win}
}

func (p *FileConverterPage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavFileConverter,
		"Konversi dokumen, gambar, ekstraksi teks OCR, kompresi berkas, serta manipulasi PDF (Merge & Split).",
		components.BadgeCyan("FILE CONVERTER PRO"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabDocConverter, theme.DocumentIcon(), p.buildDocConverterTab()),
		container.NewTabItemWithIcon(constants.TabImgConverter, theme.FileImageIcon(), p.buildImgConverterTab()),
		container.NewTabItemWithIcon(constants.TabOCR, theme.VisibilityIcon(), p.buildOCRTab()),
		container.NewTabItemWithIcon(constants.TabPDFTools, theme.FolderOpenIcon(), p.buildPDFToolsTab()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
}

func (p *FileConverterPage) getDownloadDir() string {
	home, err := os.UserHomeDir()
	if err == nil {
		dl := filepath.Join(home, "Downloads")
		if _, err := os.Stat(dl); err == nil {
			return dl
		}
		return home
	}
	return "."
}

func (p *FileConverterPage) copyToClip(txt string) {
	p.window.Clipboard().SetContent(txt)
	dialog.ShowInformation("Clipboard", constants.StatusCopied, p.window)
}

// ----------------------------------------------------------------------------
// 1. Tab Konversi Dokumen & Batch Upload
// ----------------------------------------------------------------------------
func (p *FileConverterPage) buildDocConverterTab() fyne.CanvasObject {
	var selectedFiles []string
	var convertedFiles []string

	fileListContainer := container.NewVBox()
	resListContainer := container.NewVBox()

	targetFormatSelect := widget.NewSelect([]string{"DOCX", "PDF", "TXT", "XLSX", "PPTX"}, nil)
	targetFormatSelect.SetSelected("PDF")

	compressCheck := widget.NewCheck("Kompres ukuran berkas setelah konversi", nil)
	compressCheck.SetChecked(true)

	statusLabel := widget.NewLabel("Belum ada berkas dipilih.")

	refreshLists := func() {
		fileListContainer.Objects = nil
		for i, f := range selectedFiles {
			idx := i
			path := f
			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				selectedFiles = append(selectedFiles[:idx], selectedFiles[idx+1:]...)
				fileListContainer.Refresh()
				statusLabel.SetText(fmt.Sprintf("%d berkas dipilih.", len(selectedFiles)))
			})
			delBtn.Importance = widget.LowImportance

			fi, _ := os.Stat(path)
			sizeStr := ""
			if fi != nil {
				sizeStr = fmt.Sprintf("%.2f KB", float64(fi.Size())/1024.0)
			}

			extBadge := components.BadgeYellow(strings.ToUpper(strings.TrimPrefix(filepath.Ext(path), ".")))
			row := container.NewBorder(nil, nil,
				container.NewHBox(extBadge, widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)", filepath.Base(path), sizeStr), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
				delBtn,
			)
			fileListContainer.Add(components.NewPlainCardWithAccent(row, constants.ColorAccentCyan))
		}
		fileListContainer.Refresh()
		statusLabel.SetText(fmt.Sprintf("%d berkas dipilih.", len(selectedFiles)))
	}

	uploadBtn := widget.NewButtonWithIcon(constants.BtnUploadFile, theme.FileIcon(), func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			path := closer.URI().Path()
			selectedFiles = append(selectedFiles, path)
			refreshLists()
		}, p.window)
	})
	uploadBtn.Importance = widget.HighImportance

	uploadBatchEntry := widget.NewEntry()
	uploadBatchEntry.SetPlaceHolder("Tempelkan path berkas (bisa banyak berkas dipisah baris baru)...")

	addFromEntryBtn := widget.NewButtonWithIcon("Tambah dari Path", theme.ContentAddIcon(), func() {
		lines := strings.Split(uploadBatchEntry.Text, "\n")
		added := 0
		for _, line := range lines {
			trimmed := strings.TrimSpace(strings.Trim(line, "\"'"))
			if trimmed != "" {
				if _, err := os.Stat(trimmed); err == nil {
					selectedFiles = append(selectedFiles, trimmed)
					added++
				}
			}
		}
		if added > 0 {
			uploadBatchEntry.SetText("")
			refreshLists()
		}
	})

	zipBtn := widget.NewButtonWithIcon(constants.BtnDownloadZip, theme.DocumentSaveIcon(), func() {
		if len(convertedFiles) == 0 {
			dialog.ShowInformation("Perhatian", "Belum ada berkas yang berhasil dikonversi untuk di-zip.", p.window)
			return
		}
		outZip := filepath.Join(p.getDownloadDir(), "converted_documents.zip")
		if err := converter.CreateZipArchive(convertedFiles, outZip); err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		dialog.ShowInformation("Arsip ZIP Siap", fmt.Sprintf("Semua berkas konversi berhasil dikemas ke:\n%s", outZip), p.window)
	})

	convertBtn := widget.NewButtonWithIcon(constants.BtnConvertAll, theme.ConfirmIcon(), func() {
		if len(selectedFiles) == 0 {
			dialog.ShowInformation("Pemberitahuan", "Silakan pilih setidaknya 1 berkas terlebih dahulu.", p.window)
			return
		}

		outDir := filepath.Join(p.getDownloadDir(), "IT_Toolbox_Converted")
		_ = os.MkdirAll(outDir, 0755)

		convertedFiles = nil
		resListContainer.Objects = nil

		progress := dialog.NewCustom("Sedang Mengonversi", "Mohon tunggu...", widget.NewProgressBarInfinite(), p.window)
		progress.Show()

		successCount := 0
		var errs []string

		for _, inPath := range selectedFiles {
			targetFmt := strings.ToLower(targetFormatSelect.Selected)
			outPath, err := converter.ConvertDocument(inPath, targetFmt, outDir)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: %v", filepath.Base(inPath), err))
				continue
			}

			if targetFmt == "pdf" && compressCheck.Checked {
				compPath := filepath.Join(outDir, "comp_"+filepath.Base(outPath))
				if errComp := converter.CompressPDF(outPath, "ebook", compPath); errComp == nil {
					_ = os.Remove(outPath)
					_ = os.Rename(compPath, outPath)
				}
			}

			convertedFiles = append(convertedFiles, outPath)
			successCount++

			resItem := outPath
			openBtn := widget.NewButtonWithIcon("Buka Lokasi", theme.FolderOpenIcon(), func() {
				dialog.ShowInformation("Lokasi Berkas", fmt.Sprintf("Berkas tersimpan di:\n%s", resItem), p.window)
			})
			openBtn.Importance = widget.LowImportance

			fi, _ := os.Stat(resItem)
			sizeStr := ""
			if fi != nil {
				sizeStr = fmt.Sprintf("%.2f KB", float64(fi.Size())/1024.0)
			}

			resBadge := components.BadgeSuccess("SELESAI")
			resCard := container.NewBorder(nil, nil,
				container.NewHBox(resBadge, widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)", filepath.Base(resItem), sizeStr), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
				openBtn,
			)
			resListContainer.Add(components.NewPlainCardWithAccent(resCard, constants.ColorSuccess))
		}

		progress.Hide()

		if len(errs) > 0 {
			dialog.ShowError(fmt.Errorf("beberapa konversi gagal:\n%s", strings.Join(errs, "\n")), p.window)
		} else {
			dialog.ShowInformation("Konversi Selesai", fmt.Sprintf("Berhasil mengonversi %d berkas ke folder:\n%s", successCount, outDir), p.window)
		}
		resListContainer.Refresh()
	})
	convertBtn.Importance = widget.HighImportance

	inputCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("1. Pilih Berkas Dokumen (PDF, DOCX, XLSX, PPTX, TXT):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(uploadBtn, widget.NewLabel("atau masukkan path berkas:")),
		container.NewBorder(nil, nil, nil, addFromEntryBtn, uploadBatchEntry),
		statusLabel,
		fileListContainer,
	), constants.ColorAccentCyan)

	configCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("2. Konfigurasi Konversi Dokumen:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Format Sasaran:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), targetFormatSelect),
			container.NewVBox(widget.NewLabelWithStyle("Opsi Kompresi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), compressCheck),
		),
		container.NewHBox(convertBtn, zipBtn),
	), constants.ColorAccentYellow)

	resCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("3. Hasil Konversi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		resListContainer,
	), constants.ColorSuccess)

	return container.NewVScroll(container.NewVBox(
		inputCard,
		configCard,
		resCard,
	))
}

// ----------------------------------------------------------------------------
// 2. Tab Konversi & Kompresi Gambar (JPG, PNG, WEBP)
// ----------------------------------------------------------------------------
func (p *FileConverterPage) buildImgConverterTab() fyne.CanvasObject {
	var selectedImages []string
	var convertedImages []string

	imgListContainer := container.NewVBox()
	resListContainer := container.NewVBox()

	targetFmtSelect := widget.NewSelect([]string{"JPG", "PNG", "WEBP"}, nil)
	targetFmtSelect.SetSelected("WEBP")

	qualitySlider := widget.NewSlider(10, 100)
	qualitySlider.SetValue(80)
	qualityLabel := widget.NewLabelWithStyle("Kualitas Kompresi: 80%", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	qualitySlider.OnChanged = func(val float64) {
		qualityLabel.SetText(fmt.Sprintf("Kualitas Kompresi: %d%%", int(val)))
	}

	statusLabel := widget.NewLabel("Belum ada gambar dipilih.")

	refreshImgs := func() {
		imgListContainer.Objects = nil
		for i, path := range selectedImages {
			idx := i
			pStr := path
			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				selectedImages = append(selectedImages[:idx], selectedImages[idx+1:]...)
				imgListContainer.Refresh()
				statusLabel.SetText(fmt.Sprintf("%d gambar dipilih.", len(selectedImages)))
			})
			delBtn.Importance = widget.LowImportance

			fi, _ := os.Stat(pStr)
			sizeStr := ""
			if fi != nil {
				sizeStr = fmt.Sprintf("%.2f KB", float64(fi.Size())/1024.0)
			}
			imgBadge := components.BadgeIndigo(strings.ToUpper(strings.TrimPrefix(filepath.Ext(pStr), ".")))
			row := container.NewBorder(nil, nil,
				container.NewHBox(imgBadge, widget.NewLabelWithStyle(fmt.Sprintf("%s (%s)", filepath.Base(pStr), sizeStr), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
				delBtn,
			)
			imgListContainer.Add(components.NewPlainCard(row))
		}
		imgListContainer.Refresh()
		statusLabel.SetText(fmt.Sprintf("%d gambar dipilih.", len(selectedImages)))
	}

	uploadImgBtn := widget.NewButtonWithIcon("Pilih Gambar", theme.FileImageIcon(), func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			selectedImages = append(selectedImages, closer.URI().Path())
			refreshImgs()
		}, p.window)
	})
	uploadImgBtn.Importance = widget.HighImportance

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Tempelkan path gambar...")
	addPathBtn := widget.NewButtonWithIcon("Tambah Path", theme.ContentAddIcon(), func() {
		for _, line := range strings.Split(pathEntry.Text, "\n") {
			trimmed := strings.TrimSpace(strings.Trim(line, "\"'"))
			if trimmed != "" {
				if _, err := os.Stat(trimmed); err == nil {
					selectedImages = append(selectedImages, trimmed)
				}
			}
		}
		pathEntry.SetText("")
		refreshImgs()
	})

	zipBtn := widget.NewButtonWithIcon(constants.BtnDownloadZip, theme.DocumentSaveIcon(), func() {
		if len(convertedImages) == 0 {
			dialog.ShowInformation("Perhatian", "Belum ada gambar hasil konversi.", p.window)
			return
		}
		outZip := filepath.Join(p.getDownloadDir(), "converted_images.zip")
		if err := converter.CreateZipArchive(convertedImages, outZip); err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		dialog.ShowInformation("Arsip ZIP Siap", fmt.Sprintf("Semua gambar berhasil dikemas ke:\n%s", outZip), p.window)
	})

	convertBtn := widget.NewButtonWithIcon("Konversi & Kompres Semua Gambar", theme.ConfirmIcon(), func() {
		if len(selectedImages) == 0 {
			dialog.ShowInformation("Pemberitahuan", "Pilih setidaknya 1 gambar.", p.window)
			return
		}

		outDir := filepath.Join(p.getDownloadDir(), "IT_Toolbox_Images")
		_ = os.MkdirAll(outDir, 0755)

		convertedImages = nil
		resListContainer.Objects = nil

		for _, inImg := range selectedImages {
			origStat, _ := os.Stat(inImg)
			origSize := int64(0)
			if origStat != nil {
				origSize = origStat.Size()
			}

			outImg, err := converter.ConvertImage(inImg, targetFmtSelect.Selected, int(qualitySlider.Value), outDir)
			if err != nil {
				dialog.ShowError(err, p.window)
				continue
			}

			convertedImages = append(convertedImages, outImg)
			newStat, _ := os.Stat(outImg)
			newSize := int64(0)
			if newStat != nil {
				newSize = newStat.Size()
			}

			savingPct := 0.0
			if origSize > 0 && newSize < origSize {
				savingPct = float64(origSize-newSize) / float64(origSize) * 100.0
			}

			resItem := outImg
			openBtn := widget.NewButtonWithIcon("Buka", theme.FolderOpenIcon(), func() {
				dialog.ShowInformation("Lokasi Gambar", resItem, p.window)
			})
			openBtn.Importance = widget.LowImportance

			resBadge := components.BadgeSuccess(fmt.Sprintf("HEMAT %.1f%%", savingPct))
			resRow := container.NewBorder(nil, nil,
				container.NewHBox(resBadge, widget.NewLabelWithStyle(fmt.Sprintf("%s (Awal: %.1f KB ➔ Akhir: %.1f KB)",
					filepath.Base(resItem), float64(origSize)/1024.0, float64(newSize)/1024.0),
					fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
				openBtn,
			)
			resListContainer.Add(components.NewPlainCardWithAccent(resRow, constants.ColorSuccess))
		}

		dialog.ShowInformation("Selesai", fmt.Sprintf("Konversi gambar selesai di folder:\n%s", outDir), p.window)
		resListContainer.Refresh()
	})
	convertBtn.Importance = widget.HighImportance

	inputCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("1. Pilih Gambar (JPG, PNG, WEBP):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(uploadImgBtn, widget.NewLabel("atau masukkan path:")),
		container.NewBorder(nil, nil, nil, addPathBtn, pathEntry),
		statusLabel,
		imgListContainer,
	), constants.ColorAccentCyan)

	configCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("2. Format Target & Kualitas Kompresi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Format Gambar Sasaran:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), targetFmtSelect),
			container.NewVBox(qualityLabel, qualitySlider),
		),
		container.NewHBox(convertBtn, zipBtn),
	), constants.ColorAccentYellow)

	resCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("3. Hasil Gambar & Penghematan Ukuran:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		resListContainer,
	), constants.ColorSuccess)

	return container.NewVScroll(container.NewVBox(
		inputCard,
		configCard,
		resCard,
	))
}

// ----------------------------------------------------------------------------
// 3. Tab Ekstraksi Teks (OCR)
// ----------------------------------------------------------------------------
func (p *FileConverterPage) buildOCRTab() fyne.CanvasObject {
	filePathEntry := widget.NewEntry()
	filePathEntry.SetPlaceHolder("Pilih berkas gambar atau PDF scan...")

	langSelect := widget.NewSelect([]string{"English (eng)", "Bahasa Indonesia (ind)"}, nil)
	langSelect.SetSelected("English (eng)")

	resultArea := widget.NewMultiLineEntry()
	resultArea.SetPlaceHolder("Hasil teks OCR akan ditampilkan di sini...")
	resultArea.SetMinRowsVisible(12)
	resultArea.TextStyle = fyne.TextStyle{Monospace: true}

	browseBtn := widget.NewButtonWithIcon("Pilih Berkas", theme.FileIcon(), func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			filePathEntry.SetText(closer.URI().Path())
		}, p.window)
	})
	browseBtn.Importance = widget.HighImportance

	ocrBtn := widget.NewButtonWithIcon(constants.BtnRunOCR, theme.VisibilityIcon(), func() {
		targetFile := strings.TrimSpace(filePathEntry.Text)
		if targetFile == "" {
			dialog.ShowInformation("Perhatian", "Silakan tentukan berkas gambar atau PDF yang ingin di-OCR.", p.window)
			return
		}

		langCode := "eng"
		if strings.Contains(langSelect.Selected, "ind") {
			langCode = "ind"
		}

		progress := dialog.NewCustom("Sedang Memproses OCR", "Membaca teks dari dokumen...", widget.NewProgressBarInfinite(), p.window)
		progress.Show()

		var text string
		var err error

		ext := strings.ToLower(filepath.Ext(targetFile))
		if ext == ".pdf" {
			text, err = converter.ExtractTextFromPDF(targetFile)
		} else {
			text, err = converter.ExtractTextFromImage(targetFile, langCode)
		}

		progress.Hide()

		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}

		if text == "" {
			resultArea.SetText("[Tidak ditemukan teks yang dapat dikenali dalam gambar ini]")
		} else {
			resultArea.SetText(text)
		}
	})
	ocrBtn.Importance = widget.HighImportance

	copyBtn := widget.NewButtonWithIcon("Salin Teks", theme.ContentCopyIcon(), func() {
		p.copyToClip(resultArea.Text)
	})

	saveTxtBtn := widget.NewButtonWithIcon("Simpan ke .TXT", theme.DocumentSaveIcon(), func() {
		if resultArea.Text == "" {
			return
		}
		savePath := filepath.Join(p.getDownloadDir(), "extracted_ocr.txt")
		if err := os.WriteFile(savePath, []byte(resultArea.Text), 0644); err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		dialog.ShowInformation("Tersimpan", fmt.Sprintf("Hasil OCR berhasil disimpan ke:\n%s", savePath), p.window)
	})

	cardTop := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabel("Ekstraksi teks dari gambar scan atau PDF menggunakan engine Tesseract OCR:"),
		container.NewBorder(nil, nil, widget.NewLabelWithStyle("Berkas Input:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), browseBtn, filePathEntry),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Pilihan Bahasa:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), langSelect),
			container.NewVBox(widget.NewLabelWithStyle("Aksi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), ocrBtn),
		),
	), constants.ColorAccentCyan)

	cardResult := components.NewPlainCardWithAccent(container.NewVBox(
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("Teks Hasil Ekstraksi OCR:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			container.NewHBox(copyBtn, saveTxtBtn),
		),
		resultArea,
	), constants.ColorSuccess)

	return container.NewVScroll(container.NewVBox(
		cardTop,
		cardResult,
	))
}

// ----------------------------------------------------------------------------
// 4. Tab Merge, Split & Kompresi PDF
// ----------------------------------------------------------------------------
func (p *FileConverterPage) buildPDFToolsTab() fyne.CanvasObject {
	// A. Merge PDF Section
	var mergeFiles []string
	mergeList := container.NewVBox()
	refreshMerge := func() {
		mergeList.Objects = nil
		for i, f := range mergeFiles {
			idx := i
			path := f
			delBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				mergeFiles = append(mergeFiles[:idx], mergeFiles[idx+1:]...)
				mergeList.Refresh()
			})
			delBtn.Importance = widget.LowImportance

			numBadge := components.BadgeYellow(fmt.Sprintf("#%d", idx+1))
			row := container.NewBorder(nil, nil,
				container.NewHBox(numBadge, widget.NewLabelWithStyle(filepath.Base(path), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
				delBtn,
			)
			mergeList.Add(components.NewPlainCardWithAccent(row, constants.ColorTechIndigo))
		}
		mergeList.Refresh()
	}

	addMergeBtn := widget.NewButtonWithIcon("Tambah PDF", theme.ContentAddIcon(), func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			mergeFiles = append(mergeFiles, closer.URI().Path())
			refreshMerge()
		}, p.window)
	})
	addMergeBtn.Importance = widget.HighImportance

	mergePathEntry := widget.NewEntry()
	mergePathEntry.SetPlaceHolder("Tempel path PDF untuk digabung...")
	addMergePathBtn := widget.NewButtonWithIcon("Tambah Path", theme.ContentAddIcon(), func() {
		if mergePathEntry.Text != "" {
			if _, err := os.Stat(mergePathEntry.Text); err == nil {
				mergeFiles = append(mergeFiles, mergePathEntry.Text)
				mergePathEntry.SetText("")
				refreshMerge()
			}
		}
	})

	runMergeBtn := widget.NewButtonWithIcon(constants.BtnMergePDF, theme.ConfirmIcon(), func() {
		if len(mergeFiles) < 2 {
			dialog.ShowInformation("Perhatian", "Pilih minimal 2 berkas PDF untuk digabungkan.", p.window)
			return
		}
		outPath := filepath.Join(p.getDownloadDir(), "merged_document.pdf")
		if err := converter.MergePDFs(mergeFiles, outPath); err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		dialog.ShowInformation("Penggabungan Berhasil", fmt.Sprintf("PDF berhasil digabung ke:\n%s", outPath), p.window)
	})
	runMergeBtn.Importance = widget.HighImportance

	mergeCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("📑 Penggabungan Berkas PDF (Merge):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(addMergeBtn, widget.NewLabel("atau masukkan path:")),
		container.NewBorder(nil, nil, nil, addMergePathBtn, mergePathEntry),
		mergeList,
		runMergeBtn,
	), constants.ColorTechIndigo)

	// B. Split PDF Section
	splitFileEntry := widget.NewEntry()
	splitFileEntry.SetPlaceHolder("Pilih berkas PDF yang ingin dipisah...")

	browseSplitBtn := widget.NewButtonWithIcon("Pilih PDF", theme.FileIcon(), func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			splitFileEntry.SetText(closer.URI().Path())
		}, p.window)
	})
	browseSplitBtn.Importance = widget.HighImportance

	pageFromEntry := widget.NewEntry()
	pageFromEntry.SetText("1")
	pageToEntry := widget.NewEntry()
	pageToEntry.SetText("3")

	runSplitBtn := widget.NewButtonWithIcon(constants.BtnSplitPDF, theme.ConfirmIcon(), func() {
		src := strings.TrimSpace(splitFileEntry.Text)
		if src == "" {
			dialog.ShowInformation("Perhatian", "Pilih berkas PDF sumber terlebih dahulu.", p.window)
			return
		}
		from, _ := strconv.Atoi(pageFromEntry.Text)
		to, _ := strconv.Atoi(pageToEntry.Text)

		outPath := filepath.Join(p.getDownloadDir(), fmt.Sprintf("split_hal_%d_%d.pdf", from, to))
		if err := converter.SplitPDF(src, from, to, outPath); err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		dialog.ShowInformation("Pemisahan Berhasil", fmt.Sprintf("Halaman %d s/d %d berhasil diekstrak ke:\n%s", from, to, outPath), p.window)
	})
	runSplitBtn.Importance = widget.HighImportance

	splitCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("✂️ Pemisahan Halaman PDF (Split):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, widget.NewLabelWithStyle("PDF Sumber:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), browseSplitBtn, splitFileEntry),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Dari Halaman:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), pageFromEntry),
			container.NewVBox(widget.NewLabelWithStyle("Sampai Halaman:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), pageToEntry),
		),
		runSplitBtn,
	), constants.ColorWarning)

	// C. Compress PDF Section
	compFileEntry := widget.NewEntry()
	compFileEntry.SetPlaceHolder("Pilih berkas PDF yang ingin dikompres...")

	browseCompBtn := widget.NewButtonWithIcon("Pilih PDF", theme.FileIcon(), func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			compFileEntry.SetText(closer.URI().Path())
		}, p.window)
	})
	browseCompBtn.Importance = widget.HighImportance

	compLevelSelect := widget.NewSelect([]string{"screen (Paling Kecil / Ringan)", "ebook (Sedang / Rekomendasi)", "printer (Kualitas Tinggi)"}, nil)
	compLevelSelect.SetSelected("ebook (Sedang / Rekomendasi)")

	runCompBtn := widget.NewButtonWithIcon(constants.BtnCompressPDF, theme.ViewRefreshIcon(), func() {
		src := strings.TrimSpace(compFileEntry.Text)
		if src == "" {
			dialog.ShowInformation("Perhatian", "Pilih berkas PDF terlebih dahulu.", p.window)
			return
		}

		level := "ebook"
		if strings.Contains(compLevelSelect.Selected, "screen") {
			level = "screen"
		} else if strings.Contains(compLevelSelect.Selected, "printer") {
			level = "printer"
		}

		origStat, _ := os.Stat(src)
		origSize := int64(0)
		if origStat != nil {
			origSize = origStat.Size()
		}

		outPath := filepath.Join(p.getDownloadDir(), "compressed_"+filepath.Base(src))
		if err := converter.CompressPDF(src, level, outPath); err != nil {
			dialog.ShowError(err, p.window)
			return
		}

		newStat, _ := os.Stat(outPath)
		newSize := int64(0)
		if newStat != nil {
			newSize = newStat.Size()
		}

		savingPct := 0.0
		if origSize > 0 && newSize < origSize {
			savingPct = float64(origSize-newSize) / float64(origSize) * 100.0
		}

		dialog.ShowInformation("Kompresi Berhasil",
			fmt.Sprintf("PDF berhasil dikompres!\nUkuran Awal: %.1f KB ➔ Akhir: %.1f KB (Hemat %.1f%%)\nLokasi:\n%s",
				float64(origSize)/1024.0, float64(newSize)/1024.0, savingPct, outPath), p.window)
	})
	runCompBtn.Importance = widget.HighImportance

	compCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("🗜️ Kompresi Ukuran PDF:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, widget.NewLabelWithStyle("PDF Sumber:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), browseCompBtn, compFileEntry),
		widget.NewLabelWithStyle("Tingkat Kompresi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), compLevelSelect,
		runCompBtn,
	), constants.ColorDanger)

	return container.NewVScroll(container.NewVBox(
		mergeCard,
		splitCard,
		compCard,
	))
}
