package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

var (
	serverCfg ServerCfg
	dbCfg     DBCfg
)

type DBCfg struct {
	PostgreSqlURI string `envconfig:"POSTGRESQL_URI" default:"host=localhost port=5432 user=chuongnd dbname=scrapedb password=1234abcd sslmode=disable"`
}

type ServerCfg struct {
	APIPort int `envconfig:"API_PORT" default:"8000"`
}

func InitConfig(configFile string) {
	ReadConfig(configFile)
	configs := []interface{}{
		&serverCfg,
		&dbCfg,
	}
	for _, instance := range configs {
		err := envconfig.Process("", instance)
		if err != nil {
			log.Fatalf("unable to init config: %v, err: %v", instance, err)
		}
	}
}

func ReadConfig(configFile string) {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetConfigFile(configFile)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf("fatal error config file: %w", err)
	}
	for _, env := range viper.AllKeys() {
		if viper.GetString(env) != "" {
			_ = os.Setenv(env, viper.GetString(env))
			_ = os.Setenv(strings.ToUpper(env), viper.GetString(env))
		}
	}
}

func ServerConfig() ServerCfg {
	return serverCfg
}

func DBConfig() DBCfg {
	return dbCfg
}
