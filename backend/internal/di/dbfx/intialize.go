package dbfx

import (
	"errors"

	"go.uber.org/fx"
	"gorm.io/gorm"

	"scraping-keyword-web/backend/pkg/config"
	"scraping-keyword-web/backend/pkg/db"
)

var Module = fx.Provide(
	providePostgreSQLDatabase,
)

func providePostgreSQLDatabase(lifecycle fx.Lifecycle) (pgDB *gorm.DB, err error) {
	dbCfg := config.DBConfig()
	if dbCfg.PostgreSqlURI == "" {
		return nil, errors.New("postgreSQL URI is empty")
	}

	pgDB, err = db.InitPostgreSQLClient(dbCfg.PostgreSqlURI)
	return
}
