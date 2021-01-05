package models

import (
	"github.com/jinzhu/gorm"
)

type FriendShip struct {
	gorm.Model
	UserId   uint
	FriendId uint
	State    int    `gorm:"default:0"`    // 状态
	MarkName string `gorm:"default:null"` // 备注
}
