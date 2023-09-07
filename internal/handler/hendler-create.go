package handler

import (
	"errors"
	"fmt"
	"git/ymoldabe/forum/models"
	"net/http"
	"time"
)

const dateFormat = "2006-01-02"

// postCreate обрабатывает создание нового поста.
func (h *Handler) postCreate(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	userUrlBefore[r.RemoteAddr] = r.URL.Path
	mu.Unlock()
	switch r.Method {
	case http.MethodGet:
		data := h.NewTemplateData(r)
		data.Form = models.DataTransfer{}
		h.render(w, r, http.StatusOK, "create.tmpl", data)
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			h.ClientError(w, http.StatusBadRequest)
			return
		}

		userId, err := h.GetUserAuth(r)
		if err != nil {
			h.ServerError(w, err)
			return
		}

		form := models.DataTransfer{
			UserId:     userId,
			Title:      r.PostForm.Get("title"),
			Content:    r.PostForm.Get("content"),
			Tags:       r.Form["tags"],
			CreateDate: time.Now().Format(dateFormat),
		}

		id, err := h.service.CreatePost(&form)
		if err != nil {
			if errors.Is(err, models.ErrFormNotValid) {
				data := h.NewTemplateData(r)
				data.Form = form
				h.render(w, r, http.StatusBadRequest, "create.tmpl", data)
				return
			} else {
				h.ServerError(w, err)
				return
			}
		}

		url := fmt.Sprintf("/post/view?id=%d", id)

		http.Redirect(w, r, url, http.StatusSeeOther)
	default:
		h.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
}

// createPostComment обрабатывает создание нового комментария к посту.
func (h *Handler) createPostComment(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	userUrlBefore[r.RemoteAddr] = r.URL.Path
	mu.Unlock()

	postId, status, err := h.CheckValidIdOrCommentId(r, false)
	if status == http.StatusNotFound {
		h.NotFound(w)
		return
	} else if status == http.StatusInternalServerError {
		h.ServerError(w, err)
		return
	}

	url := fmt.Sprintf("/post/view?id=%d", postId)

	if comment := (r.FormValue("comment")); comment == "" {
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}

	userId, err := h.GetUserAuth(r)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	form := models.CommentInPost{
		PostId:     postId,
		UserId:     userId,
		CreateDate: time.Now().Format(dateFormat),
		Content:    r.FormValue("comment"),
	}

	err = h.service.CreateComment(&form)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
