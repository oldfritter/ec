package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"demo/api/helpers"
	. "demo/models"
	"demo/utils"
)

func UserInfo(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	db := utils.MainDb
	var user User
	if db.Where("sn = ?", params["sn"]).First(&user).RecordNotFound() {
		return utils.BuildError("2001")
	}
	response := utils.SuccessResponse
	response.Body = user
	return c.JSON(http.StatusOK, response)
}
