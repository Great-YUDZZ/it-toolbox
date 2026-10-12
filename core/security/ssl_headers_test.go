package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCleanHost(t *testing.T) {
	cases := map[string]string{
		"https://example.com/path": "example.com",
		"http://192.168.1.1:8080/": "192.168.1.1",
		"sub.domain.co.id":          "sub.domain.co.id",
	}

	for in, expected := range cases {
		got := CleanHost(in)
		if got != expected {
			t.Errorf("CleanHost(%q) = %q; want %q", in, got, expected)
		}
	}
}

func TestInspectSecurityHeaders(t *testing.T) {
	// Mock server with A+ headers
	tsSecure := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=()")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK Secure"))
	}))
	defer tsSecure.Close()

	resSecure, err := InspectSecurityHeaders(tsSecure.URL, 2*time.Second)
	if err != nil {
		t.Fatalf("InspectSecurityHeaders failed: %v", err)
	}

	if resSecure.Score < 95 || resSecure.Grade != "A+" {
		t.Errorf("expected A+ (score >= 95), got grade %s with score %d", resSecure.Grade, resSecure.Score)
	}

	// Mock server with zero security headers
	tsInsecure := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK Insecure"))
	}))
	defer tsInsecure.Close()

	resInsecure, err := InspectSecurityHeaders(tsInsecure.URL, 2*time.Second)
	if err != nil {
		t.Fatalf("InspectSecurityHeaders failed: %v", err)
	}

	if resInsecure.Grade != "F" || resInsecure.Score != 0 {
		t.Errorf("expected F (score 0), got grade %s with score %d", resInsecure.Grade, resInsecure.Score)
	}
}

func TestCheckEmailSecurityDNSStructure(t *testing.T) {
	// Test empty
	if _, err := CheckEmailSecurityDNS(""); err == nil {
		t.Error("expected error for empty domain")
	}

	// Test valid domain structure execution (doesn't fail if offline, returns structured result)
	res, err := CheckEmailSecurityDNS("example.test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Domain != "example.test" {
		t.Errorf("expected domain cloudflare.com, got %s", res.Domain)
	}
}
