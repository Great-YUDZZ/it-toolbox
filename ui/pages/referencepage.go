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
			cmdTxt := canvas.NewText(c.Command, constants.ColorTextPrimary)
			cmdTxt.TextSize = constants.FontSizeBody
			cmdTxt.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

			copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), action)
			copyBtn.Importance = widget.LowImportance

			header := container.NewBorder(nil, nil,
				container.NewHBox(components.BadgeYellow(c.Category), cmdTxt),
				copyBtn,
			)

			descLabel := widget.NewLabel(c.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			exTxt := canvas.NewText("Contoh: "+c.Example, constants.ColorTextSecondary)
			exTxt.TextSize = constants.FontSizeSmall
			exTxt.TextStyle = fyne.TextStyle{Monospace: true}

			cardContent := container.NewVBox(
				header,
				widget.NewSeparator(),
				descLabel,
				exTxt,
			)
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorInfo))
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

			portTxt := canvas.NewText(fmt.Sprintf("Port %d", pt.Port), constants.ColorTextPrimary)
			portTxt.TextSize = constants.FontSizeH2
			portTxt.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

			protoBadge := components.BadgeYellow(pt.Protocol)
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
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, constants.ColorWarning))
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
			var accentColor color.Color
			switch {
			case cd.Code >= 500:
				statusBadge = components.BadgeDanger(fmt.Sprintf("%d", cd.Code))
				accentColor = constants.ColorDanger
			case cd.Code >= 400:
				statusBadge = components.BadgeWarning(fmt.Sprintf("%d", cd.Code))
				accentColor = constants.ColorWarning
			case cd.Code >= 300:
				statusBadge = components.BadgeIndigo(fmt.Sprintf("%d", cd.Code))
				accentColor = constants.ColorInfo
			case cd.Code >= 200:
				statusBadge = components.BadgeSuccess(fmt.Sprintf("%d", cd.Code))
				accentColor = constants.ColorSuccess
			default:
				statusBadge = components.BadgeCyan(fmt.Sprintf("%d", cd.Code))
				accentColor = constants.ColorAccentCyan
			}

			titleTxt := canvas.NewText(cd.Name, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeBody
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
			listContainer.Add(components.NewPlainCardWithAccent(cardContent, accentColor))
		}
		listContainer.Refresh()
	}

	searchBar := components.NewSearchBar(constants.SearchHTTPPlaceholder, renderList)
	renderList("")

	scroll := container.NewVScroll(listContainer)
	return container.NewBorder(container.NewPadded(searchBar.Container), nil, nil, nil, scroll)
}
