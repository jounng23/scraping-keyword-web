package userfx

import (
	"go.uber.org/fx"
	"gorm.io/gorm"

	"scraping-keyword-web/backend/pkg/user"
)

var Module = fx.Provide(
	provideUserStorage,
	provideUserRepository,
)

func provideUserStorage(db *gorm.DB) user.Storage {
	return user.NewStorage(db)
}

func provideUserRepository(dbStorage user.Storage) user.Repository {
	return user.NewRepository(dbStorage)
}
