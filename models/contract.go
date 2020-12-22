package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Contract struct {
	gorm.Model
	UserId    int
	Level     int `gorm:"default:1"`
	ExpiredAt time.Time
}
