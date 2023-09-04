package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"git/ymoldabe/forum/models"
	"git/ymoldabe/forum/pkg"

	"golang.org/x/crypto/bcrypt"
)

// AuthSqlite структура для реализации интерфейса Autorization с использованием SQLite базы данных.
type AuthSqlite struct {
	db *sql.DB
}

// NewAuthSqlite создает новый экземпляр AuthSqlite с переданной базой данных.
func NewAuthSqlite(db *sql.DB) *AuthSqlite {
	return &AuthSqlite{
		db: db,
	}
}

// InsertUser вставляет нового пользователя в базу данных.
func (s *AuthSqlite) InsertUser(form *models.UserSignUp) error {
	stmt := `INSERT INTO users (name, email, hash_password)
	VALUES(?, ?, ?)`
	_, err := s.db.Exec(stmt, form.Name, form.Email, form.HashPassword)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return models.ErrDuplicateEmail
		}
		return err
	}
	return err
}

// Authenticate выполняет аутентификацию пользователя.
func (s *AuthSqlite) Authenticate(form *models.UserSignIn) (int, error) {
	var row models.UserSignIn
	stmt := `SELECT id,  hash_password FROM users WHERE email = ?;`
	err := s.db.QueryRow(stmt, form.Email).Scan(&row.Id, &row.HashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	if err = pkg.CheckPasswordHash(form.Password, row.HashPassword); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return row.Id, nil
}

// UserSessionsAdd добавляет сессию пользователя в базу данных.
func (s *AuthSqlite) UserSessionsAdd(userId int, sessionToken string, expiresAt time.Time) error {
	stmt := `
	INSERT INTO sessions (user_id, session_token, expires_at)
	VALUES (?, ?, ?)`
	if _, err := s.db.Exec(stmt, userId, sessionToken, expiresAt); err != nil {
		return err
	}
	return nil
}

// DeleteToken удаляет сессию пользователя по токену.
func (s *AuthSqlite) DeleteToken(sessionToken string) error {
	stmt := `
	DELETE FROM	sessions
	WHERE session_token = ?`

	res, err := s.db.Exec(stmt, sessionToken)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrDeleteFailed
	}
	return err
}

// GetIdInSessions возвращает идентификатор пользователя по токену сессии.
func (s *AuthSqlite) GetIdInSessions(sessionToken string) (int, error) {
	var userId int
	stmt := `
	SELECT user_id
	FROM sessions
	WHERE session_token = ?`

	if err := s.db.QueryRow(stmt, sessionToken).Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

// CheckSessions проверяет наличие сессий пользователя.
func (s *AuthSqlite) CheckSessions(userId int) (bool, error) {
	var count int
	stmt := `
	SELECT COUNT(*)
	FROM sessions 
	WHERE user_id = ?`

	checkHasData, err := s.tableHasData("sessions")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if checkHasData {
		return false, nil
	}

	if err := s.db.QueryRow(stmt, userId).Scan(&count); err != nil {
		return false, err
	}

	return count == 1, nil
}

// UpdateToken обновляет токен сессии пользователя.
func (s *AuthSqlite) UpdateToken(sessionToken string, user_id int) error {
	stmt := `
	UPDATE sessions
	SET session_token = ?,
	expires_at = ?
	WHERE user_id = ?`

	_, err := s.db.Exec(stmt, sessionToken, time.Now().Add(24*time.Hour), user_id)
	fmt.Println("update")
	return err
}

// GetTokenSession проверяет наличие валидной сессии по токену.
func (s *AuthSqlite) GetTokenSession(cookieToken string) (bool, error) {
	var userId int
	stmt := `
	SELECT user_id 
	FROM sessions
	WHERE session_token = ? AND expires_at > ?`

	checkHasData, err := s.tableHasData("sessions")
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return false, nil
		}
		return false, err
	}

	if !checkHasData {
		return false, nil
	}

	if err := s.db.QueryRow(stmt, cookieToken, time.Now()).Scan(&userId); err != nil {
		return false, err
	}

	if userId <= 0 {
		return false, nil
	}
	return true, nil
}

// tableHasData проверяет, содержит ли таблица данные.
func (s *AuthSqlite) tableHasData(tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s LIMIT 1;", tableName)
	var result sql.NullString
	err := s.db.QueryRow(query).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = models.ErrNoRecord
			return false, err // Таблица пуста
		}
		return false, err // Ошибка выполнения запроса
	}
	return true, nil // Таблица не пуста
}
