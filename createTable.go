package main

// import (
// 	"database/sql"
// 	"log"

// 	_ "github.com/mattn/go-sqlite3"
// )

// type People struct {
// 	id         int
// 	first_name string
// 	last_name  string
// }

// func main() {
// 	db, err := sql.Open("sqlite3", "./DB/forum.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	stmt, err := db.Prepare(`CREATE TABLE users (
// 		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
// 		name VARCHAR(255) NOT NULL,
// 		email VARCHAR(255) NOT NULL,
// 		hashed_password CHAR(60) NOT NULL,
// 		created DATETIME NOT NULL
// 	);
// 	`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = stmt.Exec()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }
