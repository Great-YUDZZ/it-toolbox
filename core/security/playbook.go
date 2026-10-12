package security

import (
	"strings"
)

// SecurityPlaybookItem describes a practical cybersecurity hardening or response procedure
type SecurityPlaybookItem struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Category    string   `json:"category"` // "Linux Hardening", "Windows Security", "Incident Response", "OWASP Top 10"
	Severity    string   `json:"severity"` // "CRITICAL", "HIGH", "MEDIUM", "INFO"
	Description string   `json:"description"`
	Command     string   `json:"command,omitempty"`
	Steps       []string `json:"steps"`
	Remediation string   `json:"remediation"`
}

var playbookCatalog = []SecurityPlaybookItem{
	// 1. Linux Hardening
	{
		ID:          "LNX-01",
		Title:       "Pengerasan Layanan SSH (SSH Hardening Baseline)",
		Category:    "Linux Hardening",
		Severity:    "CRITICAL",
		Description: "Menonaktifkan login SSH langsung pengguna root dan membatasi autentikasi hanya menggunakan SSH Key kriptografi kuat.",
		Command:     "sudo sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin no/' /etc/ssh/sshd_config && sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config && sudo systemctl restart sshd",
		Steps: []string{
			"Pastikan sudah ada akun sudoer dengan SSH public key yang terpasang di ~/.ssh/authorized_keys.",
			"Ubah konfigurasi /etc/ssh/sshd_config: PermitRootLogin no",
			"Ubah konfigurasi: PasswordAuthentication no",
			"Tetapkan batas waktu Idle: ClientAliveInterval 300 dan ClientAliveCountMax 2",
			"Jalankan: sudo sshd -t untuk validasi sintaks sebelum me-restart layanan.",
			"Restart daemon SSH: sudo systemctl restart sshd",
		},
		Remediation: "Mencegah 99% serangan brute-force otomatis botnet yang menargetkan user root.",
	},
	{
		ID:          "LNX-02",
		Title:       "Konfigurasi Firewall UFW Standar Produksi",
		Category:    "Linux Hardening",
		Severity:    "HIGH",
		Description: "Menerapkan prinsip default-deny pada seluruh lalu lintas jaringan masuk dan hanya mengizinkan port layanan yang diperlukan.",
		Command:     "sudo ufw default deny incoming && sudo ufw default allow outgoing && sudo ufw allow 22/tcp && sudo ufw allow 80/tcp && sudo ufw allow 443/tcp && sudo ufw enable",
		Steps: []string{
			"Tetapkan aturan default tolak masuk: sudo ufw default deny incoming",
			"Tetapkan aturan default izinkan keluar: sudo ufw default allow outgoing",
			"Izinkan port SSH (sesuaikan jika custom port): sudo ufw allow 22/tcp",
			"Izinkan port web jika server hosting: sudo ufw allow 80/tcp && sudo ufw allow 443/tcp",
			"Aktifkan firewall: sudo ufw enable",
			"Verifikasi status: sudo ufw status verbose",
		},
		Remediation: "Membatasi port terbuka yang tidak sengaja diekspos oleh aplikasi background atau database.",
	},
	{
		ID:          "LNX-03",
		Title:       "Proteksi Brute-Force Otomatis dengan Fail2ban",
		Category:    "Linux Hardening",
		Severity:    "HIGH",
		Description: "Memasang sistem pemantauan log autentikasi untuk memblokir IP penyerang secara otomatis setelah beberapa kali gagal login.",
		Command:     "sudo apt-get update && sudo apt-get install -y fail2ban && sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local && sudo systemctl enable --now fail2ban",
		Steps: []string{
			"Pasang paket fail2ban pada sistem Debian/Ubuntu atau RHEL/Rocky.",
			"Salin template konfigurasi: sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local",
			"Konfigurasi maxretry = 3 dan bantime = 1h pada blok [sshd].",
			"Nyalakan layanan: sudo systemctl restart fail2ban",
			"Cek status jail yang aktif: sudo fail2ban-client status sshd",
		},
		Remediation: "Secara proaktif memutus serangan password guessing dan kamus credential stuffing.",
	},
	{
		ID:          "LNX-04",
		Title:       "Pengerasan Kernel Jaringan Linux (sysctl Security)",
		Category:    "Linux Hardening",
		Severity:    "MEDIUM",
		Description: "Mengaktifkan proteksi SYN Cookies dan menonaktifkan ICMP Redirect untuk mencegah serangan Denial of Service (DoS) dan Man-in-the-Middle.",
		Command:     "sudo tee -a /etc/sysctl.d/99-security.conf << 'EOF'\nnet.ipv4.tcp_syncookies = 1\nnet.ipv4.conf.all.accept_redirects = 0\nnet.ipv4.conf.all.send_redirects = 0\nnet.ipv4.conf.all.accept_source_route = 0\nnet.ipv4.conf.all.log_martians = 1\nEOF\nsudo sysctl -p /etc/sysctl.d/99-security.conf",
		Steps: []string{
			"Buat berkas konfigurasi baru /etc/sysctl.d/99-security.conf",
			"Aktifkan TCP SYN cookies untuk mitigasi SYN flood.",
			"Nonaktifkan IP source routing dan ICMP redirects.",
			"Muat ulang parameter kernel: sudo sysctl --system",
		},
		Remediation: "Melindungi network stack kernel Linux dari paket spoofing dan manipulasi rute.",
	},

	// 2. Windows Security
	{
		ID:          "WIN-01",
		Title:       "Rotasi Kata Sandi Administrator Lokal (LAPS)",
		Category:    "Windows Security",
		Severity:    "HIGH",
		Description: "Mencegah gerakan lateral (lateral movement) penyerang dengan memastikan setiap workstation memiliki kata sandi Administrator lokal yang unik dan dirotasi otomatis.",
		Command:     "Get-LapsDiagnostics; Set-LapsADComputerAccountEncryptionFlags -Identity $env:COMPUTERNAME -PasswordEncryptionRequired $true",
		Steps: []string{
			"Gunakan Windows LAPS bawaan (Windows 10/11 update terbaru) atau Azure AD LAPS.",
			"Aktifkan enkripsi password LAPS pada Active Directory atau Microsoft Entra ID.",
			"Tetapkan batas kedaluwarsa rotasi kata sandi maksimal 30 hari.",
			"Batasi izin pembacaan atribut password LAPS hanya untuk security admin.",
		},
		Remediation: "Mencegah teknik Pass-the-Hash menguasai seluruh armada PC hanya dari satu password admin yang sama.",
	},
	{
		ID:          "WIN-02",
		Title:       "Aturan Pengurangan Permukaan Serangan (ASR Rules)",
		Category:    "Windows Security",
		Severity:    "HIGH",
		Description: "Mengaktifkan Microsoft Defender Attack Surface Reduction untuk memblokir eksekusi skrip berbahaya dari aplikasi Office dan browser.",
		Command:     "Add-MpPreference -AttackSurfaceReductionRules_Ids 56a863a9-875e-4185-98a7-b882c60b3ce5 -AttackSurfaceReductionRules_Actions Enabled",
		Steps: []string{
			"Blokir aplikasi Office membuat proses anak (child process).",
			"Blokir eksekusi file executable yang terunduh mencurigakan.",
			"Blokir pemanggilan Win32 API dari macro dokumen Office.",
			"Aktifkan Script Block Logging pada PowerShell audit policy.",
		},
		Remediation: "Menahan infeksi malware tahap awal (initial access) dari lampiran email phishing.",
	},

	// 3. Incident Response
	{
		ID:          "IR-01",
		Title:       "SOP Isolasi Host & Preservasi Bukti Digital",
		Category:    "Incident Response",
		Severity:    "CRITICAL",
		Description: "Prosedur darurat saat komputer atau server terindikasi terkena malware, ransomware, atau akses tidak sah.",
		Command:     "sudo ip link set dev eth0 down || Disconnect-NetAdapter -Name * -Confirm:$false",
		Steps: []string{
			"JANGAN mematikan daya (Power Off) atau me-reboot mesin untuk menjaga bukti RAM yang mudah hilang (volatile memory).",
			"Putus koneksi jaringan fisik (cabut kabel LAN) atau matikan interface jaringan secara lokal.",
			"Catat waktu kejadian, gejala awal yang terlihat, dan akun pengguna yang sedang aktif.",
			"Jika memungkinkan, dump memori RAM menggunakan tool forensik (LiME untuk Linux atau WinPmem untuk Windows).",
			"Segera ubah kata sandi akun domain / admin dari komputer aman yang terpisah.",
		},
		Remediation: "Mencegah penyebaran infeksi malware ke segmen jaringan lain sambil menjaga integritas bukti digital.",
	},
	{
		ID:          "IR-02",
		Title:       "Investigasi Mekanisme Ketahanan Penyerang (Persistence Audit)",
		Category:    "Incident Response",
		Severity:    "HIGH",
		Description: "Menemukan backdoor atau skrip otomatis yang sengaja ditanam penyerang agar tetap bisa masuk setelah reboot.",
		Command:     "crontab -l; ls -la /etc/cron* /etc/systemd/system/; Get-ScheduledTask | Where-Object State -ne 'Disabled'",
		Steps: []string{
			"Periksa cronjobs seluruh user: /var/spool/cron/crontabs dan /etc/cron.*",
			"Periksa systemd service baru yang mencurigakan di /etc/systemd/system/ dan /lib/systemd/system/",
			"Periksa file profil shell: ~/.bashrc, ~/.profile, /etc/profile, /etc/bash.bashrc",
			"Pada Windows: Cek registry Run keys (HKLM/HKCU Software\\Microsoft\\Windows\\CurrentVersion\\Run) dan Task Scheduler.",
			"Cek SSH authorized_keys seluruh user di /home/*/.ssh/authorized_keys",
		},
		Remediation: "Memastikan penyerang benar-benar bersih dan tidak meninggalkan pintu belakang tersembunyi.",
	},

	// 4. OWASP Top 10
	{
		ID:          "OWASP-01",
		Title:       "Mitigasi Broken Access Control (A01)",
		Category:    "OWASP Top 10",
		Severity:    "CRITICAL",
		Description: "Memastikan verifikasi hak akses selalu dilakukan di sisi server (backend) pada setiap permintaan data.",
		Steps: []string{
			"Terapkan otorisasi berbasis peran (RBAC) atau atribut (ABAC) secara terpusat.",
			"Validasi kepemilikan objek (Ownership Check) pada setiap endpoint untuk mencegah IDOR (Insecure Direct Object Reference).",
			"Nonaktifkan directory listing pada web server.",
			"Catat dan pantau kegagalan otorisasi dalam log audit keamanan.",
		},
		Remediation: "Mencegah pengguna biasa membaca atau memodifikasi data milik pengguna lain dengan memanipulasi parameter ID.",
	},
	{
		ID:          "OWASP-02",
		Title:       "Mitigasi SQL Injection & Command Injection (A03)",
		Category:    "OWASP Top 10",
		Severity:    "CRITICAL",
		Description: "Mencegah eksekusi perintah SQL atau shell tak terduga dengan memisahkan kode dari data input.",
		Steps: []string{
			"Gunakan selalu Prepared Statements (Parameterized Queries) pada semua query basis data.",
			"Hindari konkatenasi string SQL mentah (misal: 'SELECT * FROM users WHERE name = ' + input).",
			"Terapkan validasi input ketat (allowlist) pada format tipe data yang diterima.",
			"Hindari memanggil shell sistem (exec, system) langsung dengan input pengguna tanpa sanitasi.",
		},
		Remediation: "Menghilangkan celah fatal pembajakan database dan pengambilalihan server dari input formulir web.",
	},
}

// GetAllPlaybooks returns all security playbook items
func GetAllPlaybooks() []SecurityPlaybookItem {
	return playbookCatalog
}

// GetPlaybookCategories returns unique category names
func GetPlaybookCategories() []string {
	return []string{
		"Semua",
		"Linux Hardening",
		"Windows Security",
		"Incident Response",
		"OWASP Top 10",
	}
}

// FilterPlaybooks searches playbooks by keyword and category
func FilterPlaybooks(query string, category string) []SecurityPlaybookItem {
	q := strings.TrimSpace(strings.ToLower(query))
	cat := strings.TrimSpace(strings.ToLower(category))

	var results []SecurityPlaybookItem
	for _, item := range playbookCatalog {
		if cat != "" && cat != "semua" && strings.ToLower(item.Category) != cat {
			continue
		}

		if q == "" {
			results = append(results, item)
			continue
		}

		if strings.Contains(strings.ToLower(item.Title), q) ||
			strings.Contains(strings.ToLower(item.ID), q) ||
			strings.Contains(strings.ToLower(item.Description), q) ||
			strings.Contains(strings.ToLower(item.Remediation), q) ||
			strings.Contains(strings.ToLower(item.Command), q) {
			results = append(results, item)
		}
	}

	return results
}
