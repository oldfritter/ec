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
