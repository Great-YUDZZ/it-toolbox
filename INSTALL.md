# Panduan Instalasi & Build Installer IT Toolbox

Aplikasi **IT Toolbox** telah dilengkapi dengan sistem installer GUI (*wizard-based*) dan pembuatan **shortcut otomatis di Desktop dan Start Menu** untuk sistem operasi **Windows** dan **Linux**.

---

## 1. Panduan Pengguna (User Installation)

### A. Di Windows (Tinggal Klik-Klik Aja)
1. Unduh atau buka file installer: **`IT-Toolbox-Setup-x64.exe`**.
2. **Double-click** file installer tersebut.
3. Ikuti wizard:
   - Pilih bahasa (tersedia Bahasa Indonesia & English).
   - Klik **Lanjut / Next**.
   - Opsi *"Buat icon di Desktop"* sudah tercentang otomatis.
   - Klik **Pasang / Install**.
   - Klik **Selesai / Finish**.
4. **Selesai!**
   - Shortcut **IT Toolbox** langsung muncul di **Desktop**.
   - Shortcut juga tersedia di **Start Menu**.
   - Aplikasi dapat langsung berjalan (tanpa jendela CMD hitam).

---

### B. Di Linux (Ubuntu, Debian, Linux Mint, Pop!_OS, Zorin, dll.)

#### Opsi 1: Menggunakan Paket `.deb` (Paling Praktis)
1. Unduh atau buka file **`dist/it-toolbox_1.0.0_amd64.deb`**.
2. **Double-click** file `.deb` tersebut di File Manager.
3. Jendela **Software Center / App Center / GDebi** akan terbuka.
4. Klik tombol **Install** dan masukkan password komputer Anda.
5. **Selesai!**
   - Shortcut **IT Toolbox** langsung muncul di **Desktop** Anda dan sudah siap diklik.
   - Aplikasi juga langsung terdaftar di menu launcher / App Drawer sistem.

*(Catatan: Jika ingin menginstal via terminal, cukup ketik: `sudo apt install ./dist/it-toolbox_1.0.0_amd64.deb`)*

#### Opsi 2: Portable 1-Click Installer (Tanpa Sudo)
Jika Anda menggunakan distro Linux non-Debian (Arch, Fedora, openSUSE) atau ingin instalasi lokal per-user:
1. Jalankan skrip:
   ```bash
   bash build/linux/install-portable.sh
   ```
2. Dialog GUI akan muncul menanyakan konfirmasi instalasi.
3. Klik **Yes / OK**, shortcut akan langsung dibuat di `~/Desktop` dan `~/.local/share/applications/`.

---

## 2. Panduan Pengembang (Cara Membuat Installer Baru)

### A. Membangun Installer Linux (.deb)
Pastikan Anda berada di direktori root proyek di Linux, lalu jalankan:
```bash
bash build/linux/package-deb.sh 1.0.0
```
Hasil installer `.deb` akan langsung tersedia di folder `dist/`:
- `dist/it-toolbox_1.0.0_amd64.deb`

Skrip ini secara otomatis:
- Mengompilasi binary Go dengan optimasi ukuran `-ldflags "-s -w"` dan `-tags x11`.
- Menyematkan desktop entry dan icon resolusi tinggi.
- Memasang script `postinst` untuk mendeteksi user aktif dan langsung meletakkan shortcut di folder `~/Desktop`.
- Memasang script `postrm` untuk membersihkan shortcut saat aplikasi di-uninstall.

---

### B. Membangun Installer Windows (.exe)
1. Buka folder proyek di komputer Windows.
2. Pastikan [Inno Setup 6](https://jrsoftware.org/isdl.php) sudah terpasang.
3. Jalankan salah satu skrip berikut:
   - **Command Prompt (CMD)**:
     ```cmd
     build\windows\build.bat
     ```
   - **PowerShell**:
     ```powershell
     .\build\windows\build.ps1
     ```
4. File installer **`dist\IT-Toolbox-Setup-x64.exe`** akan langsung dibuat.

---

### C. Build Otomatis di Cloud (GitHub Actions CI/CD)
Proyek ini telah dikonfigurasi dengan workflow [`.github/workflows/build-installers.yml`](file:///.github/workflows/build-installers.yml).
- Setiap kali Anda melakukan `git push` ke branch `main`/`master` atau membuat Git Tag (misal `v1.0.0`), GitHub Actions akan secara otomatis:
  1. Mengompilasi paket Linux `.deb` di runner Ubuntu.
  2. Mengompilasi installer Windows `.exe` di runner Windows menggunakan Inno Setup.
  3. Mengunggah kedua file installer ke menu **Artifacts** di GitHub Actions untuk langsung diunduh siapa saja.
