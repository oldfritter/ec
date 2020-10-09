package models

import (
	"github.com/jinzhu/gorm"
)

type Group struct {
	gorm.Model
	UserId       int  `json:"-"`
	User         User `gorm:"ForeignKey:UserId"`
	MaxLimit     int  `gorm:"default:5"`
	Name         string
	Desc         string
	GroupMembers []*GroupMember `gorm:"ForeignKey:GroupId"`
}
