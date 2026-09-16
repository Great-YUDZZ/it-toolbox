# Changelog IT Toolbox

Semua perubahan dan pembaruan penting pada proyek **IT Toolbox** didokumentasikan di sini.

Format berkas ini mengacu pada [Keep a Changelog](https://keepachangelog.com/id/1.0.0/) dan menganut prinsip [Semantic Versioning](https://semver.org/).

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
