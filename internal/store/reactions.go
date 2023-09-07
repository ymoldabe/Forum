package store

const (
	Like    string = "like"
	Dislike string = "dislike"
	None    string = "none"
)

/*
Этот код просто помогает поставить лайк или дизлайк все )
*/

func (s *PostStore) CheckReactionTypePost(postId, userId int) (string, error) {
	stmt := `
	SELECT reaction_type
	FROM likes_dislikes
	WHERE  post_id = ? AND user_id = ? AND comment_id IS NULL`
	var reaction_type string
	err := s.db.QueryRow(stmt, postId, userId).Scan(&reaction_type)
	return reaction_type, err
}

func (s *PostStore) UpdateReactionPost(postId, userId int, reactionType string) error {
	stmt := `
	UPDATE likes_dislikes
	SET reaction_type = ?
	WHERE post_id = ? AND user_id = ? AND comment_id IS NULL`
	_, err := s.db.Exec(stmt, reactionType, postId, userId)
	return err
}

func (s *PostStore) InsertReactionPost(postId, userId int, reaction string) error {
	stmt := `
	INSERT INTO likes_dislikes (post_id, user_id, reaction_type)
	VALUES (?, ?, ?)`
	_, err := s.db.Exec(stmt, postId, userId, reaction)
	return err
}

func (s *PostStore) CheckReactionTypeComment(postId, commentId, userId int) (string, error) {
	stmt := `
	SELECT reaction_type
	FROM likes_dislikes
	WHERE  comment_id = ? AND user_id = ? AND post_id = ?`
	var reaction_type string
	err := s.db.QueryRow(stmt, commentId, userId, postId).Scan(&reaction_type)
	return reaction_type, err
}

func (s *PostStore) UpdateReactionComment(postId, commentId, userId int, reactionType string) error {
	stmt := `
	UPDATE likes_dislikes
	SET reaction_type = ?
	WHERE comment_id = ? AND user_id = ? AND post_id = ?`
	_, err := s.db.Exec(stmt, reactionType, commentId, userId, postId)
	return err
}

func (s *PostStore) InsertReactionComment(postId, commentId, userId int, reaction string) error {
	stmt := `
	INSERT INTO likes_dislikes (post_id, comment_id, user_id, reaction_type)
	VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(stmt, postId, commentId, userId, reaction)
	return err
}
