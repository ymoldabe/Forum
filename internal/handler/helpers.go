package handler

import (
	"errors"
	"git/ymoldabe/forum/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

// TemplateData представляет структуру данных для шаблона.
type TemplateData struct {
	CurrentYear     int
	Form            interface{}
	IsAuthenticated bool
}

// NewTemplateData создает и возвращает новый экземпляр TemplateData, заполняя его данными о текущем годе и аутентификации.
func (h *Handler) NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		CurrentYear:     time.Now().Year(),
		IsAuthenticated: h.isAuthenticated(r),
	}
}

// isAuthenticated проверяет аутентификацию пользователя на основе сессионного токена.
func (h *Handler) isAuthenticated(r *http.Request) bool {
	isAuthenticated, err := r.Cookie("session_token")
	if err != nil {
		return false
	}
	cookieToken := isAuthenticated.Value

	ok, err := h.service.GetTokenSession(cookieToken)
	if err != nil {
		log.Println(err)
		return false
	}
	if !ok {
		return false
	}

	return ok
}

// redirect выполняет перенаправление на предыдущую страницу.
func redirect(w http.ResponseWriter, r *http.Request) {
	var url string

	mu.Lock()
	defer mu.Unlock()

	prevPage, ok := userUrlBefore[r.RemoteAddr]
	if ok {
		url = prevPage
	} else {
		url = "/"
	}

	defer delete(userUrlBefore, prevPage)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// CheckValidIdOrCommentId проверяет допустимость идентификатора поста или комментария.
func (h *Handler) CheckValidIdOrCommentId(r *http.Request, isComment bool) (int, int, error) {
	var idType string

	if isComment {
		idType = "comment"
	} else {
		idType = "id"
	}

	urlId := r.URL.Query().Get(idType)
	id, err := strconv.Atoi(urlId)
	if err != nil || !CheckId(urlId, id) {
		return 0, http.StatusNotFound, err
	}

	var lastId int
	var lastErr error

	if isComment {
		lastId, lastErr = h.service.CheckLastComment()
	} else {
		lastId, lastErr = h.service.CheckLastPost()
	}

	if lastErr != nil {
		if errors.Is(lastErr, models.ErrNoRowsInResSet) {
			return id, http.StatusNotFound, nil
		}
		return id, http.StatusInternalServerError, lastErr
	}

	if CheckIdReqLowerThanOne(id) || CheckPostAndLast(id, lastId) {
		return id, http.StatusNotFound, nil
	}

	return id, http.StatusOK, nil
}

// CheckId выполняет проверку идентификаторов.
func CheckId(urlId string, id int) bool {
	return urlId == strconv.Itoa(id)
}

// CheckIdReqLowerThanOne проверяет, что идентификатор поста или комментария не меньше 1.
func CheckIdReqLowerThanOne(id int) bool {
	return id < 1
}

// CheckPostAndLast выполняет сравнение идентификаторов постов или комментариев и последних идентификаторов.
func CheckPostAndLast(id, lastId int) bool {
	return lastId < id
}
