package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"ec/config"
	. "ec/models"
	"ec/services/api/helpers"
)

const (
	messagePongWait = time.Minute * 100
)

// 获取消息
func MessageListen(e echo.Context) (err error) {
	c, err := helpers.InitWsConn(e, messagePongWait)
	defer c.Close()
	user := e.Get("current_user").(User)
	ctx, cancel := context.WithCancel(context.Background())
	err = config.ListenPubSubChannels(
		ctx,
		func() error {
			return nil
		},
		func(channel string, m *[]byte) error {
			var message Message
			if channel == NotifyMessageWithRedis {
				json.Unmarshal(*m, &message)
				if message.ReceiverSn != user.Sn {
					return nil
				}
				err := c.WriteMessage(websocket.TextMessage, *m)
				if err != nil {
					log.Println("sended: %s", message)
					cancel()
				}
			}
			return nil
		},
		NotifyMessageWithRedis,
	)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// 发送消息
func MessageUpload(e echo.Context) (err error) {
	c, err := helpers.InitWsConn(e, messagePongWait)
	defer c.Close()
	user := e.Get("current_user").(User)
	for {
		_, m, e := c.ReadMessage()
		if e != nil {
			err = e
			log.Println(err)
			break
		}
		var message Message
		json.Unmarshal(m, &message)
		if message.ReceiverSn == "" {
			err = fmt.Errorf("no receiver sn")
			return
		}
		message.SenderSn = user.Sn
		b, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}
		config.PublishToPubSubChannels(NotifyMessageWithRedis, &b)
		log.Println("sended: %s", message)
	}
	return
}
