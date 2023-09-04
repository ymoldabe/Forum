package store

import (
	"database/sql"
	"git/ymoldabe/forum/configs"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// NewSqlite3 получает пути из configs и подкючается к DB
func NewSqlite3(config configs.Config) (*sql.DB, error) {
	// Подключение к DB
	db, err := sql.Open(config.DB.Driver, config.DB.Dsn)
	if err != nil {
		return nil, err
	}

	// Проверка соеденений с DB
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Создает таблицы
	if err = CreateTables(db, config); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateTables создает таблицы из миграций
func CreateTables(db *sql.DB, config configs.Config) error {
	file, err := os.ReadFile(config.Migrate)
	if err != nil {
		return err
	}
	requests := strings.Split(string(file), ";")
	for _, request := range requests {
		_, err := db.Exec(request)
		if err != nil {
			return err
		}
	}
	return err
}
