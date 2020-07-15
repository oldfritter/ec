package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var Env struct {
	Model    string `yaml:"model"`
	Newrelic struct {
		AppName    string `yaml:"app_name"`
		LicenseKey string `yaml:"license_key"`
	} `yaml:"newrelic"`

	// 经常修改的配置放在数据库，每次修改后向消息队列发送一条消息
	ConfigInDB map[string]string
}

func InitEnv() {
	path_str, _ := filepath.Abs("config/env.yml")
	content, err := ioutil.ReadFile(path_str)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = yaml.Unmarshal(content, &Env)
	if err != nil {
		log.Fatal(err)
	}
}
