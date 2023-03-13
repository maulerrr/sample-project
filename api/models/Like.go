package models

type Like struct {
	LikeID int `gorm:"primaryKey" json:"like_id"`
	UserID int `gorm:"foreignKey" json:"user_id"`
	PostID int `gorm:"foreignKey" json:"post_id"`
}
