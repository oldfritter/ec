package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"ec/config"
	. "ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

// 参数: receiver_sn, content, level, theme_id, theme_type
func MessageUpload(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(User)
	level, _ := strconv.Atoi(params["level"])
	themeId, _ := strconv.Atoi(params["theme_id"])
	message := Message{
		ThemeId:    themeId,
		ThemeType:  params["theme_type"],
		SenderSn:   user.Sn,
		ReceiverSn: params["receiver_sn"],
		Content:    params["content"],
		Level:      level,
	}

	b, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	config.PublishToPubSubChannels(NotifyMessageWithRedis, &b)

	response := utils.SuccessResponse
	return c.JSON(http.StatusOK, response)
}
