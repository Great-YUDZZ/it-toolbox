package security

import (
	"strings"
	"testing"
)

func TestGetFullASCIITable(t *testing.T) {
	table := GetFullASCIITable()
	if len(table) != 256 {
		t.Fatalf("expected 256 entries in ASCII table, got %d", len(table))
	}

	// Verify entry 0 is NUL
	if table[0].Dec != 0 || table[0].Hex != "0x00" || table[0].Category != "Control" {
		t.Errorf("unexpected entry 0: %+v", table[0])
	}

	// Verify entry 65 is 'A'
	if table[65].Dec != 65 || table[65].Hex != "0x41" || table[65].Char != "A" || table[65].Category != "Printable" {
		t.Errorf("unexpected entry 65: %+v", table[65])
	}

	// Verify entry 127 is DEL
	if table[127].Dec != 127 || table[127].Category != "Control" {
		t.Errorf("unexpected entry 127: %+v", table[127])
	}

	// Verify entry 200 is Extended
	if table[200].Category != "Extended" {
		t.Errorf("unexpected entry 200: %+v", table[200])
	}
}

func TestFilterASCIITable(t *testing.T) {
	// Search 'A'
	resA := FilterASCIITable("65", "Printable")
	foundA := false
	for _, item := range resA {
		if item.Char == "A" {
			foundA = true
			break
		}
	}
	if !foundA {
		t.Fatalf("expected to find 'A' with query '65', got %v", resA)
	}

	// Search by control category
	resControl := FilterASCIITable("", "Control")
	if len(resControl) != 33 { // 0-31 + 127 = 33
		t.Errorf("expected 33 control characters, got %d", len(resControl))
	}

	// Search by hex
	resHex := FilterASCIITable("0x41", "")
	if len(resHex) == 0 || resHex[0].Dec != 65 {
		t.Fatalf("expected to find Dec 65 with hex 0x41, got %v", resHex)
	}
}

func TestTextToASCIIAndBack(t *testing.T) {
	sample := "Hello Security 2026!"

	// 1. Decimal test
	decStr := TextToASCII(sample, "Decimal", "Spasi")
	backDec, err := ASCIIToText(decStr, "Decimal")
	if err != nil || backDec != sample {
		t.Fatalf("Decimal round-trip failed: got %q, err %v", backDec, err)
	}

	// 2. Binary test
	binStr := TextToASCII(sample, "Binary", "Spasi")
	backBin, err := ASCIIToText(binStr, "Binary")
	if err != nil || backBin != sample {
		t.Fatalf("Binary round-trip failed: got %q, err %v", backBin, err)
	}

	// 3. Contiguous binary test
	binContiguous := strings.ReplaceAll(binStr, " ", "")
	backContiguousBin, err := ASCIIToText(binContiguous, "Binary")
	if err != nil || backContiguousBin != sample {
		t.Fatalf("Contiguous binary round-trip failed: got %q, err %v", backContiguousBin, err)
	}

	// 4. Hex test with 0x prefix
	hex0xStr := TextToASCII(sample, "Hex", "Prefiks 0x")
	backHex, err := ASCIIToText(hex0xStr, "Hex")
	if err != nil || backHex != sample {
		t.Fatalf("Hex 0x round-trip failed: got %q, err %v", backHex, err)
	}

	// 5. Hex test with \x prefix
	hexSlashXStr := TextToASCII(sample, "Hex", "Prefiks \\x")
	backSlashXHex, err := ASCIIToText(hexSlashXStr, "Hex")
	if err != nil || backSlashXHex != sample {
		t.Fatalf("Hex \\x round-trip failed: got %q, err %v", backSlashXHex, err)
	}

	// 6. Octal test
	octStr := TextToASCII(sample, "Octal", "Koma")
	backOct, err := ASCIIToText(octStr, "Octal")
	if err != nil || backOct != sample {
		t.Fatalf("Octal round-trip failed: got %q, err %v", backOct, err)
	}
}

func TestASCIIToTextErrors(t *testing.T) {
	// Invalid decimal
	if _, err := ASCIIToText("999", "Decimal"); err == nil {
		t.Error("expected error for decimal 999")
	}

	// Invalid binary
	if _, err := ASCIIToText("01020304", "Binary"); err == nil {
		t.Error("expected error for invalid binary")
	}

	// Invalid hex
	if _, err := ASCIIToText("GG", "Hex"); err == nil {
		t.Error("expected error for invalid hex GG")
	}
}
