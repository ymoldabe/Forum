package service

import (
	"git/ymoldabe/forum/internal/store"
	"git/ymoldabe/forum/models"
	"git/ymoldabe/forum/validator"
)

// PostService представляет службу для обработки операций с постами и комментариями.
type PostService struct {
	store store.Post
}

// NewPostService создает новый экземпляр PostService с указанным хранилищем (store).
func NewPostService(store store.Post) *PostService {
	return &PostService{store: store}
}

// CreatePost создает новый пост.
func (s *PostService) CreatePost(form *models.DataTransfer) (int, error) {
	// Проверки на валидность полей формы создания поста.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChar(form.Title, 40), "title", "This field cannot be more than 40 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChar(form.Content, 1200), "content", "This field cannot be more than 1200 characters long")
	form.CheckField(validator.CheckTeg(form.Tags), "tags", "This field cannot be blank")

	if !form.Valid() {
		return 0, models.ErrFormNotValid
	}

	// Создание поста и возврат его идентификатора.
	id, err := s.store.CreatePost(form)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// CreateComment создает новый комментарий к посту.
func (s *PostService) CreateComment(form *models.CommentInPost) error {
	// Проверки на валидность полей формы комментария.
	form.CheckField(validator.NotBlank(form.Content), "comment", "This field cannot be blank")

	if !form.Valid() {
		return models.ErrFormNotValid
	}
	// Создание комментария.
	return s.store.CreateComment(form)
}

// GetPost получает информацию о посте по его идентификатору.
func (s *PostService) GetPost(id int) (models.GetOnePost, error) {
	return s.store.GetPost(id)
}

// GetPosts получает список всех постов.
func (s *PostService) GetPosts() ([]models.GetAllPosts, error) {
	return s.store.GetPosts()
}

// GetMyCreatedPost получает список всех постов, созданных пользователем.
func (s *PostService) GetMyCreatedPost(userId int) ([]models.GetAllPosts, error) {
	return s.store.GetMyCreatedPost(userId)
}

// ReactionPost добавляет реакцию пользователя к посту.
func (s *PostService) ReactionPost(postId, userId int, reactionType string) error {
	return s.store.ReactionPost(postId, userId, reactionType)
}

// ReactionComment добавляет реакцию пользователя к комментарию.
func (s *PostService) ReactionComment(postId, userId, commetId int, reactionType string) error {
	return s.store.ReactionComment(postId, userId, commetId, reactionType)
}

// GetMyLikesPost получает список всех постов, которые пользователь лайкнул.
func (s *PostService) GetMyLikesPost(userId int) ([]models.GetAllPosts, error) {
	return s.store.GetMyLikesPost(userId)
}

// CheckLastPost проверяет последний созданный пост.
func (s *PostService) CheckLastPost() (int, error) {
	return s.store.CheckLastPost()
}

// CheckLastComment проверяет последний созданный комментарий.
func (s *PostService) CheckLastComment() (int, error) {
	return s.store.CheckLastComment()
}
