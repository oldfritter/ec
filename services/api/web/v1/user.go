package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"ec/config"
	. "ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

func UserInfo(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := config.MainDb
	var user User
	if db.Where("sn = ?", params["sn"]).First(&user).RecordNotFound() {
		return utils.BuildError("2001")
	}
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}
