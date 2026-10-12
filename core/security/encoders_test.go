package security

import (
	"strings"
	"testing"
)

func TestDecodeJWT(t *testing.T) {
	// Sample valid JWT:
	// Header: {"alg":"HS256","typ":"JWT"}
	// Payload: {"sub":"1234567890","name":"Yuda Security","admin":true,"exp":1700000000}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ill1ZGEgU2VjdXJpdHkiLCJhZG1pbiI6dHJ1ZSwiZXhwIjoxNzAwMDAwMDAwfQ.dyt0CoTl4WoVjAHI9Q_CwSKhl6d_9rhM3NrXuJttino"

	details, err := DecodeJWT(token)
	if err != nil {
		t.Fatalf("DecodeJWT failed: %v", err)
	}

	if details.Header.Alg != "HS256" {
		t.Errorf("expected alg HS256, got %s", details.Header.Alg)
	}
	if details.Header.Typ != "JWT" {
		t.Errorf("expected typ JWT, got %s", details.Header.Typ)
	}

	if details.PayloadMap["name"] != "Yuda Security" {
		t.Errorf("expected name Yuda Security, got %v", details.PayloadMap["name"])
	}

	if !details.IsExpired {
		t.Error("expected token with exp in past to be marked as expired")
	}

	if !strings.Contains(details.StatusText, "Kedaluwarsa") {
		t.Errorf("expected expired text in status, got: %s", details.StatusText)
	}
}

func TestIdentifyHashType(t *testing.T) {
	// MD5 (32 hex)
	md5Candidates := IdentifyHashType("5d41402abc4b2a76b9719d911017c592")
	if len(md5Candidates) == 0 || md5Candidates[0].Name != "MD5" {
		t.Errorf("expected MD5 candidate, got %v", md5Candidates)
	}

	// SHA-256 (64 hex)
	sha256Candidates := IdentifyHashType("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824")
	if len(sha256Candidates) == 0 || !strings.Contains(sha256Candidates[0].Name, "SHA-256") {
		t.Errorf("expected SHA-256 candidate, got %v", sha256Candidates)
	}

	// Bcrypt
	bcryptCandidates := IdentifyHashType("$2a$12$R9h/cIPz0gi.URNNX3kh2OPST9/PgBkqquzi.Ss7KIUgO2t0jWMUW")
	if len(bcryptCandidates) == 0 || !strings.Contains(bcryptCandidates[0].Name, "bcrypt") {
		t.Errorf("expected bcrypt candidate, got %v", bcryptCandidates)
	}

	// Argon2
	argonCandidates := IdentifyHashType("$argon2id$v=19$m=65536,t=3,p=4$dummyhash")
	if len(argonCandidates) == 0 || !strings.Contains(argonCandidates[0].Name, "Argon2") {
		t.Errorf("expected Argon2 candidate, got %v", argonCandidates)
	}
}

func TestTransformSecurityText(t *testing.T) {
	orig := "AdminPassword123! <script> & /path?q=test"

	// 1. Base64
	b64, err := TransformSecurityText(orig, "base64-encode")
	if err != nil {
		t.Fatalf("base64-encode failed: %v", err)
	}
	backB64, err := TransformSecurityText(b64, "base64-decode")
	if err != nil || backB64 != orig {
		t.Errorf("base64 round-trip failed: got %q, want %q", backB64, orig)
	}

	// 2. URL
	urlEnc, err := TransformSecurityText(orig, "url-encode")
	if err != nil {
		t.Fatalf("url-encode failed: %v", err)
	}
	backURL, err := TransformSecurityText(urlEnc, "url-decode")
	if err != nil || backURL != orig {
		t.Errorf("url round-trip failed: got %q, want %q", backURL, orig)
	}

	// 3. Hex
	hexEnc, err := TransformSecurityText("HelloHex", "hex-encode")
	if err != nil {
		t.Fatalf("hex-encode failed: %v", err)
	}
	backHex, err := TransformSecurityText(hexEnc, "hex-decode")
	if err != nil || backHex != "HelloHex" {
		t.Errorf("hex round-trip failed: got %q, want HelloHex", backHex)
	}

	// 4. ROT13
	rot, err := TransformSecurityText("SecretMessage", "rot13")
	if err != nil {
		t.Fatalf("rot13 failed: %v", err)
	}
	backROT, _ := TransformSecurityText(rot, "rot13")
	if backROT != "SecretMessage" {
		t.Errorf("rot13 round-trip failed: got %q, want SecretMessage", backROT)
	}
}
