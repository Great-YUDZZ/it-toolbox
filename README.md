<div align="center">

# IT Toolbox

**All-in-One Developer, Network & System Utilities Suite**

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-blue?style=flat)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![GitHub Release](https://img.shields.io/github/v/release/Great-YUDZZ/it-toolbox?include_prereleases&color=orange)](https://github.com/Great-YUDZZ/it-toolbox/releases)

*Aplikasi desktop modern, ringan, dan cepat untuk kebutuhan teknisi IT, pengembang perangkat lunak, DevOps, dan administrator jaringan.*

</div>

---

## Panduan Unduhan & Instalasi

IT Toolbox menyediakan paket distribusi siap pakai untuk sistem operasi Linux dan Windows. Aplikasi dilengkapi dengan fitur integrasi sistem otomatis yang membuat pintasan pada Desktop dan Start Menu.

---

### 1. Sistem Operasi Linux (Ubuntu, Debian, Linux Mint, Pop!_OS, Zorin OS, dll.)

#### Metode Rekomendasi: Paket Debian (`.deb`)

1. **Unduh Paket Instalasi:**
   - **[Download it-toolbox_1.3.1_amd64.deb](https://github.com/Great-YUDZZ/it-toolbox/releases/download/v1.3.1/it-toolbox_1.3.1_amd64.deb)**

2. **Pemasangan Aplikasi:**
   - Klik ganda pada berkas `it-toolbox_1.3.1_amd64.deb` melalui File Manager.
   - Jendela App Center / Software Center / GDebi akan terbuka. Klik **Install**.
   - Atau melalui terminal:
     ```bash
     sudo apt install ./it-toolbox_1.3.1_amd64.deb
     ```

3. **Integrasi Sistem:**
   - Pintasan aplikasi akan otomatis dibuat di Desktop (`~/Desktop`) dan terdaftar pada Application Launcher / Menu Aplikasi sistem.

#### Metode Portabel (Arch, Fedora, openSUSE, atau Pengguna Non-Root)

1. Jalankan skrip instalasi portabel lokal:
   ```bash
   bash build/linux/install-portable.sh
   ```
2. Aplikasi akan terpasang di direktori `~/.local/bin/` dan mendaftarkan berkas desktop ke `~/.local/share/applications/`.

---

### 2. Sistem Operasi Windows (Windows 10 / 11 64-bit)

#### Berkas Distribusi Windows:
- **[Download it-toolbox.exe (Eksekusi Langsung)](https://github.com/Great-YUDZZ/it-toolbox/releases/download/v1.3.1/it-toolbox.exe)** *(Disarankan)*
- **[Download IT-Toolbox-Windows-Portable-x64.zip (Arsip Lengkap)](https://github.com/Great-YUDZZ/it-toolbox/releases/download/v1.3.1/IT-Toolbox-Windows-Portable-x64.zip)**

#### Cara Penggunaan & Pembuatan Pintasan:
1. Tempatkan berkas `it-toolbox.exe` pada folder pilihan Anda (contoh: `D:\Tools\IT-Toolbox\` atau `C:\Program Files\IT Toolbox\`).
2. Jalankan `it-toolbox.exe`. Aplikasi secara otomatis mendeteksi dan membuat pintasan di:
   - **Desktop**: `%USERPROFILE%\Desktop` (termasuk deteksi OneDrive Desktop).
   - **Start Menu**: `%APPDATA%\Microsoft\Windows\Start Menu\Programs`.
3. Jika menggunakan arsip ZIP, Anda juga dapat menjalankan `create-shortcut.bat` untuk memperbarui pintasan Desktop dan Start Menu secara instan.

---

## Fitur Utama IT Toolbox

| Modul | Deskripsi Fungsi |
| :--- | :--- |
| **Network Tools** | Kalkulator Subnet IP, Kalkulator CIDR, Usable Host Range Finder, dan Pelacak Diagnostik Jaringan |
| **Cisco Packet Tracer** | Direktori perintah CLI Cisco IOS (Router, Switch, PC), generator skrip konfigurasi instan, parameter kustom, dan panduan verifikasi topologi lab |
| **Security & Cryptography** | Generator Hash (MD5, SHA-1, SHA-256, SHA-512), JWT Inspector & Decoder, dan Generator Password Aman |
| **Data & Text Converter** | Base64 Encoder/Decoder, URL Encoder, JSON <-> YAML Converter, Epoch / Unix Timestamp Converter, dan Text Case Formatter |
| **Document & Media Tools** | Penggabung & Pemisah PDF, Image Converter / Optimizer, Optical Character Recognition (OCR), dan Pengemas Arsip ZIP |
| **YouTube Media Downloader** | Pengunduh video dan audio YouTube dengan pemindai resolusi dinamis (4K, 1080p, 720p, 360p, hingga Audio M4A), visual progress bar, dan auto-muxing FFmpeg |
| **System Quick Reference** | Direktori Port Jaringan standar, kamus HTTP Status Code, dan panduan perintah CLI / DevOps |
| **IT Incident Logbook** | Buku catatan aktivitas harian, pemecahan masalah (troubleshooting), checklist tugas, dan basis data error berbasis SQLite lokal (offline) |
| **Pengaturan & Pembaruan** | Pemeriksaan rilis versi terbaru GitHub secara otomatis, pengalih tema antarmuka, dan manajemen pintasan sistem |

---

## Arsitektur Antarmuka & Performa

- **Native Desktop Application**: Dikompilasi menggunakan bahasa pemrograman **Go** dengan framework GUI **Fyne v2**.
- **Multi-Theme Interface**: Mendukung tema Neo-Brutalism berdaya kontras tinggi, Neumorphism Light (Soft Glass), dan Neumorphism Dark (Midnight Slate).
- **Offline First**: Seluruh modul kalkulasi, konversi, basis data referensi, dan SQLite logbook berjalan tanpa ketergantungan koneksi internet.
- **Embedded Static Assets**: Seluruh font dan aset icon tertanam langsung di dalam berkas binary (`//go:embed`), menjaga portabilitas aplikasi.

---

## Petunjuk Kompilasi dari Source

### Prasyarat
- [Go](https://go.dev/dl/) versi 1.25 atau lebih baru
- Driver grafis dengan akselerasi OpenGL / CGO
- Windows: MinGW-w64 (`x86_64-w64-mingw32-gcc`) dan Inno Setup 6 (opsional)
- Linux: `libgl1-mesa-dev`, `xorg-dev`, `dpkg-dev`

### Menjalankan Aplikasi dalam Lingkungan Pengembangan
```bash
# Menjalankan di Linux
go run -tags x11 .

# Menjalankan di Windows
go run .
```

### Membangun Paket Distribusi
```bash
# Kompilasi Binary Linux
go build -tags x11 -ldflags="-s -w" -o it-toolbox .

# Kompilasi Binary Windows (Cross-compile dari Linux)
CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o it-toolbox.exe .

# Membangun Paket Installer Debian (.deb)
bash build/linux/package-deb.sh 1.3.1
```

---

## Catatan Rilis & Riwayat Versi

Riwayat pembaruan dan catatan rilis terdokumentasi pada berkas **[CHANGELOG.md](CHANGELOG.md)**.

---

## Lisensi

Proyek ini didistribusikan di bawah ketentuan lisensi [MIT License](LICENSE).

Hak Cipta (c) [Great-YUDZZ](https://github.com/Great-YUDZZ).
