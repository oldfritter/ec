package models

import (
	"ec/utils"
	"github.com/jinzhu/gorm"
)

type Device struct {
	gorm.Model
	UserId    int    // 设备所属用户
	IsUsed    bool   // 是否已经使用
	PublicKey string // 公钥，用于公钥验签
	Token     string `gorm:"type:varchar(64)"` // 授权token
}

func (device *Device) InitializeToken() {
	device.Token = utils.RandStringRunes(64)
}

func (device *Device) BeforeCreate(db *gorm.DB) {
	var count int
	db.Model(&Device{}).Where("token = ?", device.Token).Count(&count)
	for count > 0 {
		device.Token = utils.RandStringRunes(64)
		db.Model(&Device{}).Where("token = ?", device.Token).Count(&count)
	}
}
