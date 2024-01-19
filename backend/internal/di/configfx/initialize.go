package configfx

import (
	"go.uber.org/fx"

	"scraping-keyword-web/backend/pkg/config"
)

func Initialize(configFile string) fx.Option {
	return fx.Invoke(func() {
		config.InitConfig(configFile)
	})
}
