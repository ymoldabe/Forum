package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword []byte
	Crated         time.Time
}

type UserModel struct {
	DB *sql.DB
}

// var (
// ErrDuplicateEmail = errors.New("duplicate email")
// )
// func (m *UserModel) Insert(name, email, password string) error {
// 	hashPass, err := generateSHA256Hash(password)
// 	if err != nil {
// 		return err
// 	}

// 	stmt := `INSERT INTO users (name, email, hashed_password, created)
// 			 VALUES(?, ?, ?, ?)`

// 	_, err = m.DB.Exec(stmt, name, email, string(hashPass), time.Now().UTC())
// 	if err != nil {
// 		if strings.Contains(err.Error(), "users_uc_email") {
// 			return ErrDuplicateEmail
// 		}
// 		return err
// 	}
// 	return nil
// }

func (m *UserModel) Insert(name, email, password string) error {
	hashPass, err := generateSHA256Hash(password)
	fmt.Println(hashPass)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
			 VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(hashPass))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func generateSHA256Hash(password string) (string, error) {
	// Преобразуем пароль в срез байтов, так как SHA-256 ожидает входные данные в таком формате
	passwordBytes := []byte(password)

	// Создаем новый объект хеша SHA-256
	hasher := sha256.New()

	// Записываем пароль в хеш-функцию
	_, err := hasher.Write(passwordBytes)
	if err != nil {
		return "", err
	}

	// Получаем окончательный хеш в срезе байтов
	hashBytes := hasher.Sum(nil)

	// Преобразуем хеш в строку в формате шестнадцатеричного представления
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
