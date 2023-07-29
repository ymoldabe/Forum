package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./DB/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			hashed_password CHAR(60) NOT NULL,
			created DATETIME NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Таблица users создана успешно.")
}
