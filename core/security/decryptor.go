package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidCiphertext = errors.New("ciphertext tidak valid atau format encoding salah")
	ErrInvalidKeyLength  = errors.New("panjang kunci tidak sesuai untuk algoritma yang dipilih")
	ErrInvalidIVLength   = errors.New("panjang Initialization Vector (IV) tidak valid")
	ErrDecryptionFailed  = errors.New("dekripsi gagal, periksa kecocokan kunci, IV, atau integritas padding")
	ErrInvalidMD5Hash    = errors.New("format hash MD5 tidak valid (harus 32 karakter heksadesimal)")
	ErrHashNotFound      = errors.New("plaintext tidak ditemukan pada kamus lokal maupun basis data daring")
)

// pkcs7Unpad removes PKCS#7 padding from decrypted plaintext
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, ErrDecryptionFailed
	}
	paddingLen := int(data[len(data)-1])
	if paddingLen == 0 || paddingLen > blockSize || paddingLen > len(data) {
		return nil, ErrDecryptionFailed
	}
	for i := len(data) - paddingLen; i < len(data); i++ {
		if data[i] != byte(paddingLen) {
			return nil, ErrDecryptionFailed
		}
	}
	return data[:len(data)-paddingLen], nil
}

// DecryptAES_CBC decrypts AES-CBC ciphertext with PKCS#7 unpadding
func DecryptAES_CBC(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("%w: AES membutuhkan kunci 16, 24, atau 32 byte (didapat %d byte)", ErrInvalidKeyLength, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) == 0 || len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("%w: panjang ciphertext harus kelipatan %d byte", ErrInvalidCiphertext, aes.BlockSize)
	}
	if len(iv) != aes.BlockSize {
		return nil, fmt.Errorf("%w: IV AES-CBC harus tepat 16 byte (didapat %d byte)", ErrInvalidIVLength, len(iv))
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(ciphertext))
	mode.CryptBlocks(plain, ciphertext)

	unpadded, err := pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return unpadded, nil
}

// DecryptAES_GCM decrypts AES-GCM authenticated ciphertext
func DecryptAES_GCM(ciphertext []byte, key []byte, nonce []byte) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("%w: AES-GCM membutuhkan kunci 16, 24, atau 32 byte", ErrInvalidKeyLength)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(nonce) != gcm.NonceSize() {
		return nil, fmt.Errorf("%w: nonce AES-GCM harus tepat %d byte (didapat %d byte)", ErrInvalidIVLength, gcm.NonceSize(), len(nonce))
	}

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: autentikasi integritas GCM gagal (kunci atau tag salah)", ErrDecryptionFailed)
	}
	return plain, nil
}

// DecryptDES_CBC decrypts DES-CBC ciphertext with PKCS#7 unpadding
func DecryptDES_CBC(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	if len(key) != 8 {
		return nil, fmt.Errorf("%w: DES membutuhkan kunci tepat 8 byte (didapat %d byte)", ErrInvalidKeyLength, len(key))
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) == 0 || len(ciphertext)%des.BlockSize != 0 {
		return nil, fmt.Errorf("%w: panjang ciphertext harus kelipatan %d byte", ErrInvalidCiphertext, des.BlockSize)
	}
	if len(iv) != des.BlockSize {
		return nil, fmt.Errorf("%w: IV DES harus tepat 8 byte", ErrInvalidIVLength)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(ciphertext))
	mode.CryptBlocks(plain, ciphertext)

	return pkcs7Unpad(plain, des.BlockSize)
}

// DecryptTripleDES_CBC decrypts 3DES-CBC ciphertext with PKCS#7 unpadding
func DecryptTripleDES_CBC(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	if len(key) != 24 {
		return nil, fmt.Errorf("%w: Triple DES membutuhkan kunci tepat 24 byte (didapat %d byte)", ErrInvalidKeyLength, len(key))
	}
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) == 0 || len(ciphertext)%des.BlockSize != 0 {
		return nil, fmt.Errorf("%w: panjang ciphertext harus kelipatan %d byte", ErrInvalidCiphertext, des.BlockSize)
	}
	if len(iv) != des.BlockSize {
		return nil, fmt.Errorf("%w: IV 3DES harus tepat 8 byte", ErrInvalidIVLength)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(ciphertext))
	mode.CryptBlocks(plain, ciphertext)

	return pkcs7Unpad(plain, des.BlockSize)
}

// DecryptXOR performs repeating-key XOR decryption
func DecryptXOR(ciphertext []byte, key []byte) []byte {
	if len(key) == 0 {
		return ciphertext
	}
	res := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i++ {
		res[i] = ciphertext[i] ^ key[i%len(key)]
	}
	return res
}

// DecryptVigenere decrypts alphabetic Vigenere ciphertext preserving case and non-letters
func DecryptVigenere(ciphertext string, key string) string {
	cleanKey := strings.ToUpper(regexp.MustCompile(`[^a-zA-Z]`).ReplaceAllString(key, ""))
	if cleanKey == "" {
		return ciphertext
	}

	var sb strings.Builder
	keyIdx := 0
	keyLen := len(cleanKey)

	for _, ch := range ciphertext {
		if ch >= 'a' && ch <= 'z' {
			shift := int(cleanKey[keyIdx%keyLen] - 'A')
			dec := int(ch-'a') - shift
			if dec < 0 {
				dec += 26
			}
			sb.WriteRune(rune('a' + dec))
			keyIdx++
		} else if ch >= 'A' && ch <= 'Z' {
			shift := int(cleanKey[keyIdx%keyLen] - 'A')
			dec := int(ch-'A') - shift
			if dec < 0 {
				dec += 26
			}
			sb.WriteRune(rune('A' + dec))
			keyIdx++
		} else {
			sb.WriteRune(ch)
		}
	}
	return sb.String()
}

// parseBytes parses input string based on encoding (base64, hex, or raw)
func parseBytes(input string, encoding string) ([]byte, error) {
	trimmed := strings.TrimSpace(input)
	switch strings.ToLower(encoding) {
	case "base64":
		return base64.StdEncoding.DecodeString(trimmed)
	case "hex":
		cleanHex := strings.ReplaceAll(trimmed, " ", "")
		cleanHex = strings.ReplaceAll(cleanHex, "0x", "")
		cleanHex = strings.ReplaceAll(cleanHex, "\\x", "")
		return hex.DecodeString(cleanHex)
	default:
		return []byte(input), nil
	}
}

// DecryptWithKey is the high-level unified function for the GUI
func DecryptWithKey(algo, inputEncoding, cipherTextStr, keyStr, ivStr string) (string, error) {
	if strings.TrimSpace(cipherTextStr) == "" {
		return "", errors.New("ciphertext tidak boleh kosong")
	}

	algoNorm := strings.ToLower(strings.TrimSpace(algo))

	// Classic Vigenere is pure text-based
	if strings.Contains(algoNorm, "vigenere") {
		if strings.TrimSpace(keyStr) == "" {
			return "", errors.New("kunci Vigenere tidak boleh kosong")
		}
		return DecryptVigenere(cipherTextStr, keyStr), nil
	}

	cipherBytes, err := parseBytes(cipherTextStr, inputEncoding)
	if err != nil {
		return "", fmt.Errorf("gagal membaca ciphertext: format %s tidak valid", inputEncoding)
	}

	// Determine key bytes: if key starts with 0x or looks like hex and len matches, allow hex decode
	var keyBytes []byte
	trimmedKey := strings.TrimSpace(keyStr)
	if strings.HasPrefix(trimmedKey, "0x") {
		keyBytes, err = hex.DecodeString(strings.TrimPrefix(trimmedKey, "0x"))
		if err != nil {
			keyBytes = []byte(trimmedKey)
		}
	} else {
		keyBytes = []byte(keyStr)
	}

	// Determine IV bytes
	var ivBytes []byte
	trimmedIV := strings.TrimSpace(ivStr)
	if trimmedIV != "" {
		if strings.HasPrefix(trimmedIV, "0x") {
			ivBytes, _ = hex.DecodeString(strings.TrimPrefix(trimmedIV, "0x"))
		} else if len(trimmedIV) == 32 && regexp.MustCompile(`^[a-fA-F0-9]+$`).MatchString(trimmedIV) {
			ivBytes, _ = hex.DecodeString(trimmedIV)
		} else {
			ivBytes = []byte(trimmedIV)
		}
	}

	var resBytes []byte

	switch {
	case strings.Contains(algoNorm, "gcm"):
		if len(ivBytes) == 0 {
			ivBytes = make([]byte, 12) // Default 12 bytes nonce if empty
		}
		resBytes, err = DecryptAES_GCM(cipherBytes, keyBytes, ivBytes)
	case strings.Contains(algoNorm, "aes") || strings.Contains(algoNorm, "cbc"):
		if len(ivBytes) == 0 {
			ivBytes = make([]byte, aes.BlockSize) // Default zero IV if not provided
		}
		resBytes, err = DecryptAES_CBC(cipherBytes, keyBytes, ivBytes)
	case strings.Contains(algoNorm, "3des") || strings.Contains(algoNorm, "triple"):
		if len(ivBytes) == 0 {
			ivBytes = make([]byte, des.BlockSize)
		}
		resBytes, err = DecryptTripleDES_CBC(cipherBytes, keyBytes, ivBytes)
	case strings.Contains(algoNorm, "des"):
		if len(ivBytes) == 0 {
			ivBytes = make([]byte, des.BlockSize)
		}
		resBytes, err = DecryptDES_CBC(cipherBytes, keyBytes, ivBytes)
	case strings.Contains(algoNorm, "xor"):
		resBytes = DecryptXOR(cipherBytes, keyBytes)
	default:
		return "", fmt.Errorf("algoritma tidak dikenali: %s", algo)
	}

	if err != nil {
		return "", err
	}

	// Check if output is printable UTF-8 or raw binary
	return string(resBytes), nil
}

// ----------------------------------------------------------------------------
// MD5 Reverse Decryptor (Hybrid Lookup Engine)
// ----------------------------------------------------------------------------

// MD5ReverseResult holds the details of a found MD5 plaintext
type MD5ReverseResult struct {
	Hash       string `json:"hash"`
	Plaintext  string `json:"plaintext"`
	Source     string `json:"source"` // "Kamus Lokal", "Online API", "Custom Wordlist"
	DurationMs int64  `json:"duration_ms"`
}

// Embedded top common passwords dictionary for instant offline cracking
var commonPasswords = []string{
	"123456", "password", "123456789", "12345678", "12345", "111111", "1234567", "sunshine",
	"qwerty", "iloveyou", "princess", "admin", "welcome", "monkey", "charlie", "donald",
	"football", "dragon", "master", "michael", "shadow", "superman", "jordan", "harley",
	"secret", "killer", "matrix", "test", "root", "toor", "pass", "login", "guest", "default",
	"cisco", "class", "administrator", "system", "oracle", "postgres", "database", "security",
	"hacker", "cyber", "network", "router", "switch", "access", "terminal", "packet",
	"000000", "123123", "654321", "666666", "777777", "888888", "999999", "123321",
	"abcdef", "abc123", "letmein", "trustno1", "starwars", "batman", "liverpool", "arsenal",
	"chelsea", "barcelona", "realmadrid", "family", "computer", "internet", "cookie", "coffee",
	"summer", "spring", "winter", "autumn", "freedom", "love", "peace", "winner", "champion",
	"hello", "world", "helloworld", "god", "jesus", "jessica", "ashley", "brandon", "jennifer",
	"diamond", "gold", "silver", "platinum", "monster", "ginger", "pepper", "bandit", "baxter",
	"jakarta", "indonesia", "surabaya", "bandung", "rahasia", "merdeka", "garuda", "nusantara",
}

// ComputeMD5 calculates hex MD5 digest
func ComputeMD5(text string) string {
	h := md5.Sum([]byte(text))
	return hex.EncodeToString(h[:])
}

// LookupMD5Local checks if the hash matches any embedded common password
func LookupMD5Local(targetHash string) (string, bool) {
	norm := strings.ToLower(strings.TrimSpace(targetHash))
	for _, word := range commonPasswords {
		if ComputeMD5(word) == norm {
			return word, true
		}
	}
	return "", false
}

// LookupMD5Custom checks against user-supplied candidate words
func LookupMD5Custom(targetHash string, candidates []string) (string, bool) {
	norm := strings.ToLower(strings.TrimSpace(targetHash))
	for _, word := range candidates {
		w := strings.TrimSpace(word)
		if w == "" {
			continue
		}
		if ComputeMD5(w) == norm {
			return w, true
		}
	}
	return "", false
}

// LookupMD5Online queries free public reverse hash databases
func LookupMD5Online(ctx context.Context, targetHash string) (string, error) {
	cleanHash := strings.ToLower(strings.TrimSpace(targetHash))
	if len(cleanHash) != 32 {
		return "", ErrInvalidMD5Hash
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Strategy 1: hashtoolkit.com API
	apiURL1 := fmt.Sprintf("https://hashtoolkit.com/reverse-hash/?hash=%s", url.QueryEscape(cleanHash))
	req1, err := http.NewRequestWithContext(ctx, "GET", apiURL1, nil)
	if err == nil {
		req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) IT-Toolbox-Client")
		resp, err := client.Do(req1)
		if err == nil && resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*64))
			resp.Body.Close()

			// Check for text within res-text tag or span
			bodyStr := string(body)
			re := regexp.MustCompile(`title="decrypted md5 hash">([^<]+)</a>`)
			matches := re.FindStringSubmatch(bodyStr)
			if len(matches) > 1 {
				res := strings.TrimSpace(matches[1])
				if res != "" && ComputeMD5(res) == cleanHash {
					return res, nil
				}
			}
		}
	}

	// Strategy 2: md5decrypt.net public query API
	apiURL2 := fmt.Sprintf("https://md5decrypt.net/Api/api.php?hash=%s&hash_type=md5&email=deanna_abshire@proxymail.eu&code=1152464b80a689d2", url.QueryEscape(cleanHash))
	req2, err := http.NewRequestWithContext(ctx, "GET", apiURL2, nil)
	if err == nil {
		req2.Header.Set("User-Agent", "IT-Toolbox-Client")
		resp, err := client.Do(req2)
		if err == nil && resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
			resp.Body.Close()
			candidate := strings.TrimSpace(string(body))
			if candidate != "" && !strings.Contains(candidate, "ERROR") && ComputeMD5(candidate) == cleanHash {
				return candidate, nil
			}
		}
	}

	// Strategy 3: nitrxgen API
	apiURL3 := fmt.Sprintf("https://www.nitrxgen.net/md5db/%s", url.QueryEscape(cleanHash))
	req3, err := http.NewRequestWithContext(ctx, "GET", apiURL3, nil)
	if err == nil {
		req3.Header.Set("User-Agent", "IT-Toolbox-Client")
		resp, err := client.Do(req3)
		if err == nil && resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
			resp.Body.Close()
			candidate := strings.TrimSpace(string(body))
			if candidate != "" && ComputeMD5(candidate) == cleanHash {
				return candidate, nil
			}
		}
	}

	return "", ErrHashNotFound
}

// ReverseMD5Hybrid runs local dictionary first, then custom wordlist, then online query
func ReverseMD5Hybrid(ctx context.Context, hash string, customWords []string) (*MD5ReverseResult, error) {
	cleanHash := strings.ToLower(strings.TrimSpace(hash))
	if len(cleanHash) != 32 || !regexp.MustCompile(`^[a-fA-F0-9]{32}$`).MatchString(cleanHash) {
		return nil, ErrInvalidMD5Hash
	}

	start := time.Now()

	// 1. Check local dictionary
	if plain, found := LookupMD5Local(cleanHash); found {
		return &MD5ReverseResult{
			Hash:       cleanHash,
			Plaintext:  plain,
			Source:     "Kamus Lokal (Offline Common Passwords)",
			DurationMs: time.Since(start).Milliseconds(),
		}, nil
	}

	// 2. Check user custom wordlist if provided
	if len(customWords) > 0 {
		if plain, found := LookupMD5Custom(cleanHash, customWords); found {
			return &MD5ReverseResult{
				Hash:       cleanHash,
				Plaintext:  plain,
				Source:     "Custom Wordlist Pengguna",
				DurationMs: time.Since(start).Milliseconds(),
			}, nil
		}
	}

	// 3. Fallback to online reverse databases
	plain, err := LookupMD5Online(ctx, cleanHash)
	if err == nil && plain != "" {
		return &MD5ReverseResult{
			Hash:       cleanHash,
			Plaintext:  plain,
			Source:     "Online Reverse Database (Public API)",
			DurationMs: time.Since(start).Milliseconds(),
		}, nil
	}

	return nil, ErrHashNotFound
}
