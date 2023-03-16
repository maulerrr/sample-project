package dto

type CreateCommentDTO struct {
	UserID int    `json:"user_id"`
	PostID int    `json:"post_id"`
	Text   string `json:"text"`
}

type FindByWordsDTO struct {
	Text string `json:"text"`
}

type UpdateComment struct {
	Text string `json:"text"`
}
