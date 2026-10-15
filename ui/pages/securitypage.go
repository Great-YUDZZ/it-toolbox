package pages

import (
	"context"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/core/security"
	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

type SecurityPage struct {
	window  fyne.Window
	appTabs *container.AppTabs
}

func NewSecurityPage(win fyne.Window) *SecurityPage {
	return &SecurityPage{window: win}
}

func (p *SecurityPage) copyToClip(txt string) {
	if p.window == nil {
		return
	}
	p.window.Clipboard().SetContent(txt)
	dialog.ShowInformation("Clipboard", constants.StatusCopied, p.window)
}

func (p *SecurityPage) Build() fyne.CanvasObject {
	hero := components.NewHeroHeader(
		constants.NavSecurity,
		"Audit keamanan port jaringan, inspeksi sertifikat SSL/TLS & HTTP headers, forensik multi-hash dan analisis log anomali, dekoder JWT & tabel ASCII, serta playbook pengerasan sistem.",
		components.BadgeIndigo("CYBER SECURITY HUB PRO"),
	)

	p.appTabs = container.NewAppTabs(
		container.NewTabItemWithIcon(constants.TabSecurityScanner, theme.SearchIcon(), p.buildPortScannerTab()),
		container.NewTabItemWithIcon(constants.TabSecurityWebSSL, theme.VisibilityIcon(), p.buildWebSSLTab()),
		container.NewTabItemWithIcon(constants.TabSecurityForensics, theme.FileIcon(), p.buildForensicsTab()),
		container.NewTabItemWithIcon(constants.TabSecurityEncoders, theme.SettingsIcon(), p.buildEncodersASCIITab()),
		container.NewTabItemWithIcon(constants.TabSecurityPlaybook, theme.DocumentIcon(), p.buildPlaybookTab()),
	)

	return container.NewBorder(hero, nil, nil, nil, p.appTabs)
}

func (p *SecurityPage) SelectTab(tabName string) {
	if p.appTabs == nil || tabName == "" {
		return
	}
	for _, item := range p.appTabs.Items {
		if strings.EqualFold(item.Text, tabName) {
			p.appTabs.Select(item)
			return
		}
	}
}

// ----------------------------------------------------------------------------
// 1. Tab Port Scanner
// ----------------------------------------------------------------------------
func (p *SecurityPage) buildPortScannerTab() fyne.CanvasObject {
	hostEntry := components.NewScrollableEntry()
	hostEntry.SetPlaceHolder("IP target atau domain (contoh: 127.0.0.1, scanme.nmap.org)")
	hostEntry.SetText("127.0.0.1")

	presetSelect := widget.NewSelect([]string{
		"Top 20 Common Ports",
		"Web & DB Ports (80, 443, 3306, 5432, 6379, 8080, 8443, 27017)",
		"SSH & Remote (21, 22, 23, 3389, 5900)",
		"Custom Range",
	}, nil)

	portRangeEntry := components.NewScrollableEntry()
	portRangeEntry.SetText("21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 445, 1433, 1521, 3306, 3389, 5432, 6379, 8080, 8443")

	presetSelect.OnChanged = func(s string) {
		switch s {
		case "Top 20 Common Ports":
			portRangeEntry.SetText("21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 445, 1433, 1521, 3306, 3389, 5432, 6379, 8080, 8443")
		case "Web & DB Ports (80, 443, 3306, 5432, 6379, 8080, 8443, 27017)":
			portRangeEntry.SetText("80, 443, 3306, 5432, 6379, 8080, 8443, 27017")
		case "SSH & Remote (21, 22, 23, 3389, 5900)":
			portRangeEntry.SetText("21, 22, 23, 3389, 5900")
		}
	}
	presetSelect.SetSelected("Top 20 Common Ports")

	timeoutSelect := widget.NewSelect([]string{
		"350 ms (Standar Internet)",
		"150 ms (LAN Cepat)",
		"1000 ms (Koneksi Lambat)",
	}, nil)
	timeoutSelect.SetSelected("350 ms (Standar Internet)")

	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)

	statusLabel := widget.NewLabel("Siap memindai port jaringan.")
	statusLabel.Wrapping = fyne.TextWrapWord

	badgeCountBox := container.NewHBox(
		components.BadgeCyan("0 OPEN"),
		components.BadgeMuted("0 DIPERIKSA"),
	)

	resultsBox := container.NewVBox()
	var lastResults []security.PortResult
	var cancelFunc context.CancelFunc

	var btnScan *widget.Button
	var btnStop *widget.Button

	btnScan = widget.NewButtonWithIcon("Mulai Scan Port", theme.MediaPlayIcon(), func() {
		host := strings.TrimSpace(hostEntry.Text)
		if host == "" {
			dialog.ShowInformation("Perhatian", "Masukkan target host atau alamat IP terlebih dahulu.", p.window)
			return
		}

		ports, err := security.ParsePortRange(portRangeEntry.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}

		var timeout time.Duration
		switch timeoutSelect.Selected {
		case "150 ms (LAN Cepat)":
			timeout = 150 * time.Millisecond
		case "1000 ms (Koneksi Lambat)":
			timeout = 1000 * time.Millisecond
		default:
			timeout = 350 * time.Millisecond
		}

		btnScan.Disable()
		btnStop.Enable()
		progressBar.SetValue(0)
		resultsBox.Objects = nil
		resultsBox.Refresh()

		ctx, cancel := context.WithCancel(context.Background())
		cancelFunc = cancel

		statusLabel.SetText(fmt.Sprintf("Memindai %d port pada %s...", len(ports), host))

		go func() {
			results, err := security.ScanTarget(ctx, host, ports, timeout, 60, func(prog security.ScanProgress) {
				fyne.Do(func() {
					if prog.Total > 0 {
						progressBar.SetValue(float64(prog.Scanned) / float64(prog.Total))
					}
					badgeCountBox.Objects = []fyne.CanvasObject{
						components.BadgeSuccess(fmt.Sprintf("%d OPEN", prog.OpenCount)),
						components.BadgeMuted(fmt.Sprintf("%d / %d", prog.Scanned, prog.Total)),
					}
					badgeCountBox.Refresh()
				})
			})

			fyne.Do(func() {
				btnScan.Enable()
				btnStop.Disable()

				if err != nil {
					statusLabel.SetText(fmt.Sprintf("Pemindaian gagal: %v", err))
					return
				}

				lastResults = results
				progressBar.SetValue(1.0)

				resultsBox.Objects = nil
				openCount := 0
				for _, r := range results {
					if r.State == "OPEN" {
						openCount++
						rCopy := r
						row := container.NewBorder(
							nil, nil,
							container.NewHBox(
								components.BadgeSuccess("OPEN"),
								widget.NewLabelWithStyle(fmt.Sprintf("Port %d", rCopy.Port), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
								widget.NewLabel(fmt.Sprintf("(%s)", rCopy.Service)),
							),
							widget.NewLabel(fmt.Sprintf("%d ms", rCopy.Latency)),
						)
						if rCopy.Banner != "" {
							row = container.NewVBox(row, widget.NewLabel(fmt.Sprintf("Banner: %s", rCopy.Banner)))
						}
						resultsBox.Add(components.NewPlainCardWithAccent(row, constants.ColorSuccess))
					}
				}

				if openCount == 0 {
					resultsBox.Add(container.NewPadded(widget.NewLabel("Tidak ada port terbuka yang merespons dari rentang yang dipindai.")))
				}
				resultsBox.Refresh()

				statusLabel.SetText(fmt.Sprintf("Pemindaian selesai: Ditemukan %d port terbuka pada %s.", openCount, host))
			})
		}()
	})
	btnScan.Importance = widget.HighImportance

	btnStop = widget.NewButtonWithIcon("Hentikan", theme.MediaStopIcon(), func() {
		if cancelFunc != nil {
			cancelFunc()
			statusLabel.SetText("Pemindaian dihentikan oleh pengguna.")
			btnScan.Enable()
			btnStop.Disable()
		}
	})
	btnStop.Disable()

	btnExport := widget.NewButtonWithIcon("Ekspor Hasil", theme.DocumentSaveIcon(), func() {
		if len(lastResults) == 0 {
			dialog.ShowInformation("Perhatian", "Belum ada hasil pemindaian untuk diekspor.", p.window)
			return
		}
		host := strings.TrimSpace(hostEntry.Text)
		report := security.GenerateScanReport(host, lastResults)
		p.copyToClip(report)
		dialog.ShowInformation("Ekspor Berhasil", "Laporan hasil audit port telah disalin ke clipboard.", p.window)
	})

	btnSampleScan := widget.NewButtonWithIcon("Contoh Target", theme.HelpIcon(), func() {
		hostEntry.SetText("scanme.nmap.org")
		portRangeEntry.SetText("22, 80, 443, 8080, 1000-1024")
		timeoutSelect.SetSelected("350 ms (Standar Internet)")
	})
	btnSampleScan.Importance = widget.LowImportance

	formatGuidePortLbl := widget.NewLabel("Host/IP Target: IPv4 (contoh: 192.168.1.1, 127.0.0.1) atau Domain FQDN (contoh: scanme.nmap.org, example.com).\nFormat Port: Port tunggal (80), daftar koma (21, 22, 80, 443), atau rentang strip (8000-8080, 1-1024).")
	formatGuidePortLbl.Wrapping = fyne.TextWrapWord

	formatGuidePort := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Panduan Format Input Port Scanner:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		formatGuidePortLbl,
	), constants.ColorBorderSubtle)

	inputCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Parameter Target Pemindaian Port:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3,
			container.NewVBox(widget.NewLabelWithStyle("Host Target / IP:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), hostEntry),
			container.NewVBox(widget.NewLabelWithStyle("Preset Port:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), presetSelect),
			container.NewVBox(widget.NewLabelWithStyle("Batas Timeout:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), timeoutSelect),
		),
		widget.NewLabelWithStyle("Daftar / Rentang Port (bisa dipisah koma atau rentang strip):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewHScroll(portRangeEntry),
		container.NewHBox(
			btnScan,
			btnStop,
			btnExport,
			btnSampleScan,
		),
		progressBar,
		container.NewBorder(nil, nil, statusLabel, badgeCountBox),
	), constants.ColorAccentCyan)

	resultsCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Daftar Port Terbuka Terdeteksi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		resultsBox,
	), constants.ColorSuccess)

	return container.NewVScroll(container.NewVBox(
		formatGuidePort,
		inputCard,
		resultsCard,
	))
}

// ----------------------------------------------------------------------------
// 2. Tab SSL & Web Security
// ----------------------------------------------------------------------------
func (p *SecurityPage) buildWebSSLTab() fyne.CanvasObject {
	urlEntry := components.NewScrollableEntry()
	urlEntry.SetPlaceHolder("Masukkan URL atau domain web (contoh: https://google.com atau github.com)")
	urlEntry.SetText("https://google.com")

	statusLabel := widget.NewLabel("Masukkan domain lalu klik 'Audit Keamanan Web & SSL'.")
	statusLabel.Wrapping = fyne.TextWrapWord

	sslDetailsBox := container.NewVBox()
	headerDetailsBox := container.NewVBox()
	dnsDetailsBox := container.NewVBox()

	gradeBadgeContainer := container.NewHBox(components.BadgeMuted("GRADE -"))

	btnAudit := widget.NewButtonWithIcon("Audit Keamanan Web & SSL", theme.SearchIcon(), func() {
		target := strings.TrimSpace(urlEntry.Text)
		if target == "" {
			dialog.ShowInformation("Perhatian", "Masukkan domain atau URL web target.", p.window)
			return
		}

		statusLabel.SetText("Menghubungkan dan menganalisis SSL/TLS, Security Headers, serta DNS...")
		sslDetailsBox.Objects = nil
		headerDetailsBox.Objects = nil
		dnsDetailsBox.Objects = nil
		sslDetailsBox.Refresh()
		headerDetailsBox.Refresh()
		dnsDetailsBox.Refresh()

		go func() {
			tlsRes, tlsErr := security.InspectTLS(target, 443, 6*time.Second)
			headerRes, headerErr := security.InspectSecurityHeaders(target, 6*time.Second)
			dnsRes, dnsErr := security.CheckEmailSecurityDNS(target)

			fyne.Do(func() {
				// 1. SSL/TLS Details
				if tlsErr != nil {
					sslDetailsBox.Add(widget.NewLabel(fmt.Sprintf("Gagal inspeksi SSL/TLS: %v", tlsErr)))
				} else {
					var expiryBadge fyne.CanvasObject
					if tlsRes.IsExpired {
						expiryBadge = components.BadgeDanger("KEDALUWARSA")
					} else if tlsRes.ExpiringSoon {
						expiryBadge = components.BadgeWarning(fmt.Sprintf("KEDALUWARSA DALAM %d HARI", tlsRes.DaysLeft))
					} else {
						expiryBadge = components.BadgeSuccess(fmt.Sprintf("VALID (%d HARI LAGI)", tlsRes.DaysLeft))
					}

					sslDetailsBox.Add(container.NewHBox(
						widget.NewLabelWithStyle("Subjek Sertifikat:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabel(tlsRes.SubjectCN),
						expiryBadge,
					))
					sslDetailsBox.Add(container.NewHBox(
						widget.NewLabelWithStyle("Penerbit (CA):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabel(tlsRes.Issuer),
					))
					sslDetailsBox.Add(container.NewHBox(
						widget.NewLabelWithStyle("Masa Berlaku:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabel(fmt.Sprintf("%s s/d %s", tlsRes.ValidFrom.Format("02 Jan 2006"), tlsRes.ValidUntil.Format("02 Jan 2006"))),
					))
					sslDetailsBox.Add(container.NewHBox(
						widget.NewLabelWithStyle("Protokol TLS:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						components.BadgeIndigo(tlsRes.TLSVersion),
						widget.NewLabel(fmt.Sprintf("Cipher: %s", tlsRes.CipherSuite)),
					))
					if len(tlsRes.SANs) > 0 {
						sanPreview := strings.Join(tlsRes.SANs, ", ")
						if len(sanPreview) > 90 {
							sanPreview = sanPreview[:87] + "..."
						}
						sslDetailsBox.Add(container.NewHBox(
							widget.NewLabelWithStyle("Domain SANs:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
							widget.NewLabel(sanPreview),
						))
					}
				}
				sslDetailsBox.Refresh()

				// 2. Security Headers Details
				if headerErr != nil {
					headerDetailsBox.Add(widget.NewLabel(fmt.Sprintf("Gagal audit Security Headers: %v", headerErr)))
				} else {
					var gradeBadge fyne.CanvasObject
					switch headerRes.Grade {
					case "A+", "A":
						gradeBadge = components.BadgeSuccess("GRADE " + headerRes.Grade)
					case "B", "C":
						gradeBadge = components.BadgeWarning("GRADE " + headerRes.Grade)
					default:
						gradeBadge = components.BadgeDanger("GRADE " + headerRes.Grade)
					}
					gradeBadgeContainer.Objects = []fyne.CanvasObject{gradeBadge, components.BadgeCyan(fmt.Sprintf("SKOR %d / 100", headerRes.Score))}
					gradeBadgeContainer.Refresh()

					headerDetailsBox.Add(container.NewHBox(
						widget.NewLabelWithStyle("Server Banner:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
						widget.NewLabel(headerRes.Server),
						widget.NewLabel(fmt.Sprintf("HTTP Status: %d", headerRes.StatusCode)),
					))

					for _, h := range headerRes.Headers {
						hCopy := h
						var sevBadge fyne.CanvasObject
						if hCopy.Present {
							sevBadge = components.BadgeSuccess("PASS")
						} else if hCopy.Severity == "WARN" {
							sevBadge = components.BadgeWarning("WARN")
						} else {
							sevBadge = components.BadgeDanger("FAIL")
						}

						row := container.NewBorder(
							nil, nil,
							container.NewHBox(sevBadge, widget.NewLabelWithStyle(hCopy.HeaderName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
							nil,
							widget.NewLabel(hCopy.Recommendation),
						)
						headerDetailsBox.Add(row)
					}

					btnSnippet := widget.NewButtonWithIcon("Salin Panduan Konfigurasi Nginx & Apache", theme.ContentCopyIcon(), func() {
						p.copyToClip(headerRes.RemediationSnippets)
					})
					btnSnippet.Importance = widget.LowImportance
					headerDetailsBox.Add(btnSnippet)
				}
				headerDetailsBox.Refresh()

				// 3. DNS Email Security
				if dnsErr != nil {
					dnsDetailsBox.Add(widget.NewLabel(fmt.Sprintf("Gagal cek DNS: %v", dnsErr)))
				} else {
					dnsDetailsBox.Add(widget.NewLabelWithStyle(dnsRes.StatusSummary, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
					if dnsRes.HasSPF {
						dnsDetailsBox.Add(container.NewHBox(components.BadgeSuccess("SPF AKTIF"), widget.NewLabel(dnsRes.SPFRecord)))
					} else {
						dnsDetailsBox.Add(container.NewHBox(components.BadgeDanger("SPF HILANG"), widget.NewLabel("Belum ada record v=spf1")))
					}
					if dnsRes.HasDMARC {
						dnsDetailsBox.Add(container.NewHBox(components.BadgeSuccess("DMARC AKTIF"), widget.NewLabel(dnsRes.DMARCRecord)))
					} else {
						dnsDetailsBox.Add(container.NewHBox(components.BadgeDanger("DMARC HILANG"), widget.NewLabel("Belum ada record _dmarc v=dmarc1")))
					}
				}
				dnsDetailsBox.Refresh()

				statusLabel.SetText(fmt.Sprintf("Audit selesai untuk %s.", target))
			})
		}()
	})
	btnSampleWeb := widget.NewButtonWithIcon("Contoh Target", theme.HelpIcon(), func() {
		urlEntry.SetText("https://cloudflare.com")
	})
	btnSampleWeb.Importance = widget.LowImportance

	formatGuideWebSSLLbl := widget.NewLabel("Format URL: Masukkan protokol lengkap https://domain.com (contoh: https://google.com) atau nama domain langsung (contoh: github.com, cloudflare.com).\nCatatan: Pemeriksaan sertifikat TLS menggunakan port 443 standar, sedangkan security headers diperiksa melalui HTTP/HTTPS response.")
	formatGuideWebSSLLbl.Wrapping = fyne.TextWrapWord

	formatGuideWebSSL := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Panduan Format Input Web & SSL Scanner:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		formatGuideWebSSLLbl,
	), constants.ColorBorderSubtle)

	inputCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Target Audit Web & Keamanan Domain:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, container.NewHBox(btnAudit, btnSampleWeb), urlEntry),
		statusLabel,
	), constants.ColorAccentCyan)

	sslCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("1. Status Sertifikat SSL/TLS & Enkripsi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		sslDetailsBox,
	), constants.ColorTechIndigo)

	headerCard := components.NewPlainCardWithAccent(container.NewVBox(
		container.NewBorder(
			nil, nil,
			widget.NewLabelWithStyle("2. Evaluasi HTTP Security Headers:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			gradeBadgeContainer,
		),
		headerDetailsBox,
	), constants.ColorSuccess)

	dnsCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("3. Keamanan DNS Email (Anti-Spoofing & Phishing):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		dnsDetailsBox,
	), constants.ColorAccentCobalt)

	return container.NewVScroll(container.NewVBox(
		formatGuideWebSSL,
		inputCard,
		sslCard,
		headerCard,
		dnsCard,
	))
}

// ----------------------------------------------------------------------------
// 3. Tab Forensik & Log Audit
// ----------------------------------------------------------------------------
func (p *SecurityPage) buildForensicsTab() fyne.CanvasObject {
	// Sub-mode 1: File Hashes
	filePathEntry := components.NewScrollableEntry()
	filePathEntry.SetPlaceHolder("Pilih berkas yang ingin dihitung nilai hash & integritasnya...")

	md5Entry := components.NewScrollableEntry()
	md5Entry.Disable()
	sha1Entry := components.NewScrollableEntry()
	sha1Entry.Disable()
	sha256Entry := components.NewScrollableEntry()
	sha256Entry.Disable()
	sha512Entry := components.NewScrollableEntry()
	sha512Entry.Disable()

	magicBadgeContainer := container.NewHBox()
	integrityNoteLabel := widget.NewLabel("")
	integrityNoteLabel.Wrapping = fyne.TextWrapWord

	compareEntry := components.NewScrollableEntry()
	compareEntry.SetPlaceHolder("Tempelkan hash pembanding (misal dari situs unduhan vendor)...")
	compareResultBadge := container.NewHBox()

	updateCompare := func() {
		compareResultBadge.Objects = nil
		cmp := strings.TrimSpace(strings.ToLower(compareEntry.Text))
		if cmp == "" {
			compareResultBadge.Refresh()
			return
		}

		if strings.EqualFold(cmp, md5Entry.Text) ||
			strings.EqualFold(cmp, sha1Entry.Text) ||
			strings.EqualFold(cmp, sha256Entry.Text) ||
			strings.EqualFold(cmp, sha512Entry.Text) {
			compareResultBadge.Add(components.BadgeSuccess("VERIFIKASI COCOK (MATCH)"))
		} else {
			compareResultBadge.Add(components.BadgeDanger("TIDAK COCOK (MISMATCH)"))
		}
		compareResultBadge.Refresh()
	}
	compareEntry.OnChanged = func(_ string) { updateCompare() }

	btnBrowse := widget.NewButtonWithIcon("Pilih Berkas...", theme.FileIcon(), func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if closer == nil || err != nil {
				return
			}
			path := closer.URI().Path()
			filePathEntry.SetText(path)

			res, hErr := security.ComputeFileHashes(path)
			if hErr != nil {
				dialog.ShowError(hErr, p.window)
				return
			}

			md5Entry.SetText(res.MD5)
			sha1Entry.SetText(res.SHA1)
			sha256Entry.SetText(res.SHA256)
			sha512Entry.SetText(res.SHA512)

			magicBadgeContainer.Objects = nil
			if res.Suspicious {
				magicBadgeContainer.Add(components.BadgeDanger("ANOMALI EKSTENSI"))
			} else {
				magicBadgeContainer.Add(components.BadgeCyan(res.MagicType))
			}
			magicBadgeContainer.Refresh()
			integrityNoteLabel.SetText(fmt.Sprintf("%s (%d bytes) - %s", res.FileName, res.FileSize, res.Notes))
			updateCompare()
		}, p.window)
	})
	btnBrowse.Importance = widget.HighImportance

	vtStatusLabel := widget.NewLabel("")
	vtStatusLabel.Wrapping = fyne.TextWrapWord
	btnVT := widget.NewButtonWithIcon("Cek Reputasi VirusTotal", theme.SearchIcon(), func() {
		sha256Val := strings.TrimSpace(sha256Entry.Text)
		if sha256Val == "" {
			dialog.ShowInformation("Perhatian", "Pilih berkas terlebih dahulu untuk mendapatkan hash SHA-256.", p.window)
			return
		}

		apiKeyEntry := components.NewScrollableEntry()
		apiKeyEntry.SetPlaceHolder("Tempelkan API Key VirusTotal Anda di sini...")
		dialog.ShowCustomConfirm("VirusTotal Hash Lookup", "Cek Sekarang", "Batal", container.NewVBox(
			widget.NewLabel("Masukkan API Key VirusTotal untuk memeriksa reputasi hash berkas:"),
			apiKeyEntry,
		), func(ok bool) {
			if !ok {
				return
			}
			key := strings.TrimSpace(apiKeyEntry.Text)
			if key == "" {
				dialog.ShowInformation("Perhatian", "API key tidak boleh kosong.", p.window)
				return
			}
			vtStatusLabel.SetText("Menghubungi VirusTotal API...")
			go func() {
				vtRes, vtErr := security.LookupVirusTotalHash(key, sha256Val)
				fyne.Do(func() {
					if vtErr != nil {
						vtStatusLabel.SetText(fmt.Sprintf("Gagal cek VirusTotal: %v", vtErr))
						return
					}
					vtStatusLabel.SetText(vtRes.StatusSummary)
					dialog.ShowInformation("Hasil VirusTotal",
						fmt.Sprintf("Hash: %s\nLabel Ancaman: %s\nSkor Reputasi: %d\n\n%s",
							vtRes.SHA256, vtRes.ThreatLabel, vtRes.Reputation, vtRes.StatusSummary),
						p.window)
				})
			}()
		}, p.window)
	})
	btnVT.Importance = widget.LowImportance

	makeHashRow := func(label string, entry *widget.Entry) fyne.CanvasObject {
		copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			if entry.Text != "" {
				p.copyToClip(entry.Text)
			}
		})
		copyBtn.Importance = widget.LowImportance
		return container.NewBorder(
			nil, nil,
			container.NewHBox(widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			copyBtn,
			entry,
		)
	}

	fileContentBox := container.NewVBox(
		widget.NewLabelWithStyle("Pilih Berkas untuk Integritas & Multi-Checksum:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, btnBrowse, filePathEntry),
		container.NewHBox(magicBadgeContainer, integrityNoteLabel),
		widget.NewSeparator(),
		makeHashRow("MD5 (128-bit):   ", md5Entry),
		makeHashRow("SHA-1 (160-bit): ", sha1Entry),
		makeHashRow("SHA-256 (256-bit):", sha256Entry),
		makeHashRow("SHA-512 (512-bit):", sha512Entry),
		widget.NewSeparator(),
		container.NewBorder(nil, nil, widget.NewLabelWithStyle("Bandingkan Hash:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), compareResultBadge, compareEntry),
		container.NewHBox(btnVT, vtStatusLabel),
	)

	// Sub-mode 2: Log Anomaly Parser
	logTypeSelect := widget.NewSelect([]string{"Linux / SSH auth.log", "Web Server access.log (Nginx / Apache)"}, nil)
	logTypeSelect.SetSelected("Linux / SSH auth.log")

	logInputEntry := widget.NewMultiLineEntry()
	logInputEntry.SetPlaceHolder("Tempelkan cuplikan baris berkas log di sini...")
	logInputEntry.Wrapping = fyne.TextWrapWord
	logInputEntry.SetMinRowsVisible(6)

	btnLoadSample := widget.NewButtonWithIcon("Muat Contoh Log Uji", theme.DocumentCreateIcon(), func() {
		if logTypeSelect.Selected == "Linux / SSH auth.log" {
			logInputEntry.SetText(`Oct 12 04:02:11 web-server sshd[2841]: Failed password for root from 198.51.100.24 port 41200 ssh2
Oct 12 04:02:14 web-server sshd[2842]: Failed password for root from 198.51.100.24 port 41202 ssh2
Oct 12 04:02:18 web-server sshd[2843]: Failed password for invalid user admin from 198.51.100.24 port 41204 ssh2
Oct 12 04:03:00 web-server sshd[2844]: Failed password for invalid user test from 203.0.113.88 port 50122 ssh2
Oct 12 04:05:12 web-server sshd[2845]: Accepted publickey for yudz from 192.168.1.15 port 55432 ssh2`)
		} else {
			logInputEntry.SetText(`198.51.100.24 - - [12/Oct/2026:04:02:11 +0000] "GET /index.php?id=1%20union+select+1,2,3 HTTP/1.1" 200 452
198.51.100.24 - - [12/Oct/2026:04:02:15 +0000] "GET /../../../../etc/passwd HTTP/1.1" 404 153
203.0.113.88 - - [12/Oct/2026:04:03:01 +0000] "GET /search?q=<script>alert(1)</script> HTTP/1.1" 200 678
203.0.113.88 - - [12/Oct/2026:04:03:10 +0000] "GET /wp-login.php HTTP/1.1" 404 153`)
		}
	})
	btnLoadSample.Importance = widget.LowImportance

	logAnalysisResultBox := container.NewVBox()

	btnAnalyzeLog := widget.NewButtonWithIcon("Analisa Anomali Log", theme.SearchIcon(), func() {
		text := strings.TrimSpace(logInputEntry.Text)
		if text == "" {
			dialog.ShowInformation("Perhatian", "Tempelkan teks log terlebih dahulu.", p.window)
			return
		}

		logAnalysisResultBox.Objects = nil
		if logTypeSelect.Selected == "Linux / SSH auth.log" {
			res, err := security.ParseAuthLog(text)
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}

			summaryRow := container.NewHBox(
				components.BadgeDanger(fmt.Sprintf("%d GAGAL LOGIN", res.FailedAttempts)),
				components.BadgeWarning(fmt.Sprintf("%d USER TIDAK VALID", res.InvalidUsers)),
				components.BadgeSuccess(fmt.Sprintf("%d BERHASIL", res.SuccessLogins)),
				components.BadgeMuted(fmt.Sprintf("%d TOTAL BARIS", res.TotalLines)),
			)
			logAnalysisResultBox.Add(summaryRow)

			if len(res.TopAttackerIPs) > 0 {
				logAnalysisResultBox.Add(widget.NewLabelWithStyle("IP Penyerang Paling Agresif:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
				for _, att := range res.TopAttackerIPs {
					logAnalysisResultBox.Add(container.NewHBox(
						components.BadgeDanger(fmt.Sprintf("%d Percobaan", att.Count)),
						widget.NewLabel(att.Key),
					))
				}

				blockRules := security.GenerateBlockRules(res.AttackerIPList, "ufw")
				btnCopyBlock := widget.NewButtonWithIcon("Salin Perintah Blokir Firewall UFW", theme.ContentCopyIcon(), func() {
					p.copyToClip(blockRules)
				})
				btnCopyBlock.Importance = widget.HighImportance
				logAnalysisResultBox.Add(btnCopyBlock)
			}
		} else {
			res, err := security.ParseWebAccessLog(text)
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}

			summaryRow := container.NewHBox(
				components.BadgeDanger(fmt.Sprintf("%d SQL INJECTION", res.SqlInjections)),
				components.BadgeDanger(fmt.Sprintf("%d PATH TRAVERSAL", res.PathTraversals)),
				components.BadgeWarning(fmt.Sprintf("%d XSS ATTEMPTS", res.XssAttempts)),
				components.BadgeMuted(fmt.Sprintf("%d 404 SCANS", res.Http404Scans)),
			)
			logAnalysisResultBox.Add(summaryRow)

			if len(res.TopSuspectIPs) > 0 {
				logAnalysisResultBox.Add(widget.NewLabelWithStyle("IP Scanner / Penyerang Web:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
				for _, att := range res.TopSuspectIPs {
					logAnalysisResultBox.Add(container.NewHBox(
						components.BadgeDanger(fmt.Sprintf("%d Requests", att.Count)),
						widget.NewLabel(att.Key),
					))
				}

				blockRules := security.GenerateBlockRules(res.SuspectIPList, "ufw")
				btnCopyBlock := widget.NewButtonWithIcon("Salin Perintah Blokir Firewall UFW", theme.ContentCopyIcon(), func() {
					p.copyToClip(blockRules)
				})
				btnCopyBlock.Importance = widget.HighImportance
				logAnalysisResultBox.Add(btnCopyBlock)
			}
		}
		logAnalysisResultBox.Refresh()
	})
	btnAnalyzeLog.Importance = widget.HighImportance

	logContentBox := container.NewVBox(
		container.NewHBox(
			widget.NewLabelWithStyle("Pilih Format Log:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			logTypeSelect,
			btnLoadSample,
		),
		logInputEntry,
		btnAnalyzeLog,
		widget.NewSeparator(),
		logAnalysisResultBox,
	)
	logContentBox.Hide()

	modeRadio := widget.NewRadioGroup([]string{
		"Integritas & Multi-Hash Berkas",
		"Analisis Anomali Log (Auth & Web)",
	}, func(s string) {
		if s == "Integritas & Multi-Hash Berkas" {
			logContentBox.Hide()
			fileContentBox.Show()
		} else {
			fileContentBox.Hide()
			logContentBox.Show()
		}
	})
	modeRadio.Horizontal = false
	modeRadio.SetSelected("Integritas & Multi-Hash Berkas")

	formatGuideForensicsLbl := widget.NewLabel("1. Integritas Berkas: Semua jenis berkas biner/teks (exe, zip, iso, pdf, dll). Hash pembanding mendukung format heksadesimal 32 karakter (MD5), 40 karakter (SHA-1), 64 karakter (SHA-256), atau 128 karakter (SHA-512).\n2. Log Linux/SSH: Format syslog auth.log standar (contoh: 'Oct 12 04:02:11 host sshd[123]: Failed password for root from 1.2.3.4 port 22 ssh2').\n3. Log Web Server: Format Nginx/Apache Combined Log Format (contoh: '1.2.3.4 - - [12/Oct/2026:04:02:11 +0000] \"GET /index.php HTTP/1.1\" 200 452').")
	formatGuideForensicsLbl.Wrapping = fyne.TextWrapWord

	formatGuideForensics := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Panduan Format Input Forensik & Log Audit:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		formatGuideForensicsLbl,
	), constants.ColorBorderSubtle)

	rootCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Pilih Mode Forensik Keamanan:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		modeRadio,
		widget.NewSeparator(),
		fileContentBox,
		logContentBox,
	), constants.ColorAccentCyan)

	return container.NewVScroll(container.NewVBox(
		formatGuideForensics,
		rootCard,
	))
}

// ----------------------------------------------------------------------------
// 4. Tab Encoder, JWT & ASCII Tools
// ----------------------------------------------------------------------------
func (p *SecurityPage) buildEncodersASCIITab() fyne.CanvasObject {
	// Sub-mode 1: ASCII Table
	asciiSearchEntry := components.NewScrollableEntry()
	asciiSearchEntry.SetPlaceHolder("Cari kode desimal (65), hex (0x41), char 'A', atau nama...")

	asciiCatSelect := widget.NewSelect([]string{"Semua", "Control", "Printable", "Extended"}, nil)
	asciiCatSelect.SetSelected("Semua")

	asciiTableList := container.NewVBox()

	renderASCIITable := func() {
		asciiTableList.Objects = nil
		q := asciiSearchEntry.Text
		cat := asciiCatSelect.Selected
		items := security.FilterASCIITable(q, cat)

		// Limit rendered items to avoid sluggishness if too large
		limit := 100
		if len(items) < limit {
			limit = len(items)
		}

		for i := 0; i < limit; i++ {
			it := items[i]
			itCopy := it

			var catBadge fyne.CanvasObject
			switch itCopy.Category {
			case "Control":
				catBadge = components.BadgeWarning("CTRL")
			case "Printable":
				catBadge = components.BadgeSuccess("CHAR")
			default:
				catBadge = components.BadgeIndigo("EXT")
			}

			copyDecBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
				p.copyToClip(fmt.Sprintf("Dec: %d | Hex: %s | Bin: %s | Char: %s", itCopy.Dec, itCopy.Hex, itCopy.Bin, itCopy.Char))
			})
			copyDecBtn.Importance = widget.LowImportance

			row := container.NewBorder(
				nil, nil,
				container.NewHBox(
					catBadge,
					widget.NewLabelWithStyle(fmt.Sprintf("%-3d", itCopy.Dec), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					widget.NewLabel(fmt.Sprintf("%-6s", itCopy.Hex)),
					widget.NewLabel(fmt.Sprintf("%-8s", itCopy.Bin)),
					components.BadgeCyan(itCopy.Char),
				),
				copyDecBtn,
				widget.NewLabel(itCopy.Name),
			)
			asciiTableList.Add(components.NewPlainCardWithAccent(row, constants.ColorBorderSubtle))
		}

		if len(items) > limit {
			asciiTableList.Add(container.NewCenter(widget.NewLabel(fmt.Sprintf("... Menampilkan %d dari %d hasil. Gunakan pencarian untuk menyaring spesifik.", limit, len(items)))))
		} else if len(items) == 0 {
			asciiTableList.Add(container.NewCenter(widget.NewLabel("Tidak ada karakter ASCII yang cocok dengan filter.")))
		}
		asciiTableList.Refresh()
	}

	asciiSearchEntry.OnChanged = func(_ string) { renderASCIITable() }
	asciiCatSelect.OnChanged = func(_ string) { renderASCIITable() }
	renderASCIITable()

	asciiTableBox := container.NewVBox(
		widget.NewLabelWithStyle("Tabel Referensi Lengkap Karakter ASCII (0-255):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(2,
			asciiSearchEntry,
			container.NewHBox(widget.NewLabel("Filter:"), asciiCatSelect),
		),
		asciiTableList,
	)

	// Sub-mode 2: ASCII Translator
	transDirectionRadio := widget.NewRadioGroup([]string{"Teks Biasa -> Kode ASCII", "Kode ASCII -> Teks Biasa"}, nil)
	transDirectionRadio.SetSelected("Teks Biasa -> Kode ASCII")

	transFormatSelect := widget.NewSelect([]string{"Decimal", "Binary", "Hex", "Octal"}, nil)
	transFormatSelect.SetSelected("Decimal")

	transDelimSelect := widget.NewSelect([]string{"Spasi", "Koma", "Tanpa Spasi", "Prefiks 0x", "Prefiks \\x"}, nil)
	transDelimSelect.SetSelected("Spasi")

	transInputEntry := widget.NewMultiLineEntry()
	transInputEntry.SetPlaceHolder("Ketik teks untuk diterjemahkan ke kode ASCII...")
	transInputEntry.Wrapping = fyne.TextWrapWord
	transInputEntry.SetMinRowsVisible(4)

	transOutputEntry := widget.NewMultiLineEntry()
	transOutputEntry.SetPlaceHolder("Hasil terjemahan akan tampil di sini...")
	transOutputEntry.Wrapping = fyne.TextWrapWord
	transOutputEntry.SetMinRowsVisible(4)

	performTranslation := func() {
		in := transInputEntry.Text
		if in == "" {
			transOutputEntry.SetText("")
			return
		}

		if transDirectionRadio.Selected == "Teks Biasa -> Kode ASCII" {
			out := security.TextToASCII(in, transFormatSelect.Selected, transDelimSelect.Selected)
			transOutputEntry.SetText(out)
		} else {
			out, err := security.ASCIIToText(in, transFormatSelect.Selected)
			if err != nil {
				transOutputEntry.SetText(fmt.Sprintf("Error: %v", err))
			} else {
				transOutputEntry.SetText(out)
			}
		}
	}

	transInputEntry.OnChanged = func(_ string) { performTranslation() }
	transFormatSelect.OnChanged = func(_ string) { performTranslation() }
	transDelimSelect.OnChanged = func(_ string) { performTranslation() }
	transDirectionRadio.OnChanged = func(_ string) {
		if transDirectionRadio.Selected == "Teks Biasa -> Kode ASCII" {
			transInputEntry.SetPlaceHolder("Ketik teks untuk diterjemahkan ke kode ASCII...")
		} else {
			transInputEntry.SetPlaceHolder("Tempelkan kode angka desimal, biner, hex, atau oktal...")
		}
		performTranslation()
	}

	btnSwap := widget.NewButtonWithIcon("Tukar Arah (Swap)", theme.ViewRefreshIcon(), func() {
		currentOut := transOutputEntry.Text
		if transDirectionRadio.Selected == "Teks Biasa -> Kode ASCII" {
			transDirectionRadio.SetSelected("Kode ASCII -> Teks Biasa")
		} else {
			transDirectionRadio.SetSelected("Teks Biasa -> Kode ASCII")
		}
		transInputEntry.SetText(currentOut)
	})
	btnSwap.Importance = widget.LowImportance

	btnCopyTrans := widget.NewButtonWithIcon("Salin Hasil Terjemahan", theme.ContentCopyIcon(), func() {
		if transOutputEntry.Text != "" {
			p.copyToClip(transOutputEntry.Text)
		}
	})
	btnCopyTrans.Importance = widget.HighImportance

	asciiTranslatorBox := container.NewVBox(
		widget.NewLabelWithStyle("Penerjemah Dua Arah (ASCII Translator):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		transDirectionRadio,
		container.NewHBox(
			widget.NewLabel("Format Kode:"),
			transFormatSelect,
			widget.NewLabel("Pemisah / Delimiter:"),
			transDelimSelect,
			btnSwap,
		),
		widget.NewLabelWithStyle("Input:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		transInputEntry,
		widget.NewLabelWithStyle("Hasil Terjemahan:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		transOutputEntry,
		btnCopyTrans,
	)
	asciiTranslatorBox.Hide()

	// Sub-mode 3: JWT Inspector
	jwtInputEntry := widget.NewMultiLineEntry()
	jwtInputEntry.SetPlaceHolder("Tempelkan token JWT di sini (format: xxxxx.yyyyy.zzzzz)...")
	jwtInputEntry.Wrapping = fyne.TextWrapWord
	jwtInputEntry.SetMinRowsVisible(4)

	jwtStatusBadge := container.NewHBox(components.BadgeMuted("STATUS TOKEN"))
	jwtHeaderEntry := widget.NewMultiLineEntry()
	jwtHeaderEntry.Disable()
	jwtHeaderEntry.Wrapping = fyne.TextWrapWord

	jwtPayloadEntry := widget.NewMultiLineEntry()
	jwtPayloadEntry.Disable()
	jwtPayloadEntry.Wrapping = fyne.TextWrapWord

	btnDecodeJWT := widget.NewButtonWithIcon("Decode Token JWT", theme.SearchIcon(), func() {
		tok := strings.TrimSpace(jwtInputEntry.Text)
		if tok == "" {
			dialog.ShowInformation("Perhatian", "Tempelkan token JWT terlebih dahulu.", p.window)
			return
		}

		res, err := security.DecodeJWT(tok)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}

		jwtStatusBadge.Objects = nil
		if res.IsExpired {
			jwtStatusBadge.Add(components.BadgeDanger("EXPIRED"))
		} else {
			jwtStatusBadge.Add(components.BadgeSuccess("AKTIF / VALID"))
		}
		jwtStatusBadge.Add(components.BadgeCyan("Alg: " + res.Header.Alg))
		jwtStatusBadge.Refresh()

		jwtHeaderEntry.SetText(res.HeaderJSON)
		jwtPayloadEntry.SetText(res.PayloadJSON)
	})
	btnDecodeJWT.Importance = widget.HighImportance

	jwtBox := container.NewVBox(
		widget.NewLabelWithStyle("Inspektur Token JSON Web Token (JWT):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		jwtInputEntry,
		btnDecodeJWT,
		widget.NewSeparator(),
		container.NewBorder(nil, nil, widget.NewLabelWithStyle("Status Validasi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), jwtStatusBadge),
		widget.NewLabelWithStyle("Header Metadata:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		jwtHeaderEntry,
		widget.NewLabelWithStyle("Payload Claims (Data Pengguna / Sesi):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		jwtPayloadEntry,
	)
	jwtBox.Hide()

	// Sub-mode 4: Multi-Encoder & Hash Type Identifier
	encInputEntry := widget.NewMultiLineEntry()
	encInputEntry.SetPlaceHolder("Ketik teks atau tempelkan hash di sini...")
	encInputEntry.Wrapping = fyne.TextWrapWord
	encInputEntry.SetMinRowsVisible(4)

	encOutputEntry := widget.NewMultiLineEntry()
	encOutputEntry.SetPlaceHolder("Hasil encode / decode / analisis hash...")
	encOutputEntry.Wrapping = fyne.TextWrapWord
	encOutputEntry.SetMinRowsVisible(4)

	makeEncBtn := func(label, mode string) *widget.Button {
		b := widget.NewButton(label, func() {
			out, err := security.TransformSecurityText(encInputEntry.Text, mode)
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			encOutputEntry.SetText(out)
		})
		b.Importance = widget.LowImportance
		return b
	}

	btnIdentifyHash := widget.NewButtonWithIcon("Identifikasi Jenis Hash", theme.SearchIcon(), func() {
		candidates := security.IdentifyHashType(encInputEntry.Text)
		var sb strings.Builder
		sb.WriteString("=== HASIL ANALISIS JENIS HASH ===\n\n")
		for i, c := range candidates {
			sb.WriteString(fmt.Sprintf("%d. %s (%d-bit)\n   Deskripsi: %s\n\n", i+1, c.Name, c.BitLength, c.Description))
		}
		encOutputEntry.SetText(sb.String())
	})
	btnIdentifyHash.Importance = widget.HighImportance

	encoderBox := container.NewVBox(
		widget.NewLabelWithStyle("Multi-Encoder & Penganalisis Algoritma Hash:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		encInputEntry,
		container.NewHBox(
			makeEncBtn("Base64 Encode", "base64-encode"),
			makeEncBtn("Base64 Decode", "base64-decode"),
			makeEncBtn("URL Encode", "url-encode"),
			makeEncBtn("URL Decode", "url-decode"),
			makeEncBtn("Hex Encode", "hex-encode"),
			makeEncBtn("Hex Decode", "hex-decode"),
		),
		container.NewHBox(
			makeEncBtn("HTML Encode", "html-encode"),
			makeEncBtn("HTML Decode", "html-decode"),
			makeEncBtn("ROT13", "rot13"),
			btnIdentifyHash,
		),
		widget.NewLabelWithStyle("Output Hasil:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		encOutputEntry,
		widget.NewButtonWithIcon("Salin Hasil", theme.ContentCopyIcon(), func() {
			if encOutputEntry.Text != "" {
				p.copyToClip(encOutputEntry.Text)
			}
		}),
	)
	// Sub-mode 5: Key Decryptor
	keyAlgoSelect := widget.NewSelect([]string{
		"AES-128 (CBC PKCS7)",
		"AES-192 (CBC PKCS7)",
		"AES-256 (CBC PKCS7)",
		"AES-GCM (Authenticated)",
		"DES (CBC PKCS7)",
		"Triple DES / 3DES (CBC PKCS7)",
		"Repeating-Key XOR",
		"Vigenere Cipher",
	}, nil)
	keyAlgoSelect.SetSelected("AES-128 (CBC PKCS7)")

	keyEncodingSelect := widget.NewSelect([]string{"Base64", "Hex", "Plain Text"}, nil)
	keyEncodingSelect.SetSelected("Base64")

	keyCipherEntry := widget.NewMultiLineEntry()
	keyCipherEntry.SetPlaceHolder("Tempelkan ciphertext terenkripsi di sini...")
	keyCipherEntry.Wrapping = fyne.TextWrapWord
	keyCipherEntry.SetMinRowsVisible(4)

	keySecretEntry := components.NewScrollableEntry()
	keySecretEntry.SetPlaceHolder("Kunci rahasia (contoh: 16 char untuk AES-128, 24 untuk AES-192, 32 untuk AES-256)...")

	keyIVEntry := components.NewScrollableEntry()
	keyIVEntry.SetPlaceHolder("IV / Nonce opsional (contoh: 16-byte teks/hex untuk CBC, 12-byte untuk GCM, biarkan kosong jika nol)...")

	keyPlainOutput := widget.NewMultiLineEntry()
	keyPlainOutput.SetPlaceHolder("Hasil dekripsi plaintext akan muncul di sini...")
	keyPlainOutput.Wrapping = fyne.TextWrapWord
	keyPlainOutput.SetMinRowsVisible(4)

	btnExecuteDecrypt := widget.NewButtonWithIcon("Dekripsi Ciphertext", theme.ConfirmIcon(), func() {
		cipherText := keyCipherEntry.Text
		keyStr := keySecretEntry.Text
		ivStr := keyIVEntry.Text
		algo := keyAlgoSelect.Selected
		enc := keyEncodingSelect.Selected

		if strings.TrimSpace(cipherText) == "" {
			dialog.ShowInformation("Perhatian", "Ciphertext tidak boleh kosong.", p.window)
			return
		}

		res, err := security.DecryptWithKey(algo, enc, cipherText, keyStr, ivStr)
		if err != nil {
			dialog.ShowError(err, p.window)
			keyPlainOutput.SetText(fmt.Sprintf("Galat Dekripsi: %v", err))
			return
		}

		keyPlainOutput.SetText(res)
		dialog.ShowInformation("Sukses", "Dekripsi berhasil diselesaikan.", p.window)
	})
	btnExecuteDecrypt.Importance = widget.HighImportance

	btnSampleAES := widget.NewButton("Muat Contoh AES", func() {
		keyAlgoSelect.SetSelected("AES-128 (CBC PKCS7)")
		keyEncodingSelect.SetSelected("Base64")
		keySecretEntry.SetText("kuncirahasia1234")
		keyIVEntry.SetText("vektorinisiasi12")
		keyCipherEntry.SetText("uM4/d+q6VjQjN8YdZ6qKFAq6c7V/mR5u87c+V3g/q2E=")
	})
	btnSampleAES.Importance = widget.LowImportance

	btnSampleVig := widget.NewButton("Muat Contoh Vigenere", func() {
		keyAlgoSelect.SetSelected("Vigenere Cipher")
		keyEncodingSelect.SetSelected("Plain Text")
		keySecretEntry.SetText("LEMON")
		keyIVEntry.SetText("")
		keyCipherEntry.SetText("LXFOPV EF RNHR!")
	})
	btnSampleVig.Importance = widget.LowImportance

	keyDecryptBox := container.NewVBox(
		widget.NewLabelWithStyle("Key Decryptor (Dekripsi Simetris Berbasis Kunci):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Mendukung cipher modern AES (CBC & GCM), DES, 3DES, serta cipher klasik XOR dan Vigenere."),
		container.NewHBox(
			widget.NewLabel("Algoritma:"),
			keyAlgoSelect,
			widget.NewLabel("Format Ciphertext:"),
			keyEncodingSelect,
			btnSampleAES,
			btnSampleVig,
		),
		widget.NewLabelWithStyle("Ciphertext Input:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		keyCipherEntry,
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Secret Key:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), keySecretEntry),
			container.NewVBox(widget.NewLabelWithStyle("Initialization Vector (IV) / Nonce:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), keyIVEntry),
		),
		btnExecuteDecrypt,
		widget.NewLabelWithStyle("Plaintext Hasil Dekripsi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		keyPlainOutput,
		widget.NewButtonWithIcon("Salin Plaintext", theme.ContentCopyIcon(), func() {
			if keyPlainOutput.Text != "" {
				p.copyToClip(keyPlainOutput.Text)
			}
		}),
	)
	keyDecryptBox.Hide()

	// Sub-mode 6: MD5 Decrypt (Reverse Hash Lookup)
	md5InputEntry := components.NewScrollableEntry()
	md5InputEntry.SetPlaceHolder("Tempelkan hash MD5 32-karakter heksadesimal (contoh: 5f4dcc3b5aa765d61d8327deb882cf99)...")
	md5InputEntry.SetText("5f4dcc3b5aa765d61d8327deb882cf99")

	md5ModeSelect := widget.NewSelect([]string{
		"Hybrid (Online API + Kamus Lokal)",
		"Kamus Lokal Saja (Offline)",
		"Custom Wordlist",
	}, nil)
	md5ModeSelect.SetSelected("Hybrid (Online API + Kamus Lokal)")

	md5CustomWordsEntry := widget.NewMultiLineEntry()
	md5CustomWordsEntry.SetPlaceHolder("Daftar kata sandi kustom tambahan (satu kata per baris)...")
	md5CustomWordsEntry.Wrapping = fyne.TextWrapWord
	md5CustomWordsEntry.SetMinRowsVisible(3)
	md5CustomWordsEntry.Hide()

	md5ModeSelect.OnChanged = func(s string) {
		if s == "Custom Wordlist" {
			md5CustomWordsEntry.Show()
		} else {
			md5CustomWordsEntry.Hide()
		}
	}

	md5StatusBadgeBox := container.NewHBox(components.BadgeMuted("STATUS PENCARIAN"))
	md5ResultEntry := widget.NewMultiLineEntry()
	md5ResultEntry.SetPlaceHolder("Hasil plaintext hasil pencarian akan tampil di sini...")
	md5ResultEntry.Wrapping = fyne.TextWrapWord
	md5ResultEntry.SetMinRowsVisible(3)

	btnCrackMD5 := widget.NewButtonWithIcon("Cari Plaintext MD5", theme.SearchIcon(), func() {
		rawHash := strings.TrimSpace(md5InputEntry.Text)
		if len(rawHash) != 32 {
			dialog.ShowInformation("Format Hash Salah", "Hash MD5 harus tepat 32 karakter heksadesimal.", p.window)
			return
		}

		md5StatusBadgeBox.Objects = nil
		md5StatusBadgeBox.Add(components.BadgeWarning("SEDANG MEMERIKSA..."))
		md5StatusBadgeBox.Refresh()
		md5ResultEntry.SetText("Memproses pencarian hash...")

		var customList []string
		if md5ModeSelect.Selected == "Custom Wordlist" && md5CustomWordsEntry.Text != "" {
			lines := strings.Split(md5CustomWordsEntry.Text, "\n")
			for _, line := range lines {
				if t := strings.TrimSpace(line); t != "" {
					customList = append(customList, t)
				}
			}
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
			defer cancel()

			var plain string
			var source string
			var found bool

			if md5ModeSelect.Selected == "Kamus Lokal Saja (Offline)" {
				plain, found = security.LookupMD5Local(rawHash)
				source = "Kamus Lokal Bawaan"
			} else if md5ModeSelect.Selected == "Custom Wordlist" {
				plain, found = security.LookupMD5Custom(rawHash, customList)
				if !found {
					plain, found = security.LookupMD5Local(rawHash)
					source = "Kamus Lokal Bawaan"
				} else {
					source = "Wordlist Kustom Pengguna"
				}
			} else {
				res, err := security.ReverseMD5Hybrid(ctx, rawHash, customList)
				if err == nil && res != nil {
					plain = res.Plaintext
					source = res.Source
					found = true
				}
			}

			// Render update to GUI
			if p.window != nil {
				if found && plain != "" {
					md5StatusBadgeBox.Objects = nil
					md5StatusBadgeBox.Add(components.BadgeSuccess("DITEMUKAN"))
					md5StatusBadgeBox.Add(components.BadgeCyan("Sumber: " + source))
					md5StatusBadgeBox.Refresh()

					md5ResultEntry.SetText(fmt.Sprintf("Plaintext: %s\nHash MD5 : %s\nSumber   : %s", plain, rawHash, source))
					dialog.ShowInformation("Hash Terpecahkan!", fmt.Sprintf("Plaintext ditemukan: %s\n(Sumber: %s)", plain, source), p.window)
				} else {
					md5StatusBadgeBox.Objects = nil
					md5StatusBadgeBox.Add(components.BadgeDanger("TIDAK DITEMUKAN"))
					md5StatusBadgeBox.Refresh()

					md5ResultEntry.SetText("Hash MD5 tidak ditemukan pada kamus lokal maupun server online.\nTips: Coba gunakan Custom Wordlist atau kombinasikan dengan kamus kata sandi yang lebih spesifik.")
				}
			}
		}()
	})
	btnCrackMD5.Importance = widget.HighImportance

	btnSampleMD5 := widget.NewButtonWithIcon("Muat Contoh MD5", theme.HelpIcon(), func() {
		md5InputEntry.SetText("5f4dcc3b5aa765d61d8327deb882cf99")
	})
	btnSampleMD5.Importance = widget.LowImportance

	md5DecryptBox := container.NewVBox(
		widget.NewLabelWithStyle("MD5 Decrypt (Reverse Hash Lookup):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Pencarian plaintext dari hash MD5 menggunakan database rainbow table online dan kamus offline."),
		container.NewHBox(
			widget.NewLabel("Target Hash MD5:"),
			widget.NewLabel("Metode:"),
			md5ModeSelect,
			btnSampleMD5,
		),
		md5InputEntry,
		md5CustomWordsEntry,
		btnCrackMD5,
		widget.NewSeparator(),
		container.NewBorder(nil, nil, widget.NewLabelWithStyle("Status Hasil:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), md5StatusBadgeBox),
		md5ResultEntry,
		widget.NewButtonWithIcon("Salin Hasil Plaintext", theme.ContentCopyIcon(), func() {
			if md5ResultEntry.Text != "" {
				p.copyToClip(md5ResultEntry.Text)
			}
		}),
	)
	md5DecryptBox.Hide()

	// Sub-mode 7: Blockchain & Base Converter (Hex, BigInt, Wei, Base58, ABI)
	baseInputEntry := components.NewScrollableEntry()
	baseInputEntry.SetPlaceHolder("Masukkan angka (contoh: 0x2a, 42, 0b101010, atau angka besar 256-bit)...")
	baseInputEntry.SetText("0x2a")

	baseRadixSelect := widget.NewSelect([]string{
		"Auto Detect (0x, 0b, 0o, desimal)",
		"Hexadecimal (Basis 16)",
		"Decimal (Basis 10)",
		"Binary (Basis 2)",
		"Octal (Basis 8)",
	}, nil)
	baseRadixSelect.SetSelected("Auto Detect (0x, 0b, 0o, desimal)")

	btnSampleBase := widget.NewButtonWithIcon("Muat Contoh BigInt", theme.HelpIcon(), func() {
		baseRadixSelect.SetSelected("Hexadecimal (Basis 16)")
		baseInputEntry.SetText("0xde0b6b3a7640000") // 1 Ether in Wei (10^18)
	})
	btnSampleBase.Importance = widget.LowImportance

	baseResultDec := components.NewScrollableEntry()
	baseResultDec.Disable()
	baseResultHex := components.NewScrollableEntry()
	baseResultHex.Disable()
	baseResultBin := components.NewScrollableEntry()
	baseResultBin.Disable()
	baseResultOct := components.NewScrollableEntry()
	baseResultOct.Disable()
	baseMetaBadgeBox := container.NewHBox(
		components.BadgeCyan("0 BITS"),
		components.BadgeIndigo("0 BYTES"),
	)

	makeCopyRow := func(title string, entry *widget.Entry) fyne.CanvasObject {
		btnCopy := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
			if entry.Text != "" {
				p.copyToClip(entry.Text)
			}
		})
		btnCopy.Importance = widget.LowImportance
		return container.NewBorder(
			nil, nil,
			container.NewHBox(widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
			btnCopy,
			entry,
		)
	}

	btnConvertBase := widget.NewButtonWithIcon("Konversi Bilangan Basis", theme.ViewRefreshIcon(), func() {
		raw := strings.TrimSpace(baseInputEntry.Text)
		if raw == "" {
			dialog.ShowInformation("Perhatian", "Masukkan bilangan terlebih dahulu.", p.window)
			return
		}

		var mode string
		switch baseRadixSelect.Selected {
		case "Hexadecimal (Basis 16)":
			mode = "hex"
		case "Decimal (Basis 10)":
			mode = "dec"
		case "Binary (Basis 2)":
			mode = "bin"
		case "Octal (Basis 8)":
			mode = "oct"
		default:
			mode = "auto"
		}

		res, err := security.ConvertBaseNumber(raw, mode)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}

		baseResultDec.SetText(res.Decimal)
		baseResultHex.SetText(res.Hex)
		baseResultBin.SetText(res.Binary)
		baseResultOct.SetText(res.Octal)

		baseMetaBadgeBox.Objects = nil
		baseMetaBadgeBox.Add(components.BadgeCyan(fmt.Sprintf("%d BITS", res.BitLen)))
		baseMetaBadgeBox.Add(components.BadgeIndigo(fmt.Sprintf("%d BYTES", res.ByteLen)))
		baseMetaBadgeBox.Refresh()
	})
	btnConvertBase.Importance = widget.HighImportance

	// 2. Ethereum Unit Converter
	weiInputEntry := components.NewScrollableEntry()
	weiInputEntry.SetPlaceHolder("Jumlah nilai (contoh: 1.5 atau 1000000000000000000)...")
	weiInputEntry.SetText("1")

	weiUnitSelect := widget.NewSelect([]string{"Ether (ETH)", "Gwei", "Wei"}, nil)
	weiUnitSelect.SetSelected("Ether (ETH)")

	btnSampleWei := widget.NewButtonWithIcon("Muat 1 ETH", theme.HelpIcon(), func() {
		weiUnitSelect.SetSelected("Ether (ETH)")
		weiInputEntry.SetText("1")
	})
	btnSampleWei.Importance = widget.LowImportance

	weiResWei := components.NewScrollableEntry()
	weiResWei.Disable()
	weiResGwei := components.NewScrollableEntry()
	weiResGwei.Disable()
	weiResEth := components.NewScrollableEntry()
	weiResEth.Disable()

	btnConvertWei := widget.NewButtonWithIcon("Kalkulasi Unit Ethereum", theme.ViewRefreshIcon(), func() {
		amt := strings.TrimSpace(weiInputEntry.Text)
		if amt == "" {
			dialog.ShowInformation("Perhatian", "Masukkan jumlah saldo.", p.window)
			return
		}

		var u string
		switch weiUnitSelect.Selected {
		case "Ether (ETH)":
			u = "ether"
		case "Gwei":
			u = "gwei"
		default:
			u = "wei"
		}

		res, err := security.ConvertWei(amt, u)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}

		weiResWei.SetText(res.Wei)
		weiResGwei.SetText(res.Gwei)
		weiResEth.SetText(res.Ether)
	})
	btnConvertWei.Importance = widget.HighImportance

	// 3. Web3, ABI, & CTF Encodings
	selectorInput := components.NewScrollableEntry()
	selectorInput.SetPlaceHolder("Solidity function signature (contoh: transfer(address,uint256))...")
	selectorInput.SetText("transfer(address,uint256)")

	selectorOutput := components.NewScrollableEntry()
	selectorOutput.Disable()

	btnCalcSelector := widget.NewButtonWithIcon("Hitung 4-Byte Selector", theme.ConfirmIcon(), func() {
		sig := strings.TrimSpace(selectorInput.Text)
		if sig == "" {
			dialog.ShowInformation("Perhatian", "Masukkan signature fungsi Solidity.", p.window)
			return
		}
		sel := security.ComputeFunctionSelector(sig)
		hash := security.Keccak256Hash([]byte(sig))
		selectorOutput.SetText(fmt.Sprintf("Selector: %s | Full Keccak256: %s", sel, hash))
	})
	btnCalcSelector.Importance = widget.LowImportance

	btnCopySelector := widget.NewButtonWithIcon("Salin Selector", theme.ContentCopyIcon(), func() {
		if selectorOutput.Text != "" {
			p.copyToClip(selectorOutput.Text)
		}
	})
	btnCopySelector.Importance = widget.LowImportance

	endianInput := components.NewScrollableEntry()
	endianInput.SetPlaceHolder("Hex input untuk dibalik endianness (contoh: 0x12345678)...")
	endianInput.SetText("0x12345678")

	endianOutput := components.NewScrollableEntry()
	endianOutput.Disable()

	btnSwapEndian := widget.NewButtonWithIcon("Balik Endianness", theme.ViewRefreshIcon(), func() {
		res, err := security.ReverseEndiannessHex(endianInput.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		endianOutput.SetText(res)
	})
	btnSwapEndian.Importance = widget.LowImportance

	btnCopyEndian := widget.NewButtonWithIcon("Salin Endian", theme.ContentCopyIcon(), func() {
		if endianOutput.Text != "" {
			p.copyToClip(endianOutput.Text)
		}
	})
	btnCopyEndian.Importance = widget.LowImportance

	b58Input := components.NewScrollableEntry()
	b58Input.SetPlaceHolder("Teks atau Base58 string (contoh: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa)...")
	b58Input.SetText("Hello Web3")

	b58Output := components.NewScrollableEntry()
	b58Output.Disable()

	btnEncodeB58 := widget.NewButton("Encode Base58", func() {
		res := security.EncodeBase58([]byte(b58Input.Text))
		b58Output.SetText(res)
	})
	btnEncodeB58.Importance = widget.LowImportance

	btnDecodeB58 := widget.NewButton("Decode Base58", func() {
		bytes, err := security.DecodeBase58(b58Input.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		b58Output.SetText(string(bytes))
	})
	btnDecodeB58.Importance = widget.LowImportance

	btnCopyB58 := widget.NewButtonWithIcon("Salin Base58", theme.ContentCopyIcon(), func() {
		if b58Output.Text != "" {
			p.copyToClip(b58Output.Text)
		}
	})
	btnCopyB58.Importance = widget.LowImportance

	calldataInput := components.NewScrollableEntry()
	calldataInput.SetPlaceHolder("Hex calldata (0x48656c6c6f) atau teks ASCII biasa...")
	calldataInput.SetText("0x48656c6c6f2057656233")

	calldataOutput := components.NewScrollableEntry()
	calldataOutput.Disable()

	btnHexToAscii := widget.NewButton("Hex ke ASCII", func() {
		res, err := security.HexToASCII(calldataInput.Text)
		if err != nil {
			dialog.ShowError(err, p.window)
			return
		}
		calldataOutput.SetText(res)
	})
	btnHexToAscii.Importance = widget.LowImportance

	btnAsciiToHex := widget.NewButton("ASCII ke Hex", func() {
		res := security.ASCIIToHex(calldataInput.Text)
		calldataOutput.SetText(res)
	})
	btnAsciiToHex.Importance = widget.LowImportance

	btnCopyCalldata := widget.NewButtonWithIcon("Salin Calldata", theme.ContentCopyIcon(), func() {
		if calldataOutput.Text != "" {
			p.copyToClip(calldataOutput.Text)
		}
	})
	btnCopyCalldata.Importance = widget.LowImportance

	blockchainBox := container.NewVBox(
		widget.NewLabelWithStyle("Blockchain & Base Converter (Hex, BigInt, Wei, Base58, EVM ABI):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Peralatan komprehensif untuk analisis bilangan multi-radiks presisi tinggi, kalkulasi mata uang EVM, dan decoding CTF Blockchain."),
		widget.NewSeparator(),

		// Bagian 1
		widget.NewLabelWithStyle("1. Konversi Bilangan Multi-Radiks & BigInt (256-bit Safe):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3,
			baseInputEntry,
			baseRadixSelect,
			container.NewHBox(btnConvertBase, btnSampleBase),
		),
		baseMetaBadgeBox,
		makeCopyRow("Desimal (Dec / BigInt):", baseResultDec),
		makeCopyRow("Heksadesimal (Hex):     ", baseResultHex),
		makeCopyRow("Biner (Binary):          ", baseResultBin),
		makeCopyRow("Oktal (Octal):           ", baseResultOct),
		widget.NewSeparator(),

		// Bagian 2
		widget.NewLabelWithStyle("2. Kalkulator Satuan Ethereum (Wei, Gwei, Ether):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(3,
			weiInputEntry,
			weiUnitSelect,
			container.NewHBox(btnConvertWei, btnSampleWei),
		),
		makeCopyRow("Wei (10^0):   ", weiResWei),
		makeCopyRow("Gwei (10^9):  ", weiResGwei),
		makeCopyRow("Ether (10^18):", weiResEth),
		widget.NewSeparator(),

		// Bagian 3
		widget.NewLabelWithStyle("3. Solidity ABI & Encodings CTF Web3:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("A. 4-Byte EVM Function Selector (Keccak-256):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, container.NewHBox(btnCalcSelector, btnCopySelector), selectorInput),
		selectorOutput,

		widget.NewLabelWithStyle("B. Byte Order Swapper (Big-Endian <-> Little-Endian):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, container.NewHBox(btnSwapEndian, btnCopyEndian), endianInput),
		endianOutput,

		widget.NewLabelWithStyle("C. Base58 Encoder / Decoder (Bitcoin / IPFS):", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, container.NewHBox(btnEncodeB58, btnDecodeB58, btnCopyB58), b58Input),
		b58Output,

		widget.NewLabelWithStyle("D. Calldata Hex <-> Teks ASCII:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, container.NewHBox(btnHexToAscii, btnAsciiToHex, btnCopyCalldata), calldataInput),
		calldataOutput,
	)
	blockchainBox.Hide()

	// Main mode switcher for Tab 4
	toolsRadio := widget.NewRadioGroup([]string{
		"Tabel ASCII & Lookup (0-255)",
		"Penerjemah ASCII (ASCII Translator)",
		"JWT Inspector (Offline Token Decoder)",
		"Multi-Encoder & Hash Identifier",
		"Key Decryptor (AES, DES, Vigenere, XOR)",
		"MD5 Decrypt (Reverse Hash Lookup)",
		"Blockchain & Base Converter (Hex, BigInt, Wei, Base58, ABI)",
	}, func(s string) {
		asciiTableBox.Hide()
		asciiTranslatorBox.Hide()
		jwtBox.Hide()
		encoderBox.Hide()
		keyDecryptBox.Hide()
		md5DecryptBox.Hide()
		blockchainBox.Hide()

		switch s {
		case "Tabel ASCII & Lookup (0-255)":
			asciiTableBox.Show()
		case "Penerjemah ASCII (ASCII Translator)":
			asciiTranslatorBox.Show()
		case "JWT Inspector (Offline Token Decoder)":
			jwtBox.Show()
		case "Multi-Encoder & Hash Identifier":
			encoderBox.Show()
		case "Key Decryptor (AES, DES, Vigenere, XOR)":
			keyDecryptBox.Show()
		case "MD5 Decrypt (Reverse Hash Lookup)":
			md5DecryptBox.Show()
		case "Blockchain & Base Converter (Hex, BigInt, Wei, Base58, ABI)":
			blockchainBox.Show()
		}
	})
	toolsRadio.Horizontal = false
	toolsRadio.SetSelected("Tabel ASCII & Lookup (0-255)")

	formatGuideEncodersLbl := widget.NewLabel("1. Key Decryptor: Format input AES/DES mendukung Base64, Hex, atau Plain Text. Kunci rahasia 16 byte (AES-128), 24 byte (AES-192), 32 byte (AES-256), atau 8 byte (DES). Vigenere menggunakan kunci kata alfabet A-Z.\n2. MD5 Decrypt: Masukkan 32 karakter hex MD5 valid (contoh: 5f4dcc3b5aa765d61d8327deb882cf99 untuk 'password').\n3. Blockchain Converter: Mendukung bilangan hex prefiks 0x, biner 0b, oktal 0o, desimal integer hingga 256-bit, konversi saldo Wei/Gwei/Ether, 4-byte selector (contoh: 'transfer(address,uint256)'), dan Base58.")
	formatGuideEncodersLbl.Wrapping = fyne.TextWrapWord

	formatGuideEncoders := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Panduan Format Input Perkakas Kriptografi & Blockchain:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		formatGuideEncodersLbl,
	), constants.ColorBorderSubtle)

	rootCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Pilih Perkakas Encoder & Kriptografi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		toolsRadio,
		widget.NewSeparator(),
		asciiTableBox,
		asciiTranslatorBox,
		jwtBox,
		encoderBox,
		keyDecryptBox,
		md5DecryptBox,
		blockchainBox,
	), constants.ColorAccentCyan)

	return container.NewVScroll(container.NewVBox(
		formatGuideEncoders,
		rootCard,
	))
}

// ----------------------------------------------------------------------------
// 5. Tab Hardening Playbook
// ----------------------------------------------------------------------------
func (p *SecurityPage) buildPlaybookTab() fyne.CanvasObject {
	searchEntry := components.NewScrollableEntry()
	searchEntry.SetPlaceHolder("Cari playbook, perintah hardening, atau mitigasi OWASP...")

	catSelect := widget.NewSelect(security.GetPlaybookCategories(), nil)
	catSelect.SetSelected("Semua")

	playbookListBox := container.NewVBox()

	renderPlaybooks := func() {
		playbookListBox.Objects = nil
		items := security.FilterPlaybooks(searchEntry.Text, catSelect.Selected)

		if len(items) == 0 {
			playbookListBox.Add(container.NewCenter(widget.NewLabel("Tidak ditemukan playbook yang sesuai kata kunci.")))
			playbookListBox.Refresh()
			return
		}

		for _, item := range items {
			itCopy := item

			var sevBadge fyne.CanvasObject
			switch itCopy.Severity {
			case "CRITICAL":
				sevBadge = components.BadgeDanger("CRITICAL")
			case "HIGH":
				sevBadge = components.BadgeWarning("HIGH")
			default:
				sevBadge = components.BadgeCyan("INFO")
			}

			titleLabel := widget.NewLabelWithStyle(itCopy.Title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			titleLabel.Wrapping = fyne.TextWrapWord

			header := container.NewBorder(
				nil, nil,
				container.NewHBox(
					components.BadgeIndigo(itCopy.ID),
					sevBadge,
				),
				components.BadgeMuted(itCopy.Category),
				titleLabel,
			)

			descLabel := widget.NewLabel(itCopy.Description)
			descLabel.Wrapping = fyne.TextWrapWord

			stepsBox := container.NewVBox()
			for idx, step := range itCopy.Steps {
				sLbl := widget.NewLabel(fmt.Sprintf("%d. %s", idx+1, step))
				sLbl.Wrapping = fyne.TextWrapWord
				stepsBox.Add(sLbl)
			}

			remedLabel := widget.NewLabel(fmt.Sprintf("Dampak / Manfaat: %s", itCopy.Remediation))
			remedLabel.Wrapping = fyne.TextWrapWord

			cardItems := []fyne.CanvasObject{
				header,
				descLabel,
				widget.NewLabelWithStyle("Langkah-Langkah Implementasi:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				stepsBox,
				remedLabel,
			}

			if itCopy.Command != "" {
				cmdEntry := components.NewScrollableEntry()
				cmdEntry.SetText(itCopy.Command)
				cmdEntry.Disable()

				btnCopyCmd := widget.NewButtonWithIcon("Salin Perintah", theme.ContentCopyIcon(), func() {
					p.copyToClip(itCopy.Command)
				})
				btnCopyCmd.Importance = widget.HighImportance

				cmdRow := container.NewBorder(
					nil, nil,
					widget.NewLabelWithStyle("Perintah Cepat:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
					btnCopyCmd,
					container.NewHScroll(cmdEntry),
				)
				cardItems = append(cardItems, cmdRow)
			}

			itemCard := components.NewPlainCardWithAccent(container.NewVBox(cardItems...), constants.ColorBorderSubtle)
			playbookListBox.Add(itemCard)
		}
		playbookListBox.Refresh()
	}

	searchEntry.OnChanged = func(_ string) { renderPlaybooks() }
	catSelect.OnChanged = func(_ string) { renderPlaybooks() }
	renderPlaybooks()

	filterCard := components.NewPlainCardWithAccent(container.NewVBox(
		widget.NewLabelWithStyle("Katalog Playbook Pengerasan Sistem & Respon Insiden:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewGridWithColumns(2,
			searchEntry,
			container.NewHBox(widget.NewLabel("Kategori:"), catSelect),
		),
	), constants.ColorAccentCyan)

	return container.NewVScroll(container.NewVBox(
		filterCard,
		playbookListBox,
	))
}
