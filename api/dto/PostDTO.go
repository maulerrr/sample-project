package dto

type CreatePost struct {
	UserID int    `json:"user_id"`
	Header string `json:"header"`
	Body   string `json:"body"`
}

type UpdatePost struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}
