package crawlerfx

import (
	"github.com/gocolly/colly"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	provideCrawler,
)

func provideCrawler() *colly.Collector {
	return colly.NewCollector()
}
