package handler

import "net/http"

// home обрабатывает запросы на главную страницу.
func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	// Проверяем, если путь запроса не является корневым, вызываем NotFound.
	if r.URL.Path != "/" {
		h.NotFound(w)
		return
	}

	// Захватываем мьютекс для безопасного доступа к userUrlBefore.
	mu.Lock()
	userUrlBefore[r.RemoteAddr] = r.URL.Path
	mu.Unlock()

	// Получаем список постов от сервиса.
	form, err := h.service.GetPosts()
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Обрабатываем GET и POST запросы.
	switch r.Method {
	case http.MethodGet:
		// Создаем структуру данных для шаблона и передаем форму.
		data := h.NewTemplateData(r)
		data.Form = form

		// Рендерим шаблон и отправляем клиенту.
		h.render(w, r, http.StatusOK, "home.tmpl", data)

	case http.MethodPost:
		// Парсим форму из POST запроса.
		if err := r.ParseForm(); err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}
		tags := r.Form["tags"]

		// Фильтруем посты по выбранным тегам.
		searchForm := getTags(tags, form)

		// Создаем структуру данных для шаблона и передаем отфильтрованную форму.
		data := h.NewTemplateData(r)
		data.Form = searchForm

		// Рендерим шаблон и отправляем клиенту.
		h.render(w, r, http.StatusOK, "home.tmpl", data)

	default:
		// Если метод запроса не поддерживается, отправляем ошибку клиенту.
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
}

// MyPosts обрабатывает запросы к странице "Мои посты".
func (h *Handler) MyPosts(w http.ResponseWriter, r *http.Request) {
	// Аналогичные комментарии могут быть добавлены к функции MyPosts и MyLikes.

	// Захватываем мьютекс для безопасного доступа к userUrlBefore.
	mu.Lock()
	userUrlBefore[r.RemoteAddr] = r.URL.Path
	mu.Unlock()

	// Получаем идентификатор пользователя из аутентификации.
	user_id, err := h.GetUserAuth(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Получаем форму постов, созданных пользователем.
	form, err := h.service.GetMyCreatedPost(user_id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Обрабатываем GET и POST запросы аналогично home.
	switch r.Method {
	case http.MethodGet:
		data := h.NewTemplateData(r)
		data.Form = form

		h.render(w, r, http.StatusOK, "home.tmpl", data)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}

		tags := r.Form["tags"]

		searchForm := getTags(tags, form)

		data := h.NewTemplateData(r)
		data.Form = searchForm

		h.render(w, r, http.StatusOK, "home.tmpl", data)

	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
}

// MyLikes обрабатывает запросы к странице "Мои лайки".
func (h *Handler) MyLikes(w http.ResponseWriter, r *http.Request) {
	// Аналогичные комментарии могут быть добавлены к функции MyPosts и MyLikes.

	// Захватываем мьютекс для безопасного доступа к userUrlBefore.
	mu.Lock()
	userUrlBefore[r.RemoteAddr] = r.URL.Path
	mu.Unlock()

	// Получаем идентификатор пользователя из аутентификации.
	user_id, err := h.GetUserAuth(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Получаем форму постов, которые пользователь лайкнул.
	form, err := h.service.GetMyLikesPost(user_id)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Обрабатываем GET и POST запросы аналогично home.
	switch r.Method {
	case http.MethodGet:
		data := h.NewTemplateData(r)
		data.Form = form

		h.render(w, r, http.StatusOK, "home.tmpl", data)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}
		tags := r.Form["tags"]

		searchForm := getTags(tags, form)

		data := h.NewTemplateData(r)
		data.Form = searchForm

		h.render(w, r, http.StatusOK, "home.tmpl", data)

	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
}
