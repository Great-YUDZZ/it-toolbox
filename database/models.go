package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrDBNotInitialized = errors.New("database is not initialized, please call InitDB first")
)

func getDB() (*sql.DB, error) {
	if DB == nil {
		return InitDB()
	}
	return DB, nil
}

// ----------------------------------------------------------------------------
// 1. Error Logbook
// ----------------------------------------------------------------------------

type ErrorLog struct {
	ID           int
	Title        string
	ErrorMessage string
	Solution     string
	Tags         string
	CreatedAt    string
}

func CreateLog(log ErrorLog) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	query := `INSERT INTO error_logbook (title, error_message, solution, tags) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, log.Title, log.ErrorMessage, log.Solution, log.Tags)
	return err
}

func GetAllLogs() ([]ErrorLog, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT id, title, error_message, solution, tags, created_at FROM error_logbook ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ErrorLog
	for rows.Next() {
		var l ErrorLog
		if err := rows.Scan(&l.ID, &l.Title, &l.ErrorMessage, &l.Solution, &l.Tags, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func UpdateLog(log ErrorLog) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	query := `UPDATE error_logbook SET title = ?, error_message = ?, solution = ?, tags = ? WHERE id = ?`
	_, err = db.Exec(query, log.Title, log.ErrorMessage, log.Solution, log.Tags, log.ID)
	return err
}

func DeleteLog(id int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM error_logbook WHERE id = ?`, id)
	return err
}

func SearchLogs(query string) ([]ErrorLog, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	q := "%" + strings.TrimSpace(query) + "%"
	stmt := `SELECT id, title, error_message, solution, tags, created_at FROM error_logbook 
             WHERE title LIKE ? OR error_message LIKE ? OR solution LIKE ? OR tags LIKE ? ORDER BY id DESC`
	rows, err := db.Query(stmt, q, q, q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ErrorLog
	for rows.Next() {
		var l ErrorLog
		if err := rows.Scan(&l.ID, &l.Title, &l.ErrorMessage, &l.Solution, &l.Tags, &l.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

// ----------------------------------------------------------------------------
// 2. Code Snippets
// ----------------------------------------------------------------------------

type Snippet struct {
	ID          int
	Title       string
	Content     string
	Language    string
	Description string
	CreatedAt   string
}

func CreateSnippet(snippet Snippet) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	query := `INSERT INTO snippets (title, content, language, description) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, snippet.Title, snippet.Content, snippet.Language, snippet.Description)
	return err
}

func GetAllSnippets() ([]Snippet, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT id, title, content, language, description, created_at FROM snippets ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Snippet
	for rows.Next() {
		var s Snippet
		if err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Language, &s.Description, &s.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

func UpdateSnippet(snippet Snippet) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	query := `UPDATE snippets SET title = ?, content = ?, language = ?, description = ? WHERE id = ?`
	_, err = db.Exec(query, snippet.Title, snippet.Content, snippet.Language, snippet.Description, snippet.ID)
	return err
}

func DeleteSnippet(id int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM snippets WHERE id = ?`, id)
	return err
}

// ----------------------------------------------------------------------------
// 3. Checklists
// ----------------------------------------------------------------------------

type ChecklistItem struct {
	ID          int
	ChecklistID int
	Item        string
	Checked     bool
}

type Checklist struct {
	ID        int
	Name      string
	CreatedAt string
	Items     []ChecklistItem
}

func CreateChecklist(name string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO checklists (name) VALUES (?)`, strings.TrimSpace(name))
	return err
}

func AddItem(checklistID int, item string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO checklist_items (checklist_id, item, checked) VALUES (?, ?, FALSE)`, checklistID, strings.TrimSpace(item))
	return err
}

func ToggleItem(itemID int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE checklist_items SET checked = NOT checked WHERE id = ?`, itemID)
	return err
}

func DeleteChecklist(id int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM checklists WHERE id = ?`, id)
	return err
}

func DeleteChecklistItem(itemID int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM checklist_items WHERE id = ?`, itemID)
	return err
}

func GetChecklistsWithItems() ([]Checklist, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`SELECT id, name, created_at FROM checklists ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []Checklist
	for rows.Next() {
		var cl Checklist
		if err := rows.Scan(&cl.ID, &cl.Name, &cl.CreatedAt); err != nil {
			return nil, err
		}
		lists = append(lists, cl)
	}

	// Fetch items for each checklist
	for i := range lists {
		itemRows, err := db.Query(`SELECT id, checklist_id, item, checked FROM checklist_items WHERE checklist_id = ? ORDER BY id ASC`, lists[i].ID)
		if err != nil {
			return nil, err
		}
		var items []ChecklistItem
		for itemRows.Next() {
			var it ChecklistItem
			if err := itemRows.Scan(&it.ID, &it.ChecklistID, &it.Item, &it.Checked); err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, it)
		}
		itemRows.Close()
		lists[i].Items = items
	}

	return lists, nil
}

// ----------------------------------------------------------------------------
// 4. Tasks (Kanban)
// ----------------------------------------------------------------------------

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string // todo, inprogress, done
	Category    string // praktikum, tugas, project
	Deadline    string
	CreatedAt   string
}

func CreateTask(task Task) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	if task.Status == "" {
		task.Status = "todo"
	}
	query := `INSERT INTO tasks (title, description, status, category, deadline) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(query, task.Title, task.Description, task.Status, task.Category, task.Deadline)
	return err
}

func GetAllTasks() ([]Task, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT id, title, description, status, category, deadline, created_at FROM tasks ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Category, &t.Deadline, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func UpdateTaskStatus(id int, status string) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`UPDATE tasks SET status = ? WHERE id = ?`, status, id)
	return err
}

func UpdateTask(task Task) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	query := `UPDATE tasks SET title = ?, description = ?, status = ?, category = ?, deadline = ? WHERE id = ?`
	_, err = db.Exec(query, task.Title, task.Description, task.Status, task.Category, task.Deadline, task.ID)
	return err
}

func DeleteTask(id int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

// ----------------------------------------------------------------------------
// 5. Projects (Showcase Profiler & Portfolio Exporter)
// ----------------------------------------------------------------------------

type Project struct {
	ID           int
	Title        string
	Description  string
	Technologies string
	RepoURL      string
	LiveURL      string
	CreatedAt    string
}

func CreateProject(p Project) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	query := `INSERT INTO projects (title, description, technologies, repo_url, live_url) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(query, p.Title, p.Description, p.Technologies, p.RepoURL, p.LiveURL)
	return err
}

func GetAllProjects() ([]Project, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(`SELECT id, title, description, technologies, repo_url, live_url, created_at FROM projects ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Technologies, &p.RepoURL, &p.LiveURL, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func UpdateProject(p Project) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	query := `UPDATE projects SET title = ?, description = ?, technologies = ?, repo_url = ?, live_url = ? WHERE id = ?`
	_, err = db.Exec(query, p.Title, p.Description, p.Technologies, p.RepoURL, p.LiveURL, p.ID)
	return err
}

func DeleteProject(id int) error {
	db, err := getDB()
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM projects WHERE id = ?`, id)
	return err
}

// ExportPortfolioMarkdown generates a clean GitHub/CV portfolio.md document from all projects
func ExportPortfolioMarkdown(filepath string) error {
	projects, err := GetAllProjects()
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("# Developer Portfolio\n\n")
	sb.WriteString("Koleksi proyek perangkat lunak dan portofolio teknis yang dikelola melalui **IT Toolbox**.\n\n---\n\n")

	if len(projects) == 0 {
		sb.WriteString("*Belum ada proyek yang terdaftar di portofolio.*\n")
	} else {
		for i, p := range projects {
			sb.WriteString(fmt.Sprintf("## %d. %s\n\n", i+1, p.Title))
			if p.Description != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", p.Description))
			}
			if p.Technologies != "" {
				sb.WriteString(fmt.Sprintf("- **Teknologi**: `%s`\n", p.Technologies))
			}
			if p.RepoURL != "" {
				sb.WriteString(fmt.Sprintf("- **Repository**: [%s](%s)\n", p.RepoURL, p.RepoURL))
			}
			if p.LiveURL != "" {
				sb.WriteString(fmt.Sprintf("- **Demo / Live**: [%s](%s)\n", p.LiveURL, p.LiveURL))
			}
			sb.WriteString("\n---\n\n")
		}
	}

	return os.WriteFile(filepath, []byte(sb.String()), 0644)
}
