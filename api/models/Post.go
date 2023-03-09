package models

type Post struct {
	PostID int    `gorm:"primaryKey" json:"post_id"`
	Header string `json:"header"`
	Body   string `json:"body"`
}
