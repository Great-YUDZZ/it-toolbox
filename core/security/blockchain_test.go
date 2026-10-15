package security_test

import (
	"testing"

	"github.com/yudz/it-toolbox/core/security"
)

func TestConvertBaseNumber(t *testing.T) {
	// Tes Hex ke basis lain
	res, err := security.ConvertBaseNumber("0xFF", "hex")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Decimal != "255" {
		t.Errorf("expected decimal 255, got %s", res.Decimal)
	}
	if res.Binary != "11111111" {
		t.Errorf("expected binary 11111111, got %s", res.Binary)
	}
	if res.Octal != "0o377" {
		t.Errorf("expected octal 0o377, got %s", res.Octal)
	}

	// Tes Desimal besar (256-bit BigInt)
	bigVal := "115792089237316195423570985008687907853269984665640564039457584007913129639935" // 2^256 - 1
	resBig, err := security.ConvertBaseNumber(bigVal, "dec")
	if err != nil {
		t.Fatalf("unexpected error parsing 256-bit uint: %v", err)
	}
	expectedHex := "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	if resBig.Hex != expectedHex {
		t.Errorf("expected max uint256 hex, got %s", resBig.Hex)
	}
	if resBig.BitLen != 256 {
		t.Errorf("expected bit length 256, got %d", resBig.BitLen)
	}
}

func TestConvertWei(t *testing.T) {
	// 1 Ether = 10^18 Wei = 10^9 Gwei
	res, err := security.ConvertWei("1", "ether")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Wei != "1000000000000000000" {
		t.Errorf("expected 10^18 Wei, got %s", res.Wei)
	}
	if res.Gwei != "1000000000.000000000" {
		t.Errorf("expected 10^9 Gwei, got %s", res.Gwei)
	}

	// Konversi dari Gwei
	resGwei, err := security.ConvertWei("50", "gwei")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resGwei.Wei != "50000000000" {
		t.Errorf("expected 50*10^9 Wei, got %s", resGwei.Wei)
	}
}

func TestBase58Encoding(t *testing.T) {
	input := []byte("Hello World")
	encoded := security.EncodeBase58(input)
	decoded, err := security.DecodeBase58(encoded)
	if err != nil {
		t.Fatalf("unexpected error decoding Base58: %v", err)
	}
	if string(decoded) != string(input) {
		t.Errorf("expected %s, got %s", string(input), string(decoded))
	}
}

func TestFunctionSelector(t *testing.T) {
	// ERC-20 transfer(address,uint256) -> 0xa9059cbb
	sel := security.ComputeFunctionSelector("transfer(address,uint256)")
	expected := "0xa9059cbb"
	if sel != expected {
		t.Errorf("expected selector %s, got %s", expected, sel)
	}

	// ERC-20 approve(address,uint256) -> 0x095ea7b3
	selApprove := security.ComputeFunctionSelector("approve(address,uint256)")
	if selApprove != "0x095ea7b3" {
		t.Errorf("expected approve selector 0x095ea7b3, got %s", selApprove)
	}
}

func TestKeccak256Hash(t *testing.T) {
	// Keccak256("") = 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470
	hash := security.Keccak256Hash([]byte(""))
	expected := "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	if hash != expected {
		t.Errorf("expected empty keccak256 hash %s, got %s", expected, hash)
	}
}

func TestEndiannessAndHexASCII(t *testing.T) {
	// Test byte reversing
	rev, err := security.ReverseEndiannessHex("0x12345678")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rev != "0x78563412" {
		t.Errorf("expected reversed hex 0x78563412, got %s", rev)
	}

	// Test Hex to ASCII
	ascii, err := security.HexToASCII("0x48656c6c6f")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ascii != "Hello" {
		t.Errorf("expected 'Hello', got %s", ascii)
	}

	// Test ASCII to Hex
	hexBack := security.ASCIIToHex("Hello")
	if hexBack != "0x48656c6c6f" {
		t.Errorf("expected '0x48656c6c6f', got %s", hexBack)
	}
}
