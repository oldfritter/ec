package helpers

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func InitWsConn(e echo.Context, wait time.Duration) (c *websocket.Conn, err error) {
	upgrader := websocket.Upgrader{}
	c, err = upgrader.Upgrade(e.Response(), e.Request(), nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	c.SetWriteDeadline(time.Now().Add(wait))
	c.SetPongHandler(func(message string) error {
		c.SetReadDeadline(time.Now().Add(wait))
		return nil
	})
	return
}
