package db

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgreSQLClient(uri string) (DB *gorm.DB, err error) {
	DB, err = gorm.Open(postgres.Open(uri), &gorm.Config{})
	DB.AutoMigrate(User{})
	DB.AutoMigrate(UserKeyword{})
	DB.AutoMigrate(KeywordResult{})
	if err != nil {
		log.Error().Msg("Error connecting to the database..." + err.Error())
		return
	}
	fmt.Println("Database connection successful...", uri)
	return
}
