package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComputeFileHashes(t *testing.T) {
	// Create a temporary file with known content
	tmpDir := t.TempDir()
	dummyFile := filepath.Join(tmpDir, "testfile.txt")
	content := []byte("Hello IT Toolbox Forensics 2026")
	if err := os.WriteFile(dummyFile, content, 0644); err != nil {
		t.Fatalf("failed to write dummy file: %v", err)
	}

	res, err := ComputeFileHashes(dummyFile)
	if err != nil {
		t.Fatalf("ComputeFileHashes failed: %v", err)
	}

	if res.FileName != "testfile.txt" {
		t.Errorf("expected filename testfile.txt, got %s", res.FileName)
	}
	if res.FileSize != int64(len(content)) {
		t.Errorf("expected size %d, got %d", len(content), res.FileSize)
	}
	if len(res.MD5) != 32 {
		t.Errorf("expected 32-char MD5, got %s", res.MD5)
	}
	if len(res.SHA256) != 64 {
		t.Errorf("expected 64-char SHA256, got %s", res.SHA256)
	}
	if len(res.SHA512) != 128 {
		t.Errorf("expected 128-char SHA512, got %s", res.SHA512)
	}
}

func TestComputeFileHashesDisguisedExecutable(t *testing.T) {
	tmpDir := t.TempDir()
	disguised := filepath.Join(tmpDir, "invoice.pdf")
	// MZ header (Windows executable disguised as PDF)
	payload := []byte("MZ\x90\x00\x03\x00\x00\x00\x04\x00\x00\x00\xFF\xFF")
	if err := os.WriteFile(disguised, payload, 0644); err != nil {
		t.Fatalf("failed to write disguised file: %v", err)
	}

	res, err := ComputeFileHashes(disguised)
	if err != nil {
		t.Fatalf("ComputeFileHashes failed: %v", err)
	}

	if !res.Suspicious {
		t.Error("expected disguised executable to be flagged as Suspicious")
	}
	if !strings.Contains(res.Notes, "PERINGATAN") {
		t.Errorf("expected warning note in results, got: %s", res.Notes)
	}
}

func TestParseAuthLog(t *testing.T) {
	logSample := `Oct 12 10:00:01 server sshd[1234]: Failed password for root from 192.168.1.100 port 54321 ssh2
Oct 12 10:00:05 server sshd[1235]: Failed password for root from 192.168.1.100 port 54322 ssh2
Oct 12 10:00:10 server sshd[1236]: Failed password for invalid user admin from 10.0.0.50 port 43210 ssh2
Oct 12 10:01:00 server sshd[1237]: Accepted publickey for yudz from 192.168.1.15 port 60000 ssh2`

	analysis, err := ParseAuthLog(logSample)
	if err != nil {
		t.Fatalf("ParseAuthLog failed: %v", err)
	}

	if analysis.FailedAttempts != 3 {
		t.Errorf("expected 3 failed attempts, got %d", analysis.FailedAttempts)
	}
	if analysis.InvalidUsers != 1 {
		t.Errorf("expected 1 invalid user, got %d", analysis.InvalidUsers)
	}
	if analysis.SuccessLogins != 1 {
		t.Errorf("expected 1 success login, got %d", analysis.SuccessLogins)
	}

	if len(analysis.TopAttackerIPs) == 0 || analysis.TopAttackerIPs[0].Key != "192.168.1.100" {
		t.Errorf("expected top attacker 192.168.1.100, got %v", analysis.TopAttackerIPs)
	}
}

func TestParseWebAccessLog(t *testing.T) {
	logSample := `192.168.1.200 - - [12/Oct/2026:10:00:00 +0000] "GET /index.php?id=1%20union+select+1,2,3 HTTP/1.1" 200 1234
192.168.1.200 - - [12/Oct/2026:10:00:05 +0000] "GET /../../../../etc/passwd HTTP/1.1" 404 123
10.0.0.88 - - [12/Oct/2026:10:01:00 +0000] "GET /search?q=<script>alert(1)</script> HTTP/1.1" 200 456
10.0.0.88 - - [12/Oct/2026:10:01:10 +0000] "GET /wp-login.php HTTP/1.1" 404 123`

	analysis, err := ParseWebAccessLog(logSample)
	if err != nil {
		t.Fatalf("ParseWebAccessLog failed: %v", err)
	}

	if analysis.SqlInjections != 1 {
		t.Errorf("expected 1 SQL injection, got %d", analysis.SqlInjections)
	}
	if analysis.PathTraversals != 1 {
		t.Errorf("expected 1 path traversal, got %d", analysis.PathTraversals)
	}
	if analysis.XssAttempts != 1 {
		t.Errorf("expected 1 XSS attempt, got %d", analysis.XssAttempts)
	}
	if analysis.Http404Scans != 2 {
		t.Errorf("expected 2 404 scans, got %d", analysis.Http404Scans)
	}
}

func TestGenerateBlockRules(t *testing.T) {
	ips := []string{"192.168.1.100", "10.0.0.50"}

	ufwRules := GenerateBlockRules(ips, "ufw")
	if !strings.Contains(ufwRules, "sudo ufw deny from 192.168.1.100") {
		t.Errorf("unexpected UFW output: %s", ufwRules)
	}

	iptablesRules := GenerateBlockRules(ips, "iptables")
	if !strings.Contains(iptablesRules, "sudo iptables -A INPUT -s 192.168.1.100 -j DROP") {
		t.Errorf("unexpected iptables output: %s", iptablesRules)
	}

	psRules := GenerateBlockRules(ips, "powershell")
	if !strings.Contains(psRules, "New-NetFirewallRule") {
		t.Errorf("unexpected powershell output: %s", psRules)
	}
}
