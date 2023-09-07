package handler

import (
	"fmt"
	"git/ymoldabe/forum/pkg"
	"net/http"
)

// logRequest создает middleware для логирования информации о запросах.
func (h *Handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Записываем информацию о запросе в лог.
		pkg.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		// Передаем запрос дальше по цепочке middleware или обработчику.
		next.ServeHTTP(w, r)
	})
}

// recoverPanic создает middleware для обработки паник и восстановления после них.
func (h *Handler) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отлавливаем панику и восстанавливаемся.
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				// Отправляем серверную ошибку в ответ.
				h.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		// Передаем запрос дальше по цепочке middleware или обработчику.
		next.ServeHTTP(w, r)
	})
}

// requireAuthentication создает middleware для проверки аутентификации пользователя.
func (h *Handler) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, аутентифицирован ли пользователь.
		if !h.isAuthenticated(r) {
			// Если не аутентифицирован, перенаправляем на страницу входа.
			http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
			return
		}
		// Устанавливаем заголовок для запрета кэширования.
		w.Header().Add("Cache-Control", "no-store")
		// Передаем запрос дальше по цепочке middleware или обработчику.
		next.ServeHTTP(w, r)
	})
}
