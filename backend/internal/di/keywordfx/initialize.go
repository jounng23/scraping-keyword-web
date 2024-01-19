package keywordfx

import (
	"github.com/gocolly/colly"
	"go.uber.org/fx"
	"gorm.io/gorm"

	"scraping-keyword-web/backend/pkg/keyword"
)

var Module = fx.Provide(
	provideKeywordStorage,
	provideKeywordRepository,
)

func provideKeywordStorage(db *gorm.DB) keyword.Storage {
	return keyword.NewStorage(db)
}

func provideKeywordRepository(dbStorage keyword.Storage, crawler *colly.Collector) keyword.Repository {
	return keyword.NewRepository(dbStorage, crawler)
}
