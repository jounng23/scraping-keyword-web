package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID      string        `gorm:"primaryKey;index;unique;not null" json:"user_id"`
	Username    string        `gorm:"index;unique;not null" json:"username"`
	Password    string        `gorm:"index;not null" json:"password"`
	CreatedAt   time.Time     `json:"created_at"`
	UserKeyword []UserKeyword `gorm:"foreignKey:user_id;references:user_id"`
}

type UserKeyword struct {
	gorm.Model
	ID        string `gorm:"primaryKey;index;unique;not null" json:"id"`
	UserID    string `json:"user_id"`
	KeywordID string `json:"keyword_id"`
}

type KeywordResult struct {
	gorm.Model
	KeywordID         string        `gorm:"primaryKey;index;unique;not null" json:"keyword_id"`
	Keyword           string        `gorm:"index;unique;not null" json:"keyword"`
	AdwordTotal       int           `json:"adword_total"`
	LinkTotal         int           `json:"link_total"`
	SearchResultTotal int           `json:"search_result_total"`
	HtmlContent       string        `json:"html_content"`
	CreatedAt         time.Time     `json:"created_at"`
	UserKeyword       []UserKeyword `gorm:"foreignKey:keyword_id;references:keyword_id"`
}
