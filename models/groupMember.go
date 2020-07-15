package models

type GroupMember struct {
	GroupId  int `json:"-"`
	UserId   int `json:"-"`
	MarkName string

	Group Group `gorm:"ForeignKey:GroupId" json:"-"`
	User  User  `gorm:"ForeignKey:UserId"`
}
