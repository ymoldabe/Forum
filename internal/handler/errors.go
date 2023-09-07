package handler

import (
	"fmt"
	"git/ymoldabe/forum/pkg"
	"html/template"
	"net/http"
	"runtime/debug"
)

// HandlerError представляет ошибку обработчика.
type HandlerError struct {
	ErrorCode int
	ErrorMsg  string
}

// ServerError обрабатывает внутренние ошибки сервера, записывая их в логи и отображая пользователю страницу ошибки.
func (h *Handler) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	pkg.InfoLog.Output(2, trace)

	errorForm := HandlerError{
		ErrorCode: http.StatusInternalServerError,
		ErrorMsg:  http.StatusText(http.StatusInternalServerError),
	}

	tmpl, err := template.ParseFiles("./ui/html/error/error.html")
	if err != nil {
		pkg.ErrorLog.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(errorForm.ErrorCode)

	if err := tmpl.Execute(w, errorForm); err != nil {
		pkg.ErrorLog.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// ClientError обрабатывает ошибки, связанные с запросами клиентов, отображая пользователю страницу ошибки.
func (h *Handler) ClientError(w http.ResponseWriter, status int) {
	errorForm := HandlerError{
		ErrorCode: status,
		ErrorMsg:  http.StatusText(status),
	}

	tmpl, err := template.ParseFiles("./ui/html/error/error.html")
	if err != nil {
		pkg.ErrorLog.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if err := tmpl.Execute(w, errorForm); err != nil {
		pkg.ErrorLog.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// NotFound обрабатывает ошибку "страница не найдена" и отображает пользователю соответствующую страницу ошибки.
func (h *Handler) NotFound(w http.ResponseWriter) {
	h.ClientError(w, http.StatusNotFound)
}
