package models

import (
	envConfig "ec/config"
	"ec/utils"
)

type Config struct {
	CommonModel
	Description string `json:"description"`
	Key         string `json:"key"`
	Value       string `json:"value"`
}

func InitConfigInDB(db *utils.GormDB) {
	var configs []Config
	db.Find(&configs)
	envConfig.Env.ConfigInDB = map[string]string{}
	for _, config := range configs {
		envConfig.Env.ConfigInDB[config.Key] = config.Value
	}
}
