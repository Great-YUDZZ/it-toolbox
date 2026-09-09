<div align="center">

# 🧰 IT Toolbox

**All-in-One Developer, Network & System Utilities Suite**

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-blue?style=flat)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)]()
[![GitHub Release](https://img.shields.io/github/v/release/Great-YUDZZ/it-toolbox?include_prereleases&color=orange)](https://github.com/Great-YUDZZ/it-toolbox/releases)

*Aplikasi desktop modern, ringan, dan cepat untuk kebutuhan harian teknisi IT, programmer, DevOps, dan administrator jaringan.*

</div>

---

## 📥 Cara Download & Instalasi (Tinggal Klik-Klik Aja!)

Instalasi **IT Toolbox** dirancang sangat praktis dan ramah pengguna. Anda tidak perlu mengetik perintah terminal yang rumit. Cukup unduh installer sesuai sistem operasi Anda, ikuti wizard (Next > Next > Selesai), dan **shortcut aplikasi akan langsung muncul di Desktop serta Start Menu Anda!**

---

### 🐧 1. Untuk Pengguna Linux (Ubuntu, Debian, Linux Mint, Pop!_OS, Zorin, dll.)

Tersedia dua metode instalasi mudah untuk Linux:

#### Metode Utama: Menggunakan Paket `.deb` (Paling Direkomendasikan)

##### Langkah 1: Unduh Paket `.deb`
Klik link download langsung di bawah ini:
👉 **[Download it-toolbox_1.0.0_amd64.deb (v1.0.0)](https://github.com/Great-YUDZZ/it-toolbox/releases/download/v1.0.0/it-toolbox_1.0.0_amd64.deb)**

##### Langkah 2: Install dengan 1 Klik (Klik-Klik Aja)
1. Buka File Manager Anda dan cari file `it-toolbox_1.0.0_amd64.deb`.
2. **Double-click** file `.deb` tersebut.
3. Jendela **App Center / Software Center / GDebi** akan terbuka secara otomatis.
4. Klik tombol **Install** (masukkan password komputer Anda jika diminta).

##### 🎉 Hasil:
- Shortcut **IT Toolbox** **langsung muncul di Desktop (`~/Desktop`)** Anda dan langsung siap diklik tanpa peringatan *untrusted*.
- Aplikasi juga otomatis terdaftar di menu aplikasi / *Application Launcher* sistem Linux Anda.

> **Tips Terminal (Opsional):** Jika Anda lebih suka terminal, cukup ketik:
> ```bash
> sudo apt install ./it-toolbox_1.0.0_amd64.deb
> ```

---

### 🪟 2. Untuk Pengguna Windows (Windows 10 / 11 64-bit)

#### Langkah 1: Unduh File Windows
Klik salah satu link download langsung di bawah ini:
- 👉 **[Download it-toolbox.exe (Aplikasi Langsung Siap Pakai)](https://github.com/Great-YUDZZ/it-toolbox/releases/download/v1.0.0/it-toolbox.exe)**
- 👉 **[Download IT-Toolbox-Windows-Portable-x64.zip (Paket Zip)](https://github.com/Great-YUDZZ/it-toolbox/releases/download/v1.0.0/IT-Toolbox-Windows-Portable-x64.zip)**

#### Langkah 2: Cara Menggunakan & Membuat Shortcut Desktop
1. Simpan file `it-toolbox.exe` di folder yang Anda inginkan (misal di `C:\Program Files\IT Toolbox\` atau folder favorit Anda).
2. **Untuk membuat Shortcut di Desktop Windows (Langsung Muncul):**
   - **Klik kanan** pada file `it-toolbox.exe`
   - Pilih **Kirim ke / Send to** ➔ **Desktop (buat shortcut / create shortcut)**
   - Shortcut resmi **IT Toolbox** langsung muncul di layar Desktop Anda!
3. **Double-click shortcut tersebut**: Aplikasi akan langsung terbuka secara instan dengan tema gelap modern tanpa memunculkan jendela hitam CMD (*clean GUI mode*).

---

#### 💡 Ingin Membuat Setup Wizard (.exe Installer)?
Jika Anda ingin membuat paket installer `IT-Toolbox-Setup-x64.exe` (wizard Next -> Next -> Selesai):
1. Buka folder proyek ini di Windows.
2. Jalankan `build\windows\build.bat`.
3. Installer setup lengkap dengan wizard otomatis akan langsung ter-generate di folder `dist\`.

---

#### Metode Alternatif: Portable 1-Click Installer (Untuk Arch, Fedora, openSUSE / Tanpa Sudo)
Jika Anda menggunakan distro Linux non-Debian atau ingin instalasi portabel per-user tanpa hak akses root (`sudo`):
1. Clone atau unduh repositori ini.
2. Jalankan skrip instalasi portabel:
   ```bash
   bash build/linux/install-portable.sh
   ```
3. Dialog pop-up GUI akan muncul menanyakan konfirmasi instalasi.
4. Klik **Yes / OK**. Aplikasi akan diinstal ke folder `~/.local/bin/` dan shortcut otomatis diletakkan di **Desktop** Anda.

---

## ✨ Fitur-Fitur Utama IT Toolbox

| Kategori | Fitur & Fungsi |
| :--- | :--- |
| 🌐 **Network Tools** | IP Subnet Calculator, CIDR calculator, Host range finder, Network diagnostics & Tracker |
| 🔒 **Security & Hash** | Hash Generator (MD5, SHA-1, SHA-256, SHA-512), JWT Inspector & Decoder, Password Generator |
| 🔄 **Converter Suite** | Base64 Encoder/Decoder, URL Encoder, JSON ↔ YAML, Timestamp / Epoch Converter, Text Case Converter |
| 📄 **File & Doc Tools** | PDF Merger & Tools, Image Converter / Optimizer, Optical Character Recognition (OCR), ZIP Packager |
| 📚 **Quick Reference** | Database Port Directory, HTTP Status Codes dictionary, Cheat sheet perintah populer CLI/DevOps |
| 📝 **IT Logbook** | Pencatatan aktivitas harian teknisi IT berbasis SQLite lokal yang aman dan offline |

---

## 🎨 Antarmuka Modern & Performa

- **Native Desktop GUI**: Dibangun menggunakan bahasa **Go (Golang)** dan framework **Fyne v2**.
- **Tema Kustom Sleek Cyan & Dark Mode**: Dirancang nyaman untuk mata (*eye-friendly*) dengan tipografi modern (Google Sans & Gilroy).
- **Embedded Assets**: Semua font dan icon disematkan langsung ke dalam binary aplikasi (`//go:embed`), sehingga aplikasi tetap berjalan sempurna di mana pun shortcut diletakkan.

---

## 🛠️ Untuk Pengembang (Build dari Source)

Jika Anda ingin mengompilasi aplikasi atau membangun installer sendiri dari kode sumber:

### Prasyarat
- [Go](https://go.dev/dl/) versi 1.25 atau lebih baru
- Driver grafis dengan dukungan OpenGL / CGO
- Untuk Windows: [Inno Setup 6](https://jrsoftware.org/isdl.php) (untuk membuat setup `.exe`)
- Untuk Linux: `dpkg-dev` (untuk membuat paket `.deb`)

### 1. Menjalankan Aplikasi Secara Langsung
```bash
# Di Linux
go run -tags x11 .

# Di Windows
go run .
```

### 2. Membangun Paket Linux (.deb)
```bash
bash build/linux/package-deb.sh 1.0.0
```
*Hasil `.deb` akan disimpan di folder `dist/it-toolbox_1.0.0_amd64.deb`.*

### 3. Membangun Installer Windows (.exe)
Di komputer Windows (menggunakan Command Prompt atau PowerShell):
```cmd
build\windows\build.bat
```
*Hasil installer akan disimpan di folder `dist\IT-Toolbox-Setup-x64.exe`.*

---

## 🤝 Kontribusi

Kontribusi selalu disambut dengan baik!
1. Fork repository ini
2. Buat branch fitur baru (`git checkout -b feature/FiturKeren`)
3. Commit perubahan Anda (`git commit -m 'feat: tambah fitur keren'`)
4. Push ke branch Anda (`git push origin feature/FiturKeren`)
5. Buat **Pull Request** baru

---

## 📄 Lisensi

Proyek ini didistribusikan di bawah lisensi [MIT License](LICENSE).

Dibuat oleh [Great-YUDZZ](https://github.com/Great-YUDZZ).
