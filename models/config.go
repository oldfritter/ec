package models

import (
	"ec/config"
)

type Config struct {
	CommonModel
	Description string `json:"description"`
	Key         string `json:"key"`
	Value       string `json:"value"`
}

func InitConfigInDB(db *GormDB) {
	var configs []Config
	db.Find(&configs)
	config.Env.ConfigInDB = map[string]string{}
	for _, c := range configs {
		config.Env.ConfigInDB[c.Key] = c.Value
	}
}
