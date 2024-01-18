package router

import (
	"github.com/golang-jwt/jwt/v5"
)

type SignInRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
