package pkg

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomPassword(length int) string {
	rand.Seed(time.Now().UnixNano())

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}
