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
	stmt := `SELECT id, title, content, created, expires FROM snippets
  WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// defer should come AFTER we check for errors from query, if Query()
	/// returns an error it will panic!
	defer rows.Close()

	snippets := []*Snippet{}

	// Use rows.Next to iterate through the rows. This prepares the frst and
	// each subsequent tow to be acted on by rows.Scan. Frees up connection once its
	// complete
	for rows.Next() {
		// Crete a pointer to a new Zeroed Snippet struct
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Create, &s.Expires)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	// When the rows.Next() has finished we call rows.Err() to retrieve any errors
	// that were encountered
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
