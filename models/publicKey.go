package models

import (
	"encoding/json"
	"log"

	"ec/config"
)

type PublicKey struct {
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
