package v1

import (
	"context"
	"encoding/json"
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
	keyPongWait = time.Minute
)

// 获取新的公钥
func PublicKeyListen(e echo.Context) (err error) {
	c, err := helpers.InitWsConn(e, keyPongWait)
	defer c.Close()
	var params struct {
		UserSns []string `json:"user_sns"`
	}
	_, m, err := c.ReadMessage()
	json.Unmarshal(m, &params)
	if len(params.UserSns) == 0 {
		return
	}

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
					log.Println("sended public-key ping err: ", err)
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
				log.Println("parse public-key pong err: ", err)
				cancel()
			}
			if string(m) == timestamp {
				c.SetWriteDeadline(time.Now().Add(keyPongWait))
			}
		}
	}()

	err = config.ListenPubSubChannels(
		ctx,
		func() error {
			return nil
		},
		func(channel string, m *[]byte) error {
			var publicKey models.PublicKey
			if channel == models.NotifyPublicKeyWithRedis {
				json.Unmarshal(*m, &publicKey)
				var in bool
				for _, sn := range params.UserSns {
					if sn == publicKey.UserSn {
						in = true
					}
				}
				if !in {
					return nil
				}
				err := c.WriteMessage(websocket.TextMessage, *m)
				if err != nil {
					log.Println("refresh public key: ", publicKey)
					cancel()
				}
			}
			return nil
		},
		models.NotifyPublicKeyWithRedis,
	)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// 上传新的公钥,暂时未启用
func PublicKeyUpload(e echo.Context) (err error) {
	c, err := helpers.InitWsConn(e, keyPongWait)
	defer c.Close()
	user := e.Get("current_user").(models.User)
	for {
		_, m, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var publicKey models.PublicKey
		json.Unmarshal(m, &publicKey)
		publicKey.UserSn = user.Sn
		b, err := json.Marshal(publicKey)
		if err != nil {
			log.Println(err)
		}
		config.PublishToPubSubChannels(models.NotifyPublicKeyWithRedis, &b)
		log.Println("sended public key: ", publicKey)
		db := models.MainDbBegin()
		defer db.DbRollback()
		db.Save(&publicKey)
		db.DbCommit()
	}
	return
}
