package service

import (
	"beastpark/meetinginvitationservice/internal/model/rpc"
	"beastpark/meetinginvitationservice/pkg/config"
	"beastpark/meetinginvitationservice/pkg/database"
	"beastpark/meetinginvitationservice/pkg/http"
	"beastpark/meetinginvitationservice/pkg/log"

	"github.com/gin-gonic/gin"
)

type Other struct {
	db  *database.DB
	log *log.Logger
}

func NewOther() *Other {

	return &Other{
		db:  database.GetInstance(),
		log: log.GetInstance(),
	}
}

func (o *Other) UploadFile(ctx *gin.Context) {

	req := &rpc.UploadFileReq{}

	req.Openid = ctx.PostForm("openid")
	if !CheckOpenIdValid(ctx, req.Openid) {
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		o.log.Errorf("read file error %s", err.Error())
		http.WithResponse(ctx, http.InvalidParameter, nil)
		return
	}

	conf, ok := ctx.MustGet("config").(*config.Config)
	if !ok {
		o.log.Errorf("get config fail")
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	dst := conf.Http.SavePath + file.Filename
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		o.log.Errorf("save file error %s", err.Error())
		http.WithResponse(ctx, http.UnknownError, nil)
		return
	}

	resp := &rpc.UploadFileResp{
		Url: "https://" + conf.Http.DNS + "/download/" + file.Filename,
	}

	http.WithResponse(ctx, http.Success, resp)
}
