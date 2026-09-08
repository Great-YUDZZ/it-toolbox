package calculators

import (
	"math"
	"testing"
)

func TestDecToBin(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{10, "1010"},
		{255, "11111111"},
	}

	for _, tc := range tests {
		got := DecToBin(tc.input)
		if got != tc.expected {
			t.Errorf("DecToBin(%d) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestDecToHex(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{15, "F"},
		{255, "FF"},
		{4096, "1000"},
	}

	for _, tc := range tests {
		got := DecToHex(tc.input)
		if got != tc.expected {
			t.Errorf("DecToHex(%d) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestBinToDec(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"0", 0, false},
		{"1010", 10, false},
		{"11111111", 255, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tc := range tests {
		got, err := BinToDec(tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("BinToDec(%q) error = %v; wantErr %v", tc.input, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && got != tc.expected {
			t.Errorf("BinToDec(%q) = %d; want %d", tc.input, got, tc.expected)
		}
	}
}

func TestHexToDec(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		wantErr  bool
	}{
		{"0", 0, false},
		{"FF", 255, false},
		{"0xFF", 255, false},
		{"1000", 4096, false},
		{"invalid-hex", 0, true},
		{"", 0, true},
	}

	for _, tc := range tests {
		got, err := HexToDec(tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("HexToDec(%q) error = %v; wantErr %v", tc.input, err, tc.wantErr)
			continue
		}
		if !tc.wantErr && got != tc.expected {
			t.Errorf("HexToDec(%q) = %d; want %d", tc.input, got, tc.expected)
		}
	}
}

func TestConvertDataSize(t *testing.T) {
	// 1 MB = 1024 KB
	mbToKb := ConvertDataSize(1, "MB", "KB")
	if math.Abs(mbToKb-1024) > 1e-6 {
		t.Errorf("ConvertDataSize(1, MB, KB) = %f; want 1024", mbToKb)
	}

	// 1 Byte = 8 bits
	byteToBit := ConvertDataSize(1, "Byte", "bit")
	if math.Abs(byteToBit-8) > 1e-6 {
		t.Errorf("ConvertDataSize(1, Byte, bit) = %f; want 8", byteToBit)
	}

	// 1 GB = 1048576 KB
	gbToKb := ConvertDataSize(1, "GB", "KB")
	if math.Abs(gbToKb-1048576) > 1e-6 {
		t.Errorf("ConvertDataSize(1, GB, KB) = %f; want 1048576", gbToKb)
	}

	// Unknown unit should return 0
	invalid := ConvertDataSize(10, "XYZ", "MB")
	if invalid != 0 {
		t.Errorf("ConvertDataSize with invalid unit expected 0, got %f", invalid)
	}
}
