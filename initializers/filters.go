package initializers

import (
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"

	"ec/config"
	"ec/initializers/locale"
	. "ec/models"
)

func LimitTrafficWithIp(context echo.Context) bool {
	dataRedis := config.GetRedisConn("data")
	defer dataRedis.Close()
	key := "limit-traffic-with-ip:" + context.Path() + ":" + context.RealIP()
	timesStr, _ := redis.String(dataRedis.Do("GET", key))
	if timesStr == "" {
		dataRedis.Do("SETEX", key, 1, 60)
	} else {
		times, _ := strconv.Atoi(timesStr)
		if times > 10 {
			return false
		} else {
			dataRedis.Do("INCR", key)
		}
	}
	return true
}

func treatLanguage(context echo.Context) {
	var language string
	var lqs []locale.LangQ
	if context.QueryParam("lang") != "" {
		lqs = locale.ParseAcceptLanguage(context.QueryParam("lang"))
	} else {
		lqs = locale.ParseAcceptLanguage(context.Request().Header.Get("Accept-Language"))
	}
	if lqs[0].Lang == "en" {
		language = "en"
	} else if lqs[0].Lang == "ja" {
		language = "ja"
	} else if lqs[0].Lang == "ko" {
		language = "ko"
	} else {
		language = "zh-CN"
	}
	context.Set("language", language)
}

func checkTimestamp(context echo.Context, params map[string]string) bool {
	timestamp, _ := strconv.Atoi(params["timestamp"])
	now := time.Now()
	if int(now.Add(-time.Second*60*5).Unix()) <= timestamp && timestamp <= int(now.Add(time.Second*60*5).Unix()) {
		return true
	}
	return false
}

func checkSign(context echo.Context, params map[string]string) (allow bool) {
	mainDB := MainDbBegin()
	defer mainDB.DbRollback()
	return
}

// TODO: IP检测
func checkIP(context echo.Context, params map[string]string) bool {
	return
}
