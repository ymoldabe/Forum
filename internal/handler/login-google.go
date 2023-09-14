package handler

import (
	"fmt"
	"net/http"

	"git/ymoldabe/forum/models"
)

func (h *Handler) googleLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s", models.GoogleAuthURL, models.ClientID, "http://localhost:8080/callback", "email profile")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) collback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello i am work")
}
