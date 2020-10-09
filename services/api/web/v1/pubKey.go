package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	. "ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

// 参数: index, content
func PubKeyUpload(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := MainDbBegin()
	defer db.DbRollback()
	user := c.Get("current_user").(User)
	var publicKey PublicKey
	if db.First(&publicKey, map[string]interface{}{"index": params["index"], "user_sn": user.Sn}).RecordNotFound() {
		publicKey.UserId = int(user.ID)
		publicKey.Content = strings.Replace(params["content"], "\n", "", -1)
		db.Save(&publicKey)
	}
	index, _ := strconv.Atoi(params["index"])
	db.Model(&publicKey).Updates(PublicKey{
		Index:   index,
		UserSn:  user.Sn,
		Content: strings.Replace(params["content"], "\n", "", -1),
	})
	db.DbCommit()
	response := utils.SuccessResponse
	response.Body = publicKey
	return c.JSON(http.StatusOK, response)
}
