package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"ec/config"
	"ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

// 参数: receiver_sn, content, level, theme_id, theme_type
func MessageUpload(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(models.User)
	level, _ := strconv.Atoi(params["level"])
	themeId, _ := strconv.Atoi(params["theme_id"])
	message := models.Message{
		ThemeId:    themeId,
		ThemeType:  params["theme_type"],
		SenderSn:   user.Sn,
		ReceiverSn: params["receiver_sn"],
		Content:    params["content"],
		Level:      level,
	}
	message.Version, _ = strconv.Atoi(params["version"])

	b, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	config.PublishToPubSubChannels(models.NotifyMessageWithRedis, &b)

	config.MainDb.Save(&message)
	message.DelayNotify(config.GetRedisConn("data"))

	response := utils.SuccessResponse
	return c.JSON(http.StatusOK, response)
}

func MessageList(c echo.Context) (err error) {
	user := c.Get("current_user").(models.User)
	response := utils.SuccessResponse
	response.Body = user.NotifyMessages(config.GetRedisConn("data"))
	return c.JSON(http.StatusOK, response)
}

func MessageRead(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(models.User)
	user.ReadMessage(config.GetRedisConn("data"), params["id"])
	response := utils.SuccessResponse
	return c.JSON(http.StatusOK, response)
}
