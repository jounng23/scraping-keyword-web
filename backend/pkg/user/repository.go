package user

import (
	"context"
	"scraping-keyword-web/backend/pkg/db"
	strutil "scraping-keyword-web/backend/pkg/utils/strutils"
	"time"
)

//go:generate mockgen -source=$GOFILE -package=user_mocks -destination=$PWD/mocks/${GOFILE}
type Repository interface {
	GetUserByUsername(c context.Context, username string) (db.User, error)
	CreateUser(c context.Context, user db.User) (db.User, error)
}

type repo struct {
	dbStorage Storage
}

func (r *repo) GetUserByUsername(c context.Context, username string) (user db.User, err error) {
	user, err = r.dbStorage.GetUserByUsername(c, username)
	if err != nil {
		return
	}
	return
}

func (r *repo) CreateUser(c context.Context, user db.User) (db.User, error) {
	user.UserID = strutil.GenerateUUID()
	user.Password, _ = strutil.HashPassword(user.Password)
	user.CreatedAt = time.Now()
	return r.dbStorage.CreateUser(c, user)
}

func NewRepository(db Storage) Repository {
	return &repo{dbStorage: db}
}
