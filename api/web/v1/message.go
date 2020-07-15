package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"demo/api/helpers"
	. "demo/models"
	"demo/utils"
)

// 参数: receiver_sn, content, level
func MessageUpload(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	user := c.Get("current_user").(User)
	level, _ := strconv.Atoi(params["level"])
	message := Message{
		SenderSn:   user.Sn,
		ReceiverSn: params["receiver_sn"],
		Content:    params["content"],
		Level:      level,
	}

	b, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	utils.PublishToPubSubChannels(NotifyMessageWithRedis, &b)

	response := utils.SuccessResponse
	return c.JSON(http.StatusOK, response)
}
