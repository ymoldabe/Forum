package service

import (
	"git/ymoldabe/forum/internal/service/auth"
	"git/ymoldabe/forum/internal/service/post"
	"git/ymoldabe/forum/internal/store"
)

// Post интерфейс для методов, связанных с постами.

// Autorization интерфейс для методов, связанных с аутентификацией и управлением пользователями.

// Service структура, агрегирующая интерфейсы Post и Autorization для централизованной функциональности.
type Service struct {
	Post          post.Post
	Authorization auth.Autorization
}

// New создает новый экземпляр Service с переданным хранилищем Store.
func New(store *store.Store) *Service {
	return &Service{
		Post:          post.New(store.Post),
		Authorization: auth.New(store.Authorization),
	}
}
