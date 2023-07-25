package models

import (
	"database/sql"
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
	return nil, nil
}

func (m *ForumModel) Latest() ([]*Forum, error) {
	return nil, nil
}
