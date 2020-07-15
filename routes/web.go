package routes

import (
	"github.com/labstack/echo/v4"

	v1 "demo/api/web/v1"
)

func SetWebInterfaces(e *echo.Echo) {

	e.POST("/api/web/v1/message/upload", v1.MessageUpload)

	e.POST("/api/web/v1/pub_key/upload", v1.PubKeyUpload)

	e.POST("/api/web/v1/user/info", v1.UserInfo)
}
