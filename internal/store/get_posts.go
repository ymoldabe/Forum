package store

import (
	"git/ymoldabe/forum/models"
	"strings"
)

// GetOnePost получает информацию о конкретном посте.
func (s *PostStore) GetOnePost(id int, OnePost *models.GetOnePost) error {
	stmt := `
	SELECT 
		p.title, 
		p.content, 
		p.created_date, 
		u.name 
	FROM
		posts AS p
	JOIN 
		users u
	ON  
		u.id = p.user_id
	WHERE 
		p.id = ?`

	err := s.db.QueryRow(stmt, id).Scan(
		&OnePost.Title,
		&OnePost.Content,
		&OnePost.CreateDate,
		&OnePost.UserName,
	)

	return err
}

// GetReactionsComment получает реакции на комментарий в посте.
func (s *PostStore) GetReactionsComment(post_id int, comment *models.CommentInPost) error {
	stmtLikes := `
	SELECT 
		SUM(CASE WHEN reaction_type = 'likeComm' THEN 1 ELSE 0 END) AS likeComm_count,
    	SUM(CASE WHEN reaction_type = 'dislikeComm' THEN 1 ELSE 0 END) AS dislikeComm_count
	FROM
		likes_dislikes ld
	JOIN 
		comments c
	ON
		c.id = ld.comment_id
	WHERE 
		c.post_id = ? AND c.id = ?`

	if err := s.db.QueryRow(stmtLikes, post_id, comment.Id).Scan(&comment.Likes, &comment.Dislikes); err != nil {
		if strings.Contains(err.Error(), "converting NULL to int is unsupported") {
			comment.Likes = 0
			comment.Dislikes = 0
		} else {
			return err
		}
	}
	return nil
}

// GetReactionsPost получает реакции на пост.
func (s *PostStore) GetReactionsPost(id int, OnePost *models.GetOnePost) error {
	stmtLikes := `
	SELECT 
	SUM(CASE WHEN reaction_type = 'like' THEN 1 ELSE 0 END) AS like_count,
    SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END) AS dislike_count
	FROM
	likes_dislikes
	WHERE 
	post_id = ? AND comment_id IS NULL`

	if err := s.db.QueryRow(stmtLikes, id).Scan(&OnePost.Likes, &OnePost.Dislikes); err != nil {
		if strings.Contains(err.Error(), "converting NULL to int is unsupported") {
			OnePost.Likes = 0
			OnePost.Dislikes = 0
		} else {
			return err
		}
	}
	return nil
}

// GetTagsPosts получает теги для поста.
func (s *PostStore) GetTagsPosts(id int, OnePost *models.GetAllPosts) error {
	stmt := `
	SELECT	t.tag
	FROM 	tags t
	JOIN 	posts_tags pt
	ON 		t.id = pt.tag_id
	WHERE 	pt.post_id = ?`

	tags, err := s.db.Query(stmt, id)
	if err != nil {
		return err
	}

	defer tags.Close()

	for tags.Next() {
		var tag string

		if err = tags.Scan(&tag); err != nil {
			return err
		}

		OnePost.Tags = append(OnePost.Tags, tag)
	}
	if err = tags.Err(); err != nil {
		return nil
	}
	return nil
}

// GetTagsPost получает теги для поста.
func (s *PostStore) GetTagsPost(id int, OnePost *models.GetOnePost) error {
	stmt := `
	SELECT	t.tag
	FROM 	tags t
	JOIN 	posts_tags pt
	ON 		t.id = pt.tag_id
	WHERE 	pt.post_id = ?`

	tags, err := s.db.Query(stmt, id)
	if err != nil {
		return err
	}

	defer tags.Close()

	for tags.Next() {
		var tag string

		if err = tags.Scan(&tag); err != nil {
			return err
		}

		OnePost.Tags = append(OnePost.Tags, tag)
	}
	if err = tags.Err(); err != nil {
		return nil
	}
	return nil
}

// GetCommenstPost получает комментарии для поста.
func (s *PostStore) GetCommenstPost(id int, OnePost *models.GetOnePost) error {
	stmt := `
	SELECT		
	c.id,
	c.post_id,
	c.content,
	u.name,
	c.created_date
	FROM		comments c
	JOIN		users u
	ON			c.user_id = u.id
	WHERE 		c.post_id = ?`

	comments, err := s.db.Query(stmt, id)
	if err != nil {
		return err
	}

	defer comments.Close()

	for comments.Next() {
		var comment models.CommentInPost
		if err = comments.Scan(&comment.Id, &comment.PostId, &comment.Content, &comment.UserName, &comment.CreateDate); err != nil {
			return err
		}
		if err := s.GetReactionsComment(id, &comment); err != nil {
			return err
		}

		OnePost.Comments = append(OnePost.Comments, comment)
	}

	if err = comments.Err(); err != nil {
		return nil
	}
	return nil
}

// GetAllReactionsPost получает общие реакции на пост.
func (s *PostStore) GetAllReactionsPost(id int, OnePost *models.GetAllPosts) error {
	stmtLikes := `
	SELECT 
	SUM(CASE WHEN reaction_type = 'like' THEN 1 ELSE 0 END) AS like_count,
    SUM(CASE WHEN reaction_type = 'dislike' THEN 1 ELSE 0 END) AS dislike_count
	FROM
	likes_dislikes
	WHERE 
	post_id = ? AND comment_id IS NULL`

	if err := s.db.QueryRow(stmtLikes, id).Scan(&OnePost.Likes, &OnePost.Dislikes); err != nil {
		if strings.Contains(err.Error(), "converting NULL to int is unsupported") {
			OnePost.Likes = 0
			OnePost.Dislikes = 0
		} else {
			return err
		}
	}
	return nil
}
