package models

import (
	"database/sql"
	"errors"
	"time"
)

type Forum struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type ForumModel struct {
	DB *sql.DB
}

func (m *ForumModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO forum (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}
	// Возможно sqllite не поддержит
	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return int(id), nil
}

func (m *ForumModel) Get(id int) (*Forum, error) {

	stmt := `SELECT id, title, content, created, expires FROM forum
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Forum{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *ForumModel) Latest() ([]*Forum, error) {
	return nil, nil
}
