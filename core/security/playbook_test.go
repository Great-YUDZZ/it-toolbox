package security

import (
	"testing"
)

func TestGetAllPlaybooks(t *testing.T) {
	items := GetAllPlaybooks()
	if len(items) == 0 {
		t.Fatal("expected playbook items, got 0")
	}

	categories := GetPlaybookCategories()
	if len(categories) < 4 {
		t.Errorf("expected at least 4 categories, got %d", len(categories))
	}
}

func TestFilterPlaybooks(t *testing.T) {
	// Search SSH
	resSSH := FilterPlaybooks("SSH", "Semua")
	if len(resSSH) == 0 {
		t.Fatal("expected to find SSH playbook items")
	}

	foundSSH01 := false
	for _, p := range resSSH {
		if p.ID == "LNX-01" {
			foundSSH01 = true
			break
		}
	}
	if !foundSSH01 {
		t.Errorf("expected LNX-01 in SSH results, got %v", resSSH)
	}

	// Filter by category
	resWin := FilterPlaybooks("", "Windows Security")
	if len(resWin) == 0 {
		t.Fatal("expected Windows Security playbooks")
	}
	for _, p := range resWin {
		if p.Category != "Windows Security" {
			t.Errorf("expected category Windows Security, got %s", p.Category)
		}
	}
}
