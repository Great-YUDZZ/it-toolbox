package security

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PortResult holds the scan outcome for a single TCP port
type PortResult struct {
	Port    int    `json:"port"`
	State   string `json:"state"` // "OPEN", "CLOSED"
	Service string `json:"service"`
	Latency int64  `json:"latency_ms"`
	Banner  string `json:"banner,omitempty"`
}

// ScanProgress communicates scanning state back to caller
type ScanProgress struct {
	Scanned   int
	Total     int
	OpenCount int
	Current   int
}

var CommonPorts = map[int]string{
	20:    "FTP Data",
	21:    "FTP Control",
	22:    "SSH (Secure Shell)",
	23:    "Telnet",
	25:    "SMTP (Mail)",
	53:    "DNS (Domain Name System)",
	67:    "DHCP Server",
	68:    "DHCP Client",
	69:    "TFTP",
	80:    "HTTP (Web)",
	88:    "Kerberos",
	110:   "POP3 (Mail)",
	111:   "RPCbind",
	123:   "NTP (Network Time)",
	135:   "Microsoft EPMAP / RPC",
	137:   "NetBIOS Name Service",
	138:   "NetBIOS Datagram",
	139:   "NetBIOS Session Service",
	143:   "IMAP (Mail)",
	161:   "SNMP",
	389:   "LDAP",
	443:   "HTTPS (Secure Web)",
	445:   "Microsoft SMB / CIFS",
	465:   "SMTPS (Secure Mail)",
	514:   "Syslog",
	587:   "SMTP Submission",
	636:   "LDAPS (Secure LDAP)",
	873:   "Rsync",
	993:   "IMAPS (Secure Mail)",
	995:   "POP3S (Secure Mail)",
	1433:  "Microsoft SQL Server",
	1521:  "Oracle Database",
	2049:  "NFS (Network File System)",
	2375:  "Docker REST API (Plain)",
	2376:  "Docker REST API (TLS)",
	3306:  "MySQL / MariaDB",
	3389:  "RDP (Remote Desktop)",
	5432:  "PostgreSQL",
	5900:  "VNC Remote Desktop",
	6379:  "Redis In-Memory Database",
	8000:  "HTTP Alternate / Dev",
	8080:  "HTTP Proxy / Tomcat",
	8443:  "HTTPS Alternate",
	8888:  "HTTP Alternate",
	9000:  "SonarQube / PHP-FPM",
	9090:  "Prometheus Web UI",
	9200:  "Elasticsearch REST",
	27017: "MongoDB Database",
}

// Top20Ports provides the most frequently probed TCP ports
var Top20Ports = []int{
	21, 22, 23, 25, 53, 80, 110, 135, 139, 143,
	443, 445, 1433, 1521, 3306, 3389, 5432, 6379, 8080, 8443,
}

// GetServiceName returns common service name or unknown placeholder
func GetServiceName(port int) string {
	if name, ok := CommonPorts[port]; ok {
		return name
	}
	return "Layanan Kustom / Tidak Dikenal"
}

// ParsePortRange decodes comma and hyphen expressions into distinct sorted ports
func ParsePortRange(rangeStr string) ([]int, error) {
	clean := strings.TrimSpace(rangeStr)
	if clean == "" || strings.EqualFold(clean, "top20") || strings.EqualFold(clean, "top 20") {
		return Top20Ports, nil
	}

	seen := make(map[int]bool)
	var ports []int

	parts := strings.Split(clean, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("format range port tidak valid: %q", part)
			}
			start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err1 != nil || err2 != nil || start <= 0 || end <= 0 || start > 65535 || end > 65535 || start > end {
				return nil, fmt.Errorf("range port tidak valid (harus 1-65535): %q", part)
			}
			if end-start > 2000 {
				return nil, fmt.Errorf("maksimum 2000 port dalam satu sesi scan untuk stabilitas (dipilih: %d)", end-start+1)
			}
			for p := start; p <= end; p++ {
				if !seen[p] {
					seen[p] = true
					ports = append(ports, p)
				}
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil || p <= 0 || p > 65535 {
				return nil, fmt.Errorf("nomor port tidak valid: %q", part)
			}
			if !seen[p] {
				seen[p] = true
				ports = append(ports, p)
			}
		}
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("tidak ada port yang valid untuk dipindai")
	}

	sort.Ints(ports)
	return ports, nil
}

// ScanTarget runs a concurrent TCP port scan against the target host
func ScanTarget(ctx context.Context, host string, ports []int, timeout time.Duration, concurrency int, progressFn func(p ScanProgress)) ([]PortResult, error) {
	if host == "" {
		return nil, fmt.Errorf("target host tidak boleh kosong")
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("daftar port tidak boleh kosong")
	}
	if timeout <= 0 {
		timeout = 350 * time.Millisecond
	}
	if concurrency <= 0 {
		concurrency = 50
	}
	if concurrency > 200 {
		concurrency = 200
	}

	total := len(ports)
	var results []PortResult
	var resMu sync.Mutex

	scanned := 0
	openCount := 0
	var progMu sync.Mutex

	jobs := make(chan int, total)
	for _, p := range ports {
		jobs <- p
	}
	close(jobs)

	var wg sync.WaitGroup
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}

				addr := net.JoinHostPort(host, strconv.Itoa(port))
				start := time.Now()
				conn, err := net.DialTimeout("tcp", addr, timeout)
				latency := time.Since(start).Milliseconds()

				var res PortResult
				if err == nil {
					// Connected - port is open
					banner := ""
					_ = conn.SetReadDeadline(time.Now().Add(250 * time.Millisecond))
					buf := make([]byte, 256)
					n, _ := conn.Read(buf)
					if n > 0 {
						cleanBanner := strings.TrimSpace(string(buf[:n]))
						// Remove control chars
						cleanBanner = strings.ReplaceAll(cleanBanner, "\r", " ")
						cleanBanner = strings.ReplaceAll(cleanBanner, "\n", " ")
						if len(cleanBanner) > 60 {
							cleanBanner = cleanBanner[:57] + "..."
						}
						banner = cleanBanner
					}
					_ = conn.Close()

					res = PortResult{
						Port:    port,
						State:   "OPEN",
						Service: GetServiceName(port),
						Latency: latency,
						Banner:  banner,
					}
				} else {
					res = PortResult{
						Port:    port,
						State:   "CLOSED",
						Service: GetServiceName(port),
						Latency: latency,
					}
				}

				resMu.Lock()
				results = append(results, res)
				resMu.Unlock()

				progMu.Lock()
				scanned++
				if res.State == "OPEN" {
					openCount++
				}
				curScanned := scanned
				curOpen := openCount
				progMu.Unlock()

				if progressFn != nil {
					progressFn(ScanProgress{
						Scanned:   curScanned,
						Total:     total,
						OpenCount: curOpen,
						Current:   port,
					})
				}
			}
		}()
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results, nil
}

// GenerateScanReport produces a text format summary of port scan results
func GenerateScanReport(host string, results []PortResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== LAPORAN AUDIT PORT: %s ===\n", host))
	sb.WriteString(fmt.Sprintf("Waktu: %s\n", time.Now().Format("02 Jan 2006 15:04:05 WIB")))
	sb.WriteString(fmt.Sprintf("Total Port Diperiksa: %d\n\n", len(results)))

	openCount := 0
	for _, r := range results {
		if r.State == "OPEN" {
			openCount++
			sb.WriteString(fmt.Sprintf("[OPEN] Port %-5d (%s) - Latency: %d ms", r.Port, r.Service, r.Latency))
			if r.Banner != "" {
				sb.WriteString(fmt.Sprintf(" | Banner: %s", r.Banner))
			}
			sb.WriteString("\n")
		}
	}

	if openCount == 0 {
		sb.WriteString("Tidak ditemukan port terbuka dari rentang yang dipindai.\n")
	} else {
		sb.WriteString(fmt.Sprintf("\nRingkasan: %d port terbuka terdeteksi.\n", openCount))
	}

	return sb.String()
}
