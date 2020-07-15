package routes

import (
	"github.com/labstack/echo/v4"

	v1 "demo/api/ws/v1"
)

func SetWsInterfaces(e *echo.Echo) {

	e.GET("/api/ws/v1/message/listen", v1.MessageListen)
	e.GET("/api/ws/v1/message/upload", v1.MessageUpload)

	e.GET("/api/ws/v1/public_key/listen", v1.PublicKeyListen)
	e.GET("/api/ws/v1/public_key/upload", v1.PublicKeyUpload)

}
