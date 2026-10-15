package components_test

import (
	"testing"

	"github.com/yudz/it-toolbox/ui/components"
	"github.com/yudz/it-toolbox/ui/constants"
)

func TestGetAllTools(t *testing.T) {
	tools := components.GetAllTools()
	if len(tools) < 20 {
		t.Fatalf("expected at least 20 tools, got %d", len(tools))
	}

	// Pastikan tidak ada ID duplikat
	ids := make(map[string]bool)
	for _, tool := range tools {
		if tool.ID == "" {
			t.Errorf("tool with title %s has empty ID", tool.Title)
		}
		if ids[tool.ID] {
			t.Errorf("duplicate tool ID detected: %s", tool.ID)
		}
		ids[tool.ID] = true

		if tool.PageName == "" {
			t.Errorf("tool %s has empty PageName", tool.ID)
		}
	}
}

func TestSearchTools(t *testing.T) {
	// Query kosong harus mengembalikan daftar tool populer
	popular := components.SearchTools("")
	if len(popular) == 0 {
		t.Fatalf("expected popular tools for empty query, got 0")
	}

	// Pencarian "subnet" harus memuat Subnet & IP Calculator
	subnetResults := components.SearchTools("subnet")
	if len(subnetResults) == 0 {
		t.Fatalf("expected search 'subnet' to find items, got 0")
	}
	foundSubnet := false
	for _, item := range subnetResults {
		if item.PageName == constants.NavToolbox && item.TabName == constants.TabSubnet {
			foundSubnet = true
			break
		}
	}
	if !foundSubnet {
		t.Errorf("expected to find subnet tool in search results")
	}

	// Pencarian "decrypt" harus memuat Key Decryptor dan MD5 Decrypt
	decryptResults := components.SearchTools("decrypt")
	if len(decryptResults) < 2 {
		t.Fatalf("expected at least 2 decrypt tools, got %d", len(decryptResults))
	}

	// Pencarian "cisco"
	ciscoResults := components.SearchTools("cisco")
	if len(ciscoResults) == 0 {
		t.Fatalf("expected search 'cisco' to find items")
	}

	// Pencarian query tidak cocok
	noneResults := components.SearchTools("xyz987nonexistent")
	if len(noneResults) != 0 {
		t.Fatalf("expected 0 results for nonsense query, got %d", len(noneResults))
	}
}
