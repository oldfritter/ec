package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"ec/utils"
)

type User struct {
	CommonModel
	Sn             string `gorm:"type:varchar(16)" json:"sn"`       // 唯一编号
	PasswordDigest string `gorm:"type:varchar(64)" json:"-"`        // 经加密的密码
	Nickname       string `gorm:"type:varchar(32)" json:"nickname"` // 昵称
	State          int    `gorm:"default:null" json:"state"`        // 状态
	// GroupId int

	Tokens     []*Token     `gorm:"ForeignKey:UserId" json:"tokens"`
	PublicKeys []*PublicKey `gorm:"ForeignKey:UserId" json:"public_keys"`

	Password string `sql:"-" json:"-"`
}

func (user *User) BeforeCreate(db *gorm.DB) {
	count := 4
	for count > 0 {
		user.Sn = "DE" + utils.RandStringRunes(10) + "MO"
		db.Model(&User{}).Where("sn = ?", user.Sn).Count(&count)
	}
}

func (user *User) AfterFind(db *gorm.DB) {
	user.setPublicKeys(db)
}

func (user *User) CompareHashAndPassword() bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(user.Password))
	if err == nil {
		return true
	}
	return false
}

func (user *User) SetPasswordDigest() {
	b, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.PasswordDigest = string(b[:])
}

func (user *User) setPublicKeys(db *gorm.DB) {
	db.Where("user_sn = ?", user.Sn).Order("`index`").Find(&user.PublicKeys)
}
