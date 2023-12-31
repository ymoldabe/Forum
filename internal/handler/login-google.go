package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git/ymoldabe/forum/models"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (h *Handler) googleLogin(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=%s&prompt=select_account",
		models.GoogleAuthURL, models.GoogleClientID, models.GoogleRedirectUrl, "email profile")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) googleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	pathUrl := "/"

	state := r.FormValue("state")
	if state != "" {
		pathUrl = state
	}
	if code == "" {
		h.ClientError(w, http.StatusUnauthorized)
		return
	}

	tokenRes, err := GetGoogleOauthToken(code)
	if err != nil {
		h.ClientError(w, http.StatusBadGateway)
		return
	}

	google_user, err := GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)
	if err != nil {
		h.ClientError(w, http.StatusBadGateway)
	}

	now := time.Now()

	email := strings.ToLower(google_user.Email)

	user_data := models.GoogleLoginUserData{
		Name:      google_user.Name,
		Email:     email,
		Password:  "",
		Photo:     google_user.Picture,
		Provider:  "Google",
		Role:      "user",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := h.service.Authorization.GoogleAuthUser(&user_data)
	if err != nil {
		http.Redirect(w, r, "/auth/login", http.StatusBadGateway)
		return
	}
	h.NewCookieFile(w, r, id)

	http.Redirect(w, r, pathUrl, http.StatusTemporaryRedirect)
}

type GoogleUserResult struct {
	Id             string
	Email          string
	Verified_email bool
	Name           string
	Given_name     string
	Family_name    string
	Picture        string
	Locale         string
}

func GetGoogleUser(access_token string, id_token string) (*GoogleUserResult, error) {
	rootUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", access_token)

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", id_token))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer

	if _, err = io.Copy(&resBody, res.Body); err != nil {
		return nil, err
	}
	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &GoogleUserResult{
		Id:             GoogleUserRes["id"].(string),
		Email:          GoogleUserRes["email"].(string),
		Verified_email: GoogleUserRes["verified_email"].(bool),
		Name:           GoogleUserRes["name"].(string),
		Given_name:     GoogleUserRes["given_name"].(string),
		Picture:        GoogleUserRes["picture"].(string),
		Locale:         GoogleUserRes["locale"].(string),
	}
	return userBody, nil
}

type GoogleOauthToken struct {
	Access_token string
	Id_token     string
}

func GetGoogleOauthToken(code string) (*GoogleOauthToken, error) {
	const rootURl = "https://oauth2.googleapis.com/token"

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", models.GoogleClientID)
	values.Add("client_secret", models.GoogleClientSecret)
	values.Add("redirect_uri", models.GoogleRedirectUrl)

	query := values.Encode()

	req, err := http.NewRequest("POST", rootURl, bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleOauthTokenRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleOauthTokenRes); err != nil {
		return nil, err
	}

	tokenBody := &GoogleOauthToken{
		Access_token: GoogleOauthTokenRes["access_token"].(string),
		Id_token:     GoogleOauthTokenRes["id_token"].(string),
	}
	return tokenBody, nil
}
