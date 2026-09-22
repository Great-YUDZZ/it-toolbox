package calculators

import (
	"strings"
	"testing"
)

func TestGenerateMD5(t *testing.T) {
	// md5("hello") = 5d41402abc4b2a76b9719d911017c592
	got := GenerateMD5("hello")
	want := "5d41402abc4b2a76b9719d911017c592"
	if got != want {
		t.Errorf("GenerateMD5(\"hello\") = %q; want %q", got, want)
	}
}

func TestGenerateSHA256(t *testing.T) {
	// sha256("hello") = 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
	got := GenerateSHA256("hello")
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if got != want {
		t.Errorf("GenerateSHA256(\"hello\") = %q; want %q", got, want)
	}
}

func TestGenerateRandomToken(t *testing.T) {
	token := GenerateRandomToken(16)
	if len(token) != 16 {
		t.Errorf("GenerateRandomToken(16) length = %d; want 16", len(token))
	}
	token2 := GenerateRandomToken(16)
	if token == token2 {
		t.Errorf("Two random tokens should not be identical: %q vs %q", token, token2)
	}
}

func TestGeneratePassword(t *testing.T) {
	pwdNoSym := GeneratePassword(12, false)
	if len(pwdNoSym) != 12 {
		t.Errorf("GeneratePassword(12, false) length = %d; want 12", len(pwdNoSym))
	}

	for _, c := range pwdNoSym {
		if strings.ContainsRune(symbols, c) {
			t.Errorf("Password should not contain symbols when useSymbols is false, found %c", c)
		}
	}

	pwdWithSym := GeneratePassword(24, true)
	if len(pwdWithSym) != 24 {
		t.Errorf("GeneratePassword(24, true) length = %d; want 24", len(pwdWithSym))
	}
}

func TestGenerateSHA1(t *testing.T) {
	// sha1("hello") = aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
	got := GenerateSHA1("hello")
	want := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
	if got != want {
		t.Errorf("GenerateSHA1(\"hello\") = %q; want %q", got, want)
	}
}

func TestGenerateSHA512(t *testing.T) {
	// sha512("hello") = 9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043
	got := GenerateSHA512("hello")
	want := "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"
	if got != want {
		t.Errorf("GenerateSHA512(\"hello\") = %q; want %q", got, want)
	}
}

func TestGenerateAndVerifyBcrypt(t *testing.T) {
	password := "SecretP@ssword2026"
	hash, err := GenerateBcrypt(password, 10)
	if err != nil {
		t.Fatalf("GenerateBcrypt error: %v", err)
	}
	if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
		t.Errorf("Expected bcrypt hash prefix, got %q", hash)
	}

	if !VerifyBcrypt(hash, password) {
		t.Errorf("VerifyBcrypt failed for matching password")
	}

	if VerifyBcrypt(hash, "WrongPassword") {
		t.Errorf("VerifyBcrypt should fail for wrong password")
	}
}

func TestVerifyHash(t *testing.T) {
	password := "admin123"

	// bcrypt
	bHash, _ := GenerateBcrypt(password, 10)
	if !VerifyHash("bcrypt", bHash, password, "") {
		t.Errorf("VerifyHash bcrypt failed")
	}
	if VerifyHash("bcrypt", bHash, "wrong", "") {
		t.Errorf("VerifyHash bcrypt matched wrong password")
	}

	// SHA-256
	sha256Hash := GenerateSHA256(password)
	if !VerifyHash("SHA-256", sha256Hash, password, "") {
		t.Errorf("VerifyHash SHA-256 failed")
	}

	// SHA-512 with salt
	salt := "myRandomSalt99"
	sha512Hash := GenerateSHA512(password + salt)
	if !VerifyHash("SHA-512", sha512Hash, password, salt) {
		t.Errorf("VerifyHash SHA-512 with salt failed")
	}

	// MD5
	md5Hash := GenerateMD5(password)
	if !VerifyHash("MD5", md5Hash, password, "") {
		t.Errorf("VerifyHash MD5 failed")
	}

	// SHA-1
	sha1Hash := GenerateSHA1(password)
	if !VerifyHash("SHA-1", sha1Hash, password, "") {
		t.Errorf("VerifyHash SHA-1 failed")
	}
}
