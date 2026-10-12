package security

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// FileHashResult holds cryptographic checksums and file type integrity inspection
type FileHashResult struct {
	FilePath   string `json:"file_path"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	MD5        string `json:"md5"`
	SHA1       string `json:"sha1"`
	SHA256     string `json:"sha256"`
	SHA512     string `json:"sha512"`
	MagicType  string `json:"magic_type"`
	MagicExt   string `json:"magic_ext"`
	Suspicious bool   `json:"suspicious"`
	Notes      string `json:"notes,omitempty"`
}

// AuthAttempt tracks failed authentication counts by IP or user
type AuthAttempt struct {
	Key      string `json:"key"`
	Count    int    `json:"count"`
	Category string `json:"category"`
}

// AuthLogAnalysis summarizes Linux/SSH authentication security telemetry
type AuthLogAnalysis struct {
	TotalLines     int           `json:"total_lines"`
	FailedAttempts int           `json:"failed_attempts"`
	SuccessLogins  int           `json:"success_logins"`
	InvalidUsers   int           `json:"invalid_users"`
	TopAttackerIPs []AuthAttempt `json:"top_attacker_ips"`
	TopTargetUsers []AuthAttempt `json:"top_target_users"`
	AttackerIPList []string      `json:"attacker_ip_list"`
}

// WebLogAnalysis summarizes web server access log scanning telemetry
type WebLogAnalysis struct {
	TotalRequests  int           `json:"total_requests"`
	Http404Scans   int           `json:"http_404_scans"`
	PathTraversals int           `json:"path_traversals"`
	SqlInjections  int           `json:"sql_injections"`
	XssAttempts    int           `json:"xss_attempts"`
	TopSuspectIPs  []AuthAttempt `json:"top_suspect_ips"`
	SuspiciousUris []string      `json:"suspicious_uris"`
	SuspectIPList  []string      `json:"suspect_ip_list"`
}

// VTLookupResult holds VirusTotal file hash reputation summary
type VTLookupResult struct {
	SHA256        string `json:"sha256"`
	Malicious     int    `json:"malicious"`
	Suspicious    int    `json:"suspicious"`
	Undetected    int    `json:"undetected"`
	Harmless      int    `json:"harmless"`
	Reputation    int    `json:"reputation"`
	ThreatLabel   string `json:"threat_label"`
	StatusSummary string `json:"status_summary"`
}

var ipRegex = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

// ComputeFileHashes generates multi-hash digests and inspects magic byte headers
func ComputeFileHashes(filePath string) (*FileHashResult, error) {
	cleanPath := strings.TrimSpace(filePath)
	if cleanPath == "" {
		return nil, fmt.Errorf("jalur berkas tidak boleh kosong")
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("berkas tidak ditemukan: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("jalur adalah direktori, silakan pilih berkas")
	}

	f, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka berkas: %w", err)
	}
	defer f.Close()

	// Multi-hasher setup
	hMD5 := md5.New()
	hSHA1 := sha1.New()
	hSHA256 := sha256.New()
	hSHA512 := sha512.New()
	multiWriter := io.MultiWriter(hMD5, hSHA1, hSHA256, hSHA512)

	// Read magic bytes header first
	header := make([]byte, 512)
	nHeader, _ := f.Read(header)
	if nHeader > 0 {
		_, _ = multiWriter.Write(header[:nHeader])
	}

	// Copy rest of file
	buf := make([]byte, 64*1024)
	for {
		n, rErr := f.Read(buf)
		if n > 0 {
			_, _ = multiWriter.Write(buf[:n])
		}
		if rErr == io.EOF {
			break
		}
		if rErr != nil {
			return nil, fmt.Errorf("gagal membaca stream berkas: %w", rErr)
		}
	}

	magicType, magicExt := detectMagic(header[:nHeader])
	ext := strings.ToLower(filepath.Ext(cleanPath))

	isSuspicious := false
	var noteStr string

	// Check for extension disguise
	if (magicExt == ".exe" || magicExt == ".elf") && (ext == ".jpg" || ext == ".png" || ext == ".pdf" || ext == ".txt") {
		isSuspicious = true
		noteStr = fmt.Sprintf("PERINGATAN: Berkas memiliki ekstensi %s tetapi header biner adalah %s (%s). Kemungkinan payload berbahaya tersamar!", ext, magicType, magicExt)
	} else if magicType != "Tidak Dikenal / Teks Biasa" {
		noteStr = fmt.Sprintf("Header biner terverifikasi sesuai tipe: %s", magicType)
	}

	return &FileHashResult{
		FilePath:   cleanPath,
		FileName:   filepath.Base(cleanPath),
		FileSize:   info.Size(),
		MD5:        hex.EncodeToString(hMD5.Sum(nil)),
		SHA1:       hex.EncodeToString(hSHA1.Sum(nil)),
		SHA256:     hex.EncodeToString(hSHA256.Sum(nil)),
		SHA512:     hex.EncodeToString(hSHA512.Sum(nil)),
		MagicType:  magicType,
		MagicExt:   magicExt,
		Suspicious: isSuspicious,
		Notes:      noteStr,
	}, nil
}

func detectMagic(header []byte) (string, string) {
	if len(header) < 4 {
		return "Berkas Terlalu Kecil / Kosong", ""
	}

	switch {
	case bytes.HasPrefix(header, []byte("MZ")):
		return "Windows PE Executable / DLL", ".exe"
	case bytes.HasPrefix(header, []byte{0x7F, 'E', 'L', 'F'}):
		return "Linux ELF Executable / Binary", ".elf"
	case bytes.HasPrefix(header, []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}):
		return "PNG Image", ".png"
	case bytes.HasPrefix(header, []byte{0xFF, 0xD8, 0xFF}):
		return "JPEG / JPG Image", ".jpg"
	case bytes.HasPrefix(header, []byte("%PDF")):
		return "PDF Document", ".pdf"
	case bytes.HasPrefix(header, []byte("PK\x03\x04")):
		return "ZIP Archive / OpenOffice / APK", ".zip"
	case bytes.HasPrefix(header, []byte{0x1F, 0x8B}):
		return "GZIP Compressed Archive", ".gz"
	case bytes.HasPrefix(header, []byte("GIF87a")) || bytes.HasPrefix(header, []byte("GIF89a")):
		return "GIF Image", ".gif"
	case bytes.HasPrefix(header, []byte("Rar!\x1A\x07")):
		return "RAR Archive", ".rar"
	case bytes.HasPrefix(header, []byte{0x25, 0x21, 0x50, 0x53}):
		return "PostScript Document", ".ps"
	case bytes.HasPrefix(header, []byte("7z\xBC\xAF\x27\x1C")):
		return "7-Zip Archive", ".7z"
	default:
		return "Tidak Dikenal / Teks Biasa", ""
	}
}

// ParseAuthLog analyzes Linux SSH/auth.log lines for authentication brute force
func ParseAuthLog(logText string) (*AuthLogAnalysis, error) {
	lines := strings.Split(logText, "\n")
	analysis := &AuthLogAnalysis{
		TotalLines: len(lines),
	}

	ipCounts := make(map[string]int)
	userCounts := make(map[string]int)

	userRegex := regexp.MustCompile(`for (?:invalid user )?([a-zA-Z0-9_\-\.]+) from`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)
		isFailed := strings.Contains(lower, "failed password") || strings.Contains(lower, "authentication failure")
		isInvalid := strings.Contains(lower, "invalid user")
		isSuccess := strings.Contains(lower, "accepted password") || strings.Contains(lower, "accepted publickey")

		if isFailed {
			analysis.FailedAttempts++
		}
		if isInvalid {
			analysis.InvalidUsers++
		}
		if isSuccess {
			analysis.SuccessLogins++
		}

		if isFailed || isInvalid {
			// Extract IP
			ipMatch := ipRegex.FindString(line)
			if ipMatch != "" && ipMatch != "127.0.0.1" && ipMatch != "0.0.0.0" {
				ipCounts[ipMatch]++
			}

			// Extract User
			uMatches := userRegex.FindStringSubmatch(line)
			if len(uMatches) > 1 {
				userCounts[uMatches[1]]++
			}
		}
	}

	// Sort attacker IPs
	for ip, count := range ipCounts {
		analysis.TopAttackerIPs = append(analysis.TopAttackerIPs, AuthAttempt{
			Key:      ip,
			Count:    count,
			Category: "IP Penyerang",
		})
		analysis.AttackerIPList = append(analysis.AttackerIPList, ip)
	}
	sort.Slice(analysis.TopAttackerIPs, func(i, j int) bool {
		return analysis.TopAttackerIPs[i].Count > analysis.TopAttackerIPs[j].Count
	})

	// Sort targeted users
	for user, count := range userCounts {
		analysis.TopTargetUsers = append(analysis.TopTargetUsers, AuthAttempt{
			Key:      user,
			Count:    count,
			Category: "Username Target",
		})
	}
	sort.Slice(analysis.TopTargetUsers, func(i, j int) bool {
		return analysis.TopTargetUsers[i].Count > analysis.TopTargetUsers[j].Count
	})

	return analysis, nil
}

// ParseWebAccessLog analyzes HTTP server access logs for common attack patterns
func ParseWebAccessLog(logText string) (*WebLogAnalysis, error) {
	lines := strings.Split(logText, "\n")
	analysis := &WebLogAnalysis{
		TotalRequests: len(lines),
	}

	ipCounts := make(map[string]int)
	seenUris := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lower := strings.ToLower(line)
		is404 := strings.Contains(line, " 404 ")
		isTraversal := strings.Contains(lower, "../") || strings.Contains(lower, "..\\") || strings.Contains(lower, "%2e%2e")
		isSQLi := strings.Contains(lower, "union select") || strings.Contains(lower, "union+select") || strings.Contains(lower, "' or 1=1") || strings.Contains(lower, "information_schema")
		isXSS := strings.Contains(lower, "<script") || strings.Contains(lower, "%3cscript") || strings.Contains(lower, "javascript:")

		if is404 {
			analysis.Http404Scans++
		}
		if isTraversal {
			analysis.PathTraversals++
		}
		if isSQLi {
			analysis.SqlInjections++
		}
		if isXSS {
			analysis.XssAttempts++
		}

		if is404 || isTraversal || isSQLi || isXSS {
			ipMatch := ipRegex.FindString(line)
			if ipMatch != "" && ipMatch != "127.0.0.1" {
				ipCounts[ipMatch]++
			}

			// Extract request URI snippet
			if idx1 := strings.Index(line, "\""); idx1 != -1 {
				rest := line[idx1+1:]
				if idx2 := strings.Index(rest, "\""); idx2 != -1 {
					reqLine := rest[:idx2]
					parts := strings.Fields(reqLine)
					if len(parts) >= 2 {
						uri := parts[1]
						if !seenUris[uri] && len(analysis.SuspiciousUris) < 20 {
							seenUris[uri] = true
							analysis.SuspiciousUris = append(analysis.SuspiciousUris, uri)
						}
					}
				}
			}
		}
	}

	for ip, count := range ipCounts {
		analysis.TopSuspectIPs = append(analysis.TopSuspectIPs, AuthAttempt{
			Key:      ip,
			Count:    count,
			Category: "IP Scanner / Penyerang Web",
		})
		analysis.SuspectIPList = append(analysis.SuspectIPList, ip)
	}
	sort.Slice(analysis.TopSuspectIPs, func(i, j int) bool {
		return analysis.TopSuspectIPs[i].Count > analysis.TopSuspectIPs[j].Count
	})

	return analysis, nil
}

// GenerateBlockRules produces firewall rules to block suspicious IPs
func GenerateBlockRules(ips []string, tool string) string {
	if len(ips) == 0 {
		return "# Tidak ada IP untuk diblokir."
	}

	var sb strings.Builder
	tool = strings.ToLower(strings.TrimSpace(tool))

	switch tool {
	case "iptables":
		sb.WriteString("# Perintah Pemblokiran Linux iptables:\n")
		for _, ip := range ips {
			sb.WriteString(fmt.Sprintf("sudo iptables -A INPUT -s %s -j DROP\n", ip))
		}
	case "firewalld":
		sb.WriteString("# Perintah Pemblokiran Linux firewalld:\n")
		for _, ip := range ips {
			sb.WriteString(fmt.Sprintf("sudo firewall-cmd --add-rich-rule='rule family=\"ipv4\" source address=\"%s\" reject' --permanent\n", ip))
		}
		sb.WriteString("sudo firewall-cmd --reload\n")
	case "powershell", "windows":
		sb.WriteString("# Perintah Pemblokiran Windows Firewall (PowerShell Run as Administrator):\n")
		for i, ip := range ips {
			sb.WriteString(fmt.Sprintf("New-NetFirewallRule -Name \"Block_Attacker_%d\" -DisplayName \"Block IP %s\" -Direction Inbound -Action Block -RemoteAddress %s\n", i+1, ip, ip))
		}
	default: // UFW
		sb.WriteString("# Perintah Pemblokiran Linux UFW (Uncomplicated Firewall):\n")
		for _, ip := range ips {
			sb.WriteString(fmt.Sprintf("sudo ufw deny from %s to any\n", ip))
		}
	}

	return sb.String()
}

// LookupVirusTotalHash queries VirusTotal API v3 for known hash reputation
func LookupVirusTotalHash(apiKey string, sha256Hash string) (*VTLookupResult, error) {
	key := strings.TrimSpace(apiKey)
	if key == "" {
		return nil, fmt.Errorf("API key VirusTotal diperlukan untuk query online")
	}
	hash := strings.TrimSpace(sha256Hash)
	if len(hash) != 64 {
		return nil, fmt.Errorf("hash harus berupa SHA-256 (64 karakter heksadesimal)")
	}

	apiURL := fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s", hash)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apikey", key)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi VirusTotal: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return &VTLookupResult{
			SHA256:        hash,
			StatusSummary: "Berkas Tidak Ditemukan: Hash belum pernah diunggah atau dianalisis di VirusTotal.",
		}, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("VirusTotal mengembalikan status HTTP %d", resp.StatusCode)
	}

	var payload struct {
		Data struct {
			Attributes struct {
				LastAnalysisStats struct {
					Malicious  int `json:"malicious"`
					Suspicious int `json:"suspicious"`
					Undetected int `json:"undetected"`
					Harmless   int `json:"harmless"`
				} `json:"last_analysis_stats"`
				Reputation                int    `json:"reputation"`
				PopularThreatClassification struct {
					SuggestedThreatLabel string `json:"suggested_threat_label"`
				} `json:"popular_threat_classification"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("gagal membaca respons JSON VirusTotal: %w", err)
	}

	stats := payload.Data.Attributes.LastAnalysisStats
	label := payload.Data.Attributes.PopularThreatClassification.SuggestedThreatLabel
	if label == "" {
		label = "Tidak Tersedia"
	}

	summary := fmt.Sprintf("Analisis VT: %d engine mendeteksi Berbahaya (Malicious), %d Mencurigakan, %d Bersih.", stats.Malicious, stats.Suspicious, stats.Harmless)

	return &VTLookupResult{
		SHA256:        hash,
		Malicious:     stats.Malicious,
		Suspicious:    stats.Suspicious,
		Undetected:    stats.Undetected,
		Harmless:      stats.Harmless,
		Reputation:    payload.Data.Attributes.Reputation,
		ThreatLabel:   label,
		StatusSummary: summary,
	}, nil
}
