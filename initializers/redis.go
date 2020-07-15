package initializers

import (
	"ec/config"
)

func InitRedisPools() {
	config.DatePool = config.NewRedisPool("data")
	config.LimitPool = config.NewRedisPool("limit")
	config.PublishPool = config.NewRedisPool("publish")
}

func CloseRedisPools() {
	config.DatePool.Close()
	config.LimitPool.Close()
	config.PublishPool.Close()
}
