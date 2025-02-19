package cache

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type SQLiteCache struct {
	db *sql.DB
}

func NewSQLiteCache(dbPath string) (*SQLiteCache, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Create table if not exists
	//_, err = db.Exec(`
	//    CREATE TABLE IF NOT EXISTS repo_descriptions (
	//        id INTEGER PRIMARY KEY AUTOINCREMENT,
	//        repo_name TEXT UNIQUE,
	//        description LONG TEXT,
	//        created_at DATETIME,
	//        updated_at DATETIME
	//    )
	//`)
	//if err != nil {
	//	return nil, fmt.Errorf("error creating table: %w", err)
	//}

	return &SQLiteCache{db: db}, nil
}

func (c *SQLiteCache) GetDescription(repoName string) (string, error) {
	var description string
	err := c.db.QueryRow(
		"SELECT description FROM repo_descriptions WHERE repo_name = ?",
		repoName,
	).Scan(&description)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("error getting description: %w", err)
	}

	return description, nil
}

func (c *SQLiteCache) SetDescripton(repoName, description string) error {
	now := time.Now()
	_, err := c.db.Exec(`
        INSERT INTO repo_descriptions (repo_name, description, created_at, updated_at)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(repo_name) DO UPDATE SET
            description = excluded.description,
            updated_at = excluded.updated_at
    `, repoName, description, now, now)

	if err != nil {
		return fmt.Errorf("error setting description: %w", err)
	}

	return nil
}

func (c *SQLiteCache) Close() error {
	return c.db.Close()
}
