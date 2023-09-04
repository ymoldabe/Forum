package service

import (
	"errors"
	"time"

	"git/ymoldabe/forum/internal/store"
	"git/ymoldabe/forum/models"
	"git/ymoldabe/forum/pkg"
	"git/ymoldabe/forum/validator"
)

// AuthService представляет службу для обработки операций аутентификации и управления пользователями.
type AuthService struct {
	store store.Autorization
}

// NewAuthService создает новый экземпляр AuthService с указанным хранилищем (store).
func NewAuthService(store store.Autorization) *AuthService {
	return &AuthService{
		store: store,
	}
}

// InsertUser вставляет нового пользователя в базу данных.
func (s *AuthService) InsertUser(form *models.UserSignUp) error {
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

	err = s.store.InsertUser(form)
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
func (s *AuthService) Authenticate(form *models.UserSignIn) (int, error) {
	// Проверки на валидность полей формы входа.
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This fiel must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChar(form.Password, 8), "password", "This field must be at laest 8 characters long")
	if !form.Valid() {
		return 0, models.ErrFormNotValid
	}

	// Попытка аутентификации пользователя и возврат его идентификатора.
	id, err := s.store.Authenticate(form)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			return 0, models.ErrDuplicateEmail
		} else if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrectfd")
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

// UserSessionsAdd добавляет сессию пользователя в базу данных.
func (s *AuthService) UserSessionsAdd(userId int, sessionToken string, expiresAt time.Time) error {
	return s.store.UserSessionsAdd(userId, sessionToken, expiresAt)
}

// DeleteToken удаляет токен сессии пользователя.
func (s *AuthService) DeleteToken(sessionToken string) error {
	return s.store.DeleteToken(sessionToken)
}

// GetIdInSessions получает идентификатор пользователя по токену сессии.
func (s *AuthService) GetIdInSessions(sessionToken string) (int, error) {
	return s.store.GetIdInSessions(sessionToken)
}

// CheckSessions проверяет наличие сессии у пользователя.
func (s *AuthService) CheckSessions(userId int) (bool, error) {
	return s.store.CheckSessions(userId)
}

// UpdateToken обновляет токен сессии пользователя.
func (s *AuthService) UpdateToken(sessionToken string, user_id int) error {
	return s.store.UpdateToken(sessionToken, user_id)
}

// GetTokenSession проверяет наличие токена сессии.
func (s *AuthService) GetTokenSession(cookieToken string) (bool, error) {
	return s.store.GetTokenSession(cookieToken)
}
