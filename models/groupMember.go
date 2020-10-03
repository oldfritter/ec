package models

type GroupMember struct {
	GroupId  int   `json:"-"`
	Group    Group `gorm:"ForeignKey:GroupId" json:"-"`
	MemberId int   `json:"-"`
	Member   User  `gorm:"ForeignKey:MemberId"`
}
