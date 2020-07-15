package models

import (
	"encoding/json"
	"log"

	"github.com/jinzhu/gorm"

	"ec/config"
)

type PublicKey struct {
	gorm.Model
	Index   int
	UserId  int
	Version int `gorm:"default:0"`
	UserSn  string
	Content string `gorm:"type:text"`
}

func (pk *PublicKey) AfterSave() {
	b, err := json.Marshal(pk)
	if err != nil {
		log.Println(err)
	}
	config.PublishToPubSubChannels(NotifyPublicKeyWithRedis, &b)
}

func (pk *PublicKey) BeforeSave() {
	pk.Version = pk.Version + 1
}
