package models

type FriendShip struct {
	OwnerId  int
	Owner    User `gorm:"ForeignKey:OwnerId"`
	FriendId int
	Friend   User `gorm:"ForeignKey:FriendId"`
}
