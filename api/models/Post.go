package models

type Post struct {
	PostID int    `gorm:"primaryKey" json:"post_id"`
	UserID int    `gorm:"foreignKey" json:"user_id"`
	Header string `json:"header"`
	Body   string `json:"body"`
}
