package models

import (
	"database/sql"
	"time"
)

// Snippet type
type Snippet struct {
	ID      int
	Title   string
	Content string
	Create  time.Time
	Expires time.Time
}

// SnippetModel a SnippetsModel type which wraps sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert snippet
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// Get snippet by id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Latest - get last 10 snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
