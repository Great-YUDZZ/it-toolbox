package pages

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2/canvas"
	"github.com/yudz/it-toolbox/core/converter"
	"github.com/yudz/it-toolbox/core/qrcode"
	"github.com/yudz/it-toolbox/core/youtube"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

type FileConverterPage struct {
	window       fyne.Window
	qrFileServer *qrcode.FileShareServer
}

func NewFileConverterPage(win fyne.Window) *FileConverterPage {
	return &FileConverterPage{window: win}
}

func (p *FileConverterPage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavFileConverter,
		"Konversi dokumen, gambar, ekstraksi teks OCR, kompresi berkas, manipulasi PDF, generator QR code, serta unduh video YouTube dengan resolusi pilihan.",
		components.BadgeCyan("FILE CONVERTER PRO"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabDocConverter, theme.DocumentIcon(), p.buildDocConverterTab()),
		container.NewTabItemWithIcon(constants.TabImgConverter, theme.FileImageIcon(), p.buildImgConverterTab()),
		container.NewTabItemWithIcon(constants.TabOCR, theme.VisibilityIcon(), p.buildOCRTab()),
		container.NewTabItemWithIcon(constants.TabPDFTools, theme.FolderOpenIcon(), p.buildPDFToolsTab()),
		container.NewTabItemWithIcon(constants.TabQRCode, theme.ViewFullScreenIcon(), p.buildQRCodeTab()),
		container.NewTabItemWithIcon(constants.TabYouTube, theme.DownloadIcon(), p.buildYouTubeDownloaderTab()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
}

func (p *FileConverterPage) BuildYouTubePage() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavYouTube,
		"Unduh video YouTube dengan resolusi pilihan yang tersedia (4K 2160p, 1080p Full HD, 720p HD, 360p, hingga Audio M4A).",
		components.BadgeCyan("YOUTUBE ENGINE PRO"),
	)
	return container.NewBorder(hero, nil, nil, nil, p.buildYouTubeDownloaderTab())
}

func (p *FileConverterPage) BuildQRCodePage() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.TabQRCode,
		"Buat QR Code instan dari tautan URL website, teks bebas, atau berbagi foto, video, dan berkas apapun langsung via jaringan Wi-Fi lokal.",
		components.BadgeCyan("QR CODE ENGINE PRO"),
	)
	return container.NewBorder(hero, nil, nil, nil, p.buildQRCodeTab())
}

func (p *FileConverterPage) openTarget(target string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", target)
	case "darwin":
		cmd = exec.Command("open", target)
	default:
		cmd = exec.Command("xdg-open", target)
	}
	_ = cmd.Start()
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
				container.NewHBox(resBadge, widget.NewLabelWithStyle(fmt.Sprintf("%s (Awal: %.1f KB -> Akhir: %.1f KB)",
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
		widget.NewLabelWithStyle("Penggabungan Berkas PDF (Merge):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
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
		widget.NewLabelWithStyle("Pemisahan Halaman PDF (Split):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
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
			fmt.Sprintf("PDF berhasil dikompres!\nUkuran Awal: %.1f KB -> Akhir: %.1f KB (Hemat %.1f%%)\nLokasi:\n%s",
				float64(origSize)/1024.0, float64(newSize)/1024.0, savingPct, outPath), p.window)
	})
	runCompBtn.Importance = widget.HighImportance

	compCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Kompresi Ukuran PDF:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
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

// ----------------------------------------------------------------------------
// 5. Tab Unduh Video YouTube dengan Pilihan Resolusi
// ----------------------------------------------------------------------------
func (p *FileConverterPage) buildYouTubeDownloaderTab() fyne.CanvasObject {
	ytDownloader := youtube.NewDownloader()

	// 1. Input URL Section
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Tempel link YouTube (cth: https://www.youtube.com/watch?v=... atau youtu.be/...)")

	pasteBtn := widget.NewButtonWithIcon("Tempel Link", theme.ContentPasteIcon(), func() {
		clip := p.window.Clipboard().Content()
		if clip != "" {
			urlEntry.SetText(strings.TrimSpace(clip))
		}
	})
	pasteBtn.Importance = widget.MediumImportance

	checkBtn := widget.NewButtonWithIcon(constants.BtnFetchVideo, theme.SearchIcon(), nil)
	checkBtn.Importance = widget.HighImportance

	inputCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Masukkan Tautan Video YouTube:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, container.NewHBox(pasteBtn, checkBtn), urlEntry),
	), constants.ColorAccentCyan)

	// 2. Video Details & Resolution Selection Section
	var currentDetails *youtube.VideoDetails
	var selectedOption *youtube.ResolutionOption

	titleLabel := widget.NewLabelWithStyle("Belum ada video yang dimuat.", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	titleLabel.Wrapping = fyne.TextWrapWord

	badgeRow := container.NewHBox()
	resolutionSelect := widget.NewSelect([]string{"-- Masukkan URL & Klik 'Periksa Video' --"}, nil)

	destDirEntry := widget.NewEntry()
	destDirEntry.SetText(p.getDownloadDir())

	chooseDirBtn := widget.NewButtonWithIcon("Pilih Folder", theme.FolderOpenIcon(), func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil && err == nil {
				destDirEntry.SetText(uri.Path())
			}
		}, p.window)
	})

	destDirRow := container.NewBorder(nil, nil,
		widget.NewLabelWithStyle("Simpan ke:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		chooseDirBtn,
		destDirEntry,
	)

	downloadProgressBar := widget.NewProgressBar()
	downloadProgressBar.SetValue(0)
	downloadProgressBar.Hide()

	statusLabel := widget.NewLabel("Tempelkan tautan video YouTube di atas, lalu tekan tombol 'Periksa Video'.")
	statusLabel.Wrapping = fyne.TextWrapWord

	openFileBtn := widget.NewButtonWithIcon("Buka Berkas", theme.MediaPlayIcon(), nil)
	openFileBtn.Hide()
	openFolderBtn := widget.NewButtonWithIcon("Buka Folder", theme.FolderOpenIcon(), nil)
	openFolderBtn.Hide()
	actionBtns := container.NewHBox(openFileBtn, openFolderBtn)

	downloadBtn := widget.NewButtonWithIcon(constants.BtnDownloadVideo, theme.DownloadIcon(), nil)
	downloadBtn.Importance = widget.HighImportance
	downloadBtn.Disable()

	openTarget := func(path string) {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", "", path)
		case "darwin":
			cmd = exec.Command("open", path)
		default:
			cmd = exec.Command("xdg-open", path)
		}
		_ = cmd.Start()
	}

	detailsCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Informasi Video & Resolusi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		titleLabel,
		badgeRow,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Pilihan Kualitas / Resolusi yang Tersedia:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		resolutionSelect,
		destDirRow,
		container.NewHBox(downloadBtn, actionBtns),
		downloadProgressBar,
		statusLabel,
	), constants.ColorTechIndigo)

	// Action: Check video
	checkBtn.OnTapped = func() {
		rawURL := strings.TrimSpace(urlEntry.Text)
		if rawURL == "" {
			dialog.ShowInformation("Perhatian", "Harap masukkan tautan video YouTube terlebih dahulu.", p.window)
			return
		}

		checkBtn.Disable()
		downloadBtn.Disable()
		openFileBtn.Hide()
		openFolderBtn.Hide()
		titleLabel.SetText("Sedang memeriksa data video dari YouTube...")
		badgeRow.Objects = nil
		badgeRow.Refresh()
		statusLabel.SetText("Menghubungi server YouTube dan mengekstrak daftar resolusi...")

		go func() {
			details, err := ytDownloader.FetchVideoDetails(rawURL)
			if err != nil {
				fyne.Do(func() {
					checkBtn.Enable()
					titleLabel.SetText("Gagal memuat data video.")
					statusLabel.SetText(fmt.Sprintf("Kesalahan: %v", err))
					dialog.ShowError(err, p.window)
				})
				return
			}

			fyne.Do(func() {
				checkBtn.Enable()
				currentDetails = details
				titleLabel.SetText(details.Title)

				badgeRow.Objects = nil
				badgeRow.Add(components.BadgeCyan("" + details.Author))
				badgeRow.Add(components.BadgeYellow("" + details.DurationStr))
				badgeRow.Add(components.BadgeSuccess(fmt.Sprintf("%d Pilihan Kualitas", len(details.Options))))
				if youtube.HasFFmpeg() {
					badgeRow.Add(components.BadgeSuccess("FFmpeg Aktif (Auto-Mux)"))
				} else {
					badgeRow.Add(components.BadgeWarning("Unduh Direct (Tanpa FFmpeg)"))
				}
				badgeRow.Refresh()

				var optLabels []string
				for _, opt := range details.Options {
					optLabels = append(optLabels, opt.DisplayLabel)
				}

				resolutionSelect.Options = optLabels
				if len(optLabels) > 0 {
					resolutionSelect.SetSelected(optLabels[0])
					firstOpt := details.Options[0]
					selectedOption = &firstOpt
					downloadBtn.Enable()
				}

				resolutionSelect.OnChanged = func(selected string) {
					for _, opt := range currentDetails.Options {
						if opt.DisplayLabel == selected {
							copyOpt := opt
							selectedOption = &copyOpt
							break
						}
					}
				}

				statusLabel.SetText(fmt.Sprintf("Siap! Ditemukan %d pilihan resolusi sesuai video ini. Pilih kualitas dan klik Unduh.", len(details.Options)))
			})
		}()
	}

	// Action: Download
	downloadBtn.OnTapped = func() {
		if currentDetails == nil || selectedOption == nil {
			dialog.ShowInformation("Perhatian", "Pilih resolusi video terlebih dahulu.", p.window)
			return
		}

		destDir := strings.TrimSpace(destDirEntry.Text)
		if destDir == "" {
			destDir = p.getDownloadDir()
		}

		downloadBtn.Disable()
		checkBtn.Disable()
		downloadProgressBar.SetValue(0)
		downloadProgressBar.Show()
		openFileBtn.Hide()
		openFolderBtn.Hide()
		statusLabel.SetText("Memulai streaming unduhan...")

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 45*time.Minute)
			defer cancel()

			finalPath, err := ytDownloader.Download(ctx, currentDetails, *selectedOption, destDir,
				func(downloaded, total int64, pct, speed float64, status string) {
					fyne.Do(func() {
						downloadProgressBar.SetValue(pct / 100.0)
						statusLabel.SetText(fmt.Sprintf("[%s] %.1f%% — %s", selectedOption.QualityLabel, pct, status))
					})
				})

			fyne.Do(func() {
				downloadBtn.Enable()
				checkBtn.Enable()

				if err != nil {
					statusLabel.SetText(fmt.Sprintf("Unduhan gagal: %v", err))
					dialog.ShowError(err, p.window)
					return
				}

				downloadProgressBar.SetValue(1.0)
				statusLabel.SetText(fmt.Sprintf("Berhasil disimpan ke:\n%s", finalPath))

				openFileBtn.OnTapped = func() {
					openTarget(finalPath)
				}
				openFolderBtn.OnTapped = func() {
					openTarget(filepath.Dir(finalPath))
				}
				openFileBtn.Show()
				openFolderBtn.Show()

				dialog.ShowInformation("Unduhan Selesai!",
					fmt.Sprintf("Video YouTube berhasil diunduh:\n\nJudul: %s\nKualitas: %s\nLokasi:\n%s",
						currentDetails.Title, selectedOption.QualityLabel, finalPath),
					p.window)
			})
		}()
	}

	return container.NewVScroll(container.NewVBox(
		inputCard,
		detailsCard,
	))
}

// ----------------------------------------------------------------------------
// Generator QR Code & Berbagi Berkas Wi-Fi
// ----------------------------------------------------------------------------
func (p *FileConverterPage) buildQRCodeTab() fyne.CanvasObject {
	currentMode := "text" // "text" or "file"
	var currentPNGBytes []byte
	var currentContent string

	schemes := qrcode.GetAvailableColorSchemes()
	var schemeNames []string
	for _, s := range schemes {
		schemeNames = append(schemeNames, s.Name)
	}

	// Preview Image Widget
	qrImg := canvas.NewImageFromResource(theme.ViewFullScreenIcon())
	qrImg.FillMode = canvas.ImageFillContain
	qrImg.SetMinSize(fyne.NewSize(220, 220))

	// Status & Content Labels
	statusLabel := widget.NewLabel("Masukkan tautan atau pilih berkas untuk membuat QR code.")
	statusLabel.Wrapping = fyne.TextWrapWord
	statsLabel := widget.NewLabel("")

	contentPreviewEntry := components.NewScrollableEntry()
	contentPreviewEntry.Disable()
	contentPreviewEntry.SetPlaceHolder("Konten atau tautan QR code akan tampil di sini...")

	// Options
	sizeSelect := widget.NewSelect([]string{"256 px", "512 px (Standar)", "1024 px (HD)"}, nil)
	sizeSelect.SetSelected("512 px (Standar)")

	levelSelect := widget.NewSelect([]string{
		"Medium (15% - Rekomendasi)",
		"Low (7%)",
		"High (25%)",
		"Highest (30%)",
	}, nil)
	levelSelect.SetSelected("Medium (15% - Rekomendasi)")

	paletteSelect := widget.NewSelect(schemeNames, nil)
	if len(schemeNames) > 0 {
		paletteSelect.SetSelected(schemeNames[0])
	}

	var textEntry *widget.Entry

	// Function to regenerate QR code image
	updateQRCode := func() {
		var rawText string
		if currentMode == "file" {
			if p.qrFileServer != nil && p.qrFileServer.IsRunning() {
				rawText = p.qrFileServer.GetShareURL()
			} else {
				rawText = ""
			}
		} else {
			if textEntry != nil {
				rawText = strings.TrimSpace(textEntry.Text)
			}
		}

		if rawText == "" {
			if currentMode == "file" {
				statusLabel.SetText("Pilih berkas lalu klik 'Mulai Berbagi (Wi-Fi)' untuk mengaktifkan server dan membuat QR code.")
			} else {
				statusLabel.SetText("Ketikkan teks atau URL website untuk membuat QR code.")
			}
			return
		}

		currentContent = rawText
		contentPreviewEntry.SetText(rawText)

		// Determine size
		size := 512
		if strings.HasPrefix(sizeSelect.Selected, "256") {
			size = 256
		} else if strings.HasPrefix(sizeSelect.Selected, "1024") {
			size = 1024
		}

		// Determine level
		recLevel := qrcode.ParseRecoveryLevel(levelSelect.Selected)

		// Determine palette
		selectedScheme := schemes[0]
		for _, s := range schemes {
			if s.Name == paletteSelect.Selected {
				selectedScheme = s
				break
			}
		}

		pngBytes, err := qrcode.GenerateCustomQRCode(rawText, size, recLevel, selectedScheme.FG, selectedScheme.BG)
		if err != nil {
			statusLabel.SetText(fmt.Sprintf("Gagal membuat QR code: %v", err))
			return
		}

		currentPNGBytes = pngBytes
		res := fyne.NewStaticResource("qrcode.png", pngBytes)
		qrImg.Resource = res
		qrImg.Refresh()
		statusLabel.SetText(fmt.Sprintf("QR Code siap (%d x %d px) - Level %s", size, size, levelSelect.Selected))
	}

	sizeSelect.OnChanged = func(_ string) { updateQRCode() }
	levelSelect.OnChanged = func(_ string) { updateQRCode() }
	paletteSelect.OnChanged = func(_ string) { updateQRCode() }

	// Mode 1: Text / URL inputs
	textEntry = components.NewScrollableEntry()
	textEntry.SetPlaceHolder("Contoh: https://mywebsite.com atau pesan teks...")
	textEntry.SetText("https://github.com/Great-YUDZZ/it-toolbox")
	textEntry.OnChanged = func(_ string) {
		if currentMode == "text" {
			updateQRCode()
		}
	}

	presetWebBtn := widget.NewButtonWithIcon("Web HTTPS", theme.NavigateNextIcon(), func() {
		textEntry.SetText("https://")
		updateQRCode()
	})
	presetWifiBtn := widget.NewButtonWithIcon("Wi-Fi Config", theme.SettingsIcon(), func() {
		textEntry.SetText("WIFI:T:WPA;S:NamaSSID;P:PasswordWiFi;;")
		updateQRCode()
	})
	presetWaBtn := widget.NewButtonWithIcon("WhatsApp", theme.MailComposeIcon(), func() {
		textEntry.SetText("https://wa.me/6281234567890?text=Halo%20IT%20Toolbox")
		updateQRCode()
	})
	presetMailBtn := widget.NewButtonWithIcon("Email", theme.MailSendIcon(), func() {
		textEntry.SetText("mailto:admin@example.com?subject=Halo%20IT%20Toolbox")
		updateQRCode()
	})

	presetRow := container.NewVBox(
		widget.NewLabelWithStyle("Preset Cepat:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHBox(presetWebBtn, presetWifiBtn, presetWaBtn, presetMailBtn),
	)

	textConfigContainer := container.NewVBox(
		widget.NewLabelWithStyle("Masukkan Tautan URL atau Teks Bebas:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		textEntry,
		presetRow,
	)

	// Mode 2: File Sharing inputs
	filePathEntry := components.NewScrollableEntry()
	filePathEntry.SetPlaceHolder("Pilih berkas foto, video, dokumen, atau berkas apapun...")

	fileBadgeContainer := container.NewHBox()
	fileInfoLabel := widget.NewLabel("Belum ada berkas dipilih.")

	updateFileInfo := func(path string) {
		fileBadgeContainer.Objects = nil
		if path == "" {
			fileInfoLabel.SetText("Belum ada berkas dipilih.")
			fileBadgeContainer.Refresh()
			return
		}
		fi, err := os.Stat(path)
		if err != nil {
			fileInfoLabel.SetText(fmt.Sprintf("Berkas tidak valid: %v", err))
			fileBadgeContainer.Refresh()
			return
		}

		sizeStr := qrcode.FormatFileSize(fi.Size())
		ext := strings.ToLower(filepath.Ext(path))
		var badge fyne.CanvasObject
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp":
			badge = components.BadgeCyan("Foto / Gambar")
		case ".mp4", ".mkv", ".avi", ".mov", ".webm", ".flv":
			badge = components.BadgeIndigo("Video")
		case ".mp3", ".wav", ".ogg", ".flac", ".m4a":
			badge = components.BadgeYellow("Audio")
		case ".pdf", ".docx", ".xlsx", ".pptx", ".txt", ".zip", ".tar", ".gz":
			badge = components.BadgeSuccess("Dokumen / Arsip")
		default:
			badge = components.BadgeMuted("Berkas Umum")
		}

		fileBadgeContainer.Add(badge)
		fileBadgeContainer.Refresh()
		fileInfoLabel.SetText(fmt.Sprintf("%s (%s)", filepath.Base(path), sizeStr))
	}

	browseFileBtn := widget.NewButtonWithIcon("Pilih Berkas...", theme.FileIcon(), func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			path := closer.URI().Path()
			filePathEntry.SetText(path)
			updateFileInfo(path)
		}, p.window)
	})
	browseFileBtn.Importance = widget.HighImportance

	filePathEntry.OnChanged = func(path string) {
		updateFileInfo(strings.TrimSpace(path))
	}

	// Network IP & Port selection
	localIPs, _ := qrcode.GetLocalIPv4Addresses()
	var ipOptions []string
	for _, ip := range localIPs {
		ipOptions = append(ipOptions, ip)
	}
	if len(ipOptions) == 0 {
		ipOptions = append(ipOptions, "127.0.0.1")
	}

	ipSelect := widget.NewSelect(ipOptions, nil)
	ipSelect.SetSelected(ipOptions[0])

	refreshIPBtn := widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		freshIPs, _ := qrcode.GetLocalIPv4Addresses()
		if len(freshIPs) > 0 {
			ipSelect.Options = freshIPs
			ipSelect.SetSelected(freshIPs[0])
		}
	})
	refreshIPBtn.Importance = widget.LowImportance

	portEntry := components.NewScrollableEntry()
	portEntry.SetText("8765")

	serverBadgeContainer := container.NewHBox(components.BadgeMuted("Server Nonaktif"))
	var toggleServerBtn *widget.Button

	stopServerFunc := func() {
		if p.qrFileServer != nil && p.qrFileServer.IsRunning() {
			_ = p.qrFileServer.Stop()
		}
		toggleServerBtn.SetText("Mulai Berbagi (Wi-Fi)")
		toggleServerBtn.SetIcon(theme.MediaPlayIcon())
		toggleServerBtn.Importance = widget.HighImportance
		serverBadgeContainer.Objects = []fyne.CanvasObject{components.BadgeMuted("Server Nonaktif")}
		serverBadgeContainer.Refresh()
		statsLabel.SetText("")
		updateQRCode()
	}

	toggleServerBtn = widget.NewButtonWithIcon("Mulai Berbagi (Wi-Fi)", theme.MediaPlayIcon(), func() {
		if p.qrFileServer != nil && p.qrFileServer.IsRunning() {
			stopServerFunc()
			return
		}

		filePath := strings.TrimSpace(filePathEntry.Text)
		if filePath == "" {
			dialog.ShowInformation("Perhatian", "Pilih berkas yang ingin dibagikan terlebih dahulu.", p.window)
			return
		}
		if _, err := os.Stat(filePath); err != nil {
			dialog.ShowError(fmt.Errorf("berkas tidak ditemukan: %w", err), p.window)
			return
		}

		ip := ipSelect.Selected
		if ip == "" {
			ip = "0.0.0.0"
		}
		port, err := strconv.Atoi(strings.TrimSpace(portEntry.Text))
		if err != nil || port <= 0 || port > 65535 {
			port = 8765
			portEntry.SetText("8765")
		}

		srv := qrcode.NewFileShareServer()
		shareURL, err := srv.Start(filePath, port, ip)
		if err != nil {
			dialog.ShowError(fmt.Errorf("gagal memulai server Wi-Fi: %w", err), p.window)
			return
		}

		p.qrFileServer = srv
		toggleServerBtn.SetText("Hentikan Berbagi")
		toggleServerBtn.SetIcon(theme.MediaStopIcon())
		toggleServerBtn.Importance = widget.DangerImportance
		serverBadgeContainer.Objects = []fyne.CanvasObject{components.BadgeSuccess("Server Aktif (Berbagi Berkas)")}
		serverBadgeContainer.Refresh()

		updateQRCode()

		// Periodic status update while server is active
		go func(activeSrv *qrcode.FileShareServer, expectedURL string) {
			for activeSrv.IsRunning() {
				time.Sleep(1 * time.Second)
				_, _, _, count := activeSrv.GetFileInfo()
				fyne.Do(func() {
					if activeSrv.IsRunning() {
						statsLabel.SetText(fmt.Sprintf("Tautan: %s | Total Akses / Unduhan: %d kali", expectedURL, count))
					}
				})
			}
		}(srv, shareURL)
	})
	toggleServerBtn.Importance = widget.HighImportance

	networkRow := container.NewHBox(
		widget.NewLabelWithStyle("IP Host:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		ipSelect,
		refreshIPBtn,
		widget.NewLabelWithStyle("Port:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		portEntry,
		serverBadgeContainer,
	)

	fileConfigContainer := container.NewVBox(
		widget.NewLabelWithStyle("Pilih Berkas yang Ingin Dibagikan:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, browseFileBtn, filePathEntry),
		container.NewHBox(fileBadgeContainer, fileInfoLabel),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Pengaturan Jaringan Server Wi-Fi Lokal:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		networkRow,
		toggleServerBtn,
		widget.NewLabel("Pastikan smartphone terhubung pada jaringan Wi-Fi / Hotspot yang sama."),
	)
	fileConfigContainer.Hide()

	// Mode switch radio
	modeRadio := widget.NewRadioGroup([]string{
		"Tautan URL / Teks Bebas",
		"Berbagi Berkas (Foto, Video, File) via Wi-Fi",
	}, func(selected string) {
		if selected == "Tautan URL / Teks Bebas" {
			currentMode = "text"
			fileConfigContainer.Hide()
			textConfigContainer.Show()
		} else {
			currentMode = "file"
			textConfigContainer.Hide()
			fileConfigContainer.Show()
		}
		updateQRCode()
	})
	modeRadio.Horizontal = false
	modeRadio.SetSelected("Tautan URL / Teks Bebas")

	// Customization card
	customCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Pengaturan Ukuran & Gaya QR Code:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3,
			container.NewVBox(canvas.NewText("Ukuran Gambar:", constants.ColorTextPrimary), sizeSelect),
			container.NewVBox(canvas.NewText("Tingkat Koreksi Error:", constants.ColorTextPrimary), levelSelect),
			container.NewVBox(canvas.NewText("Tema Palet Warna:", constants.ColorTextPrimary), paletteSelect),
		),
	), constants.ColorAccentCyan)

	// Action buttons
	copyBtn := widget.NewButtonWithIcon("Salin Tautan / Teks", theme.ContentCopyIcon(), func() {
		if currentContent == "" {
			dialog.ShowInformation("Perhatian", "Belum ada konten QR code untuk disalin.", p.window)
			return
		}
		p.copyToClip(currentContent)
	})

	savePNGBtn := widget.NewButtonWithIcon("Simpan Gambar (PNG)", theme.DocumentSaveIcon(), func() {
		if len(currentPNGBytes) == 0 {
			dialog.ShowInformation("Perhatian", "Belum ada gambar QR code untuk disimpan.", p.window)
			return
		}
		saveDlg := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
			if err != nil || uc == nil {
				return
			}
			defer uc.Close()
			_, err = uc.Write(currentPNGBytes)
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			dialog.ShowInformation("Tersimpan", fmt.Sprintf("Berkas QR Code PNG berhasil disimpan:\n%s", uc.URI().Path()), p.window)
		}, p.window)
		saveDlg.SetFileName("qrcode.png")
		saveDlg.Show()
	})
	savePNGBtn.Importance = widget.HighImportance

	openBrowserBtn := widget.NewButtonWithIcon("Buka di Browser", theme.NavigateNextIcon(), func() {
		if currentContent == "" {
			dialog.ShowInformation("Perhatian", "Belum ada URL untuk dibuka.", p.window)
			return
		}
		if !strings.HasPrefix(currentContent, "http://") && !strings.HasPrefix(currentContent, "https://") {
			dialog.ShowInformation("Perhatian", "Konten bukan merupakan URL web (http/https).", p.window)
			return
		}
		p.openTarget(currentContent)
	})

	actionButtonsRow := container.NewHBox(
		savePNGBtn,
		copyBtn,
		openBrowserBtn,
	)

	qrBox := container.NewCenter(
		container.NewPadded(qrImg),
	)

	previewCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Pratinjau QR Code & Ekspor:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		qrBox,
		statusLabel,
		statsLabel,
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle("Konten Terkode:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			nil,
			contentPreviewEntry,
		),
		actionButtonsRow,
	), constants.ColorSuccess)

	// Initial render
	updateQRCode()

	return container.NewVScroll(container.NewVBox(
		components.NewPlainCardWithAccent(container.NewVBox(
			widget.NewLabelWithStyle("Pilih Mode Pembuatan QR Code:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			modeRadio,
			widget.NewSeparator(),
			textConfigContainer,
			fileConfigContainer,
		), constants.ColorAccentCyan),
		customCard,
		previewCard,
	))
}
