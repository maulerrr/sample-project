package models

import "time"

type User struct {
	UserID    int       `gorm:"primaryKey" json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
