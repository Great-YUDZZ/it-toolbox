package calculators

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenerateMD5 computes the MD5 hex digest of the input string
func GenerateMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateSHA1 computes the SHA1 hex digest of the input string
func GenerateSHA1(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateSHA256 computes the SHA256 hex digest of the input string
func GenerateSHA256(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateSHA512 computes the SHA512 hex digest of the input string
func GenerateSHA512(input string) string {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateBcrypt hashes password using bcrypt with the specified cost (default 10 if <= 0)
func GenerateBcrypt(password string, cost int) (string, error) {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyBcrypt checks if a plaintext password matches the bcrypt hash
func VerifyBcrypt(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// VerifyHash verifies whether a password matches the hashValue using the specified algorithm
func VerifyHash(algorithm, hashValue, password, salt string) bool {
	normAlgo := strings.ToUpper(strings.TrimSpace(algorithm))
	target := strings.TrimSpace(hashValue)

	switch normAlgo {
	case "BCRYPT":
		return VerifyBcrypt(target, password)
	case "SHA-256", "SHA256":
		expected := GenerateSHA256(password + salt)
		return strings.EqualFold(target, expected)
	case "SHA-512", "SHA512":
		expected := GenerateSHA512(password + salt)
		return strings.EqualFold(target, expected)
	case "MD5":
		expected := GenerateMD5(password + salt)
		return strings.EqualFold(target, expected)
	case "SHA-1", "SHA1":
		expected := GenerateSHA1(password + salt)
		return strings.EqualFold(target, expected)
	default:
		return false
	}
}

// GenerateRandomToken generates a cryptographically secure random hexadecimal token
func GenerateRandomToken(length int) string {
	if length <= 0 {
		return ""
	}
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
			password[i] = charPool[i%len(charPool)]
			continue
		}
		password[i] = charPool[randomIndex.Int64()]
	}

	return string(password)
}
