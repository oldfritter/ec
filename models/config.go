package models

import (
	"github.com/jinzhu/gorm"

	"ec/config"
)

type Config struct {
	gorm.Model
	Description string
	Key         string
	Value       string
}

func InitConfigInDB(db *GormDB) {
	var configs []Config
	db.Find(&configs)
	config.Env.ConfigInDB = map[string]string{}
	for _, c := range configs {
		config.Env.ConfigInDB[c.Key] = c.Value
	}
}
