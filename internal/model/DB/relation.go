package DB

import (
	"beastpark/meetinginvitationservice/pkg/utils"

	"gorm.io/gorm"
)

type Relation struct {
	gorm.Model
	OpenID     string `gorm:"primary_key"`
	MeetingID  string
	Action     string //signUp signIn create view
	Name       string
	Phone      int
	PartnerNum int
	EndTime    utils.StringTime
	// IsNotify   bool
}

func (*Relation) TableName() string {
	return "relation"
}

func (*Relation) GetUserSignUp(meetingID string) ([]string, error) {

	return nil, nil
}
