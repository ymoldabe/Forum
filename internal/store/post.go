package store

import (
	"database/sql"
	"errors"
	"git/ymoldabe/forum/models"
	"log"
)

// PostStore представляет хранилище для операций с постами и комментариями.
type PostStore struct {
	db *sql.DB
}

// NewPostSqlite создает новый экземпляр PostStore, используя переданную базу данных.
func NewPostSqlite(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

// CreatePost создает новый пост в базе данных на основе данных из формы.
func (s *PostStore) CreatePost(form *models.DataTransfer) (int, error) {
	stmtPost := `
	INSERT INTO posts (title, content, created_date, user_id)
	VALUES (?, ?, ?, ?)`

	res, err := s.db.Exec(stmtPost, form.Title, form.Content, form.CreateDate, form.UserId)
	if err != nil {
		return 0, err
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmtTags := `
	INSERT INTO posts_tags (post_id, tag_id)
	VALUES (?, ?)`

	stmtTagsExec, err := s.db.Prepare(stmtTags)
	if err != nil {
		return 0, err
	}
	defer stmtTagsExec.Close()

	for _, v := range form.Tags {
		if _, err = stmtTagsExec.Exec(postID, v); err != nil {
			return 0, err
		}
	}

	return int(postID), nil
}

// CreateComment создает новый комментарий к посту в базе данных.
func (s *PostStore) CreateComment(form *models.CommentInPost) error {
	stmt := `INSERT INTO comments(content, post_id, user_id, created_date)
			VALUES(?, ?, ?, ?)`
	if _, err := s.db.Exec(stmt, form.Content, form.PostId, form.UserId, form.CreateDate); err != nil {
		return err
	}

	return nil
}

// GetPost получает информацию о конкретном посте.
func (s *PostStore) GetPost(id int) (models.GetOnePost, error) {
	var OnePost models.GetOnePost

	if err := s.GetOnePost(id, &OnePost); err != nil {
		if err == sql.ErrNoRows {
			return models.GetOnePost{}, models.ErrNoRowsInResSet
		}
		return models.GetOnePost{}, err
	}

	OnePost.Id = id

	if err := s.GetTagsPost(id, &OnePost); err != nil {
		return models.GetOnePost{}, err
	}

	if err := s.GetReactionsPost(id, &OnePost); err != nil {
		return models.GetOnePost{}, err
	}

	if err := s.GetCommenstPost(id, &OnePost); err != nil {
		return models.GetOnePost{}, err
	}

	return OnePost, nil
}

// GetPosts возвращает список всех постов.
func (s *PostStore) GetPosts() ([]models.GetAllPosts, error) {
	var posts []models.GetAllPosts
	stmt := `
	SELECT p.id, p.title, p.created_date, u.name
	FROM posts p
	JOIN users u 
	ON u.id = p.user_id
	ORDER BY p.id DESC
	`
	postsRow, err := s.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer postsRow.Close()

	for postsRow.Next() {
		post := models.GetAllPosts{}

		if err := postsRow.Scan(&post.Id, &post.Title, &post.CreateDate, &post.UserName); err != nil {
			return nil, err
		}
		if err := s.GetTagsPosts(post.Id, &post); err != nil {
			log.Println(err)
			return nil, err
		}

		if err := s.GetAllReactionsPost(post.Id, &post); err != nil {
			log.Println(err)
			return nil, err
		}

		posts = append(posts, post)
	}
	if err := postsRow.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetMyLikesPost возвращает список всех постов с пометкой like || dislike
func (s *PostStore) GetMyLikesPost(userId int) ([]models.GetAllPosts, error) {
	// Запрос для получения постов, которые пользователь лайкнул
	stmt := `
	SELECT p.id, p.title, p.created_date, u.name
	FROM posts p
	JOIN likes_dislikes ld
	ON ld.post_id = p.id
	JOIN users u 
	ON u.id = p.user_id
	WHERE ld.user_id = ? AND ld.reaction_type = 'like'
	ORDER BY p.id DESC
	`

	// Выполняем запрос к базе данных
	postsRow, err := s.db.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer postsRow.Close()

	var posts []models.GetAllPosts

	// Итерируем по результатам запроса
	for postsRow.Next() {
		post := models.GetAllPosts{}

		// Сканируем результаты запроса в структуру post
		if err := postsRow.Scan(&post.Id, &post.Title, &post.CreateDate, &post.UserName); err != nil {
			return nil, err
		}

		// Получаем теги и реакции для данного поста
		if err := s.GetTagsPosts(post.Id, &post); err != nil {
			log.Println(err)
			return nil, err
		}

		if err := s.GetAllReactionsPost(post.Id, &post); err != nil {
			log.Println(err)
			return nil, err
		}

		posts = append(posts, post)
	}

	// Проверяем наличие ошибок после завершения итерации
	if err := postsRow.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetMyCreatedPost возвращает список постов, созданных указанным пользователем.
func (s *PostStore) GetMyCreatedPost(userId int) ([]models.GetAllPosts, error) {
	var posts []models.GetAllPosts

	// SQL-запрос для получения постов, созданных пользователем
	stmt := `
	SELECT p.id, p.title, p.created_date, u.name
	FROM posts p
	JOIN users u 
	ON u.id = p.user_id
	WHERE p.user_id = ?
	ORDER BY p.id DESC
	`

	// Выполняем SQL-запрос к базе данных
	postsRow, err := s.db.Query(stmt, userId)
	if err != nil {
		return nil, err
	}
	defer postsRow.Close()

	// Итерируем по результатам запроса
	for postsRow.Next() {
		post := models.GetAllPosts{}

		// Сканируем результаты запроса в структуру post
		if err := postsRow.Scan(&post.Id, &post.Title, &post.CreateDate, &post.UserName); err != nil {
			return nil, err
		}

		// Получаем теги и реакции для данного поста
		if err := s.GetTagsPosts(post.Id, &post); err != nil {
			log.Println(err)
			return nil, err
		}

		if err := s.GetAllReactionsPost(post.Id, &post); err != nil {
			log.Println(err)
			return nil, err
		}

		// Добавляем текущий пост в список постов
		posts = append(posts, post)
	}

	// Проверяем наличие ошибок после завершения итерации
	if err := postsRow.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// ReactionPost обрабатывает реакцию пользователя на определенный пост.
func (s *PostStore) ReactionPost(postId, userId int, reactionType string) error {
	// Проверяем существующую реакцию пользователя.
	reaction, err := s.CheckReactionTypePost(postId, userId)
	if errors.Is(err, sql.ErrNoRows) {
		// Если реакции нет, вставляем новую реакцию.
		return s.InsertReactionPost(postId, userId, reactionType)
	} else if err != nil {
		return err
	}

	if reaction == reactionType {
		// Если реакция совпадает с существующей, обновляем в "нет реакции".
		if err = s.UpdateReactionPost(postId, userId, None); err != nil {
			return err
		}
	} else {
		// В противном случае, обновляем реакцию пользователя.
		if err := s.UpdateReactionPost(postId, userId, reactionType); err != nil {
			return err
		}
	}
	return nil
}

// ReactionComment добовляет рекакцию к комменту.
func (s *PostStore) ReactionComment(postId, userId, commentId int, reactionType string) error {
	reaction, err := s.CheckReactionTypeComment(postId, commentId, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return s.InsertReactionComment(postId, commentId, userId, reactionType)
	} else if err != nil {
		return err
	}

	if reaction == reactionType {
		if err = s.UpdateReactionComment(postId, commentId, userId, None); err != nil {
			return err
		}
	} else {
		if err := s.UpdateReactionComment(postId, commentId, userId, reactionType); err != nil {
			return err
		}
	}
	return nil
}

// CheckLastPost проверяет последний созданный пост.
func (s *PostStore) CheckLastPost() (int, error) {
	var lastPostID int64

	row := s.db.QueryRow("SELECT id FROM posts ORDER BY id DESC LIMIT 1;")
	err := row.Scan(&lastPostID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRowsInResSet
		}
		return 0, err

	}

	if lastPostID == 0 {
		return 0, models.ErrNoRowsInResSet
	}

	return int(lastPostID), nil
}

// CheckLastComment проверяет последний созданный комментарий.
func (s *PostStore) CheckLastComment() (int, error) {
	var lastCommentID int64

	row := s.db.QueryRow("SELECT id FROM comments ORDER BY id DESC LIMIT 1;")
	err := row.Scan(&lastCommentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRowsInResSet
		}
		return 0, err

	}

	if lastCommentID == 0 {
		return 0, models.ErrNoRowsInResSet
	}

	return int(lastCommentID), nil
}
