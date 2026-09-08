package reference

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Command struct {
	Category    string `json:"category,omitempty"`
	Command     string `json:"command"`
	Description string `json:"description"`
	Example     string `json:"example"`
}

type CommandCategory struct {
	Name     string    `json:"name"`
	Commands []Command `json:"commands"`
}

type CommandData struct {
	Categories []CommandCategory `json:"categories"`
}

type Port struct {
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	Service     string `json:"service"`
	Description string `json:"description"`
}

type PortData struct {
	Ports []Port `json:"ports"`
}

type HTTPCode struct {
	Code        int    `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Cause       string `json:"cause"`
}

type HTTPCodeData struct {
	Codes []HTTPCode `json:"codes"`
}

// readFileWithFallback attempts to locate the data file relative to current working directory
func readFileWithFallback(filename string) ([]byte, error) {
	candidates := []string{
		filepath.Join("data", filename),
		filepath.Join("..", "data", filename),
		filepath.Join("..", "..", "data", filename),
	}

	for _, path := range candidates {
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("could not locate data file %q in candidate paths", filename)
}

// LoadCommands loads command cheat sheets from commands.json
func LoadCommands() (CommandData, error) {
	var result CommandData
	data, err := readFileWithFallback("commands.json")
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal commands.json: %w", err)
	}

	// Propagate category name to individual commands for easy flat filtering
	for i := range result.Categories {
		catName := result.Categories[i].Name
		for j := range result.Categories[i].Commands {
			result.Categories[i].Commands[j].Category = catName
		}
	}

	return result, nil
}

// LoadPorts loads network ports from ports.json
func LoadPorts() (PortData, error) {
	var result PortData
	data, err := readFileWithFallback("ports.json")
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal ports.json: %w", err)
	}
	return result, nil
}

// LoadHTTPCodes loads HTTP status code glossary from httpcodes.json
func LoadHTTPCodes() (HTTPCodeData, error) {
	var result HTTPCodeData
	data, err := readFileWithFallback("httpcodes.json")
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal httpcodes.json: %w", err)
	}
	return result, nil
}

// SearchCommands filters commands across categories matching the search query
func SearchCommands(query string) []Command {
	data, err := LoadCommands()
	if err != nil {
		return nil
	}

	q := strings.ToLower(strings.TrimSpace(query))
	var matches []Command

	for _, cat := range data.Categories {
		for _, cmd := range cat.Commands {
			if q == "" ||
				strings.Contains(strings.ToLower(cat.Name), q) ||
				strings.Contains(strings.ToLower(cmd.Command), q) ||
				strings.Contains(strings.ToLower(cmd.Description), q) ||
				strings.Contains(strings.ToLower(cmd.Example), q) {
				matches = append(matches, cmd)
			}
		}
	}
	return matches
}

// SearchPorts filters ports by number, protocol, service, or description
func SearchPorts(query string) []Port {
	data, err := LoadPorts()
	if err != nil {
		return nil
	}

	q := strings.ToLower(strings.TrimSpace(query))
	var matches []Port

	for _, p := range data.Ports {
		portStr := strconv.Itoa(p.Port)
		if q == "" ||
			strings.Contains(portStr, q) ||
			strings.Contains(strings.ToLower(p.Protocol), q) ||
			strings.Contains(strings.ToLower(p.Service), q) ||
			strings.Contains(strings.ToLower(p.Description), q) {
			matches = append(matches, p)
		}
	}
	return matches
}

// SearchHTTPCodes filters HTTP status codes by code number, name, description, or cause
func SearchHTTPCodes(query string) []HTTPCode {
	data, err := LoadHTTPCodes()
	if err != nil {
		return nil
	}

	q := strings.ToLower(strings.TrimSpace(query))
	var matches []HTTPCode

	for _, c := range data.Codes {
		codeStr := strconv.Itoa(c.Code)
		if q == "" ||
			strings.Contains(codeStr, q) ||
			strings.Contains(strings.ToLower(c.Name), q) ||
			strings.Contains(strings.ToLower(c.Description), q) ||
			strings.Contains(strings.ToLower(c.Cause), q) {
			matches = append(matches, c)
		}
	}
	return matches
}
