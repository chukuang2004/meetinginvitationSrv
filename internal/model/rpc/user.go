package rpc

type LoginReq struct {
	Wxcode string `json:"wxCode" binding:"required"`
}

type LoginResp struct {
	Openid              string `json:"openid"`
	WXName              string `json:"wxName"`
	WXAvatar            string `json:"wxAvatar"`
	EnableMsgPush       bool   `json:"enableMsgPush"`
	EnableMeetingNotify bool   `json:"enableMeetingNotify"`
}

type NotifyReq struct {
	Openid              string `json:"openid" binding:"required"`
	EnableMsgPush       bool   `json:"enableMsgPush"`
	EnableMeetingNotify bool   `json:"enableMeetingNotify"`
}

type PushReq struct {
	Openid    string `json:"openid" binding:"required"`
	MeetingID string `json:"meetingID" binding:"required"`
}

type UserUpdateReq struct {
	Openid   string `json:"openid" binding:"required"`
	WXName   string `json:"wxName"`
	WXAvatar string `json:"wxAvatar"`
}
