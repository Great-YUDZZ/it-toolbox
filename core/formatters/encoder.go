package formatters

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

// EncodeBase64 encodes input string to standard Base64
func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// DecodeBase64 decodes standard Base64 string back to plaintext
func DecodeBase64(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("invalid Base64 payload: %w", err)
	}
	return string(data), nil
}

// EncodeURL applies percent-encoding to a URL query parameter string
func EncodeURL(input string) string {
	return url.QueryEscape(input)
}

// DecodeURL unescapes a percent-encoded URL string
func DecodeURL(input string) (string, error) {
	decoded, err := url.QueryUnescape(input)
	if err != nil {
		return "", fmt.Errorf("invalid URL encoding: %w", err)
	}
	return decoded, nil
}
