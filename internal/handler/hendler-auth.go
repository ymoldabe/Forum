package handler

import (
	"errors"
	"git/ymoldabe/forum/models"
	"net/http"
)

// signUp обрабатывает регистрацию пользователей.
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := h.NewTemplateData(r)
		data.Form = models.UserSignUp{}
		h.render(w, r, http.StatusOK, "signup.tmpl", data)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}

		form := models.UserSignUp{
			Name:      r.PostForm.Get("name"),
			Email:     r.PostForm.Get("email"),
			Password1: r.PostForm.Get("password1"),
			Password2: r.PostForm.Get("password2"),
		}

		err = h.service.InsertUser(&form)
		if err != nil {
			if errors.Is(err, models.ErrFormNotValid) {
				data := h.NewTemplateData(r)
				data.Form = form
				h.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
				return
			} else if errors.Is(err, models.ErrDuplicateEmail) {
				data := h.NewTemplateData(r)
				data.Form = form
				h.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
				return
			} else {
				h.ServerError(w, err)
			}
		}
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
}

// signIn обрабатывает вход пользователей.
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := h.NewTemplateData(r)
		data.Form = models.UserSignIn{}
		h.render(w, r, http.StatusOK, "login.tmpl", data)
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}

		form := models.UserSignIn{
			Email:    r.PostForm.Get("email"),
			Password: r.PostForm.Get("password"),
		}

		id, err := h.service.Authenticate(&form)
		if err != nil {
			if errors.Is(err, models.ErrFormNotValid) ||
				errors.Is(err, models.ErrInvalidCredentials) {
				data := h.NewTemplateData(r)
				data.Form = form
				h.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
				return
			} else {
				h.ServerError(w, err)
			}
		}

		h.NewCookieFile(w, r, id)
		redirect(w, r)

	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
}

// logout обрабатывает выход пользователя из системы.
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	err := h.CloseCookieFile(w, r)
	if err != nil {
		if err == http.ErrNoCookie {
			h.ClientError(w, http.StatusBadRequest)
			return
		} else {
			h.ClientError(w, http.StatusBadRequest)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
