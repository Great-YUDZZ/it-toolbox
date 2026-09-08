package calculators

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidBinary = errors.New("invalid binary string")
	ErrInvalidHex    = errors.New("invalid hexadecimal string")
	ErrInvalidUnit   = errors.New("unsupported data size unit")
)

// DecToBin converts a decimal integer to binary string representation
func DecToBin(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		// Use standard two's complement or signed representation
		return "-" + strconv.FormatUint(uint64(-n), 2)
	}
	return strconv.FormatUint(uint64(n), 2)
}

// DecToHex converts a decimal integer to uppercase hexadecimal string
func DecToHex(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + strings.ToUpper(strconv.FormatUint(uint64(-n), 16))
	}
	return strings.ToUpper(strconv.FormatUint(uint64(n), 16))
}

// BinToDec converts a binary string to a decimal integer
func BinToDec(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, ErrInvalidBinary
	}
	isNegative := false
	if strings.HasPrefix(s, "-") {
		isNegative = true
		s = s[1:]
	}
	val, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidBinary, err)
	}
	result := int(val)
	if isNegative {
		result = -result
	}
	return result, nil
}

// HexToDec converts a hexadecimal string to a decimal integer
func HexToDec(s string) (int, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	if s == "" {
		return 0, ErrInvalidHex
	}
	isNegative := false
	if strings.HasPrefix(s, "-") {
		isNegative = true
		s = s[1:]
	}
	val, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrInvalidHex, err)
	}
	result := int(val)
	if isNegative {
		result = -result
	}
	return result, nil
}

// DataSize multipliers relative to bytes (1 Byte = 8 bits, 1 KB = 1024 Bytes, etc.)
var unitMultipliers = map[string]float64{
	"bit":  0.125, // 1/8 byte
	"byte": 1.0,
	"kb":   1024.0,
	"mb":   1024.0 * 1024.0,
	"gb":   1024.0 * 1024.0 * 1024.0,
	"tb":   1024.0 * 1024.0 * 1024.0 * 1024.0,
}

// ConvertDataSize converts a data size value between units: bit, Byte, KB, MB, GB, TB
func ConvertDataSize(value float64, from, to string) float64 {
	fromKey := strings.ToLower(strings.TrimSpace(from))
	toKey := strings.ToLower(strings.TrimSpace(to))

	fromMult, fromOk := unitMultipliers[fromKey]
	toMult, toOk := unitMultipliers[toKey]

	if !fromOk || !toOk || toMult == 0 {
		return 0
	}

	// Normalize to bytes first, then convert to target unit
	bytes := value * fromMult
	return bytes / toMult
}
