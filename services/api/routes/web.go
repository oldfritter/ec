package routes

import (
	"github.com/labstack/echo/v4"

	v1 "ec/services/api/web/v1"
)

func SetWebInterfaces(e *echo.Echo) {

	e.POST("/api/web/v1/message/upload", v1.MessageUpload)

	e.POST("/api/web/v1/pub_key/upload", v1.PubKeyUpload)

	e.GET("/api/web/v1/user/me", v1.UserMe)
	e.POST("/api/web/v1/user/me", v1.UserMe)
	e.GET("/api/web/v1/user/info", v1.UserInfo)
	e.POST("/api/web/v1/user/info", v1.UserInfo)
	e.POST("/api/web/v1/user/login", v1.UserLogin)
	e.POST("/api/web/v1/user/register", v1.UserRegister)

	e.GET("/api/web/v1/cloud_storage/upload/auth", v1.CloudStorageUploadAuth)
}
