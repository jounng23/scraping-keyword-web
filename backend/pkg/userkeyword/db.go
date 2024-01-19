package userkeyword

import (
	"context"
	"scraping-keyword-web/backend/pkg/db"

	"gorm.io/gorm"
)

type Storage interface {
	CreateMultipleUserKeywordByUserID(c context.Context, newUserKeywor []*db.UserKeyword) error
	GetUserKeywordByUserID(c context.Context, userID string, limit, offset int, sort string) ([]db.UserKeyword, int64, error)
}

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{db}
}

func (s *storage) GetUserKeywordByUserID(c context.Context, userID string, limit, offset int, sort string) (userKeywords []db.UserKeyword, total int64, err error) {
	res := s.db.Where("user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order(sort).
		Find(&userKeywords).
		Count(&total)
	if res.Error != nil {
		return
	}
	return
}

func (s *storage) CreateMultipleUserKeywordByUserID(c context.Context, newUserKeywor []*db.UserKeyword) error {
	res := s.db.Create(&newUserKeywor)
	return res.Error
}
