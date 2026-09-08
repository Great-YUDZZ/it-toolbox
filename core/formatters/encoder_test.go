package formatters

import (
	"testing"
)

func TestBase64(t *testing.T) {
	raw := "Antigravity IT Toolbox 2026!"
	encoded := EncodeBase64(raw)
	if encoded == "" || encoded == raw {
		t.Fatalf("EncodeBase64 produced invalid output: %q", encoded)
	}

	decoded, err := DecodeBase64(encoded)
	if err != nil {
		t.Fatalf("DecodeBase64 error: %v", err)
	}
	if decoded != raw {
		t.Errorf("DecodeBase64 = %q; want %q", decoded, raw)
	}

	_, errInv := DecodeBase64("!!!not-valid-base64!!!")
	if errInv == nil {
		t.Errorf("Expected error for invalid base64, got nil")
	}
}

func TestURLEncodeDecode(t *testing.T) {
	raw := "key=hello world & special/chars=100%"
	encoded := EncodeURL(raw)
	if encoded == raw {
		t.Fatalf("EncodeURL produced identical output: %q", encoded)
	}

	decoded, err := DecodeURL(encoded)
	if err != nil {
		t.Fatalf("DecodeURL error: %v", err)
	}
	if decoded != raw {
		t.Errorf("DecodeURL = %q; want %q", decoded, raw)
	}
}
