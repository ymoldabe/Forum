package auth

import (
	"git/ymoldabe/forum/models"
	"git/ymoldabe/forum/pkg"
	"git/ymoldabe/forum/validator"
)

func (a *AuthService) GithubAuthUser(form *models.GithubLoginUserData) (int, error) {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	if !form.Valid() {
		return 0, models.ErrFormNotValid
	}

	// Генерация хеша пароля и вставка пользователя в хранилище.
	hash := pkg.GenerateRandomPassword(30)

	form.Password = hash

	form.Provider = models.ProviderGit

	id, err := a.store.GithubAuthUser(form)
	if err != nil {
		return 0, err
	}
	return id, nil
}
