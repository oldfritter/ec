package helpers

import (
	"log"
	"net/http"
	"strconv"
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
	c.SetPingHandler(func(message string) error {
		err := c.WriteMessage(websocket.TextMessage, []byte(strconv.FormatInt(time.Now().UnixNano()/1000000, 10)))
		if err != nil {
			log.Println("sended ping err: ", err)
		}
		return nil
	})
	return
}
