package store

import (
	"database/sql"
	"git/ymoldabe/forum/models"
	"time"
)

// Post определяет интерфейс для операций с постами и комментариями.
type Post interface {
	CreatePost(form *models.DataTransfer) (int, error)
	CreateComment(form *models.CommentInPost) error
	GetPost(id int) (models.GetOnePost, error)
	GetPosts() ([]models.GetAllPosts, error)
	GetMyCreatedPost(userId int) ([]models.GetAllPosts, error)
	GetMyLikesPost(userId int) ([]models.GetAllPosts, error)
	ReactionPost(postId, userId int, reactionType string) error
	ReactionComment(postId, userId, commetId int, reactionType string) error
	CheckLastPost() (int, error)
	CheckLastComment() (int, error)
}

// Authorization определяет интерфейс для операций аутентификации и управления пользователями.
type Autorization interface {
	InsertUser(form *models.UserSignUp) error
	Authenticate(form *models.UserSignIn) (int, error)
	UserSessionsAdd(userId int, sessionToken string, expiresAt time.Time) error
	DeleteToken(sessionToken string) error
	GetIdInSessions(sessionToken string) (int, error)
	CheckSessions(userId int) (bool, error)
	UpdateToken(sessionToken string, user_id int) error
	GetTokenSession(cookieToken string) (bool, error)
}

// Store объединяет интерфейсы Post и Authorization.
type Store struct {
	Post
	Autorization
}

// New создает новый экземпляр хранилища Store.
func New(db *sql.DB) *Store {
	return &Store{
		Post:         NewPostSqlite(db),
		Autorization: NewAuthSqlite(db),
	}
}
