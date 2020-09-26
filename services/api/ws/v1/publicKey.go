package v1

import (
	"context"
	"encoding/json"
	// "fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	. "ec/models"
	"ec/services/api/helpers"
	"ec/utils"
)

const (
	keyPongWait = time.Minute * 10
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
	err = utils.ListenPubSubChannels(
		ctx,
		func() error {
			return nil
		},
		func(channel string, m *[]byte) error {
			var publicKey PublicKey
			if channel == NotifyPublicKeyWithRedis {
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
					log.Println("refresh public key: %s", publicKey)
					cancel()
				}
			}
			return nil
		},
		NotifyPublicKeyWithRedis,
	)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// 上传新的公钥
func PublicKeyUpload(e echo.Context) (err error) {
	c, err := helpers.InitWsConn(e, keyPongWait)
	defer c.Close()
	user := e.Get("current_user").(User)
	for {
		_, m, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		var publicKey PublicKey
		json.Unmarshal(m, &publicKey)
		publicKey.UserSn = user.Sn
		b, err := json.Marshal(publicKey)
		if err != nil {
			log.Println(err)
		}
		utils.PublishToPubSubChannels(NotifyPublicKeyWithRedis, &b)
		log.Println("sended public key: %s", publicKey)
		db := utils.MainDbBegin()
		defer db.DbRollback()
		db.Save(&publicKey)
		db.DbCommit()
	}
	return
}
