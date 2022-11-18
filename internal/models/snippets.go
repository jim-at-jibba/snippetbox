package models

import (
	"database/sql"
	"errors"
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
	stmt := `INSERT INTO snippets (title, content, created, expires)
  VALUES(?,?,UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// Get snippet by id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
  WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	// Initialise a pointer to a new zeroes Snippet struct
	s := &Snippet{}

	// Notice the arguments to row.Scan are *pointers* to the place you want
	// to copy data into, the number of augments must be exactly the same as the
	// number of columns returned
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Create, &s.Expires)
	if err != nil {
		//If the query returns no rows, then row.Scan will return sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// We return our own error here to help encapsulate the model. This means
			// our app is not concerned with datastore specific errors
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// Latest - get last 10 snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
