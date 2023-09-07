package pkg

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	salt       = "auidsyfhqewgrsaodu"
	signingKey = "jahdljha8weyrjasdf"
	tokenTTL   = 12 * time.Hour
)

// GeneratePasswordHash хэширует пароль.
func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), 14)
	return string(bytes), err
}

// CheckPasswordHash проверяет пароль на соответсвие.
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))
	return err
}
