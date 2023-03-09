package models

import "github.com/dgrijalva/jwt-go"

type TokenResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	jwt.StandardClaims
}
