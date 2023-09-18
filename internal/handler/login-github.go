package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"git/ymoldabe/forum/models"
)

func (h *Handler) githubLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s",
		models.GitHubAuthURL,
		models.GithubClientID,
		models.GithubRedirectUrl,
	)

	http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
}

func (h *Handler) githubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken, err := getGithubAccessToken(code)
	if err != nil {
		h.ClientError(w, http.StatusBadGateway)
		return
	}

	githubData, err := getGithubData(githubAccessToken)
	if err != nil {
		h.ClientError(w, http.StatusBadGateway)
	}

	userData, err := getUserData(githubData)
	if err != nil {
		h.ServerError(w, err)
		return
	}
	id, err := h.service.Authorization.GithubAuthUser(&userData)
	if err != nil {
		h.ClientError(w, http.StatusBadGateway)
		return
	}
	h.NewCookieFile(w, r, id)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getUserData(data string) (models.GithubLoginUserData, error) {
	userData := models.GithubLoginUserData{}
	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return models.GithubLoginUserData{}, err
	}

	return userData, nil
}

func getGithubAccessToken(code string) (string, error) {
	requestBodyMap := map[string]string{
		"client_id":     models.GithubClientID,
		"client_secret": models.GithubClientSecret,
		"code":          code,
	}
	requestJSON, err := json.Marshal(requestBodyMap)
	if err != nil {
		return "", err
	}

	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		return "", reqerr
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	if err := json.Unmarshal(respbody, &ghresp); err != nil {
		return "", err
	}

	return ghresp.AccessToken, nil
}

func getGithubData(accessToken string) (string, error) {
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		return "", reqerr
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		return "", resperr
	}

	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respbody), nil
}
