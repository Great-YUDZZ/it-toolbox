package components

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/yudz/it-toolbox/ui/constants"
)

// ToolItem merepresentasikan entri fitur dalam direktori pencarian cepat
type ToolItem struct {
	ID          string
	Title       string
	Description string
	Category    string
	Keywords    []string
	PageName    string
	TabName     string
	Icon        fyne.Resource
	IsPopular   bool
}

// GetAllTools mengembalikan daftar lengkap seluruh tool dan fitur dalam IT Toolbox
func GetAllTools() []ToolItem {
	return []ToolItem{
		// ALAT & KALKULATOR
		{
			ID:          "calc_subnet",
			Title:       "Subnet & IP Calculator",
			Description: "Rekomendasi prefix CIDR, kalkulasi host usable, netmask, wildcard, dan broadcast.",
			Category:    "NETWORK",
			Keywords:    []string{"subnet", "ip", "cidr", "netmask", "broadcast", "host", "usable", "wildcard", "kalkulator", "jaringan"},
			PageName:    constants.NavToolbox,
			TabName:     constants.TabSubnet,
			Icon:        theme.ComputerIcon(),
			IsPopular:   true,
		},
		{
			ID:          "calc_converter",
			Title:       "Konversi Bilangan & Data Storage",
			Description: "Konversi Desimal, Biner, Heksadesimal, Oktal, serta konverter satuan kapasitas memori.",
			Category:    "TOOLS",
			Keywords:    []string{"konversi", "biner", "hex", "desimal", "oktal", "byte", "mb", "gb", "storage", "kapasitas"},
			PageName:    constants.NavToolbox,
			TabName:     constants.TabConverter,
			Icon:        theme.MenuIcon(),
			IsPopular:   true,
		},
		{
			ID:          "calc_hashgen",
			Title:       "Hash Generator & Password Vault",
			Description: "Hasilkan hash MD5, SHA-1, SHA-256, SHA-512, generate password acak aman, dan audit entropi.",
			Category:    "SECURITY",
			Keywords:    []string{"hash", "md5", "sha256", "sha1", "sha512", "password", "generator", "entropi", "vault", "keamanan"},
			PageName:    constants.NavToolbox,
			TabName:     constants.TabHashGen,
			Icon:        theme.VisibilityIcon(),
			IsPopular:   true,
		},
		{
			ID:          "calc_formatter",
			Title:       "Format & Encode (JSON, XML, SQL, Base64)",
			Description: "Formatter dan beautifier JSON, XML, SQL, URL encoder, Base64 encoder, dan diff comparator.",
			Category:    "CODE",
			Keywords:    []string{"json", "xml", "sql", "format", "beautify", "base64", "urlencode", "diff", "minify"},
			PageName:    constants.NavToolbox,
			TabName:     constants.TabFormatter,
			Icon:        theme.DocumentCreateIcon(),
			IsPopular:   true,
		},

		// CYBER SECURITY HUB
		{
			ID:          "sec_portscanner",
			Title:       "Port Scanner & Audit Service",
			Description: "Pemindaian port TCP lokal dan remote, inspeksi latency, serta identifikasi service aktif.",
			Category:    "SECURITY",
			Keywords:    []string{"port", "scanner", "scan", "tcp", "nmap", "open port", "latency", "host"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityScanner,
			Icon:        theme.SearchIcon(),
			IsPopular:   true,
		},
		{
			ID:          "sec_webssl",
			Title:       "SSL/TLS Certificate & Header Inspector",
			Description: "Inspeksi sertifikat HTTPS, tanggal kadaluarsa, issuer, dan evaluasi HTTP Security Headers.",
			Category:    "SECURITY",
			Keywords:    []string{"ssl", "tls", "https", "sertifikat", "header", "csp", "hsts", "web security"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityWebSSL,
			Icon:        theme.VisibilityIcon(),
			IsPopular:   true,
		},
		{
			ID:          "sec_forensics",
			Title:       "Multi-Hash Forensics & Log Anomaly Audit",
			Description: "Hitung checksum berkas (MD5, SHA1, SHA256) serta analisis anomali pola serangan log server.",
			Category:    "SECURITY",
			Keywords:    []string{"forensik", "checksum", "file hash", "log", "anomali", "bruteforce", "sqli", "audit"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityForensics,
			Icon:        theme.FileIcon(),
			IsPopular:   false,
		},
		{
			ID:          "sec_keydecryptor",
			Title:       "Key Decryptor (AES, DES, 3DES, XOR, Vigenere)",
			Description: "Dekripsi simetris modern (AES-128/192/256 CBC & GCM), legacy (DES, 3DES), dan cipher klasik.",
			Category:    "SECURITY",
			Keywords:    []string{"key decryptor", "decryptor", "aes", "des", "3des", "xor", "vigenere", "dekripsi", "kunci", "cipher"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityEncoders,
			Icon:        theme.SettingsIcon(),
			IsPopular:   true,
		},
		{
			ID:          "sec_md5decrypt",
			Title:       "MD5 Decrypt (Reverse Hash Lookup)",
			Description: "Reverse lookup hash MD5 menggunakan kamus offline bawaan, custom wordlist, dan API daring.",
			Category:    "SECURITY",
			Keywords:    []string{"md5 decrypt", "md5", "reverse hash", "lookup", "rainbow table", "crack", "wordlist", "hash decrypt"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityEncoders,
			Icon:        theme.SettingsIcon(),
			IsPopular:   true,
		},
		{
			ID:          "sec_blockchain",
			Title:       "Blockchain & Base Converter (Hex, BigInt, Wei, Base58, ABI)",
			Description: "Konversi bilangan multi-radiks 256-bit (Hex, Dec, Bin, Oct), satuan Ether/Wei, 4-byte ABI selector, dan Base58 untuk CTF Blockchain.",
			Category:    "SECURITY",
			Keywords:    []string{"blockchain", "ctf", "hex", "bigint", "wei", "gwei", "ether", "base58", "keccak", "solidity", "endian", "calldata", "abi"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityEncoders,
			Icon:        theme.SettingsIcon(),
			IsPopular:   true,
		},
		{
			ID:          "sec_jwt",
			Title:       "JWT Token Inspector & Decoder",
			Description: "Dekode JSON Web Token header, payload, klaim waktu kadaluarsa (exp), dan validasi struktur.",
			Category:    "SECURITY",
			Keywords:    []string{"jwt", "json web token", "token", "auth", "bearer", "claims", "decoder"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityEncoders,
			Icon:        theme.SettingsIcon(),
			IsPopular:   false,
		},
		{
			ID:          "sec_ascii",
			Title:       "Interactive ASCII Table & Translator",
			Description: "Tabel referensi ASCII interaktif lengkap (0-127), konverter teks ke biner, oktal, hex, dan desimal.",
			Category:    "SECURITY",
			Keywords:    []string{"ascii", "ascii table", "tabel ascii", "ascii translate", "karakter", "biner ke teks", "hex to text"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityEncoders,
			Icon:        theme.SettingsIcon(),
			IsPopular:   true,
		},
		{
			ID:          "sec_playbook",
			Title:       "Security Hardening Playbook",
			Description: "Checklist pengerasan keamanan sistem operasi Linux, SSH, Nginx/Apache, firewall, dan database.",
			Category:    "SECURITY",
			Keywords:    []string{"playbook", "hardening", "keamanan", "ssh", "linux", "firewall", "iptables", "checklist"},
			PageName:    constants.NavSecurity,
			TabName:     constants.TabSecurityPlaybook,
			Icon:        theme.DocumentIcon(),
			IsPopular:   false,
		},

		// KONVERTER BERKAS & TOOLS MEDIA
		{
			ID:          "file_qrcode",
			Title:       "Generator QR Code & Berbagi Berkas",
			Description: "Buat kode QR dari tautan teks atau unggah file, foto, video dengan server unduhan HTTP lokal instan.",
			Category:    "TOOLS",
			Keywords:    []string{"qr code", "qr", "barcode", "generator qr", "bagikan berkas", "file share", "foto", "video"},
			PageName:    constants.NavFileConverter,
			TabName:     constants.TabQRCode,
			Icon:        theme.ViewFullScreenIcon(),
			IsPopular:   true,
		},
		{
			ID:          "file_doc",
			Title:       "Konverter Dokumen (PDF, DOCX, TXT, MD, HTML)",
			Description: "Konversi berkas dokumen antar format secara cepat, aman, dan tanpa koneksi internet.",
			Category:    "CONVERTER",
			Keywords:    []string{"dokumen", "pdf", "docx", "txt", "markdown", "html", "konversi dokumen"},
			PageName:    constants.NavFileConverter,
			TabName:     constants.TabDocConverter,
			Icon:        theme.DocumentIcon(),
			IsPopular:   false,
		},
		{
			ID:          "file_img",
			Title:       "Konverter & Kompresi Gambar (PNG, JPG, WEBP)",
			Description: "Ubah format gambar, sesuaikan dimensi, rotasi, kompresi kualitas, dan konversi multi-file.",
			Category:    "CONVERTER",
			Keywords:    []string{"gambar", "image", "png", "jpg", "jpeg", "webp", "kompresi", "resize", "foto"},
			PageName:    constants.NavFileConverter,
			TabName:     constants.TabImgConverter,
			Icon:        theme.FileImageIcon(),
			IsPopular:   false,
		},
		{
			ID:          "file_ocr",
			Title:       "Ekstraksi Teks dari Gambar (OCR)",
			Description: "Ekstraksi teks otomatis dari tangkapan layar, foto dokumen, atau infografis.",
			Category:    "CONVERTER",
			Keywords:    []string{"ocr", "ekstraksi teks", "scan teks", "gambar ke teks", "tesseract"},
			PageName:    constants.NavFileConverter,
			TabName:     constants.TabOCR,
			Icon:        theme.VisibilityIcon(),
			IsPopular:   false,
		},
		{
			ID:          "file_pdf",
			Title:       "Manipulasi PDF (Merge & Split)",
			Description: "Gabungkan beberapa berkas PDF menjadi satu atau pisahkan halaman tertentu menjadi file terpisah.",
			Category:    "CONVERTER",
			Keywords:    []string{"pdf", "merge pdf", "split pdf", "gabung pdf", "pisah pdf"},
			PageName:    constants.NavFileConverter,
			TabName:     constants.TabPDFTools,
			Icon:        theme.FolderOpenIcon(),
			IsPopular:   false,
		},
		{
			ID:          "file_youtube",
			Title:       "Unduh Video & Audio YouTube",
			Description: "Unduh video YouTube dalam berbagai kualitas (4K, 1080p, 720p, 360p) atau ekstrak audio M4A.",
			Category:    "TOOLS",
			Keywords:    []string{"youtube", "download youtube", "unduh video", "mp4", "m4a", "audio", "video downloader"},
			PageName:    constants.NavYouTube,
			TabName:     "",
			Icon:        theme.DownloadIcon(),
			IsPopular:   true,
		},

		// CISCO PACKET TRACER
		{
			ID:          "cisco_library",
			Title:       "Cisco IOS Command Library & Resep Topologi",
			Description: "Koleksi resep konfigurasi CLI Cisco IOS untuk VLAN, Trunking, OSPF, DHCP, NAT, dan ACL.",
			Category:    "CISCO",
			Keywords:    []string{"cisco", "packet tracer", "vlan", "trunking", "ospf", "dhcp", "nat", "router", "switch", "cli"},
			PageName:    constants.NavCisco,
			TabName:     constants.TabCiscoLibrary,
			Icon:        theme.ComputerIcon(),
			IsPopular:   true,
		},
		{
			ID:          "cisco_verify",
			Title:       "Cisco Verification & Troubleshooting Guide",
			Description: "Panduan perintah verifikasi status link, show ip interface brief, tabel routing, dan debug Cisco.",
			Category:    "CISCO",
			Keywords:    []string{"cisco show", "troubleshoot", "show ip route", "ping", "traceroute", "verifikasi"},
			PageName:    constants.NavCisco,
			TabName:     constants.TabCiscoVerify,
			Icon:        theme.ConfirmIcon(),
			IsPopular:   false,
		},

		// SCHEMA DATABASE & DDL
		{
			ID:          "db_generator",
			Title:       "Pembuat Tabel Kustom & Generator DDL",
			Description: "Rancang skema tabel database secara visual dan generate skrip SQL untuk MySQL, PostgreSQL, SQLite.",
			Category:    "DATABASE",
			Keywords:    []string{"database", "sql", "ddl", "tabel", "create table", "generator sql", "mysql", "postgres", "sqlite"},
			PageName:    constants.NavDatabase,
			TabName:     constants.TabDBGenerator,
			Icon:        theme.DocumentCreateIcon(),
			IsPopular:   true,
		},
		{
			ID:          "db_schemas",
			Title:       "Perpustakaan Skema Database Siap Pakai",
			Description: "Template skema database lengkap untuk E-Commerce, Autentikasi Pengguna, Blog, dan Manajemen Sekolah.",
			Category:    "DATABASE",
			Keywords:    []string{"skema database", "template sql", "ecommerce schema", "auth schema", "relasi tabel"},
			PageName:    constants.NavDatabase,
			TabName:     constants.TabDBSchemas,
			Icon:        theme.FolderOpenIcon(),
			IsPopular:   false,
		},
		{
			ID:          "db_datatypes",
			Title:       "Kamus Tipe Data Antar-Database",
			Description: "Tabel komparasi tipe data SQL antara MySQL, MariaDB, PostgreSQL, SQLite, dan SQL Server.",
			Category:    "DATABASE",
			Keywords:    []string{"tipe data", "data types", "varchar", "int", "timestamp", "komparasi sql"},
			PageName:    constants.NavDatabase,
			TabName:     constants.TabDBDataTypes,
			Icon:        theme.InfoIcon(),
			IsPopular:   false,
		},

		// REFERENSI & DOKUMENTASI
		{
			ID:          "ref_commands",
			Title:       "Katalog Perintah Linux & Windows CLI",
			Description: "Cheat-sheet praktis perintah terminal harian untuk Linux sysadmin, Git, Docker, dan Windows PowerShell.",
			Category:    "DOCS",
			Keywords:    []string{"perintah linux", "command line", "bash", "powershell", "git", "docker", "sysadmin", "cli"},
			PageName:    constants.NavReference,
			TabName:     constants.TabCommands,
			Icon:        theme.DocumentIcon(),
			IsPopular:   true,
		},
		{
			ID:          "ref_ports",
			Title:       "Daftar Port Jaringan Standar IANA",
			Description: "Referensi nomor port protokol jaringan populer (HTTP 80, HTTPS 443, SSH 22, DNS 53, MySQL 3306).",
			Category:    "DOCS",
			Keywords:    []string{"port", "port jaringan", "iana", "nomor port", "ssh port", "http port", "dns"},
			PageName:    constants.NavReference,
			TabName:     constants.TabPorts,
			Icon:        theme.RadioButtonIcon(),
			IsPopular:   false,
		},
		{
			ID:          "ref_http",
			Title:       "Katalog HTTP Status Codes",
			Description: "Arti dan penjelasan kode status HTTP (200 OK, 301 Redirect, 400 Bad Request, 404, 500 Server Error).",
			Category:    "DOCS",
			Keywords:    []string{"http status", "http codes", "404", "500", "200 ok", "rest api"},
			PageName:    constants.NavReference,
			TabName:     constants.TabHTTPCodes,
			Icon:        theme.InfoIcon(),
			IsPopular:   false,
		},

		// LOGBOOK & MANAJEMEN TUGAS
		{
			ID:          "log_error",
			Title:       "Error Logbook & Solusi Masalah",
			Description: "Catat histori pesan error teknis, solusi pemecahan masalah, dan penanganan bug berulang.",
			Category:    "LOGS",
			Keywords:    []string{"error log", "debugging", "catatan error", "solusi bug", "troubleshooting log"},
			PageName:    constants.NavLogbook,
			TabName:     constants.TabErrorLog,
			Icon:        theme.DocumentCreateIcon(),
			IsPopular:   false,
		},
		{
			ID:          "log_snippets",
			Title:       "Penyimpan Cuplikan Kode (Code Snippets)",
			Description: "Simpan dan kelola potongan kode favorit dalam berbagai bahasa pemrograman.",
			Category:    "LOGS",
			Keywords:    []string{"snippet", "code snippet", "cuplikan kode", "simpan kode"},
			PageName:    constants.NavLogbook,
			TabName:     constants.TabSnippets,
			Icon:        theme.DocumentCreateIcon(),
			IsPopular:   false,
		},
		{
			ID:          "tracker_kanban",
			Title:       "Papan Kanban & Manajemen Tugas",
			Description: "Visualisasikan alur kerja tugas To-Do, In Progress, dan Done dengan papan Kanban interaktif.",
			Category:    "TASKS",
			Keywords:    []string{"kanban", "tugas", "pelacak tugas", "todo", "in progress", "done", "manajemen proyek"},
			PageName:    constants.NavTracker,
			TabName:     constants.TabKanban,
			Icon:        theme.ListIcon(),
			IsPopular:   false,
		},

		// SISTEM & PENGATURAN
		{
			ID:          "sys_settings",
			Title:       "Pengaturan Tampilan & Periksa Update",
			Description: "Ganti tema visual antarmuka, kelola pintasan sistem desktop, dan cek rilis pembaruan GitHub.",
			Category:    "SYSTEM",
			Keywords:    []string{"pengaturan", "tema", "update", "pembaruan", "github update", "versi", "pintasan"},
			PageName:    constants.NavSettings,
			TabName:     "",
			Icon:        theme.SettingsIcon(),
			IsPopular:   false,
		},
	}
}

// SearchTools menyaring daftar tool berdasarkan kata kunci pencarian
func SearchTools(query string) []ToolItem {
	all := GetAllTools()
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		// Kembalikan tool populer jika query kosong
		var popular []ToolItem
		for _, t := range all {
			if t.IsPopular {
				popular = append(popular, t)
			}
		}
		return popular
	}

	terms := strings.Fields(q)
	var matches []ToolItem

	for _, item := range all {
		haystack := strings.ToLower(fmt.Sprintf("%s %s %s %s %s",
			item.Title,
			item.Description,
			item.Category,
			item.PageName,
			strings.Join(item.Keywords, " "),
		))

		allMatched := true
		for _, term := range terms {
			if !strings.Contains(haystack, term) {
				allMatched = false
				break
			}
		}

		if allMatched {
			matches = append(matches, item)
		}
	}

	return matches
}

// ShowCommandPalette menampilkan modal Command Palette pencarian tools
func ShowCommandPalette(win fyne.Window, onSelect func(pageName, tabName string)) {
	var d dialog.Dialog

	searchEntry := NewScrollableEntry()
	searchEntry.SetPlaceHolder("Ketik nama fitur, fungsi, atau kata kunci (contoh: decrypt, subnet, qr, ssl, cisco)...")

	resultsBox := container.NewVBox()
	resultsScroll := container.NewVScroll(resultsBox)
	resultsScroll.SetMinSize(fyne.NewSize(580, 360))

	renderResults := func(query string) {
		resultsBox.Objects = nil
		items := SearchTools(query)

		headerTxt := "Rekomendasi Fitur Populer:"
		if strings.TrimSpace(query) != "" {
			headerTxt = fmt.Sprintf("Hasil Pencarian (%d ditemukan):", len(items))
		}

		lblHead := canvas.NewText(headerTxt, constants.ColorTextSecondary)
		lblHead.TextSize = constants.FontSizeLabel
		lblHead.TextStyle = fyne.TextStyle{Bold: true}
		resultsBox.Add(lblHead)
		resultsBox.Add(widget.NewSeparator())

		if len(items) == 0 {
			emptyCard := container.NewPadded(container.NewVBox(
				canvas.NewText("Tidak ditemukan tool yang cocok dengan kata kunci.", constants.ColorTextMuted),
				canvas.NewText("Coba kata kunci lain seperti: ip, aes, qr, port, cisco, json, hash.", constants.ColorTextMuted),
			))
			resultsBox.Add(emptyCard)
			resultsBox.Refresh()
			return
		}

		for _, item := range items {
			it := item // capture closure

			titleTxt := canvas.NewText(it.Title, constants.ColorTextPrimary)
			titleTxt.TextSize = constants.FontSizeBody
			titleTxt.TextStyle = fyne.TextStyle{Bold: true}

			descTxt := canvas.NewText(it.Description, constants.ColorTextMuted)
			descTxt.TextSize = constants.FontSizeSmall

			var catBadge fyne.CanvasObject
			switch it.Category {
			case "SECURITY":
				catBadge = BadgeIndigo("SECURITY")
			case "NETWORK":
				catBadge = BadgeCyan("NETWORK")
			case "CONVERTER":
				catBadge = BadgeSuccess("CONVERTER")
			case "DATABASE":
				catBadge = BadgeYellow("DATABASE")
			case "CISCO":
				catBadge = BadgeCyan("CISCO")
			case "CODE":
				catBadge = BadgeSuccess("CODE")
			default:
				catBadge = BadgeCyan(it.Category)
			}

			locInfo := it.PageName
			if it.TabName != "" {
				locInfo += " > " + it.TabName
			}
			locTxt := canvas.NewText(locInfo, constants.ColorTextSecondary)
			locTxt.TextSize = constants.FontSizeLabel

			textCol := container.NewVBox(
				container.NewHBox(titleTxt, catBadge),
				descTxt,
				locTxt,
			)

			iconWidget := widget.NewIcon(it.Icon)
			rowContent := container.NewBorder(nil, nil, iconWidget, nil, container.NewPadded(textCol))

			btn := widget.NewButton("", func() {
				if d != nil {
					d.Hide()
				}
				if onSelect != nil {
					onSelect(it.PageName, it.TabName)
				}
			})
			btn.Importance = widget.LowImportance

			rowStack := container.NewStack(btn, rowContent)
			resultsBox.Add(rowStack)
			resultsBox.Add(widget.NewSeparator())
		}
		resultsBox.Refresh()
	}

	searchEntry.OnChanged = func(s string) {
		renderResults(s)
	}

	// Render awal dengan tool populer
	renderResults("")

	// Header modal
	titleModal := canvas.NewText("Quick Tool Finder (Command Palette)", constants.ColorTextPrimary)
	titleModal.TextSize = constants.FontSizeH2
	titleModal.TextStyle = fyne.TextStyle{Bold: true}

	shortcutBadge := BadgeCyan("Ctrl+K")

	headerTop := container.NewBorder(nil, nil,
		container.NewHBox(widget.NewIcon(theme.SearchIcon()), titleModal),
		shortcutBadge,
	)

	subHeader := canvas.NewText("Temukan dan buka modul atau tab spesifik langsung dalam 1-klik.", constants.ColorTextMuted)
	subHeader.TextSize = constants.FontSizeSmall

	headerBox := container.NewVBox(
		headerTop,
		subHeader,
		widget.NewSeparator(),
		searchEntry,
		widget.NewSeparator(),
	)

	// Tombol tutup
	btnClose := widget.NewButtonWithIcon("Tutup (Esc)", theme.CancelIcon(), func() {
		if d != nil {
			d.Hide()
		}
	})
	btnClose.Importance = widget.LowImportance

	footerBox := container.NewBorder(nil, nil,
		canvas.NewText("Tekan tombol atau Enter untuk membuka fitur", constants.ColorTextMuted),
		btnClose,
	)

	dialogBody := container.NewBorder(
		headerBox,
		container.NewVBox(widget.NewSeparator(), footerBox),
		nil,
		nil,
		resultsScroll,
	)

	card := NewPlainCardWithAccent(dialogBody, constants.ColorAccentCobalt)
	d = dialog.NewCustomWithoutButtons("", card, win)
	d.Resize(fyne.NewSize(620, 480))
	d.Show()

	// Fokuskan input pencarian
	if canvasObj := win.Canvas(); canvasObj != nil {
		canvasObj.Focus(searchEntry)
	}
}

// QuickSearchTriggerButton membuat tombol pemicu pencarian bergaya pencarian modern
func QuickSearchTriggerButton(onTap func()) *widget.Button {
	btn := widget.NewButtonWithIcon("Cari Tools / Fitur (Ctrl+K)...", theme.SearchIcon(), onTap)
	btn.Importance = widget.LowImportance
	return btn
}
