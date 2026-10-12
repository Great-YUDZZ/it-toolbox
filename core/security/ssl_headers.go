package security

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TLSInspectionResult contains the cryptographic certificate inspection details
type TLSInspectionResult struct {
	Domain       string    `json:"domain"`
	Port         int       `json:"port"`
	SubjectCN    string    `json:"subject_cn"`
	Issuer       string    `json:"issuer"`
	SANs         []string  `json:"sans"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidUntil   time.Time `json:"valid_until"`
	DaysLeft     int       `json:"days_left"`
	IsExpired    bool      `json:"is_expired"`
	ExpiringSoon bool      `json:"expiring_soon"`
	TLSVersion   string    `json:"tls_version"`
	CipherSuite  string    `json:"cipher_suite"`
	InsecureTLS  bool      `json:"insecure_tls"`
}

// HeaderCheckItem detail for a single HTTP security header
type HeaderCheckItem struct {
	HeaderName     string `json:"header_name"`
	Present        bool   `json:"present"`
	Value          string `json:"value"`
	Recommendation string `json:"recommendation"`
	Severity       string `json:"severity"` // "PASS", "WARN", "FAIL"
}

// HeaderAuditResult aggregates security header checks and compliance grade
type HeaderAuditResult struct {
	URL                 string            `json:"url"`
	StatusCode          int               `json:"status_code"`
	Grade               string            `json:"grade"` // "A+", "A", "B", "C", "D", "F"
	Score               int               `json:"score"` // 0 - 100
	Headers             []HeaderCheckItem `json:"headers"`
	Server              string            `json:"server"`
	RawHeaders          map[string]string `json:"raw_headers"`
	RemediationSnippets string            `json:"remediation_snippets"`
}

// EmailSecurityResult records SPF, DMARC, and MX DNS security hygiene
type EmailSecurityResult struct {
	Domain        string   `json:"domain"`
	HasSPF        bool     `json:"has_spf"`
	SPFRecord     string   `json:"spf_record"`
	HasDMARC      bool     `json:"has_dmarc"`
	DMARCRecord   string   `json:"dmarc_record"`
	MXRecords     []string `json:"mx_records"`
	StatusSummary string   `json:"status_summary"`
}

func CleanHost(input string) string {
	clean := strings.TrimSpace(input)
	clean = strings.TrimPrefix(clean, "https://")
	clean = strings.TrimPrefix(clean, "http://")
	if idx := strings.Index(clean, "/"); idx != -1 {
		clean = clean[:idx]
	}
	if host, _, err := net.SplitHostPort(clean); err == nil {
		clean = host
	}
	return clean
}

func tlsVersionString(ver uint16) (string, bool) {
	switch ver {
	case tls.VersionTLS13:
		return "TLS 1.3 (Modern Secure)", false
	case tls.VersionTLS12:
		return "TLS 1.2 (Standard Secure)", false
	case tls.VersionTLS11:
		return "TLS 1.1 (Deprecated / Insecure)", true
	case tls.VersionTLS10:
		return "TLS 1.0 (Deprecated / Insecure)", true
	default:
		return fmt.Sprintf("Unknown (0x%04x)", ver), true
	}
}

// InspectTLS connects to domain via TLS and evaluates certificate health
func InspectTLS(domain string, port int, timeout time.Duration) (*TLSInspectionResult, error) {
	host := CleanHost(domain)
	if host == "" {
		return nil, fmt.Errorf("domain atau host tidak boleh kosong")
	}
	if port <= 0 {
		port = 443
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	dialer := &net.Dialer{Timeout: timeout}

	// First attempt standard validation
	conf := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", addr, conf)
	if err != nil {
		// Attempt fallback with InsecureSkipVerify to still fetch certificate details if self-signed or expired
		insecureConf := &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: true,
		}
		insecureConn, err2 := tls.DialWithDialer(dialer, "tcp", addr, insecureConf)
		if err2 != nil {
			return nil, fmt.Errorf("gagal terhubung ke %s via TLS: %w", addr, err)
		}
		conn = insecureConn
	}
	defer conn.Close()

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return nil, fmt.Errorf("tidak ada sertifikat publik yang dikembalikan oleh host %s", host)
	}

	cert := state.PeerCertificates[0]
	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)

	tlsVer, isInsecure := tlsVersionString(state.Version)
	cipherName := tls.CipherSuiteName(state.CipherSuite)

	issuerStr := cert.Issuer.CommonName
	if issuerStr == "" && len(cert.Issuer.Organization) > 0 {
		issuerStr = cert.Issuer.Organization[0]
	}
	if issuerStr == "" {
		issuerStr = cert.Issuer.String()
	}

	subjectCN := cert.Subject.CommonName
	if subjectCN == "" {
		subjectCN = host
	}

	return &TLSInspectionResult{
		Domain:       host,
		Port:         port,
		SubjectCN:    subjectCN,
		Issuer:       issuerStr,
		SANs:         cert.DNSNames,
		ValidFrom:    cert.NotBefore,
		ValidUntil:   cert.NotAfter,
		DaysLeft:     daysLeft,
		IsExpired:    now.After(cert.NotAfter),
		ExpiringSoon: daysLeft >= 0 && daysLeft <= 30,
		TLSVersion:   tlsVer,
		CipherSuite:  cipherName,
		InsecureTLS:  isInsecure,
	}, nil
}

// InspectSecurityHeaders queries HTTP endpoints and grades security headers
func InspectSecurityHeaders(targetURL string, timeout time.Duration) (*HeaderAuditResult, error) {
	rawURL := strings.TrimSpace(targetURL)
	if rawURL == "" {
		return nil, fmt.Errorf("URL tidak boleh kosong")
	}
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("format URL tidak valid: %w", err)
	}

	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", parsed.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "IT-Toolbox-SecurityAuditor/1.5")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan request ke %s: %w", parsed.String(), err)
	}
	defer resp.Body.Close()

	rawMap := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			rawMap[k] = v[0]
		}
	}

	var checks []HeaderCheckItem
	score := 0

	// 1. Strict-Transport-Security (20 pts)
	hsts := resp.Header.Get("Strict-Transport-Security")
	if hsts != "" {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "Strict-Transport-Security (HSTS)",
			Present:        true,
			Value:          hsts,
			Recommendation: "Sudah aktif dan melindungi koneksi HTTPS dari downgrade attacks.",
			Severity:       "PASS",
		})
		score += 20
	} else {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "Strict-Transport-Security (HSTS)",
			Present:        false,
			Value:          "-",
			Recommendation: "Tambahkan: Strict-Transport-Security: max-age=31536000; includeSubDomains",
			Severity:       "FAIL",
		})
	}

	// 2. Content-Security-Policy (25 pts)
	csp := resp.Header.Get("Content-Security-Policy")
	if csp != "" {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "Content-Security-Policy (CSP)",
			Present:        true,
			Value:          csp,
			Recommendation: "CSP aktif, membatasi pemuatan skrip eksternal dan serangan XSS.",
			Severity:       "PASS",
		})
		score += 25
	} else {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "Content-Security-Policy (CSP)",
			Present:        false,
			Value:          "-",
			Recommendation: "Definisikan aturan CSP minimal: default-src 'self'",
			Severity:       "FAIL",
		})
	}

	// 3. X-Frame-Options (15 pts)
	xfo := resp.Header.Get("X-Frame-Options")
	if xfo != "" {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "X-Frame-Options (Clickjacking Protection)",
			Present:        true,
			Value:          xfo,
			Recommendation: "Mencegah situs disematkan ke dalam iframe pihak ketiga.",
			Severity:       "PASS",
		})
		score += 15
	} else {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "X-Frame-Options",
			Present:        false,
			Value:          "-",
			Recommendation: "Tambahkan: X-Frame-Options: SAMEORIGIN atau DENY",
			Severity:       "FAIL",
		})
	}

	// 4. X-Content-Type-Options (15 pts)
	xcto := resp.Header.Get("X-Content-Type-Options")
	if strings.EqualFold(xcto, "nosniff") {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "X-Content-Type-Options",
			Present:        true,
			Value:          xcto,
			Recommendation: "Mencegah browser melakukan MIME-sniffing berkas berbahaya.",
			Severity:       "PASS",
		})
		score += 15
	} else {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "X-Content-Type-Options",
			Present:        false,
			Value:          xcto,
			Recommendation: "Tambahkan: X-Content-Type-Options: nosniff",
			Severity:       "FAIL",
		})
	}

	// 5. Referrer-Policy (15 pts)
	refPol := resp.Header.Get("Referrer-Policy")
	if refPol != "" {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "Referrer-Policy",
			Present:        true,
			Value:          refPol,
			Recommendation: "Mencegah kebocoran URL sensitif pada request outbound.",
			Severity:       "PASS",
		})
		score += 15
	} else {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "Referrer-Policy",
			Present:        false,
			Value:          "-",
			Recommendation: "Tambahkan: Referrer-Policy: strict-origin-when-cross-origin",
			Severity:       "WARN",
		})
	}

	// 6. Permissions-Policy (10 pts)
	permPol := resp.Header.Get("Permissions-Policy")
	if permPol != "" {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "Permissions-Policy",
			Present:        true,
			Value:          permPol,
			Recommendation: "Membatasi akses API browser (kamera, mikrofon, geolokasi).",
			Severity:       "PASS",
		})
		score += 10
	} else {
		checks = append(checks, HeaderCheckItem{
			HeaderName:     "Permissions-Policy",
			Present:        false,
			Value:          "-",
			Recommendation: "Tambahkan: Permissions-Policy: camera=(), microphone=(), geolocation=()",
			Severity:       "WARN",
		})
	}

	// Calculate Grade
	var grade string
	switch {
	case score >= 95:
		grade = "A+"
	case score >= 80:
		grade = "A"
	case score >= 65:
		grade = "B"
	case score >= 50:
		grade = "C"
	case score >= 30:
		grade = "D"
	default:
		grade = "F"
	}

	serverBanner := resp.Header.Get("Server")
	if serverBanner == "" {
		serverBanner = "Tersembunyi / Tidak Ditampilkan"
	}

	remediation := `### Rekomendasi Konfigurasi Nginx:
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Permissions-Policy "camera=(), microphone=(), geolocation=()" always;
add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; object-src 'none';" always;

### Rekomendasi Konfigurasi Apache (.htaccess):
<IfModule mod_headers.c>
  Header always set Strict-Transport-Security "max-age=31536000; includeSubDomains"
  Header always set X-Frame-Options "SAMEORIGIN"
  Header always set X-Content-Type-Options "nosniff"
  Header always set Referrer-Policy "strict-origin-when-cross-origin"
  Header always set Permissions-Policy "camera=(), microphone=(), geolocation=()"
</IfModule>`

	return &HeaderAuditResult{
		URL:                 parsed.String(),
		StatusCode:          resp.StatusCode,
		Grade:               grade,
		Score:               score,
		Headers:             checks,
		Server:              serverBanner,
		RawHeaders:          rawMap,
		RemediationSnippets: remediation,
	}, nil
}

// CheckEmailSecurityDNS performs DNS queries to analyze SPF, DMARC, and MX with timeout
func CheckEmailSecurityDNS(domain string) (*EmailSecurityResult, error) {
	clean := CleanHost(domain)
	if clean == "" {
		return nil, fmt.Errorf("nama domain tidak boleh kosong")
	}

	res := &EmailSecurityResult{
		Domain: clean,
	}

	resolver := &net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 1. Check MX records
	mxList, err := resolver.LookupMX(ctx, clean)
	if err == nil && len(mxList) > 0 {
		for _, mx := range mxList {
			res.MXRecords = append(res.MXRecords, fmt.Sprintf("%s (Pref: %d)", mx.Host, mx.Pref))
		}
	}

	// 2. Check SPF
	txtRecords, err := resolver.LookupTXT(ctx, clean)
	if err == nil {
		for _, txt := range txtRecords {
			if strings.HasPrefix(strings.ToLower(txt), "v=spf1") {
				res.HasSPF = true
				res.SPFRecord = txt
				break
			}
		}
	}

	// 3. Check DMARC
	dmarcDomain := "_dmarc." + clean
	dmarcTxt, err := resolver.LookupTXT(ctx, dmarcDomain)
	if err == nil {
		for _, txt := range dmarcTxt {
			if strings.HasPrefix(strings.ToLower(txt), "v=dmarc1") {
				res.HasDMARC = true
				res.DMARCRecord = txt
				break
			}
		}
	}

	// Summary
	if res.HasSPF && res.HasDMARC {
		res.StatusSummary = "Proteksi Email Lengkap: SPF dan DMARC terkonfigurasi dengan baik."
	} else if res.HasSPF && !res.HasDMARC {
		res.StatusSummary = "Proteksi Parsial: SPF aktif, namun DMARC belum dipasang."
	} else if !res.HasSPF && res.HasDMARC {
		res.StatusSummary = "Proteksi Parsial: DMARC aktif, namun SPF belum dipasang."
	} else {
		res.StatusSummary = "Rentan Spoofing: Domain tidak memiliki record SPF maupun DMARC."
	}

	return res, nil
}
