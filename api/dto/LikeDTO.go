package dto

type AddLike struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}
