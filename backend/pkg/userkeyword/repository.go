package userkeyword

import (
	"context"
	"scraping-keyword-web/backend/pkg/db"
	strutil "scraping-keyword-web/backend/pkg/utils/strutils"
)

//go:generate mockgen -source=$GOFILE -package=userkeyword_mocks -destination=$PWD/mocks/${GOFILE}
type Repository interface {
	GetUserKeywordByUserID(c context.Context, userID string, limit, offset int, sort string) ([]db.UserKeyword, error)
	CreateMultipleUserKeywordByUserID(c context.Context, userID string, keywordIDs []string) ([]db.UserKeyword, error)
}

type repo struct {
	dbStorage Storage
}

func (r *repo) CreateMultipleUserKeywordByUserID(c context.Context, userID string, keywordIDs []string) (newUserKeyword []db.UserKeyword, err error) {
	userKeywords := make([]*db.UserKeyword, 0, len(keywordIDs))
	for _, kID := range keywordIDs {
		userKeywords = append(userKeywords, &db.UserKeyword{
			ID:        strutil.GenerateUUID(),
			KeywordID: kID,
			UserID:    userID,
		})
	}

	err = r.dbStorage.CreateMultipleUserKeywordByUserID(c, userKeywords)
	if err != nil {
		return
	}

	for _, userKeyword := range userKeywords {
		newUserKeyword = append(newUserKeyword, *userKeyword)
	}
	return
}

func (r *repo) GetUserKeywordByUserID(c context.Context, userID string, limit, offset int, sort string) (userKeywords []db.UserKeyword, err error) {
	return r.dbStorage.GetUserKeywordByUserID(c, userID, limit, offset, sort)
}

func NewRepository(db Storage) Repository {
	return &repo{dbStorage: db}
}
