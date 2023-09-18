package auth

import (
	"git/ymoldabe/forum/models"
	"git/ymoldabe/forum/pkg"
	"git/ymoldabe/forum/validator"
)

func (a *AuthService) GoogleAuthUser(form *models.GoogleLoginUserData) (int, error) {
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	if !form.Valid() {
		return 0, models.ErrFormNotValid
	}

	// Генерация хеша пароля и вставка пользователя в хранилище.
	hash := pkg.GenerateRandomPassword(30)

	form.Password = hash

	form.Provider = models.ProviderGoogle

	id, err := a.store.GoogleAuthUser(form)
	if err != nil {
		return 0, err
	}
	return id, nil
}
