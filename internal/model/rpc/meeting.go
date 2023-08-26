package rpc

import (
	"beastpark/meetinginvitationservice/pkg/utils"
)

type QueryReq struct {
	Openid    string `form:"openid" binding:"required"`
	MeetingID string `form:"meetingID" binding:"required"`
}

type QueryResp struct {
	Meeting *MeetingInfo `json:"meeting"`
}

type CreateReq struct {
	Openid    string           `json:"openid"`
	Name      string           `json:"name"`
	Place     string           `json:"place"`
	Meeting   string           `json:"meeting"`
	Poster    string           `json:"poster"`
	StartTime utils.StringTime `json:"startTime"`
	EndTime   utils.StringTime `json:"endTime"`
}

type CreateResp struct {
	ID string `json:"ID"`
}

type DeleteReq struct {
	Openid    string `form:"openid" binding:"required"`
	MeetingID string `form:"meetingID" binding:"required"`
}

type UpdateReq struct {
	Openid    string           `json:"openid" binding:"required"`
	MeetingID string           `form:"meetingID" binding:"required"`
	Name      string           `json:"name"`
	Place     string           `json:"place"`
	Meeting   string           `json:"meeting"`
	Poster    string           `json:"poster"`
	StartTime utils.StringTime `json:"startTime"`
	EndTime   utils.StringTime `json:"endTime"`
}

type FollowingReq struct {
	Openid string `form:"openid" binding:"required"`
}

type MeetingInfo struct {
	ID        string `json:"ID"`
	Name      string `json:"name"`
	Place     string `json:"place"`
	IsOwner   bool   `json:"isOwner"`
	Poster    string `json:"poster"`
	Thumbnail string `json:"thumbnail"`
	Meeting   string `json:"meeting"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
type FollowingResp struct {
	Meeting []*MeetingInfo `json:"meeting"`
}

type HistoryReq struct {
	Openid string `form:"openid" binding:"required"`
}

type HistoryResp struct {
	Meeting []*MeetingInfo `json:"meeting"`
}

type QR4inviteReq struct {
	Openid    string `form:"openid" binding:"required"`
	MeetingID string `form:"meetingID" binding:"required"`
}

type QR4inviteResp struct {
	QRCode string `json:"QRCode"`
}

type ViewReq struct {
	Openid    string `json:"openid" binding:"required"`
	MeetingID string `form:"meetingID" binding:"required"`
}

type SignUpReq struct {
	Openid     string `json:"openid" binding:"required"`
	MeetingID  string `form:"meetingID" binding:"required"`
	Name       string `json:"name"`
	PhoneNum   int    `json:"phoneNum"`
	PartnerNum int    `json:"partnerNum"`
}

type QuitReq struct {
	Openid    string `json:"openid" binding:"required"`
	MeetingID string `form:"meetingID" binding:"required"`
}

type ReportReq struct {
	Openid    string `form:"openid" binding:"required"`
	MeetingID string `form:"meetingID" binding:"required"`
}

type ReportResp struct {
	View   int `json:"View"`
	SignUp int `json:"signUp"`
	SignIn int `json:"signIn"`
}

type QR4SignInReq struct {
	Openid    string `form:"openid" binding:"required"`
	MeetingID string `form:"meetingID" binding:"required"`
}

type QR4SignInResp struct {
	QRCode string `json:"QRCode"`
}

type SignInReq struct {
	Openid    string `form:"openid" binding:"required"`
	MeetingID string `form:"meetingID" binding:"required"`
}
