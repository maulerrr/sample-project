package models

type Comment struct {
	CommentID int    `gorm:"primaryKey" json:"comment_id"`
	UserID    int    `gorm:"foreignKey" json:"user_id"`
	PostID    int    `gorm:"foreignKey" json:"post_id"`
	Text      string `json:"text"`
}
