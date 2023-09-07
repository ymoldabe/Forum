package handler

import (
	"bytes"
	"fmt"
	"git/ymoldabe/forum/pkg"
	"log"
	"net/http"
)

// render отвечает за рендеринг HTML-шаблона и отправку его клиенту.
func (h *Handler) render(w http.ResponseWriter, r *http.Request, status int, page string, data *TemplateData) {
	// Создаем кэш шаблонов.
	tmplCache, err := pkg.NewTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// Извлекаем шаблон из кэша по его имени.
	ts, ok := tmplCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.ServerError(w, err)
		return
	}

	// Создаем буфер для записи рендеринга.
	buf := new(bytes.Buffer)

	// Выполняем рендеринг шаблона и записываем результат в буфер.
	err = ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// Устанавливаем статус ответа и отправляем рендеринг клиенту.
	w.WriteHeader(status)
	buf.WriteTo(w)
}
