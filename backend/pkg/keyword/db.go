package keyword

import (
	"context"
	"scraping-keyword-web/backend/pkg/db"

	"gorm.io/gorm"
)

type Storage interface {
	CreateKeywordResults(c context.Context, keywordResults []*db.KeywordResult) error
	GetKeywordResultByKeywords(c context.Context, keywords []string) (keywordResults []db.KeywordResult, err error)
	GetKeywordResultByIDs(c context.Context, ids []string) (keywordResults []db.KeywordResult, err error)
}

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{db}
}

func (s *storage) CreateKeywordResults(c context.Context, keywordResults []*db.KeywordResult) (err error) {
	res := s.db.Create(&keywordResults)
	if res.Error != nil {
		err = res.Error
	}
	return
}

func (s *storage) GetKeywordResultByIDs(c context.Context, ids []string) (keywordResults []db.KeywordResult, err error) {
	res := s.db.Where("keyword_id IN ?", ids).Find(&keywordResults)
	if res.Error != nil {
		err = res.Error
	}
	return
}

func (s *storage) GetKeywordResultByKeywords(c context.Context, keywords []string) (keywordResults []db.KeywordResult, err error) {
	res := s.db.Where("keyword IN ?", keywords).Find(&keywordResults)
	if res.Error != nil {
		err = res.Error
	}
	return
}
