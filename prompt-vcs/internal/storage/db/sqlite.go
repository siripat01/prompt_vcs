package storage

import (
	"database/sql"
	"encoding/json"
	"time"

	"prompt-vcs/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.init(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) init() error {
	schema := `
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS prompts (
		id TEXT PRIMARY KEY,
		version INTEGER DEFAULT 1,
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		variables TEXT,
		tags TEXT,
		collection TEXT DEFAULT '',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS version_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		prompt_id TEXT NOT NULL,
		version INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(prompt_id) REFERENCES prompts(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS commits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		prompt_id TEXT NOT NULL,
		hash TEXT NOT NULL UNIQUE,
		message TEXT,
		author TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(prompt_id) REFERENCES prompts(id)
	);

	CREATE TABLE IF NOT EXISTS collections (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) CreatePrompt(p *models.Prompt) error {
	_, err := db.conn.Exec(
		`INSERT INTO prompts (id, name, content, variables, tags, collection, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Content,
		encodeJSON(p.Variables), encodeJSON(p.Tags),
		p.Collection, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (db *DB) GetPrompt(id string) (*models.Prompt, error) {
	row := db.conn.QueryRow(
		`SELECT id, version, name, content, variables, tags, collection, created_at, updated_at
		FROM prompts WHERE id = ?`, id,
	)

	p := &models.Prompt{}
	var varsJSON, tagsJSON string
	err := row.Scan(&p.ID, &p.Version, &p.Name, &p.Content,
		&varsJSON, &tagsJSON, &p.Collection, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	p.Variables = decodeJSON(varsJSON)
	p.Tags = decodeJSON(tagsJSON)
	return p, nil
}

func (db *DB) UpdatePromptVersion(id string, newContent string) error {
	prompt, err := db.GetPrompt(id)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(
		`INSERT INTO version_history (prompt_id, version, content) VALUES (?, ?, ?)`,
		id, prompt.Version, prompt.Content,
	)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(
		`UPDATE prompts
		SET content = ?, version = version + 1, updated_at = ?
		WHERE id = ?`,
		newContent, time.Now(), id,
	)
	return err
}

func (db *DB) GetVersion(promptID string, version int) (string, error) {
	var content string
	err := db.conn.QueryRow(
		`SELECT content FROM version_history WHERE prompt_id = ? AND version = ?`,
		promptID, version,
	).Scan(&content)
	return content, err
}

func (db *DB) AddCommit(promptID, hash, message, author string) error {
	_, err := db.conn.Exec(
		`INSERT INTO commits (prompt_id, hash, message, author) VALUES (?, ?, ?, ?)`,
		promptID, hash, message, author,
	)
	return err
}

func (db *DB) GetCommits(promptID string) ([]models.Commit, error) {
	rows, err := db.conn.Query(
		`SELECT hash, message, author, timestamp FROM commits WHERE prompt_id = ? ORDER BY timestamp DESC`,
		promptID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commits []models.Commit
	for rows.Next() {
		var c models.Commit
		err := rows.Scan(&c.Hash, &c.Message, &c.Author, &c.Timestamp)
		if err != nil {
			return nil, err
		}
		commits = append(commits, c)
	}
	return commits, nil
}

func (db *DB) ListPrompts(collection string) ([]models.Prompt, error) {
	rows, err := db.conn.Query(
		`SELECT id, version, name, content, variables, tags, collection, created_at, updated_at
		FROM prompts WHERE collection = ? ORDER BY updated_at DESC`,
		collection,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []models.Prompt
	for rows.Next() {
		var p models.Prompt
		var varsJSON, tagsJSON string
		err := rows.Scan(&p.ID, &p.Version, &p.Name, &p.Content,
			&varsJSON, &tagsJSON, &p.Collection, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		p.Variables = decodeJSON(varsJSON)
		p.Tags = decodeJSON(tagsJSON)
		prompts = append(prompts, p)
	}
	return prompts, nil
}

func (db *DB) DeletePrompt(id string) error {
	_, err := db.conn.Exec(`DELETE FROM prompts WHERE id = ?`, id)
	return err
}

func encodeJSON(data []string) string {
	b, err := json.Marshal(data)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func decodeJSON(data string) []string {
	var result []string
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return []string{}
	}
	return result
}

func (db *DB) Close() error {
	return db.conn.Close()
}
