package dto

type CreatePost struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}

type UpdatePost struct {
	Header string `json:"header"`
	Body   string `json:"body"`
}
