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

// 	stmt, err := db.Prepare(`CREATE TABLE sessions (
//     token CHAR(43) PRIMARY KEY,
//     data BLOB NOT NULL,
//     expiry TIMESTAMP(6) NOT NULL
// );`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = stmt.Exec()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }
