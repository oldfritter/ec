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

func ParsePong(c *websocket.Conn, timestamp *string, wait time.Duration) (err error) {
	_, m, err := c.ReadMessage()
	if err != nil {
		log.Println("parse pong err: ", err)
		return
	}
	if string(m) == *timestamp {
		log.Println("parse pong : ", string(m))
		c.SetWriteDeadline(time.Now().Add(wait))
		err = ParsePong(c, timestamp, wait)
	}
	return
}

func SendPing(c *websocket.Conn, timestamp *string) (err error) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			*timestamp = strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
			ping := c.PingHandler()
			err := ping(*timestamp)
			if err != nil {
				log.Println("sended message ping err: ", err)
			}
		}
	}
	return
}
