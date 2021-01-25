package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"ec/config"
	"ec/models"
	"ec/services/api/helpers"
)

const (
	messagePongWait = time.Minute
)

// 获取消息
func MessageListen(e echo.Context) (err error) {
	c, err := helpers.InitWsConn(e, messagePongWait)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	user := e.Get("current_user").(models.User)
	ctx, cancel := context.WithCancel(context.Background())
	var timestamp string

	// 定时发出ping
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				timestamp = strconv.FormatInt(time.Now().UnixNano()/1000000, 10)
				ping := c.PingHandler()
				err := ping(timestamp)
				if err != nil {
					log.Println("sended message ping err: ", err)
					cancel()
				}
			}
		}
	}()

	// 读取pong
	go func() {
		for {
			_, m, err := c.ReadMessage()
			if err != nil {
				log.Println("parse message pong err: ", err)
				cancel()
			}
			if string(m) == timestamp {
				c.SetWriteDeadline(time.Now().Add(messagePongWait))
			}
		}
	}()

	err = config.ListenPubSubChannels(
		ctx,
		func() error {
			return nil
		},
		func(channel string, m *[]byte) error {
			var message models.Message
			if channel == models.NotifyMessageWithRedis {
				json.Unmarshal(*m, &message)
				if message.ReceiverSn != user.Sn {
					return nil
				}
				err := c.WriteMessage(websocket.TextMessage, *m)
				if err != nil {
					log.Println("sended: ", message)
					cancel()
				}
			}
			return nil
		},
		models.NotifyMessageWithRedis,
	)
	if err != nil {
		log.Println(err)
	}
	return
}

// 发送消息,暂时未启用
func MessageUpload(e echo.Context) (err error) {
	c, err := helpers.InitWsConn(e, messagePongWait)
	defer c.Close()
	user := e.Get("current_user").(models.User)
	for {
		_, m, e := c.ReadMessage()
		if e != nil {
			err = e
			log.Println(err)
			break
		}
		var message models.Message
		json.Unmarshal(m, &message)
		if message.ReceiverSn == "" {
			err = fmt.Errorf("no receiver sn")
			return
		}
		if message.Content == "" {
			err = fmt.Errorf("no content")
			return
		}
		message.SenderSn = user.Sn
		b, err := json.Marshal(message)
		if err != nil {
			log.Println(err)
		}
		config.PublishToPubSubChannels(models.NotifyMessageWithRedis, &b)
		log.Println("sended: ", message)
	}
	return
}
