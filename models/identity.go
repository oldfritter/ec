package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Identity struct {
	gorm.Model
	UserId      int       `json:"-"`
	User        User      `gorm:"ForeignKey:UserId" jons:"-"`
	Source      string    `gorm:"type:varchar(32)"` // Email or Phone, Wechat, Alipay
	Symbol      string    `gorm:"type:varchar(64)"` // Email address or Phone number, openid, uid
	AccessToken string    `gorm:"type:varchar(64)"` // 授权token
	ExpiredAt   time.Time `gorm:"default:null"`     // 过期时间
}
