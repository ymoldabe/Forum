package auth

import (
	"errors"
	"time"

	store_auth "git/ymoldabe/forum/internal/store/auth"
	"git/ymoldabe/forum/models"
	"git/ymoldabe/forum/pkg"
	"git/ymoldabe/forum/validator"
)

// AuthService представляет службу для обработки операций аутентификации и управления пользователями.
type AuthService struct {
	store store_auth.Authorization
}

// NewAuthService создает новый экземпляр AuthService с указанным хранилищем (store).
func New(store store_auth.Authorization) *AuthService {
	return &AuthService{
		store: store,
	}
}

// InsertUser вставляет нового пользователя в базу данных.
func (a *AuthService) InsertUser(form *models.UserSignUp) error {
	// Проверки на валидность полей формы регистрации.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This fiel must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password1), "password", "This field cannot be blank")
	form.CheckField(validator.MinChar(form.Password1, 8), "password", "This field must be at laest 8 characters long")
	form.CheckField(validator.NotBlank(form.Password2), "password", "This field cannot be blank")
	form.CheckField(validator.MinChar(form.Password2, 8), "password", "This field must be at laest 8 characters long")
	form.CheckField(validator.CheckPassword(form.Password1, form.Password2), "password", "flogged do not match")
	if !form.Valid() {
		return models.ErrFormNotValid
	}

	// Генерация хеша пароля и вставка пользователя в хранилище.
	hash, err := pkg.GeneratePasswordHash(form.Password1)
	if err != nil {
		return err
	}

	form.HashPassword = hash

	form.Provider = models.ProviderDefult

	err = a.store.InsertUser(form)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			return models.ErrDuplicateEmail
		} else {
			return err
		}
	}

	return nil
}

// Authenticate аутентифицирует пользователя.
func (a *AuthService) Authenticate(form *models.UserSignIn) (int, error) {
	// Проверки на валидность полей формы входа.
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This fiel must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChar(form.Password, 8), "password", "This field must be at laest 8 characters long")
	if !form.Valid() {
		return 0, models.ErrFormNotValid
	}

	// Попытка аутентификации пользователя и возврат его идентификатора.
	id, err := a.store.Authenticate(form)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			return 0, models.ErrDuplicateEmail
		} else if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrectfd")
			return 0, models.ErrInvalidCredentials
		} else if errors.Is(err, models.ErrProviderGoogle) {
			form.AddNonFieldError("Please log in via Google")
			return 0, models.ErrInvalidCredentials
		} else if errors.Is(err, models.ErrProviderGithub) {
			form.AddNonFieldError("Please log in via Github")
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

// UserSessionsAdd добавляет сессию пользователя в базу данных.
func (a *AuthService) UserSessionsAdd(userId int, sessionToken string, expiresAt time.Time) error {
	return a.store.UserSessionsAdd(userId, sessionToken, expiresAt)
}

// DeleteToken удаляет токен сессии пользователя.
func (a *AuthService) DeleteToken(sessionToken string) error {
	return a.store.DeleteToken(sessionToken)
}

// GetIdInSessions получает идентификатор пользователя по токену сессии.
func (a *AuthService) GetIdInSessions(sessionToken string) (int, error) {
	return a.store.GetIdInSessions(sessionToken)
}

// CheckSessions проверяет наличие сессии у пользователя.
func (a *AuthService) CheckSessions(userId int) (bool, error) {
	return a.store.CheckSessions(userId)
}

// UpdateToken обновляет токен сессии пользователя.
func (a *AuthService) UpdateToken(sessionToken string, user_id int) error {
	return a.store.UpdateToken(sessionToken, user_id)
}

// GetTokenSession проверяет наличие токена сессии.
func (a *AuthService) GetTokenSession(cookieToken string) (bool, error) {
	return a.store.GetTokenSession(cookieToken)
}
