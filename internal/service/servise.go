package service

import (
	"time"

	"git/ymoldabe/forum/internal/store"
	"git/ymoldabe/forum/models"
)

// Post интерфейс для методов, связанных с постами.
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

// Autorization интерфейс для методов, связанных с аутентификацией и управлением пользователями.
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

// Service структура, агрегирующая интерфейсы Post и Autorization для централизованной функциональности.
type Service struct {
	Post
	Autorization
}

// New создает новый экземпляр Service с переданным хранилищем Store.
func New(store *store.Store) *Service {
	return &Service{
		Post:         NewPostService(store.Post),
		Autorization: NewAuthService(store.Autorization),
	}
}
