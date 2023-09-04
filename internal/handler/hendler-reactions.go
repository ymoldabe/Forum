package handler

import (
	"fmt"
	"git/ymoldabe/forum/models"
	"net/http"
)

// reactionComment обрабатывает POST-запросы для реакций на комментарии к постам.
func (h *Handler) reactionComment(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса - POST.
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	// Получаем идентификатор поста, статус и ошибку из функции CheckValidIdOrCommentId для поста.
	postId, status, err := h.CheckValidIdOrCommentId(r, false)
	if status == http.StatusNotFound {
		h.NotFound(w)
		return
	} else if status == http.StatusInternalServerError {
		h.ServerError(w, err)
		return
	}

	// Получаем идентификатор комментария, статус и ошибку из функции CheckValidIdOrCommentId для комментария.
	commentId, status, err := h.CheckValidIdOrCommentId(r, true)
	if status == http.StatusNotFound {
		h.NotFound(w)
		return
	} else if status == http.StatusInternalServerError {
		h.ServerError(w, err)
		return
	}

	// Получаем идентификатор пользователя из аутентификации.
	userId, err := h.GetUserAuth(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Получаем значение реакции из формы запроса.
	reaction := r.FormValue("reaction")

	// Обрабатываем реакцию на комментарий в зависимости от значения reaction.
	switch reaction {
	case models.LikeComment, models.DislikeComment:
		// Вызываем сервис для обработки реакции на комментарий.
		if err = h.service.ReactionComment(postId, userId, commentId, reaction); err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}
	default:
		// Если значение reaction не соответствует ожидаемым значениям, возвращаем ошибку клиента.
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	// Формируем URL для перенаправления на страницу просмотра поста.
	url := fmt.Sprintf("/post/view?id=%d", postId)
	// Осуществляем перенаправление пользователя на данную страницу.
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// reaction обрабатывает POST-запросы для реакций на посты.
func (h *Handler) reaction(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса - POST.
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Получаем идентификатор поста, статус и ошибку из функции CheckValidIdOrCommentId.
	postId, status, err := h.CheckValidIdOrCommentId(r, false)
	if status == http.StatusNotFound {
		h.NotFound(w)
		return
	} else if status == http.StatusInternalServerError {
		h.ServerError(w, err)
		return
	}

	// Получаем идентификатор пользователя из аутентификации.
	userId, err := h.GetUserAuth(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Получаем значение реакции из формы запроса.
	reaction := r.FormValue("reaction")

	// Обрабатываем реакцию на пост в зависимости от значения reaction.
	switch reaction {
	case models.LikePost, models.DislikePost:
		// Вызываем сервис для обработки реакции на пост.
		if err = h.service.ReactionPost(postId, userId, reaction); err != nil {
			h.ServerError(w, err)
			return
		}
	default:
		// Если значение reaction не соответствует ожидаемым значениям, возвращаем ошибку клиента.
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	// Формируем URL для перенаправления на страницу просмотра поста.
	url := fmt.Sprintf("/post/view?id=%d", postId)
	// Осуществляем перенаправление пользователя на данную страницу.
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// reactionFromHome обрабатывает POST-запросы для реакций на посты на домашней странице.
func (h *Handler) reactionFromHome(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса - POST.
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	// Получаем идентификатор поста, статус и ошибку из функции CheckValidIdOrCommentId.
	postId, status, err := h.CheckValidIdOrCommentId(r, false)
	if status == http.StatusNotFound {
		h.NotFound(w)
		return
	} else if status == http.StatusInternalServerError {
		h.ServerError(w, err)
		return
	}

	// Получаем идентификатор пользователя из аутентификации.
	userId, err := h.GetUserAuth(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Получаем значение реакции из формы запроса.
	reaction := r.FormValue("reaction")

	// Обрабатываем реакцию на пост в зависимости от значения reaction.
	switch reaction {
	case models.LikePost:
		// Вызываем сервис для обработки лайка поста.
		if err = h.service.ReactionPost(postId, userId, reaction); err != nil {
			h.ServerError(w, err)
			return
		}

	case models.DislikePost:
		// Вызываем сервис для обработки дизлайка поста.
		if err = h.service.ReactionPost(postId, userId, reaction); err != nil {
			h.ServerError(w, err)
			return
		}
	default:
		// Если значение reaction не соответствует ожидаемым значениям, возвращаем ошибку клиента.
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	// После обработки реакции, перенаправляем пользователя на текущую страницу.
	redirect(w, r)
}
