package DB

import "gorm.io/gorm"

type User struct {
	gorm.Model
	OpenID              string `gorm:"primary_key"`
	UnionID             string
	Name                string
	Phone               string
	WXName              string
	WXAvatar            string
	EnableMsgPush       bool
	EnableMeetingNotify bool
}

func (*User) TableName() string {
	return "user"
}
