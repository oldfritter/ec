package models

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	ThemeId    int
	ThemeType  string `gorm:"size:16"`
	SenderSn   string `gorm:"size32"`
	ReceiverSn string `gorm:"size:32"`
	Content    string `gorm:"type:text"`
	Level      int    // 多层加密，标明层数
	Version    int
}

func (self *Message) receiverNotifyKey() string {
	return fmt.Sprintf("ec:message:notify:%v", self.ReceiverSn)
}

func (self *Message) DelayNotify(redisCon redis.Conn) {
	redisCon.Do("RPUSH", self.receiverNotifyKey(), self.ID)
	redisCon.Do("EXPIRE", self.receiverNotifyKey(), 3600*24*30)
}
