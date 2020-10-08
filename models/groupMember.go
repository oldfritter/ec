package models

type GroupMember struct {
	GroupId  int    `json:"-"`
	MemberId int    `json:"-"`
	MarkName string `json:"mark_name"`

	Group  Group `gorm:"ForeignKey:GroupId" json:"-"`
	Member User  `gorm:"ForeignKey:MemberId" json:"member"`
}
