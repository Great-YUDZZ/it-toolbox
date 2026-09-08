package calculators

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// GenerateMD5 computes the MD5 hex digest of the input string
func GenerateMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateSHA256 computes the SHA256 hex digest of the input string
func GenerateSHA256(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateRandomToken generates a cryptographically secure random hexadecimal token
func GenerateRandomToken(length int) string {
	if length <= 0 {
		return ""
	}
	// length in hex characters requires length/2 bytes, rounded up
	byteCount := (length + 1) / 2
	b := make([]byte, byteCount)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	s := hex.EncodeToString(b)
	if len(s) > length {
		return s[:length]
	}
	return s
}

const (
	lettersUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lettersLower = "abcdefghijklmnopqrstuvwxyz"
	digits       = "0123456789"
	symbols      = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

// GeneratePassword creates a secure random password of given length with or without symbols
func GeneratePassword(length int, useSymbols bool) string {
	if length <= 0 {
		return ""
	}

	charPool := lettersUpper + lettersLower + digits
	if useSymbols {
		charPool += symbols
	}

	poolLen := big.NewInt(int64(len(charPool)))
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, poolLen)
		if err != nil {
			// Fallback index
			password[i] = charPool[i%len(charPool)]
			continue
		}
		password[i] = charPool[randomIndex.Int64()]
	}

	return string(password)
}
