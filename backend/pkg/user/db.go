package user

import (
	"context"
	"errors"
	"scraping-keyword-web/backend/pkg/db"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("user not found")

type Storage interface {
	GetUserByID(c context.Context, ID string) (user db.User, err error)
	GetUserByAuthentication(c context.Context, username, hashedPassword string) (db.User, error)
	CreateUser(c context.Context, user db.User) (db.User, error)
}

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{db}
}

func (s *storage) GetUserByAuthentication(c context.Context, username, hashedPassword string) (user db.User, err error) {
	var users []db.User
	res := s.db.Select("user_id", "created_at").
		Where("username = ? AND password = ?", username, hashedPassword).
		Find(&users)
	if res.Error != nil {
		return
	}
	user = users[0]
	return
}

func (s *storage) GetUserByID(c context.Context, ID string) (user db.User, err error) {
	var users []db.User
	res := s.db.Select("user_id", "username", "created_at").
		Where("user_id = ?", ID).
		Find(&users)
	if res.Error != nil {
		return
	}
	user = users[0]
	return
}

func (s *storage) CreateUser(c context.Context, newUser db.User) (user db.User, err error) {
	res := s.db.Create(&newUser)
	if res.Error != nil {
		err = res.Error
		return user, err
	}
	return s.GetUserByID(c, newUser.UserID)
}
