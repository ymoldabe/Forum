package models

import (
	"git/ymoldabe/forum/validator"
)

var (
	LikePost       = "like"
	DislikePost    = "dislike"
	LikeComment    = "likeComm"
	DislikeComment = "dislikeComm"
)

type DataTransfer struct {
	PostId     int
	UserId     int
	UserName   string
	Title      string
	Content    string
	CreateDate string
	Tags       []string
	validator.Validator
}

type Teg struct {
	Teg string
}

type GetOnePost struct {
	Id         int
	Title      string
	Content    string
	CreateDate string
	UserName   string
	Tags       []string
	Comments   []CommentInPost
	validator.Validator
	Likes    int
	Dislikes int
}

type LikeANDDislike struct {
	Id      int
	UserId  int
	Like    int
	Dislike int
}

type CommentInPost struct {
	Id         int
	PostId     int
	UserId     int
	Content    string
	UserName   string
	CreateDate string
	Likes      int
	Dislikes   int
	validator.Validator
}

type GetAllPosts struct {
	Id         int
	CreateDate string
	Title      string
	UserName   string
	Likes      int
	Dislikes   int
	Tags       []string
	ValidTags  string
}
