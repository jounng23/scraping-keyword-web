package keyword

import (
	"context"
	"scraping-keyword-web/backend/pkg/db"
	strutil "scraping-keyword-web/backend/pkg/utils/strutils"
	"sync"
	"time"

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

func (repo *repo) crawlKeywordResults(kw string) (result db.KeywordResult, err error) {
	result.Keyword = kw

	// Visit the Google search page
	repo.crawler.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	})

	repo.crawler.OnHTML("html", func(e *colly.HTMLElement) {
		// Extracting the total number of AdWords advertisers on the page
		result.AdwordTotal = e.DOM.Find(".qvfQJe").Length()

		// Extracting the total number of links on the page
		result.LinkTotal = e.DOM.Find("a").Length()

		// Extracting the total search results for this keyword
		result.SearchResultTotal = e.DOM.Find("#result-stats").Text()

	})

	repo.crawler.OnResponse(func(r *colly.Response) {
		result.HtmlContent = string(r.Body)
	})

	url, _ := strutil.AddQueryParamsToRawUrl("https://www.google.com/search", map[string]string{
		"q": kw,
	})
	err = repo.crawler.Visit(url)
	if err != nil {
		return
	}

	// Start scraping
	repo.crawler.Wait()
	return
}

func (repo *repo) CrawlKeywordResults(keywords []string) ([]db.KeywordResult, error) {
	keywordResults := make([]db.KeywordResult, 0, len(keywords))

	var mutex sync.Mutex
	var wg errgroup.Group
	for _, kw := range keywords {
		keyword := kw
		wg.Go(func() error {
			kwRes, err := repo.crawlKeywordResults(keyword)
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
