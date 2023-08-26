package rpc

import (
	"beastpark/meetinginvitationservice/internal/model/DB"
)

type TemplatesReq struct {
	Openid string `form:"openid" binding:"required"`
}

type Template struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type TemplatesResp struct {
	Templates []*Template `json:"templates"`
}

type UploadTemplatesReq struct {
	Openid string               `form:"openid" binding:"required"`
	Name   string               `json:"name"`
	Data   []*DB.InvitationPage `json:"data"`
}

type DeleteTemplatesReq struct {
	Openid     string `form:"openid" binding:"required"`
	TemplateID int    `form:"templateID" binding:"required"`
}
