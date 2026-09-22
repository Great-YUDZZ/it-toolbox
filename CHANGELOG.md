# Changelog IT Toolbox

Semua perubahan dan pembaruan penting pada proyek **IT Toolbox** didokumentasikan di sini.

Format berkas ini mengacu pada [Keep a Changelog](https://keepachangelog.com/id/1.0.0/) dan menganut prinsip [Semantic Versioning](https://semver.org/).

---

## [1.5.0] - 2026-09-22

### Fitur Baru
- **Hashed Password Vault & Manajemen Keamanan Sandi**:
  - Modul penyimpanan lokal persisten SQLite (`~/.it-toolbox/it-toolbox.db`) untuk sandi yang telah di-hash secara aman.
  - Dukungan algoritma hashing lengkap: **bcrypt** (standar industri via `golang.org/x/crypto/bcrypt`), **SHA-256**, **SHA-512**, **MD5**, dan **SHA-1**.
  - Generator Hash Multi-Digest instan dengan tombol cepat **Simpan ke Vault** dan **Salin**.
  - Integrasi generator password entropy tinggi dengan opsi simpan langsung ke Vault.
  - Modal penyimpanan Neo-Brutalist dengan live preview kalkulasi hash saat mengetik, opsi salt kustom, dan pilihan penyimpanan sandi plaintext opsional.
  - Pencarian langsung (*live search*), filter berbasis algoritma, dan badge warna tematik (BCRYPT, SHA-256, SHA-512, MD5, SHA-1).
  - Masking sandi plaintext (`********`) dengan tombol intip/sembunyi dan salin sandi.
  - Alat verifikasi kecocokan sandi plaintext terhadap hash tersimpan secara real-time.
  - Ekspor seluruh data kredensial Vault ke format file **JSON** dan **CSV** yang otomatis tersalin ke clipboard.

### Perbaikan & Optimasi Tampilan (Zero-Overflow Assurance)
- **Eliminasi Pemotongan UI (Horizontal Overflow)**:
  - Memperbaiki `NewScrollableEntry` dengan scroller internal berpelindung `scrollShield` sehingga lebar minimum entry tidak membengkak saat menampilkan string digest panjang (seperti SHA-512 128-karakter).
  - Restrukturisasi tata letak header Vault dan preset chip panjang karakter, menurunkan lebar minimum halaman hingga lebih dari 60% (dari 1.288px ke 516px).
  - Memastikan seluruh tombol aksi, badge, dan konten tetap proporsional dan tidak terpotong pada berbagai resolusi layar.

---

## [1.4.1] - 2026-09-20

### Peningkatan & Perbaikan Kinerja
- **Scroll Transparan di Seluruh Area Input**:
  - Memperbaiki perilaku mouse wheel pada input pencarian (*search bar*), kotak isian parameter perintah (seperti `{{HOSTNAME}}`, `{{IP}}`), dan kotak teks multi-baris (*CLI command box* dan *SQL preview*).
  - Scrolling roda mouse saat kursor berada tepat di atas area input teks kini secara transparan langsung meneruskan pergerakan scroll ke kontainer halaman utama tanpa tertahan.
- **Eliminasi Freeze pada Modul Cisco Packet Tracer**:
  - Mengeliminasi *redundant rendering* berulang pada tab Resep dan menerapkan *progressive batching* (8 item awal dengan tombol muat semua) serta *view caching*.
  - Waktu pembukaan halaman Cisco terpangkas hingga 96%, membuka seketika (<25ms) dan instan (170ns) saat kembali dari cache.
- **Animasi Buka/Tutup Sidebar Terkunci di 60 FPS**:
  - Menerapkan `mainSlidingLayout` yang memisahkan translasi posisi kontainer dari rekalkulasi ukuran teks (*text reflow*).
  - Selama transisi animasi 180ms, elemen halaman hanya bergeser posisi tanpa kalkulasi ulang pemotongan baris, menghasilkan animasi buka/tutup sidebar yang konsisten mulus 60 FPS di seluruh halaman termasuk Cisco dan Schema Database.

---

## [1.4.0] - 2026-09-20

### Fitur Baru
- **Modul Baru: Schema Database & SQL Architect**:
  - Menu navigasi khusus **"Schema Database"** di sidebar dengan 4 tab komprehensif:
    1. **Perpustakaan Skema Siap Pakai**: DDL lengkap dengan relasi PK-FK untuk Autentikasi & RBAC, E-Commerce & Toko Online, Sistem Akademik & Lab TKJ, serta Infrastruktur Jaringan & IPAM (tersedia untuk MySQL/MariaDB, PostgreSQL, dan SQLite) dengan tombol salin dan ekspor file `.sql`.
    2. **Katalog Perintah DDL & Administrasi**: Direktori perintah terlengkap untuk pembuatan DB, hak akses user & privileges, foreign key cascading, backup/dump CLI (mysqldump, pg_dump, sqlite3), serta pemeliharaan & optimasi (vacuum, reindex, explain) disertai modal penyesuaian parameter interaktif.
    3. **Pembuat Tabel Kustom (Interactive DDL Builder)**: Generator tabel interaktif dengan pilihan preset kolom (Standar, Akun Pengguna, Produk E-Commerce), kustomisasi tipe data, nullability, unique, dan live SQL preview.
    4. **Kamus & Komparasi Tipe Data Lintas DBMS**: Matriks perbandingan tipe data antara MySQL, MariaDB, PostgreSQL, SQLite, SQL Server (T-SQL), dan Oracle Database lengkap dengan rekomendasi best practice.

---

## [1.3.2] - 2026-09-16

### Peningkatan & Fitur Baru
- **Pilihan Lokasi Pintasan di Awal Peluncuran / Instalasi**:
  - Dialog interaktif otomatis saat pertama kali aplikasi dibuka untuk memilih pembuatan pintasan di Desktop, Start Menu, atau keduanya, dengan tombol "Nanti Saja" dan "Pasang Pintasan".
  - Dukungan pemilihan pintasan di wizard instalasi Inno Setup Windows (`desktopicon` dan `startmenuicon`).
  - Menu interaktif pada skrip `create-shortcut.bat` untuk memilih Desktop, Start Menu, atau keduanya secara fleksibel.
  - Opsi pembuatan pintasan fleksibel pada kartu Integrasi Sistem di halaman Pengaturan.
- **Sistem Pembaruan Otomatis Terhubung GitHub Releases**:
  - Pengecekan pembaruan langsung terhadap endpoint GitHub Releases API.
  - Dialog notifikasi pembaruan dengan tombol "Buka Halaman Rilis", "Catatan Rilis", dan "Tutup".
- **Desain Dialog Modal Baru & Standar Visual Bebas Emoji**:
  - Efek latar belakang modal dialog dengan scrim bayangan transparan (translucent glass/scrim effect).
  - Penghapusan seluruh emoji di seluruh antarmuka aplikasi, dokumentasi README, dan skrip instalasi untuk tampilan teknis yang bersih dan profesional.
- **Perbaikan Scroll pada Area Formulir Cisco**:
  - Event scroll wheel mouse pada area teks multi-baris (seperti kolom komentar) kini diteruskan secara mulus ke kontainer scroll induk.

---

## [1.3.1] - 2026-09-15

### Fitur Baru & Peningkatan
- **Pemeriksaan Pembaruan Versi di Pengaturan**:
  - Kartu pembaruan aplikasi dengan tombol cek manual dan indikator versi saat ini vs versi rilis terbaru di GitHub.
- **Dukungan Pintasan Start Menu**:
  - Dukungan pembuatan pintasan terpadu di Desktop dan Start Menu pada Windows dan Linux.

---

## [1.3.0] - 2026-09-15

### Desain & Fitur Baru (Multi-Theme Support: Neumorphism)
- **Tema Neumorphism (Soft UI) — Terang & Gelap**:
  - Penambahan fitur penggantian tema multi-mode fleksibel melalui **Modal Dialog Pemilih Tema** dan **Quick Cycle Button** di sidebar.
  - **Neo-Brutalism (Signature)**: Mode ikonik tunggal tanpa pemisahan gelap/terang, mempertahankan estetika kontras tinggi retro paper (`#FFFDF8`), border tebal solid hitam 2.5px, bayangan tajam (*hard offset shadow +4px*), font 100% *pitch black*, dan *sticker badges* warna-warni.
  - **Neumorphism Glass — Mode Terang (Glass-Neumorphism UI)**: Perpaduan visual sejati antara *Soft UI tactile 3D extrusion* dan *Glassmorphism luminous transparency*. Menggunakan latar belakang kanvas ambient luminous gradient (gradasi lembut *sky cyan* `#D6E6FD` ke *dreamy lavender* `#EEE2FD`), kartu *frosted glass* putih berkilau (`#FFFFFF`) dengan tepian kristal reflektif 1.5px murni (`#FFFFFF`), bayangan ganda lembut (*top-left pure white light halo* + *soft cool blue-slate depth* `#C2D0E2`), sidebar bernuansa *frosted ice vertical gradient*, serta pil status dan badge kristal pastel dengan tipografi tajam berdaya kontras tinggi.
  - **Neumorphism — Mode Gelap (Dark Soft UI)**: Estetika modern slate gelap (`#21242B`) yang sangat nyaman di mata, dilengkapi *embossed dual shadow*, border halus 0.8px, sudut lengkung 14px, dan tipografi *crisp soft white* (`#F1F5F9`).
  - Adaptasi otomatis seluruh komponen (*cards, badges, KPI stat cards, nav items, search bars, dividers, and dialogs*) secara instan tanpa perlu memuat ulang aplikasi.
  - Penyimpanan preferensi tema secara persisten (`active_theme`).

---

## [1.2.0] - 2026-09-15

### Fitur Baru
- **YouTube Video Downloader**:
  - Deteksi dan ekstraksi otomatis daftar pilihan resolusi video yang benar-benar tersedia (mulai dari 4K 2160p, 2K 1440p, 1080p Full HD, 720p HD, 480p, 360p, 240p, 144p, hingga format Audio M4A/MP3 saja).
  - Ekstraksi informasi video lengkap: Judul, Channel/Author badge, Durasi badge, Jumlah pilihan kualitas, dan Status FFmpeg.
  - Estimasi ukuran berkas (MB) yang akurat pada setiap opsi resolusi di dropdown pemilih.
  - Streaming berkecepatan tinggi (7 - 10+ MB/s) dengan mekanisme bypass 403 Forbidden YouTube InnerTube.
  - Dukungan penggabungan otomatis (*auto-muxing*) stream video 1080p/4K dengan audio berformat AAC jika sistem memiliki FFmpeg.
  - Indikator kemajuan unduhan (*real-time progress bar*) menampilkan persentase, ukuran MB terunduh, dan estimasi kecepatan unduh.
  - Tombol aksi cepat pasca-unduh: **"Buka Berkas"** (memutar langsung dengan media player default) dan **"Buka Folder"** (membuka lokasi file di file manager).
- **Akses Navigasi Cepat (Quick Nav)**:
  - Penambahan menu navigasi baru **"Unduh YouTube"** langsung di sidebar sebelah kiri di bawah grup *Alat & Kalkulator*, memudahkan akses dengan 1 kali klik.
  - Tetap terintegrasi sebagai tab ke-5 pada menu **"Konverter Berkas"**.

### Perbaikan & Peningkatan
- Pembersihan karakter ilegal secara otomatis pada penamaan berkas hasil unduhan (`:`, `/`, `\`, `|`, `?`, `*`, `<`, `>`).
- Optimalisasi kompilasi Windows dengan MinGW cross-compiler untuk mode GUI murni tanpa jendela CMD hitam.

---

## [1.1.0] - 2026-09-10

### Desain & Tampilan (Redesign Neo-Brutalism)
- Pembaruan antarmuka secara menyeluruh mengadopsi estetika **Neo-Brutalism**:
  - Garis tepi tebal (*2.5px solid border*) dengan sudut membulat modern.
  - *Hard offset zero-blur shadow* khas Gumroad.
  - *Sticker badges* berwarna kontras tinggi (Gumroad Yellow, Cyber Cyan, Neo Mint, Amber, Neo Red).
  - Palet warna semantik untuk kartu status kalkulasi dan ringkasan.
- Ikon aplikasi resmi baru beresolusi tinggi (512x512) dengan integrasi multi-format `.ico` dan `.png`.
- Peningkatan tata letak agar responsif dan aman pada monitor berskala DPI tinggi (2.75x scaling).

### Fitur Cisco Packet Tracer
- Penambahan fitur kustomisasi parameter topologi jaringan (Hostname, VLAN ID, IP Address, Subnet Mask).
- Output perintah CLI otomatis menyesuaikan input pengguna.
- Fitur pencatatan langkah-langkah pembuatan topologi (*Topology Notes*) dengan kemampuan tambah, edit, dan hapus catatan.

---

## [1.0.0] - 2026-09-08

### Rilis Perdana
- **Toolbox & Kalkulator**:
  - Subnet & IP Calculator (CIDR, VLSM, Subnet Recommender).
  - Number & Data Converter (Desimal, Biner, Heksadesimal, Data Size Converter).
  - Hash & Password Generator (MD5, SHA-256, Random Token, Password).
  - Format & Encoder (JSON, YAML, XML, Base64, URL Encode).
- **Konverter Berkas**:
  - Konversi Dokumen (DOCX, PDF, TXT, XLSX).
  - Konversi Gambar (PNG, JPG, WEBP, BMP) & Kompresi.
  - Ekstraksi Teks (OCR) berbasis Tesseract.
  - PDF Merger & Splitter.
- **Katalog Referensi**:
  - Database perintah CLI populer (Docker, Git, Linux, Kubernetes).
  - Direktori port jaringan & protokol standar.
  - Direktori kode status HTTP (1xx - 5xx).
- **Catatan & Debugging**:
  - Error Logbook dengan pencarian dan filter tag.
  - Code Snippets manager.
  - Checklist teknisi IT.
- **Pelacak Tugas**:
  - Papan Kanban visual (To Do, In Progress, Done).
  - Portofolio pelacak proyek IT.
