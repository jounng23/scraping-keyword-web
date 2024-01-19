package userkeywordfx

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	"scraping-keyword-web/backend/pkg/userkeyword"
)

var Module = fx.Provide(
	provideUserStorage,
	provideUserRepository,
)

func provideUserStorage(db *gorm.DB) userkeyword.Storage {
	return userkeyword.NewStorage(db)
}

func provideUserRepository(dbStorage userkeyword.Storage) userkeyword.Repository {
	return userkeyword.NewRepository(dbStorage)
}
