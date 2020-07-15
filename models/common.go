package models

import (
	"time"

	"demo/utils"
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

func AutoMigrations() {
	mainDB := utils.MainDbBegin()
	defer mainDB.DbRollback()

	// config
	mainDB.AutoMigrate(&Config{})

	// device
	mainDB.AutoMigrate(&Device{})

	// identity
	mainDB.AutoMigrate(&Identity{})
	mainDB.Model(&Identity{}).AddUniqueIndex("index_identity_on_source_and_symbol", "source", "symbol")

	// public_key
	mainDB.AutoMigrate(&PublicKey{})
	mainDB.Model(&PublicKey{}).AddUniqueIndex("index_public_key_on_user_sn_and_index", "user_sn", "index")

	// token
	mainDB.AutoMigrate(&Token{})

	// user
	mainDB.AutoMigrate(&User{})

}
