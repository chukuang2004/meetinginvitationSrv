package DB

import (
	"beastpark/meetinginvitationservice/pkg/utils"

	"gorm.io/gorm"
)

type Meeting struct {
	gorm.Model
	ID        string `gorm:"primary_key"`
	CreaterID string
	Name      string
	Place     string
	StartTime utils.StringTime
	EndTime   utils.StringTime
	Data      string `gorm:"size:1024*1024"`
	Poster    string `gorm:"size:1024"`
	Thumbnail string `gorm:"size:1024*1024"`
	QR4Invite string `gorm:"size:1024*1024"`
	QR4SignIn string `gorm:"size:1024*1024"`
}

func (*Meeting) TableName() string {
	return "meeting"
}
