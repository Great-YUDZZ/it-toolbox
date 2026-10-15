package security

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

// BaseNumberResult menampung representasi konversi bilangan multi-radiks
type BaseNumberResult struct {
	Decimal string `json:"decimal"`
	Hex     string `json:"hex"`
	Binary  string `json:"binary"`
	Octal   string `json:"octal"`
	BigInt  string `json:"bigint"`
	BitLen  int    `json:"bit_len"`
	ByteLen int    `json:"byte_len"`
}

// ConvertBaseNumber mengonversi input angka ke Dec, Hex, Binary, Octal, dan BigInt
func ConvertBaseNumber(input string, fromBase string) (*BaseNumberResult, error) {
	clean := strings.TrimSpace(input)
	if clean == "" {
		return nil, errors.New("masukan bilangan tidak boleh kosong")
	}

	bi := new(big.Int)
	var base int

	switch strings.ToLower(fromBase) {
	case "hex", "heksadesimal", "16":
		base = 16
		clean = strings.TrimPrefix(clean, "0x")
		clean = strings.TrimPrefix(clean, "0X")
		clean = strings.ReplaceAll(clean, " ", "")
	case "bin", "biner", "binary", "2":
		base = 2
		clean = strings.TrimPrefix(clean, "0b")
		clean = strings.TrimPrefix(clean, "0B")
		clean = strings.ReplaceAll(clean, " ", "")
	case "oct", "oktal", "octal", "8":
		base = 8
		clean = strings.TrimPrefix(clean, "0o")
		clean = strings.TrimPrefix(clean, "0O")
		clean = strings.ReplaceAll(clean, " ", "")
	case "dec", "desimal", "decimal", "10":
		base = 10
		clean = strings.ReplaceAll(clean, " ", "")
		clean = strings.ReplaceAll(clean, ",", "")
	default: // Otomatis deteksi dari format
		if strings.HasPrefix(clean, "0x") || strings.HasPrefix(clean, "0X") {
			base = 16
			clean = clean[2:]
		} else if strings.HasPrefix(clean, "0b") || strings.HasPrefix(clean, "0B") {
			base = 2
			clean = clean[2:]
		} else if strings.HasPrefix(clean, "0o") || strings.HasPrefix(clean, "0O") {
			base = 8
			clean = clean[2:]
		} else {
			base = 10
		}
	}

	if _, ok := bi.SetString(clean, base); !ok {
		return nil, fmt.Errorf("format input tidak valid untuk basis %d: %s", base, clean)
	}

	bitLen := bi.BitLen()
	byteLen := (bitLen + 7) / 8

	// Format Hex 0x
	hexStr := bi.Text(16)
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	hexWithPrefix := "0x" + hexStr

	// Format Biner terformat (kelipatan 8 bit)
	binStr := bi.Text(2)
	pad := (8 - (len(binStr) % 8)) % 8
	if pad > 0 {
		binStr = strings.Repeat("0", pad) + binStr
	}

	return &BaseNumberResult{
		Decimal: bi.Text(10),
		Hex:     hexWithPrefix,
		Binary:  binStr,
		Octal:   "0o" + bi.Text(8),
		BigInt:  bi.String(),
		BitLen:  bitLen,
		ByteLen: byteLen,
	}, nil
}

// WeiConverterResult menampung denominasi mata uang Ethereum presisi tinggi
type WeiConverterResult struct {
	Wei   string `json:"wei"`
	Gwei  string `json:"gwei"`
	Ether string `json:"ether"`
}

// ConvertWei mengonversi nilai antara Wei, Gwei, dan Ether menggunakan big.Float
func ConvertWei(amountStr string, fromUnit string) (*WeiConverterResult, error) {
	clean := strings.TrimSpace(amountStr)
	clean = strings.ReplaceAll(clean, ",", "")
	if clean == "" {
		return nil, errors.New("jumlah saldo tidak boleh kosong")
	}

	val, _, err := big.ParseFloat(clean, 10, 256, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("nilai angka tidak valid: %w", err)
	}

	weiMultiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)) // 10^18
	gweiMultiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil))  // 10^9

	var weiVal *big.Float

	switch strings.ToLower(fromUnit) {
	case "ether", "eth":
		weiVal = new(big.Float).Mul(val, weiMultiplier)
	case "gwei":
		weiVal = new(big.Float).Mul(val, gweiMultiplier)
	default: // "wei"
		weiVal = val
	}

	gweiVal := new(big.Float).Quo(weiVal, gweiMultiplier)
	etherVal := new(big.Float).Quo(weiVal, weiMultiplier)

	// Ubah Wei ke representasi integer bulat tanpa notasi ilmiah jika bilangan bulat
	weiInt, _ := weiVal.Int(nil)

	return &WeiConverterResult{
		Wei:   weiInt.String(),
		Gwei:  gweiVal.Text('f', 9),
		Ether: etherVal.Text('f', 18),
	}, nil
}

// Alfabet Base58 standar Bitcoin / IPFS (tanpa 0, O, I, l)
const b58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// EncodeBase58 melakukan encoding byte slice ke representasi Base58
func EncodeBase58(input []byte) string {
	if len(input) == 0 {
		return ""
	}

	var zeros int
	for zeros < len(input) && input[zeros] == 0 {
		zeros++
	}

	bi := new(big.Int).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)

	var res []byte
	for bi.Cmp(zero) > 0 {
		bi.DivMod(bi, base, mod)
		res = append(res, b58Alphabet[mod.Int64()])
	}

	// Balikkan urutan karakter
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	// Tambahkan leading zero '1'
	prefix := strings.Repeat("1", zeros)
	return prefix + string(res)
}

// DecodeBase58 melakukan decoding string Base58 menjadi slice byte
func DecodeBase58(input string) ([]byte, error) {
	clean := strings.TrimSpace(input)
	if clean == "" {
		return nil, nil
	}

	var zeros int
	for zeros < len(clean) && clean[zeros] == '1' {
		zeros++
	}

	bi := big.NewInt(0)
	base := big.NewInt(58)

	for _, ch := range clean {
		idx := strings.IndexRune(b58Alphabet, ch)
		if idx == -1 {
			return nil, fmt.Errorf("karakter tidak valid pada Base58: %c", ch)
		}
		bi.Mul(bi, base)
		bi.Add(bi, big.NewInt(int64(idx)))
	}

	decoded := bi.Bytes()
	res := make([]byte, zeros+len(decoded))
	copy(res[zeros:], decoded)
	return res, nil
}

// Keccak256Hash menghitung hash Keccak-256 khas Ethereum
func Keccak256Hash(data []byte) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return "0x" + hex.EncodeToString(hash.Sum(nil))
}

// ComputeFunctionSelector menghitung 4-byte EVM method selector dari signature fungsi Solidity
// Contoh: "transfer(address,uint256)" -> "0xa9059cbb"
func ComputeFunctionSelector(signature string) string {
	clean := strings.TrimSpace(signature)
	clean = strings.ReplaceAll(clean, " ", "")
	fullHash := Keccak256Hash([]byte(clean))
	if len(fullHash) >= 10 { // "0x" + 8 hex chars (4 bytes)
		return fullHash[:10]
	}
	return fullHash
}

// ReverseEndiannessHex membalikkan urutan byte dalam string heksadesimal (Big Endian <-> Little Endian)
func ReverseEndiannessHex(hexStr string) (string, error) {
	clean := strings.TrimSpace(hexStr)
	clean = strings.TrimPrefix(clean, "0x")
	clean = strings.TrimPrefix(clean, "0X")
	clean = strings.ReplaceAll(clean, " ", "")

	if len(clean)%2 != 0 {
		clean = "0" + clean
	}

	bytes, err := hex.DecodeString(clean)
	if err != nil {
		return "", fmt.Errorf("hex tidak valid: %w", err)
	}

	// Balik urutan byte
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return "0x" + hex.EncodeToString(bytes), nil
}

// HexToASCII mengonversi string heksadesimal menjadi string teks ASCII/UTF-8
func HexToASCII(hexStr string) (string, error) {
	clean := strings.TrimSpace(hexStr)
	clean = strings.TrimPrefix(clean, "0x")
	clean = strings.TrimPrefix(clean, "0X")
	clean = strings.ReplaceAll(clean, " ", "")

	if len(clean)%2 != 0 {
		clean = "0" + clean
	}

	bytes, err := hex.DecodeString(clean)
	if err != nil {
		return "", fmt.Errorf("format heksadesimal tidak valid: %w", err)
	}

	return string(bytes), nil
}

// ASCIIToHex mengonversi string teks biasa menjadi heksadesimal dengan prefiks 0x
func ASCIIToHex(text string) string {
	return "0x" + hex.EncodeToString([]byte(text))
}

// UnpackUint256Hex membaca 32-byte (64-char) slot data EVM ABI menjadi BigInt desimal
func UnpackUint256Hex(slotHex string) (*big.Int, error) {
	clean := strings.TrimSpace(slotHex)
	clean = strings.TrimPrefix(clean, "0x")
	clean = strings.TrimPrefix(clean, "0X")

	if len(clean) > 64 {
		clean = clean[:64]
	}

	bi := new(big.Int)
	if _, ok := bi.SetString(clean, 16); !ok {
		return nil, errors.New("gagal membedah slot uint256 hex")
	}
	return bi, nil
}

// IntToBigEndianBytes mengembalikan byte representation Big-Endian
func IntToBigEndianBytes(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return buf
}
