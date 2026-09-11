package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schemaDDL = `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS error_logbook (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    title         TEXT NOT NULL,
    error_message TEXT,
    solution      TEXT,
    tags          TEXT,
    created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS snippets (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    content     TEXT NOT NULL,
    language    TEXT,
    description TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS checklists (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS checklist_items (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    checklist_id INTEGER,
    item         TEXT NOT NULL,
    checked      BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (checklist_id) REFERENCES checklists(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tasks (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    description TEXT,
    status      TEXT DEFAULT 'todo',
    category    TEXT,
    deadline    TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS projects (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    title        TEXT NOT NULL,
    description  TEXT,
    technologies TEXT,
    repo_url     TEXT,
    live_url     TEXT,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cisco_topologies (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    description TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cisco_topology_steps (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    topology_id  INTEGER NOT NULL,
    step_number  INTEGER NOT NULL,
    title        TEXT NOT NULL,
    detail       TEXT,
    is_completed BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (topology_id) REFERENCES cisco_topologies(id) ON DELETE CASCADE
);
`

// InitDB initializes SQLite database connection and runs migration schemas
func InitDB() (*sql.DB, error) {
	if DB != nil {
		return DB, nil
	}

	homeDir, err := os.UserHomeDir()
	var dbPath string
	if err == nil {
		appDir := filepath.Join(homeDir, ".it-toolbox")
		_ = os.MkdirAll(appDir, 0755)
		dbPath = filepath.Join(appDir, "it-toolbox.db")
	} else {
		dbPath = "it-toolbox.db"
	}

	return InitDBWithPath(dbPath)
}

// InitDBWithPath initializes the database with a custom path (useful for testing e.g. :memory:)
func InitDBWithPath(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if _, err := db.Exec(schemaDDL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to execute database migrations: %w", err)
	}

	DB = db
	return db, nil
}
