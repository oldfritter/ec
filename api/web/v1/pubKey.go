package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"demo/api/helpers"
	. "demo/models"
	"demo/utils"
)

// 参数: index, content
func PubKeyUpload(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := utils.MainDbBegin()
	defer db.DbRollback()
	user := c.Get("current_user").(User)
	var publicKey PublicKey
	if db.Where("`index` = ?", params["index"]).Where("user_sn = ?", user.Sn).First(&publicKey).RecordNotFound() {
		publicKey.Index, _ = strconv.Atoi(params["index"])
		publicKey.UserSn = user.Sn
	}
	publicKey.Content = params["content"]
	db.Save(&publicKey)
	db.DbCommit()
	response := utils.SuccessResponse
	response.Body = publicKey
	return c.JSON(http.StatusOK, response)
}
