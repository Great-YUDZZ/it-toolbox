package security

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// JWTHeader represents the decoded JWT header metadata
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid,omitempty"`
}

// JWTDetails holds the fully inspected JWT components and expiration state
type JWTDetails struct {
	RawToken    string                 `json:"raw_token"`
	Header      JWTHeader              `json:"header"`
	HeaderJSON  string                 `json:"header_json"`
	PayloadJSON string                 `json:"payload_json"`
	PayloadMap  map[string]interface{} `json:"payload_map"`
	Signature   string                 `json:"signature"`
	IssuedAt    *time.Time             `json:"issued_at,omitempty"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	NotBefore   *time.Time             `json:"not_before,omitempty"`
	IsExpired   bool                   `json:"is_expired"`
	StatusText  string                 `json:"status_text"`
}

// HashCandidate describes a plausible hash algorithm deduced from pattern
type HashCandidate struct {
	Name        string `json:"name"`
	BitLength   int    `json:"bit_length"`
	Description string `json:"description"`
}

// decodeBase64URL safely decodes standard or unpadded URL-safe Base64
func decodeBase64URL(seg string) ([]byte, error) {
	s := strings.TrimSpace(seg)
	// Pad with '=' if needed
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

// DecodeJWT parses a JSON Web Token offline into readable header, payload, and status
func DecodeJWT(tokenStr string) (*JWTDetails, error) {
	clean := strings.TrimSpace(tokenStr)
	if clean == "" {
		return nil, fmt.Errorf("token JWT tidak boleh kosong")
	}

	parts := strings.Split(clean, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("format JWT tidak valid (harus terdiri dari 3 bagian dipisah titik: header.payload.signature)")
	}

	// 1. Decode Header
	hdrBytes, err := decodeBase64URL(parts[0])
	if err != nil {
		return nil, fmt.Errorf("gagal men-decode header JWT: %w", err)
	}

	var header JWTHeader
	_ = json.Unmarshal(hdrBytes, &header)

	var prettyHdr bytes.Buffer
	if err := json.Indent(&prettyHdr, hdrBytes, "", "  "); err != nil {
		prettyHdr.WriteString(string(hdrBytes))
	}

	// 2. Decode Payload
	payloadBytes, err := decodeBase64URL(parts[1])
	if err != nil {
		return nil, fmt.Errorf("gagal men-decode payload JWT: %w", err)
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
		payloadMap = make(map[string]interface{})
	}

	var prettyPayload bytes.Buffer
	if err := json.Indent(&prettyPayload, payloadBytes, "", "  "); err != nil {
		prettyPayload.WriteString(string(payloadBytes))
	}

	// 3. Parse timestamp claims
	now := time.Now()
	var expTime, iatTime, nbfTime *time.Time
	isExpired := false
	statusText := "Token Aktif (Valid)"

	parseTimestamp := func(val interface{}) *time.Time {
		if val == nil {
			return nil
		}
		var sec int64
		switch v := val.(type) {
		case float64:
			sec = int64(v)
		case int64:
			sec = v
		case int:
			sec = int64(v)
		case string:
			s, _ := strconv.ParseInt(v, 10, 64)
			sec = s
		}
		if sec > 0 {
			t := time.Unix(sec, 0)
			return &t
		}
		return nil
	}

	if expVal, ok := payloadMap["exp"]; ok {
		expTime = parseTimestamp(expVal)
		if expTime != nil {
			if now.After(*expTime) {
				isExpired = true
				statusText = fmt.Sprintf("Token Telah Kedaluwarsa (Expired pada: %s)", expTime.Format("02 Jan 2006 15:04:05 WIB"))
			} else {
				statusText = fmt.Sprintf("Token Masih Berlaku (Kedaluwarsa pada: %s)", expTime.Format("02 Jan 2006 15:04:05 WIB"))
			}
		}
	} else {
		statusText = "Token Tanpa Batas Waktu Kedaluwarsa (No 'exp' claim)"
	}

	if iatVal, ok := payloadMap["iat"]; ok {
		iatTime = parseTimestamp(iatVal)
	}
	if nbfVal, ok := payloadMap["nbf"]; ok {
		nbfTime = parseTimestamp(nbfVal)
		if nbfTime != nil && now.Before(*nbfTime) {
			statusText = fmt.Sprintf("Token Belum Aktif (Not Before: %s)", nbfTime.Format("02 Jan 2006 15:04:05 WIB"))
		}
	}

	return &JWTDetails{
		RawToken:    clean,
		Header:      header,
		HeaderJSON:  prettyHdr.String(),
		PayloadJSON: prettyPayload.String(),
		PayloadMap:  payloadMap,
		Signature:   parts[2],
		IssuedAt:    iatTime,
		ExpiresAt:   expTime,
		NotBefore:   nbfTime,
		IsExpired:   isExpired,
		StatusText:  statusText,
	}, nil
}



// IdentifyHashType analyzes string properties to determine probable hash algorithm
func IdentifyHashType(hashStr string) []HashCandidate {
	clean := strings.TrimSpace(hashStr)
	if clean == "" {
		return nil
	}

	var candidates []HashCandidate
	isHex := regexp.MustCompile(`^[a-fA-F0-9]+$`).MatchString(clean)

	switch {
	case strings.HasPrefix(clean, "$2a$") || strings.HasPrefix(clean, "$2b$") || strings.HasPrefix(clean, "$2y$"):
		candidates = append(candidates, HashCandidate{
			Name:        "bcrypt (Blowfish Key Derivation)",
			BitLength:   184,
			Description: "Algoritma hash adaptif aman standar industri untuk kata sandi.",
		})
	case strings.HasPrefix(clean, "$argon2"):
		candidates = append(candidates, HashCandidate{
			Name:        "Argon2 (Argon2id / Argon2i)",
			BitLength:   256,
			Description: "Pemenang Password Hashing Competition (PHC), sangat tahan serangan GPU/ASIC.",
		})
	case strings.HasPrefix(clean, "$6$"):
		candidates = append(candidates, HashCandidate{
			Name:        "SHA-512 Crypt (Linux Shadow)",
			BitLength:   512,
			Description: "Format hash standar `/etc/shadow` Linux modern.",
		})
	case strings.HasPrefix(clean, "$1$"):
		candidates = append(candidates, HashCandidate{
			Name:        "MD5 Crypt (Legacy Linux Shadow)",
			BitLength:   128,
			Description: "Format hash lawas `/etc/shadow`, rentan dan usang.",
		})
	case isHex:
		switch len(clean) {
		case 32:
			candidates = append(candidates,
				HashCandidate{Name: "MD5", BitLength: 128, Description: "Algoritma hash 128-bit usang, rentan tabrakan (collisions)."},
				HashCandidate{Name: "NTLM (Windows NT LanMan)", BitLength: 128, Description: "Format hash kredensial Windows Active Directory dan SAM."},
				HashCandidate{Name: "MD4", BitLength: 128, Description: "Algoritma pendahulu MD5, sudah tidak aman."},
			)
		case 40:
			candidates = append(candidates,
				HashCandidate{Name: "SHA-1", BitLength: 160, Description: "Hash 160-bit banyak digunakan pada Git commit lawas, rentan tabrakan."},
				HashCandidate{Name: "RIPEMD-160", BitLength: 160, Description: "Hash kriptografi 160-bit yang umum dipakai di alamat Bitcoin."},
			)
		case 64:
			candidates = append(candidates,
				HashCandidate{Name: "SHA-256 (SHA-2 Family)", BitLength: 256, Description: "Standar industri hash 256-bit kuat untuk integritas dan sertifikat."},
				HashCandidate{Name: "HMAC-SHA256", BitLength: 256, Description: "Hash message authentication code berbasis kunci rahasia."},
				HashCandidate{Name: "Keccak-256 (Ethereum)", BitLength: 256, Description: "Hash varian SHA-3 yang digunakan pada smart contract Ethereum."},
			)
		case 96:
			candidates = append(candidates,
				HashCandidate{Name: "SHA-384", BitLength: 384, Description: "Hash 384-bit dari keluarga SHA-2 untuk sertifikasi tingkat tinggi."},
			)
		case 128:
			candidates = append(candidates,
				HashCandidate{Name: "SHA-512", BitLength: 512, Description: "Hash 512-bit sangat kuat untuk keamanan komputasi tinggi."},
				HashCandidate{Name: "Whirlpool", BitLength: 512, Description: "Fungsi hash 512-bit berbasis cipher block W."},
			)
		}
	}

	if len(candidates) == 0 {
		candidates = append(candidates, HashCandidate{
			Name:        "Format Kustom / Tidak Dikenali",
			BitLength:   len(clean) * 4,
			Description: fmt.Sprintf("Panjang karakter: %d. Tidak cocok dengan pola hash standar umum.", len(clean)),
		})
	}

	return candidates
}

// TransformSecurityText performs encoding and decoding operations
func TransformSecurityText(input string, mode string) (string, error) {
	switch strings.ToLower(mode) {
	case "base64-encode":
		return base64.StdEncoding.EncodeToString([]byte(input)), nil
	case "base64-decode":
		b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(input))
		if err != nil {
			return "", fmt.Errorf("gagal decode Base64: %w", err)
		}
		return string(b), nil

	case "base64url-encode":
		return base64.URLEncoding.EncodeToString([]byte(input)), nil
	case "base64url-decode":
		b, err := decodeBase64URL(input)
		if err != nil {
			return "", fmt.Errorf("gagal decode Base64 URL: %w", err)
		}
		return string(b), nil

	case "url-encode":
		return url.QueryEscape(input), nil
	case "url-decode":
		res, err := url.QueryUnescape(input)
		if err != nil {
			return "", fmt.Errorf("gagal decode URL: %w", err)
		}
		return res, nil

	case "hex-encode":
		return hex.EncodeToString([]byte(input)), nil
	case "hex-decode":
		clean := strings.ReplaceAll(strings.TrimSpace(input), " ", "")
		clean = strings.ReplaceAll(clean, "0x", "")
		b, err := hex.DecodeString(clean)
		if err != nil {
			return "", fmt.Errorf("gagal decode Hex: %w", err)
		}
		return string(b), nil

	case "html-encode":
		return html.EscapeString(input), nil
	case "html-decode":
		return html.UnescapeString(input), nil

	case "rot13":
		var sb strings.Builder
		for _, r := range input {
			switch {
			case 'a' <= r && r <= 'z':
				sb.WriteRune('a' + (r-'a'+13)%26)
			case 'A' <= r && r <= 'Z':
				sb.WriteRune('A' + (r-'A'+13)%26)
			default:
				sb.WriteRune(r)
			}
		}
		return sb.String(), nil

	default:
		return input, nil
	}
}
