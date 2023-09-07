package handler

import (
	"errors"
	"git/ymoldabe/forum/models"
	"net/http"
	"strconv"
)

// postView обрабатывает GET-запросы для просмотра поста на форуме.
func (h *Handler) postView(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса - GET.
	if r.Method != http.MethodGet {
		h.ClientError(w, http.StatusBadRequest)
		return
	}

	// Захватываем мьютекс для безопасного доступа к userUrlBefore.
	mu.Lock()
	userUrlBefore[r.RemoteAddr] = r.URL.Path
	mu.Unlock()

	// Получаем идентификатор поста из параметра запроса.
	urlId := r.URL.Query().Get("id")

	// Преобразуем идентификатор в целое число.
	id, err := strconv.Atoi(urlId)
	if err != nil || !CheckId(urlId, id) {
		h.NotFound(w)
		return
	}

	// Получаем данные поста из сервиса.
	form, err := h.service.GetPost(id)
	if err != nil {
		// Обрабатываем различные ошибки, которые могут возникнуть при запросе данных поста.
		if errors.Is(err, models.ErrNoRecord) || errors.Is(err, models.ErrNoRowsInResSet) {
			h.NotFound(w)
			return
		} else {
			h.ServerError(w, err)
			return
		}
	}

	// Создаем структуру данных для шаблона и передаем данные поста.
	data := h.NewTemplateData(r)
	data.Form = form

	// Рендерим шаблон и отправляем клиенту.
	h.render(w, r, http.StatusOK, "view.tmpl", data)
}
