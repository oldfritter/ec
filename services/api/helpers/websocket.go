package helpers

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func InitWsConn(e echo.Context, wait time.Duration) (c *websocket.Conn, err error) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(request *http.Request) bool {
		// TODO: 检测请求的Origin
		return true
	}
	c, err = upgrader.Upgrade(e.Response(), e.Request(), nil)
	if err != nil {
		log.Println("upgrade err:", err)
		return
	}
	c.SetWriteDeadline(time.Now().Add(wait))
	c.SetPongHandler(func(message string) error {
		c.SetReadDeadline(time.Now().Add(wait))
		return nil
	})
	return
}
