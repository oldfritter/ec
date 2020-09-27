package config

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

var (
	LogDb  *gorm.DB
	MainDb *gorm.DB
)

func GetConnectionString(config *ConfigEnv, name string) string {
	host := config.Get(Env.Model+"."+name+".host", "")
	port := config.Get(Env.Model+"."+name+".port", "3306")
	user := config.Get(Env.Model+"."+name+".username", "")
	pass := config.Get(Env.Model+"."+name+".password", "")
	dbname := config.Get(Env.Model+"."+name+".database", "")
	protocol := config.Get(Env.Model+"."+name+".protocol", "tcp")
	dbargs := config.Get(Env.Model+"."+name+".dbargs", " ")
	if strings.Trim(dbargs, " ") != "" {
		dbargs = "?" + dbargs
	} else {
		dbargs = ""
	}
	return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s", user, pass, protocol, host, port, dbname, dbargs)
}
