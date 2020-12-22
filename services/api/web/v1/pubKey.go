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
	var re struct{ Max int64 }
	db.Model(&models.Contract{}).Where("user_id = ?", user.ID).Where("expired_at > NOW()").Select("Max(level) as max").Scan(&re)
	index, _ := strconv.Atoi(params["index"])
	if int64(index) > re.Max {
		return utils.BuildError("2004")
	}
	var publicKey models.PublicKey
	db.FirstOrInit(&publicKey, models.PublicKey{
		UserId: int(user.ID),
		Index:  index,
		UserSn: user.Sn,
	})
	publicKey.Content = params["content"]
	db.Save(&publicKey)
	db.DbCommit()
	response := utils.SuccessResponse
	response.Body = publicKey
	return c.JSON(http.StatusOK, response)
}
