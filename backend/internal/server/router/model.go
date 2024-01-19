package router

import (
	"scraping-keyword-web/backend/pkg/db"

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

type GetKeywordsResponse struct {
	Metadata       GetKeywordsResponseMetaData `json:"metadata"`
	KeywordResults []db.KeywordResult          `json:"keyword_results"`
}

type GetKeywordsResponseMetaData struct {
	Total int64 `json:"total"`
}
