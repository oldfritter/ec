package models

type FriendShip struct {
	UserId   int
	FriendId int
	MarkName string
	State    int `gorm:"default:null"` // 状态
}
