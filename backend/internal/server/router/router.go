package router

import (
	"scraping-keyword-web/backend/pkg/keyword"
	"scraping-keyword-web/backend/pkg/user"
	"scraping-keyword-web/backend/pkg/userkeyword"

	"github.com/gin-gonic/gin"
)

type Router interface {
	Register(g gin.IRouter)
}

type router struct {
	userRepo        user.Repository
	keywordRepo     keyword.Repository
	userkeywordRepo userkeyword.Repository
}

func NewRouter(userRepo user.Repository, keywordRepo keyword.Repository, userkeywordRepo userkeyword.Repository) Router {
	return &router{userRepo: userRepo, keywordRepo: keywordRepo, userkeywordRepo: userkeywordRepo}
}

func (r *router) Register(g gin.IRouter) {
	userGroup := g.Group("/users")
	{
		userGroup.POST("/signin", r.signin)
		userGroup.POST("/signup", r.signup)
		userGroup.POST("/verify", r.verify)
	}

	keywordGroup := g.Group("/keywords")
	{
		keywordGroup.GET("", r.getKeywords)
		keywordGroup.POST("/upload", r.uploadKeywordFile)
	}
}
