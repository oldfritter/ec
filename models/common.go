package models

import (
	"time"
)

const (
	NotifyMessageWithRedis   = "notify:message:"
	NotifyPublicKeyWithRedis = "notify:public_key:"
)

type CommonModel struct {
	Id        int       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
