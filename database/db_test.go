package database

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDB(t *testing.T) {
	// Re-initialize DB using in-memory SQLite for testing
	db, err := InitDBWithPath("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("InitDBWithPath failed: %v", err)
	}
	DB = db
}

func TestErrorLogCRUD(t *testing.T) {
	setupTestDB(t)

	// Create
	err := CreateLog(ErrorLog{
		Title:        "Null Pointer Exception",
		ErrorMessage: "runtime error: invalid memory address",
		Solution:     "Check if object is nil before accessing fields",
		Tags:         "go,debug",
	})
	if err != nil {
		t.Fatalf("CreateLog error: %v", err)
	}

	// Read
	logs, err := GetAllLogs()
	if err != nil {
		t.Fatalf("GetAllLogs error: %v", err)
	}
	if len(logs) == 0 {
		t.Fatalf("Expected at least 1 log, got 0")
	}

	first := logs[0]
	if first.Title != "Null Pointer Exception" {
		t.Errorf("Expected title 'Null Pointer Exception', got %q", first.Title)
	}

	// Update
	first.Title = "Fixed NPE"
	if err := UpdateLog(first); err != nil {
		t.Fatalf("UpdateLog error: %v", err)
	}

	// Search
	results, err := SearchLogs("Fixed")
	if err != nil || len(results) == 0 {
		t.Fatalf("SearchLogs('Fixed') failed: %v, count: %d", err, len(results))
	}

	// Delete
	if err := DeleteLog(first.ID); err != nil {
		t.Fatalf("DeleteLog error: %v", err)
	}
}

func TestSnippetCRUD(t *testing.T) {
	setupTestDB(t)

	err := CreateSnippet(Snippet{
		Title:       "HTTP Get Helper",
		Content:     "resp, err := http.Get(url)",
		Language:    "go",
		Description: "Simple GET request snippet",
	})
	if err != nil {
		t.Fatalf("CreateSnippet error: %v", err)
	}

	snippets, err := GetAllSnippets()
	if err != nil || len(snippets) == 0 {
		t.Fatalf("GetAllSnippets error: %v", err)
	}

	s := snippets[0]
	s.Title = "Updated HTTP Get"
	if err := UpdateSnippet(s); err != nil {
		t.Fatalf("UpdateSnippet error: %v", err)
	}

	if err := DeleteSnippet(s.ID); err != nil {
		t.Fatalf("DeleteSnippet error: %v", err)
	}
}

func TestChecklistCRUD(t *testing.T) {
	setupTestDB(t)

	if err := CreateChecklist("Pre-flight Deployment"); err != nil {
		t.Fatalf("CreateChecklist error: %v", err)
	}

	lists, err := GetChecklistsWithItems()
	if err != nil || len(lists) == 0 {
		t.Fatalf("GetChecklistsWithItems error: %v", err)
	}

	clID := lists[0].ID
	if err := AddItem(clID, "Run tests"); err != nil {
		t.Fatalf("AddItem error: %v", err)
	}

	// Re-fetch to check item
	lists, _ = GetChecklistsWithItems()
	if len(lists[0].Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(lists[0].Items))
	}

	itemID := lists[0].Items[0].ID
	if err := ToggleItem(itemID); err != nil {
		t.Fatalf("ToggleItem error: %v", err)
	}

	lists, _ = GetChecklistsWithItems()
	if !lists[0].Items[0].Checked {
		t.Errorf("Expected item to be checked after toggle")
	}

	if err := DeleteChecklist(clID); err != nil {
		t.Fatalf("DeleteChecklist error: %v", err)
	}
}

func TestTaskAndProjectCRUD(t *testing.T) {
	setupTestDB(t)

	// Task
	err := CreateTask(Task{
		Title:       "Praktikum Jaringan 1",
		Description: "Subnetting Cisco Packet Tracer",
		Status:      "todo",
		Category:    "praktikum",
		Deadline:    "2026-09-15",
	})
	if err != nil {
		t.Fatalf("CreateTask error: %v", err)
	}

	tasks, err := GetAllTasks()
	if err != nil || len(tasks) == 0 {
		t.Fatalf("GetAllTasks error: %v", err)
	}

	if err := UpdateTaskStatus(tasks[0].ID, "inprogress"); err != nil {
		t.Fatalf("UpdateTaskStatus error: %v", err)
	}

	// Project & Portfolio Export
	err = CreateProject(Project{
		Title:        "IT Toolbox",
		Description:  "Desktop toolkit for IT students in Go + Fyne",
		Technologies: "Go, Fyne, SQLite",
		RepoURL:      "https://github.com/yudz/it-toolbox",
		LiveURL:      "",
	})
	if err != nil {
		t.Fatalf("CreateProject error: %v", err)
	}

	tmpDir := t.TempDir()
	exportPath := filepath.Join(tmpDir, "portfolio.md")
	if err := ExportPortfolioMarkdown(exportPath); err != nil {
		t.Fatalf("ExportPortfolioMarkdown failed: %v", err)
	}

	data, err := os.ReadFile(exportPath)
	if err != nil || len(data) == 0 {
		t.Fatalf("Exported portfolio file is empty or missing: %v", err)
	}
}
