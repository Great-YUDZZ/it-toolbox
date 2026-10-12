package security

import (
	"fmt"
	"strconv"
	"strings"
)

// ASCIIEntry represents a single ASCII table item with multiple representations
type ASCIIEntry struct {
	Dec      int    `json:"dec"`
	Hex      string `json:"hex"`
	Oct      string `json:"oct"`
	Bin      string `json:"bin"`
	Char     string `json:"char"`
	Name     string `json:"name"`
	Category string `json:"category"` // "Control", "Printable", "Extended"
}

var controlNames = map[int]string{
	0:   "NUL (Null character)",
	1:   "SOH (Start of Heading)",
	2:   "STX (Start of Text)",
	3:   "ETX (End of Text)",
	4:   "EOT (End of Transmission)",
	5:   "ENQ (Enquiry)",
	6:   "ACK (Acknowledge)",
	7:   "BEL (Bell / Alert)",
	8:   "BS (Backspace)",
	9:   "TAB (Horizontal Tab)",
	10:  "LF (Line Feed / Newline)",
	11:  "VT (Vertical Tab)",
	12:  "FF (Form Feed)",
	13:  "CR (Carriage Return)",
	14:  "SO (Shift Out)",
	15:  "SI (Shift In)",
	16:  "DLE (Data Link Escape)",
	17:  "DC1 (Device Control 1 / XON)",
	18:  "DC2 (Device Control 2)",
	19:  "DC3 (Device Control 3 / XOFF)",
	20:  "DC4 (Device Control 4)",
	21:  "NAK (Negative Acknowledge)",
	22:  "SYN (Synchronous Idle)",
	23:  "ETB (End of Transmission Block)",
	24:  "CAN (Cancel)",
	25:  "EM (End of Medium)",
	26:  "SUB (Substitute)",
	27:  "ESC (Escape)",
	28:  "FS (File Separator)",
	29:  "GS (Group Separator)",
	30:  "RS (Record Separator)",
	31:  "US (Unit Separator)",
	32:  "Space",
	127: "DEL (Delete)",
}

// GetFullASCIITable generates the complete 256-entry ASCII reference table
func GetFullASCIITable() []ASCIIEntry {
	table := make([]ASCIIEntry, 256)
	for i := 0; i < 256; i++ {
		hexStr := fmt.Sprintf("%02X", i)
		octStr := fmt.Sprintf("%03o", i)
		binStr := fmt.Sprintf("%08b", i)

		var charStr string
		var nameStr string
		var category string

		switch {
		case i < 32:
			category = "Control"
			charStr = fmt.Sprintf("[%02d]", i)
			nameStr = controlNames[i]
		case i == 32:
			category = "Printable"
			charStr = "[SPACE]"
			nameStr = "Space"
		case i < 127:
			category = "Printable"
			charStr = string(rune(i))
			nameStr = fmt.Sprintf("Character '%c'", rune(i))
		case i == 127:
			category = "Control"
			charStr = "[DEL]"
			nameStr = controlNames[127]
		default:
			category = "Extended"
			charStr = string(rune(i))
			nameStr = fmt.Sprintf("Extended ASCII #%d", i)
		}

		table[i] = ASCIIEntry{
			Dec:      i,
			Hex:      "0x" + hexStr,
			Oct:      octStr,
			Bin:      binStr,
			Char:     charStr,
			Name:     nameStr,
			Category: category,
		}
	}
	return table
}

// FilterASCIITable filters ASCII table items by search query and category
func FilterASCIITable(query string, category string) []ASCIIEntry {
	full := GetFullASCIITable()
	q := strings.TrimSpace(strings.ToLower(query))
	cat := strings.TrimSpace(strings.ToLower(category))

	var results []ASCIIEntry
	for _, item := range full {
		if cat != "" && cat != "semua" && strings.ToLower(item.Category) != cat {
			continue
		}

		if q == "" {
			results = append(results, item)
			continue
		}

		// Match dec, hex, oct, bin, char, or name
		decStr := strconv.Itoa(item.Dec)
		if strings.Contains(decStr, q) ||
			strings.Contains(strings.ToLower(item.Hex), q) ||
			strings.Contains(item.Oct, q) ||
			strings.Contains(item.Bin, q) ||
			strings.Contains(strings.ToLower(item.Char), q) ||
			strings.Contains(strings.ToLower(item.Name), q) {
			results = append(results, item)
		}
	}
	return results
}

// TextToASCII converts plain text into specified representation with delimiter
func TextToASCII(text string, format string, delimiter string) string {
	if text == "" {
		return ""
	}

	bytes := []byte(text)
	parts := make([]string, len(bytes))

	sep := delimiter
	switch strings.ToLower(delimiter) {
	case "spasi", "space", " ":
		sep = " "
	case "koma", "comma", ",":
		sep = ", "
	case "tanpa spasi", "none", "":
		sep = ""
	case "prefiks 0x":
		sep = " "
	case "prefiks \\x":
		sep = ""
	}

	useHex0x := strings.EqualFold(delimiter, "prefiks 0x")
	useHexSlashX := strings.EqualFold(delimiter, "prefiks \\x")

	for i, b := range bytes {
		switch strings.ToLower(format) {
		case "binary", "bin", "biner":
			parts[i] = fmt.Sprintf("%08b", b)
		case "hex", "hexadecimal", "heksadesimal":
			if useHex0x {
				parts[i] = fmt.Sprintf("0x%02X", b)
			} else if useHexSlashX {
				parts[i] = fmt.Sprintf("\\x%02X", b)
			} else {
				parts[i] = fmt.Sprintf("%02X", b)
			}
		case "octal", "oct", "oktal":
			parts[i] = fmt.Sprintf("%03o", b)
		default: // Decimal
			parts[i] = strconv.Itoa(int(b))
		}
	}

	return strings.Join(parts, sep)
}

// ASCIIToText converts formatted ASCII numbers back into readable text string
func ASCIIToText(input string, format string) (string, error) {
	clean := strings.TrimSpace(input)
	if clean == "" {
		return "", nil
	}

	var bytes []byte

	switch strings.ToLower(format) {
	case "binary", "bin", "biner":
		// Can be separated by spaces, commas, or contiguous
		clean = strings.ReplaceAll(clean, ",", " ")
		fields := strings.Fields(clean)
		if len(fields) == 1 && len(fields[0])%8 == 0 && len(fields[0]) > 8 {
			// Chunk into 8-bit slices
			s := fields[0]
			fields = nil
			for i := 0; i < len(s); i += 8 {
				fields = append(fields, s[i:i+8])
			}
		}
		for _, f := range fields {
			val, err := strconv.ParseUint(f, 2, 8)
			if err != nil {
				return "", fmt.Errorf("byte biner tidak valid: %q", f)
			}
			bytes = append(bytes, byte(val))
		}

	case "hex", "hexadecimal", "heksadesimal":
		// Remove common prefixes
		clean = strings.ReplaceAll(clean, "\\x", " ")
		clean = strings.ReplaceAll(clean, "0x", " ")
		clean = strings.ReplaceAll(clean, ",", " ")
		fields := strings.Fields(clean)
		if len(fields) == 1 && len(fields[0])%2 == 0 && len(fields[0]) > 2 {
			// Chunk into 2-character hex slices
			s := fields[0]
			fields = nil
			for i := 0; i < len(s); i += 2 {
				fields = append(fields, s[i:i+2])
			}
		}
		for _, f := range fields {
			val, err := strconv.ParseUint(f, 16, 8)
			if err != nil {
				return "", fmt.Errorf("kode hex tidak valid: %q", f)
			}
			bytes = append(bytes, byte(val))
		}

	case "octal", "oct", "oktal":
		clean = strings.ReplaceAll(clean, ",", " ")
		fields := strings.Fields(clean)
		for _, f := range fields {
			val, err := strconv.ParseUint(f, 8, 8)
			if err != nil {
				return "", fmt.Errorf("kode oktal tidak valid: %q", f)
			}
			bytes = append(bytes, byte(val))
		}

	default: // Decimal
		clean = strings.ReplaceAll(clean, ",", " ")
		fields := strings.Fields(clean)
		for _, f := range fields {
			val, err := strconv.Atoi(f)
			if err != nil || val < 0 || val > 255 {
				return "", fmt.Errorf("angka desimal tidak valid (harus 0-255): %q", f)
			}
			bytes = append(bytes, byte(val))
		}
	}

	return string(bytes), nil
}
