package pages

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/core/reference"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

type ReferencePage struct {
	window fyne.Window
}

func NewReferencePage(win fyne.Window) *ReferencePage {
	return &ReferencePage{window: win}
}

func (p *ReferencePage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavReference,
		"Kumpulan referensi cepat untuk Terminal Commands, Port Jaringan, dan HTTP Status Codes.",
		components.BadgeCyan("DOCS & REFERENCE"),
	)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabCommands, theme.DocumentIcon(), p.buildCommandsTab()),
		container.NewTabItemWithIcon(constants.TabPorts, theme.RadioButtonIcon(), p.buildPortsTab()),
		container.NewTabItemWithIcon(constants.TabHTTPCodes, theme.InfoIcon(), p.buildHTTPCodesTab()),
	)

	return container.NewBorder(hero, nil, nil, nil, tabs)
}

func (p *ReferencePage) copyToClip(txt string) {
	p.window.Clipboard().SetContent(txt)
	dialog.ShowInformation("Clipboard", constants.StatusCopied, p.window)
}

// 1. Tab Commands
func (p *ReferencePage) buildCommandsTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

	renderList := func(query string) {
		listContainer.Objects = nil
		items := reference.SearchCommands(query)
		if len(items) == 0 {
			listContainer.Add(widget.NewLabel(constants.NoDataMessage))
			listContainer.Refresh()
			return
		}

		for _, cmd := range items {
			c := cmd
			action := func() {
				p.copyToClip(c.Command)
			}

			// Code box for command
			cmdTxt := canvas.NewText(c.Command, color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0xFF})
			cmdTxt.TextSize = 13.5
			cmdTxt.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

			copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), action)
			copyBtn.Importance = widget.LowImportance

			header := container.NewBorder(nil, nil,
				container.NewHBox(components.BadgeIndigo(c.Category), cmdTxt),
				copyBtn,
			)

			descLabel := widget.NewLabel(c.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			exTxt := canvas.NewText("Contoh: "+c.Example, color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF})
			exTxt.TextSize = 11.5
			exTxt.TextStyle = fyne.TextStyle{Monospace: true}

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				descLabel,
				exTxt,
			)
			listContainer.Add(components.NewPlainCard(cardContent))
		}
		listContainer.Refresh()
	}

	searchBar := components.NewSearchBar(constants.SearchCommandPlaceholder, renderList)
	renderList("")

	scroll := container.NewVScroll(listContainer)
	return container.NewBorder(container.NewPadded(searchBar.Container), nil, nil, nil, scroll)
}

// 2. Tab Ports
func (p *ReferencePage) buildPortsTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

	renderList := func(query string) {
		listContainer.Objects = nil
		items := reference.SearchPorts(query)
		if len(items) == 0 {
			listContainer.Add(widget.NewLabel(constants.NoDataMessage))
			listContainer.Refresh()
			return
		}

		for _, port := range items {
			pt := port
			action := func() {
				p.copyToClip(fmt.Sprintf("%d", pt.Port))
			}

			portTxt := canvas.NewText(fmt.Sprintf("Port %d", pt.Port), color.RGBA{R: 0x00, G: 0xD4, B: 0xFF, A: 0xFF})
			portTxt.TextSize = 14
			portTxt.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

			protoBadge := components.BadgeIndigo(pt.Protocol)
			serviceBadge := components.BadgeCyan(pt.Service)

			copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), action)
			copyBtn.Importance = widget.LowImportance

			header := container.NewBorder(nil, nil,
				container.NewHBox(portTxt, protoBadge, serviceBadge),
				copyBtn,
			)

			descLabel := widget.NewLabel(pt.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				descLabel,
			)
			listContainer.Add(components.NewPlainCard(cardContent))
		}
		listContainer.Refresh()
	}

	searchBar := components.NewSearchBar(constants.SearchPortPlaceholder, renderList)
	renderList("")

	scroll := container.NewVScroll(listContainer)
	return container.NewBorder(container.NewPadded(searchBar.Container), nil, nil, nil, scroll)
}

// 3. Tab HTTP Status Codes
func (p *ReferencePage) buildHTTPCodesTab() fyne.CanvasObject {
	listContainer := container.NewVBox()

	renderList := func(query string) {
		listContainer.Objects = nil
		items := reference.SearchHTTPCodes(query)
		if len(items) == 0 {
			listContainer.Add(widget.NewLabel(constants.NoDataMessage))
			listContainer.Refresh()
			return
		}

		for _, code := range items {
			cd := code
			action := func() {
				p.copyToClip(fmt.Sprintf("%d %s", cd.Code, cd.Name))
			}

			var statusBadge fyne.CanvasObject
			switch {
			case cd.Code >= 500:
				statusBadge = components.BadgeError(fmt.Sprintf("%d", cd.Code))
			case cd.Code >= 400:
				statusBadge = components.BadgeWarning(fmt.Sprintf("%d", cd.Code))
			case cd.Code >= 300:
				statusBadge = components.BadgeIndigo(fmt.Sprintf("%d", cd.Code))
			case cd.Code >= 200:
				statusBadge = components.BadgeSuccess(fmt.Sprintf("%d", cd.Code))
			default:
				statusBadge = components.BadgeCyan(fmt.Sprintf("%d", cd.Code))
			}

			titleTxt := canvas.NewText(cd.Name, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
			titleTxt.TextSize = 13.5
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), action)
			copyBtn.Importance = widget.LowImportance

			header := container.NewBorder(nil, nil,
				container.NewHBox(statusBadge, titleTxt),
				copyBtn,
			)

			descLabel := widget.NewLabel(cd.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			causeLabel := widget.NewLabel(fmt.Sprintf("Penyebab Umum: %s", cd.Cause))
			causeLabel.Wrapping = fyne.TextWrapWord

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				descLabel,
				causeLabel,
			)
			listContainer.Add(components.NewPlainCard(cardContent))
		}
		listContainer.Refresh()
	}

	searchBar := components.NewSearchBar(constants.SearchHTTPPlaceholder, renderList)
	renderList("")

	scroll := container.NewVScroll(listContainer)
	return container.NewBorder(container.NewPadded(searchBar.Container), nil, nil, nil, scroll)
}
