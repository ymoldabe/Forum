package handler

import (
	"errors"
	"git/ymoldabe/forum/models"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// NewCookieFile создает новую куку (cookie) для аутентификации пользователя.
func (h *Handler) NewCookieFile(w http.ResponseWriter, r *http.Request, UserId int) {
	// Генерируем уникальный идентификатор сессии.
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(1 * time.Hour)

	// Проверяем, есть ли у пользователя уже существующая сессия.
	ok, err := h.service.CheckSessions(UserId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			// Ничего не делаем, так как это означает, что у пользователя еще нет сессии.
		} else {
			h.ServerError(w, err)
			return
		}
	}
	if !ok {
		// Удаляем существующую куку (если есть).
		http.SetCookie(w, &http.Cookie{
			Value:   "session_token",
			Name:    "",
			Expires: time.Now(),
		})
		// Обновляем токен в базе данных.
		if err := h.service.UpdateToken(sessionToken, UserId); err != nil {
			h.ServerError(w, err)
			return
		}
	}

	// Добавляем новую сессию пользователя в базу данных.
	if err := h.service.UserSessionsAdd(UserId, sessionToken, expiresAt); err != nil {
		h.ServerError(w, err)
		return
	}

	// Устанавливаем куку с токеном сессии для пользователя.
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
	})
}

// CloseCookieFile закрывает куку (cookie) сессии пользователя.
func (h *Handler) CloseCookieFile(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie("session_token")
	if err != nil {
		return err
	}
	sessionToken := c.Value

	// Удаляем токен сессии из базы данных.
	if err := h.service.DeleteToken(sessionToken); err != nil {
		log.Println(err)
		return err
	}

	// Устанавливаем куку сессии в значении пустой строки и нулевым временем жизни.
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: false,
		MaxAge:   -1,
	})
	return nil
}

// GetUserAuth получает идентификатор аутентифицированного пользователя из куки (cookie) сессии.
func (h *Handler) GetUserAuth(r *http.Request) (int, error) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}
	sessionToken := c.Value
	return h.service.GetIdInSessions(sessionToken)
}
