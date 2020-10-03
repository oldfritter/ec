package models

type Group struct {
	CommonModel
	OwnerId  int    `json:"-"`
	Owner    User   `gorm:"ForeignKey:OwnerId" json:"owner"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	MaxLimit int    `gorm:"default:5"`

	GroupMembers []*GroupMember `gorm:"ForeignKey:GroupId" json:"group_members"`
}
