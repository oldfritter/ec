package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

// 参数: name 文件名
func CloudStorageUploadAuth(c echo.Context) (err error) {
	params := helpers.StringParams(c)
	cf := models.CloudFile{OriginName: params["name"]}
	cf.Initialize()
	response := utils.SuccessResponse
	s := cf.Attrs()
	response.Body = s
	return c.JSON(http.StatusOK, response)
}
