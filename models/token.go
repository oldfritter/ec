package models

import (
	"time"

	"github.com/jinzhu/gorm"

	"ec/utils"
)

const (
	DefaultExpire = time.Hour * 24 // 1天有效期
)

type Token struct {
	gorm.Model
	UserId       int
	IsUsed       bool
	Type         string    `gorm:"type:varchar(32)"` // 令牌类型 Tokens::Login, Tokens::AccessToken
	Token        string    `gorm:"type:varchar(64)"`
	RemoteIp     string    `gorm:"varchar(64)" json:"-"`
	ExpireAt     time.Time `gorm:"default:null"`
	LastVerifyAt time.Time `gorm:"default:null" json:"-"`
}

func (token *Token) BeforeCreate(db *gorm.DB) {
	if token.Type == "" {
		token.Type = "Tokens::Login"
	}
	token.ExpireAt = time.Now().Add(DefaultExpire)
	count := 9
	for count > 0 {
		token.Token = utils.RandStringRunes(64)
		db.Model(&Token{}).Where("token = ?", token.Token).Count(&count)
	}
}
