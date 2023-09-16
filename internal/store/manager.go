package store

import (
	"database/sql"

	store_auth "git/ymoldabe/forum/internal/store/auth"
	store_post "git/ymoldabe/forum/internal/store/post"
)

// Post определяет интерфейс для операций с постами и комментариями.

// Authorization определяет интерфейс для операций аутентификации и управления пользователями.

// Store объединяет интерфейсы Post и Authorization.
type Store struct {
	Post          store_post.Post
	Authorization store_auth.Authorization
}

// New создает новый экземпляр хранилища Store.
func New(db *sql.DB) *Store {
	return &Store{
		Post:          store_post.New(db),
		Authorization: store_auth.New(db),
	}
}
