package store_auth

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"git/ymoldabe/forum/models"
	"git/ymoldabe/forum/pkg"

	"golang.org/x/crypto/bcrypt"
)

// Auth структура для реализации интерфейса Autorization с использованием SQLite базы данных.
type Auth struct {
	db *sql.DB
}

// New создает новый экземпляр AuthSqlite с переданной базой данных.
func New(db *sql.DB) *Auth {
	return &Auth{
		db: db,
	}
}

// InsertUser вставляет нового пользователя в базу данных.
func (a *Auth) InsertUser(form *models.UserSignUp) error {
	stmt := `INSERT INTO users (name, email, hash_password, provider)
	VALUES(?, ?, ?, ?)`
	_, err := a.db.Exec(stmt, form.Name, form.Email, form.HashPassword, form.Provider)
	if err != nil {
		if strings.Contains(models.ProviderGoogle, form.Provider) {
			return models.ErrProviderGoogle
		} else if strings.Contains(models.ProviderGit, form.Provider) {
			return models.ErrProviderGithub
		}

		if err.Error() == "UNIQUE constraint failed: users.email" {
			return models.ErrDuplicateEmail
		}
		return err
	}
	return err
}

// Authenticate выполняет аутентификацию пользователя.
func (a *Auth) Authenticate(form *models.UserSignIn) (int, error) {
	var user models.UserSignIn
	stmt := `SELECT id, hash_password, provider FROM users WHERE email = ?;`
	err := a.db.QueryRow(stmt, form.Email).Scan(&user.Id, &user.HashPassword, &user.Provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	if strings.Contains(models.ProviderGoogle, user.Provider) {
		return 0, models.ErrProviderGoogle
	} else if strings.Contains(models.ProviderGit, user.Provider) {
		return 0, models.ErrProviderGithub
	}

	if err = pkg.CheckPasswordHash(form.Password, user.HashPassword); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return user.Id, nil
}

// UserSessionsAdd добавляет сессию пользователя в базу данных.
func (a *Auth) UserSessionsAdd(userId int, sessionToken string, expiresAt time.Time) error {
	stmt := `
	INSERT INTO sessions (user_id, session_token, expires_at)
	VALUES (?, ?, ?)`
	if _, err := a.db.Exec(stmt, userId, sessionToken, expiresAt); err != nil {
		return err
	}
	return nil
}

// DeleteToken удаляет сессию пользователя по токену.
func (a *Auth) DeleteToken(sessionToken string) error {
	stmt := `
	DELETE FROM	sessions
	WHERE session_token = ?`

	res, err := a.db.Exec(stmt, sessionToken)
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
func (a *Auth) GetIdInSessions(sessionToken string) (int, error) {
	var userId int
	stmt := `
	SELECT user_id
	FROM sessions
	WHERE session_token = ?`

	if err := a.db.QueryRow(stmt, sessionToken).Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

// CheckSessions проверяет наличие сессий пользователя.
func (a *Auth) CheckSessions(userId int) (bool, error) {
	var count int
	stmt := `
	SELECT COUNT(*)
	FROM sessions 
	WHERE user_id = ?`

	checkHasData, err := a.tableHasData("sessions")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if checkHasData {
		return false, nil
	}

	if err := a.db.QueryRow(stmt, userId).Scan(&count); err != nil {
		return false, err
	}

	return count == 1, nil
}

// UpdateToken обновляет токен сессии пользователя.
func (a *Auth) UpdateToken(sessionToken string, user_id int) error {
	stmt := `
	UPDATE sessions
	SET session_token = ?,
	expires_at = ?
	WHERE user_id = ?`

	_, err := a.db.Exec(stmt, sessionToken, time.Now().Add(24*time.Hour), user_id)
	return err
}

// GetTokenSession проверяет наличие валидной сессии по токену.
func (a *Auth) GetTokenSession(cookieToken string) (bool, error) {
	var userId int
	stmt := `
	SELECT user_id 
	FROM sessions
	WHERE session_token = ? AND expires_at > ?`

	checkHasData, err := a.tableHasData("sessions")
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			return false, nil
		}
		return false, err
	}

	if !checkHasData {
		return false, nil
	}

	if err := a.db.QueryRow(stmt, cookieToken, time.Now()).Scan(&userId); err != nil {
		return false, err
	}

	if userId <= 0 {
		return false, nil
	}
	return true, nil
}

// tableHasData проверяет, содержит ли таблица данные.
func (a *Auth) tableHasData(tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT 1 FROM %s LIMIT 1;", tableName)
	var result sql.NullString
	err := a.db.QueryRow(query).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = models.ErrNoRecord
			return false, err // Таблица пуста
		}
		return false, err // Ошибка выполнения запроса
	}
	return true, nil // Таблица не пуста
}
