package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	"ec/utils"
)

type User struct {
	gorm.Model
	Sn             string `gorm:"type:varchar(16)"`
	PasswordDigest string `gorm:"type:varchar(64)" json:"-"`
	Nickname       string `gorm:"type:varchar(32)"`
	State          int    `gorm:"default:null"`

	Tokens     []*Token     `gorm:"ForeignKey:UserId"`
	PublicKeys []*PublicKey `gorm:"ForeignKey:UserId"`
	Groups     []*Group     `gorm:"many2many:group_members"`
	Friends    []*User      `gorm:"many2many:friend_ships"`

	Password string `sql:"-" json:"-"`
}

func (user *User) BeforeCreate(db *gorm.DB) {
	if user.Password != "" {
		user.SetPasswordDigest()
	}
	count := 4
	for count > 0 {
		user.Sn = "DE" + utils.RandStringRunes(10) + "MO"
		db.Model(&User{}).Where("sn = ?", user.Sn).Count(&count)
	}
}

func (user *User) AfterFind(db *gorm.DB) {
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
