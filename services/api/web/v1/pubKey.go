package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

// 参数: index, content
func PubKeyUpload(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := models.MainDbBegin()
	defer db.DbRollback()
	user := c.Get("current_user").(models.User)
	var publicKey models.PublicKey
	if db.First(&publicKey, map[string]interface{}{"index": params["index"], "user_sn": user.Sn}).RecordNotFound() {
		publicKey.UserId = int(user.ID)
		publicKey.Content = params["content"]
		db.Save(&publicKey)
	}
	index, _ := strconv.Atoi(params["index"])
	db.FirstOrInit(&publicKey, models.PublicKey{
		Index:   index,
		UserSn:  user.Sn,
		Content: params["content"],
	})
	db.Save(&publicKey)
	db.DbCommit()
	response := utils.SuccessResponse
	response.Body = publicKey
	return c.JSON(http.StatusOK, response)
}
