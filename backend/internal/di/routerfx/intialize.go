package routerfx

import (
	"go.uber.org/fx"

	"scraping-keyword-web/backend/internal/server/router"
	"scraping-keyword-web/backend/pkg/keyword"
	"scraping-keyword-web/backend/pkg/user"
	"scraping-keyword-web/backend/pkg/userkeyword"
)

var Module = fx.Provide(provideRouter)

func provideRouter(userRepo user.Repository, keywordRepo keyword.Repository, userkeywordRepo userkeyword.Repository) router.Router {
	return router.NewRouter(userRepo, keywordRepo, userkeywordRepo)
}
