package security

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestParsePortRange(t *testing.T) {
	// Top20 default
	p1, err := ParsePortRange("")
	if err != nil || len(p1) != 20 {
		t.Fatalf("expected 20 ports from empty string, got %d (err: %v)", len(p1), err)
	}

	// List + Range
	p2, err := ParsePortRange("80, 443, 8080-8083")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []int{80, 443, 8080, 8081, 8082, 8083}
	if len(p2) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(p2))
	}
	for i := range expected {
		if p2[i] != expected[i] {
			t.Errorf("at %d: expected %d, got %d", i, expected[i], p2[i])
		}
	}

	// Invalid
	if _, err := ParsePortRange("abc"); err == nil {
		t.Error("expected error for non-numeric port")
	}
	if _, err := ParsePortRange("70000"); err == nil {
		t.Error("expected error for port > 65535")
	}
	if _, err := ParsePortRange("100-50"); err == nil {
		t.Error("expected error for reversed range")
	}
}

func TestScanTargetLocalhost(t *testing.T) {
	// Start a local mock TCP listener
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on local port: %v", err)
	}
	defer ln.Close()

	_, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("failed to get port: %v", err)
	}
	openPort, _ := strconv.Atoi(portStr)

	// Accept in background to simulate an open port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_, _ = conn.Write([]byte("SSH-2.0-MockServer_1.0\r\n"))
			_ = conn.Close()
		}
	}()

	// Scan openPort and a dummy high port (64999)
	targetPorts := []int{openPort, 64999}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results, err := ScanTarget(ctx, "127.0.0.1", targetPorts, 300*time.Millisecond, 2, nil)
	if err != nil {
		t.Fatalf("ScanTarget failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	var foundOpen bool
	for _, r := range results {
		if r.Port == openPort {
			if r.State != "OPEN" {
				t.Errorf("expected port %d to be OPEN, got %s", openPort, r.State)
			}
			foundOpen = true
		}
	}

	if !foundOpen {
		t.Errorf("mock port %d was not detected as OPEN", openPort)
	}

	report := GenerateScanReport("127.0.0.1", results)
	if !strings.Contains(report, "LAPORAN AUDIT PORT") {
		t.Errorf("expected report header in output, got: %s", report)
	}
}
