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
