package models

type FriendShip struct {
	OwnerId  int
	FriendId int
	MarkName string `json:"mark_name"`

	Owner  User `gorm:"ForeignKey:OwnerId"`
	Friend User `gorm:"ForeignKey:FriendId"`
}
