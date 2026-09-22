package pages

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/yudz/it-toolbox/core/calculators"
	"github.com/yudz/it-toolbox/database"
)

func TestCalculatorPageBuild(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	win := test.NewWindow(nil)
	defer win.Close()

	page := NewCalculatorPage(win)
	obj := page.Build()

	if obj == nil {
		t.Fatal("expected non-nil CanvasObject from CalculatorPage.Build()")
	}
}

func TestAlgorithmBadges(t *testing.T) {
	algos := []string{"bcrypt", "BCRYPT", "SHA-256", "sha256", "SHA-512", "MD5", "SHA-1", "UNKNOWN"}
	for _, a := range algos {
		badge := getAlgorithmBadge(a)
		if badge == nil {
			t.Fatalf("expected non-nil badge for algorithm: %s", a)
		}
	}
}

func TestSaveAndVerifyHashedPasswordWorkflow(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	win := test.NewWindow(nil)
	defer win.Close()

	// 1. Generate bcrypt hash
	rawPass := "SecureVaultPass2026!"
	bHash, err := calculators.GenerateBcrypt(rawPass, 10)
	if err != nil {
		t.Fatalf("failed to generate bcrypt hash: %v", err)
	}

	// 2. Save to database
	hp := &database.HashedPassword{
		Title:         "Test Vault Entry",
		Algorithm:     "bcrypt",
		HashValue:     bHash,
		PlainPassword: rawPass,
		Notes:         "Automated unit test record",
	}
	id, err := database.CreateHashedPassword(hp)
	if err != nil {
		t.Fatalf("failed to create hashed password in database: %v", err)
	}
	defer func() {
		_ = database.DeleteHashedPassword(id)
	}()

	// 3. Search database
	results, err := database.SearchHashedPasswords("Test Vault", "bcrypt")
	if err != nil {
		t.Fatalf("failed to search hashed passwords: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least 1 search result, got 0")
	}

	found := false
	for _, r := range results {
		if r.ID == id {
			found = true
			if r.Title != hp.Title {
				t.Errorf("expected title %s, got %s", hp.Title, r.Title)
			}
			if r.HashValue != hp.HashValue {
				t.Errorf("expected hash %s, got %s", hp.HashValue, r.HashValue)
			}
			break
		}
	}
	if !found {
		t.Errorf("created record ID %d not found in search results", id)
	}

	// 4. Verify password matches hash
	if !calculators.VerifyHash("bcrypt", bHash, rawPass, "") {
		t.Error("expected VerifyHash to return true for correct password")
	}
	if calculators.VerifyHash("bcrypt", bHash, "WrongPassword123", "") {
		t.Error("expected VerifyHash to return false for incorrect password")
	}
}

func TestHashGenTabMinSizeBounded(t *testing.T) {
	app := test.NewApp()
	defer app.Quit()

	win := test.NewWindow(nil)
	defer win.Close()

	page := NewCalculatorPage(win)
	tab := page.buildHashGenTab()

	// Ensure layout does not overflow beyond standard window width
	if tab.MinSize().Width > 600 {
		t.Fatalf("HashGenTab MinSize width (%f) exceeds 600px, which causes cutoffs on standard screens", tab.MinSize().Width)
	}
}
