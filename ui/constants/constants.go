package constants

// UI Text Constants to avoid hardcoded strings across the application
const (
	AppTitle       = "IT Toolbox"
	AppSubtitle    = "All-in-one Toolkit for IT Students & Developers"
	AppVersion     = "1.1.0"
	DefaultWinW    = 1050
	DefaultWinH    = 680
	SidebarWidth   = 200

	// Navigation Menu Labels
	NavToolbox     = "Toolbox & Kalkulator"
	NavReference   = "Katalog Referensi"
	NavLogbook     = "Catatan & Debugging"
	NavTracker     = "Pelacak Tugas"
	NavFileConverter = "Konverter Berkas"
	NavCisco         = "Cisco Packet Tracer"

	// Calculator Tabs
	TabConverter   = "Number & Data"
	TabSubnet      = "Subnet & IP"
	TabHashGen     = "Hash & Password"
	TabFormatter   = "Format & Encode"

	// Cisco Packet Tracer Tabs
	TabCiscoAll        = "Semua Perintah"
	TabCiscoBasic      = "Dasar & Keamanan"
	TabCiscoIP         = "IP & Interface"
	TabCiscoVLAN       = "VLAN & Switching"
	TabCiscoRouting    = "Routing (OSPF/Static)"
	TabCiscoServices   = "DHCP & Services"
	TabCiscoSecurity   = "ACL & Port Security"
	TabCiscoNAT        = "NAT & PAT"
	TabCiscoDiag       = "Show & Troubleshoot"
	TabCiscoVerify     = "Verifikasi Topologi"

	// File Converter Tabs
	TabDocConverter = "Konversi Dokumen"
	TabImgConverter = "Konversi Gambar"
	TabOCR          = "Ekstraksi Teks (OCR)"
	TabPDFTools     = "Merge & Split PDF"

	// Reference Tabs
	TabCommands    = "Commands"
	TabPorts       = "Ports"
	TabHTTPCodes   = "HTTP Codes"

	// Logbook Tabs
	TabErrorLog    = "Error Logbook"
	TabSnippets    = "Code Snippets"
	TabChecklists  = "Checklists"

	// Tracker Tabs
	TabKanban      = "Papan Kanban"
	TabProjects    = "Portofolio Proyek"

	// File Converter Buttons
	BtnUploadFile   = "Pilih Berkas"
	BtnUploadBatch  = "Pilih Banyak Berkas (Batch)"
	BtnConvertAll   = "Konversi Semua"
	BtnDownloadZip  = "Unduh Semua (.ZIP)"
	BtnMergePDF     = "Gabungkan PDF"
	BtnSplitPDF     = "Pisahkan PDF"
	BtnCompressPDF  = "Kompres Ukuran PDF"
	BtnRunOCR       = "Ekstrak Teks OCR"

	// Common Buttons & Labels
	BtnCalculate   = "Hitung"
	BtnConvert     = "Konversi"
	BtnGenerate    = "Generate"
	BtnFormat      = "Format"
	BtnClear       = "Bersihkan"
	BtnCopy        = "Salin ke Clipboard"
	BtnAdd         = "Tambah"
	BtnSave        = "Simpan"
	BtnDelete      = "Hapus"
	BtnExportMD    = "Ekspor ke Markdown"
	BtnSearch      = "Cari"

	// Search Placeholders
	SearchPlaceholder        = "Ketik untuk mencari..."
	SearchCommandPlaceholder = "Cari command (misal: docker, git, ls)..."
	SearchPortPlaceholder    = "Cari port atau protokol (misal: 80, SSH, HTTPS)..."
	SearchHTTPPlaceholder    = "Cari status code (misal: 404, Unauthorized)..."
	SearchCiscoPlaceholder   = "Cari perintah Cisco (misal: vlan, ospf, dhcp, nat, roas, ping)..."

	// Subnet Recommender
	SubnetRecommenderTitle   = "Rekomendasi Subnet Sesuai Kebutuhan Host"
	SubnetRecommenderDesc    = "Kalkulasi otomatis prefix CIDR, subnet mask, dan rentang alamat IP paling efisien"
	BtnRecommend             = "Cari Subnet Paling Cocok"
	HostInputPlaceholder     = "Jumlah host (cth: 50)..."
	BaseIPPlaceholder        = "IP Jaringan (cth: 192.168.1.0)..."

	// Status Messages
	StatusCopied   = "Berhasil disalin ke clipboard!"
	StatusSuccess  = "Operasi berhasil"
	StatusError    = "Terjadi kesalahan"
	NoDataMessage  = "Tidak ada data yang ditemukan"
)
