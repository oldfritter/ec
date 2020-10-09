package initializers

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"

	"ec/config"
	"ec/models"
	// "ec/utils"
)

func initDb(name string) *gorm.DB {
	dbConfig := config.GetDatabaseConfig()
	connString := config.GetConnectionString(dbConfig, name)
	db, err := gorm.Open("mysql", connString)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(dbConfig.GetInt(config.Env.Model+"."+name+".pool", 5))
	db.DB().SetMaxOpenConns(dbConfig.GetInt(config.Env.Model+"."+name+".maxopen", 0))
	du, _ := time.ParseDuration(dbConfig.Get(config.Env.Model+"."+name+".timeout", "3600") + "s")
	db.DB().SetConnMaxLifetime(du)
	db.Exec("set transaction isolation level repeatable read")
	// file := utils.GetLogFile("gorm")
	// db.SetLogger(gorm.Logger{LogWriter: log.New(file, "\r\n", 0)})
	db.LogMode(true)
	return db
}

func InitMainDB() {
	config.MainDb = initDb("main")

	models.MainMigrations()
}

func InitLogDB() {
	config.LogDb = initDb("log")

	models.LogMigrations()
}

func CloseMainDB() {
	err := config.MainDb.Close()
	if err != nil {
		log.Println(err)
	}
}

func CloseLogDB() {
	err := config.LogDb.Close()
	if err != nil {
		log.Println(err)
	}
}
