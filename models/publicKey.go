package models

import (
	"encoding/json"
	"log"

	"ec/config"
)

type PublicKey struct {
	CommonModel
	Index   int    `json:"index"`
	UserId  int    `json:"user_id"`
	UserSn  string `json:"user_sn"`
	Content string `json:"content" gorm:"type:text"`
}

func (pk *PublicKey) AfterSave() {
	b, err := json.Marshal(pk)
	if err != nil {
		log.Println(err)
	}
	config.PublishToPubSubChannels(NotifyPublicKeyWithRedis, &b)
}
