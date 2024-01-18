package keyword

import (
	"context"
	"scraping-keyword-web/backend/pkg/db"
	strutil "scraping-keyword-web/backend/pkg/utils/strutils"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -source=$GOFILE -package=keyword_mocks -destination=$PWD/mocks/${GOFILE}
type Repository interface {
	GetKeywordResultByIDs(c context.Context, ids []string) ([]db.KeywordResult, error)
	GetKeywordResultByKeywords(c context.Context, keywords []string) ([]db.KeywordResult, error)
	CreateKeywordResults(c context.Context, keywordResults []db.KeywordResult) ([]db.KeywordResult, error)
	CrawlKeywordResults(keywords []string) ([]db.KeywordResult, error)
}

type repo struct {
	dbStorage Storage
	crawler   *colly.Collector
}

func NewRepository(db Storage, crawler *colly.Collector) Repository {
	return &repo{dbStorage: db, crawler: crawler}
}

func (repo *repo) crawlKeywordResult(kw string) (result db.KeywordResult, err error) {
	result.Keyword = kw

	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	url, _ := strutil.AddQueryParamsToRawUrl("https://www.google.com/search", map[string]string{
		"q": kw,
	})

	var htmlContent string
	err = chromedp.Run(ctx,
		// visit the target page
		chromedp.Navigate(url),
		// wait for the page to load
		chromedp.Sleep(2000*time.Millisecond),
		// extract the raw HTML from the page
		chromedp.ActionFunc(func(ctx context.Context) error {
			// select the root node on the page
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			htmlContent, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return
	}

	result.HtmlContent = htmlContent
	result.AdwordTotal = doc.Find(".pla-unit").Length()
	result.LinkTotal = doc.Find("a[href]").Length()
	result.SearchResultTotal = strutil.CollectTotalSearchResultsFromStats(doc.Find("#result-stats").Text())
	return
}

func (repo *repo) CrawlKeywordResults(keywords []string) ([]db.KeywordResult, error) {
	keywordResults := make([]db.KeywordResult, 0, len(keywords))

	var mutex sync.Mutex
	var wg errgroup.Group
	for _, kw := range keywords {
		keyword := kw
		wg.Go(func() error {
			kwRes, err := repo.crawlKeywordResult(keyword)
			if err != nil {
				return err
			}

			mutex.Lock()
			keywordResults = append(keywordResults, kwRes)
			mutex.Unlock()
			return nil
		})
	}

	err := wg.Wait()
	if err != nil {
		return nil, err
	}

	return keywordResults, nil
}

func (r *repo) CreateKeywordResults(c context.Context, keywordInfos []db.KeywordResult) (newKeywordResults []db.KeywordResult, err error) {
	keywordResults := make([]*db.KeywordResult, 0, len(keywordInfos))
	for _, k := range keywordInfos {
		kwRes := k
		kwRes.KeywordID = strutil.GenerateUUID()
		kwRes.CreatedAt = time.Now()
		keywordResults = append(keywordResults, &kwRes)
	}

	err = r.dbStorage.CreateKeywordResults(c, keywordResults)
	if err != nil {
		return newKeywordResults, err
	}

	newKeywordResults = make([]db.KeywordResult, 0, len(keywordInfos))
	for _, k := range keywordResults {
		newKeywordResults = append(newKeywordResults, *k)
	}
	return
}

func (r *repo) GetKeywordResultByKeywords(c context.Context, keywords []string) ([]db.KeywordResult, error) {
	return r.dbStorage.GetKeywordResultByKeywords(c, keywords)
}

func (r *repo) GetKeywordResultByIDs(c context.Context, ids []string) ([]db.KeywordResult, error) {
	return r.dbStorage.GetKeywordResultByIDs(c, ids)
}
